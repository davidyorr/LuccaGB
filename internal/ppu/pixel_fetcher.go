package ppu

type PixelFetcher struct {
	ppu                        *PPU
	state                      FetcherState
	counter                    uint16
	fetchedTileNumber          uint8
	fetchedTileDataLow         uint8
	fetchedTileDataHigh        uint8
	xPositionCounter           uint8
	pixelsToDiscard            uint8
	backgroundScrollingPenalty uint8
	isFetchingWindow           bool
	currentX                   uint8
	windowLineCounter          uint8
	isFirstFetchOfScanline     bool
	scanlineHadWindowPixels    bool
	wyEqualedLyDuringFrame     bool

	isFetchingSprite bool
	spriteIndex      uint8

	// background and window FIFO
	backgroundFifo PixelFifo
	// sprite (object) FIFO
	spriteFifo PixelFifo
}

type FetcherState int

const (
	StateGetTile FetcherState = iota
	StateGetTileDataLow
	StateGetTileDataHigh
	StateSleep
	StatePush
)

func newPixelFetcher(ppu *PPU) *PixelFetcher {
	pixelFetcher := &PixelFetcher{
		ppu: ppu,
	}

	return pixelFetcher
}

func (fetcher *PixelFetcher) prepareForScanline() {
	fetcher.state = StateGetTile
	fetcher.backgroundFifo.Reset()
	fetcher.spriteFifo.Reset()
	fetcher.isFirstFetchOfScanline = true
	fetcher.isFetchingSprite = false
	fetcher.isFetchingWindow = false
	fetcher.currentX = 0
	fetcher.counter = 0
	fetcher.xPositionCounter = 0

	// SCX
	fetcher.pixelsToDiscard = fetcher.ppu.scx % 8
	if fetcher.pixelsToDiscard == 0 {
		fetcher.backgroundScrollingPenalty = 0
	} else if fetcher.pixelsToDiscard < 5 {
		fetcher.backgroundScrollingPenalty = 4
	} else {
		fetcher.backgroundScrollingPenalty = 8
	}
}

func (fetcher *PixelFetcher) step() {
	if !fetcher.isFetchingSprite {
		fetcher.attemptToPushPixel()
	}
	fetcher.tick()
}

func (fetcher *PixelFetcher) tick() {
	// Only check for new sprites if we aren't already fetching a sprite
	if !fetcher.isFetchingSprite {
		// If the X-Position of any sprite in the sprite buffer is less than or
		// equal to the current Pixel-X-Position + 8, a sprite fetch is initiated.
		// See: https://ashiepaws.github.io/GBEDG/ppu/#sprite-fetching
		for i, oamIndex := range fetcher.ppu.spriteBuffer {
			// each sprite is 4 bytes long
			baseAddress := oamIndex * 4
			spriteX := fetcher.ppu.oam[baseAddress+1]

			if spriteX <= fetcher.currentX+8 {
				fetcher.isFetchingSprite = true
				fetcher.state = StateGetTile
				fetcher.spriteIndex = oamIndex
				fetcher.ppu.spriteBuffer = append(fetcher.ppu.spriteBuffer[:i], fetcher.ppu.spriteBuffer[i+1:]...)

				// reset background fetcher
				fetcher.counter = 0
				fetcher.fetchedTileNumber = 0
				fetcher.fetchedTileDataLow = 0
				fetcher.fetchedTileDataHigh = 0
				return
			}
		}
	}

	windowEnabled := (fetcher.ppu.lcdc>>5)&1 == 1
	if !fetcher.isFetchingWindow && (windowEnabled) && (fetcher.wyEqualedLyDuringFrame) && (fetcher.currentX+7 >= fetcher.ppu.wx) {
		fetcher.state = StateGetTile
		fetcher.backgroundFifo.Reset()
		fetcher.isFetchingWindow = true
		fetcher.scanlineHadWindowPixels = true
		fetcher.xPositionCounter = 0
		fetcher.counter = 0
		fetcher.pixelsToDiscard = 0
		fetcher.backgroundScrollingPenalty = 0
		return
	}

	fetcher.counter++

	switch fetcher.state {
	case StateGetTile:
		if fetcher.counter < 2 {
			return
		}
		fetcher.counter = 0

		if fetcher.isFetchingSprite {
			oamIndex := fetcher.spriteIndex
			// each sprite is 4 bytes, byte 2 is the tile index
			tileIndex := (oamIndex * 4) + 2
			fetcher.fetchedTileNumber = fetcher.ppu.oam[tileIndex]
		} else {
			var tileMapAreaStart uint16 = 0x9800
			var yTile uint8 = 0
			var xTile uint8 = fetcher.xPositionCounter

			// during window fetching we ignore the SCX and SCY values completely
			if fetcher.isFetchingWindow {
				if (fetcher.ppu.lcdc>>6)&1 == 1 {
					tileMapAreaStart = 0x9C00
				}
				yTile = (fetcher.windowLineCounter / 8)
			} else {
				if (fetcher.ppu.lcdc>>3)&1 == 1 {
					tileMapAreaStart = 0x9C00
				}
				yTile = ((fetcher.ppu.ly + fetcher.ppu.scy) & 0xFF) / 8
				xTile += (fetcher.ppu.scx / 8)
				xTile &= 0x1f // for wrap-around
			}

			address := (tileMapAreaStart - 0x8000) + ((uint16(yTile)*32)+uint16(xTile))&0x3FF
			fetcher.fetchedTileNumber = fetcher.ppu.videoRam[address]
		}

		fetcher.state = StateGetTileDataLow
	case StateGetTileDataLow:
		if fetcher.counter < 2 {
			return
		}
		fetcher.counter = 0

		fetcher.fetchedTileDataLow = fetcher.fetchTileData(0)
		fetcher.state = StateGetTileDataHigh
	case StateGetTileDataHigh:
		if fetcher.counter < 2 {
			return
		}
		fetcher.counter = 0

		fetcher.fetchedTileDataHigh = fetcher.fetchTileData(1)

		// Note: The first time the background fetcher completes this step on a scanline the status is fully reset and operation restarts at Step 1.
		// See: https://ashiepaws.github.io/GBEDG/ppu/#background-pixel-fetching
		if fetcher.isFirstFetchOfScanline {
			fetcher.isFirstFetchOfScanline = false
			fetcher.state = StateGetTile
		} else {
			fetcher.state = StatePush
		}
	case StatePush:
		if fetcher.counter < 2 {
			return
		}

		pushedToFifo := false

		if fetcher.isFetchingSprite {
			oamIndex := fetcher.spriteIndex
			spriteX := fetcher.ppu.oam[oamIndex*4+1]
			spriteFlags := fetcher.ppu.oam[oamIndex*4+3]
			spriteFlipX := (spriteFlags>>5)&1 == 1
			var tempBuffer [8]FIFO

			for i := 7; i >= 0; i-- {
				bit := i
				if !spriteFlipX {
					bit = 7 - i
				}
				lowBit := (fetcher.fetchedTileDataLow >> bit) & 1
				highBit := (fetcher.fetchedTileDataHigh >> bit) & 1
				color := (highBit << 1) | lowBit
				pixel := FIFO{
					colorId:            color,
					palette:            (spriteFlags >> 4) & 1,
					backgroundPriority: (spriteFlags >> 7) & 1,
				}
				tempBuffer[i] = pixel
			}

			// only pixels which are actually visible on the screen are loaded into the FIFO
			// pixels can only be loaded into FIFO slots if there is no pixel in the given slot already
			// See: https://ashiepaws.github.io/GBEDG/ppu/#sprite-fetching
			var pixelsToDiscard uint8 = 0
			if spriteX < 8 {
				pixelsToDiscard = 8 - spriteX
			}
			for i := 0; i <= 7; i++ {
				if i < int(pixelsToDiscard) {
					continue
				}

				fifoIndex := i - int(pixelsToDiscard)

				// if the FIFO is not big enough to hold this pixel we need to expand it
				if fifoIndex >= fetcher.spriteFifo.size {
					fetcher.spriteFifo.Push(tempBuffer[i])
				} else {
					// if the FIFO already has a pixel in this slot we only overwrite it if
					// 1. the existing pixel is transparent (color ID 0)
					// 2. the new pixel is not transparent (color ID not 0)
					slot := fetcher.spriteFifo.Peek(fifoIndex)
					if slot.colorId == 0 && tempBuffer[i].colorId != 0 {
						*slot = tempBuffer[i]
					}
				}
			}
			fetcher.state = StateGetTile
			fetcher.isFetchingSprite = false
			pushedToFifo = true
		} else {
			// Note: While fetching background pixels, this step is only executed if the background FIFO is fully empty.
			// If it is not, this step repeats every cycle until it succeeds.
			// See: https://ashiepaws.github.io/GBEDG/ppu/#background-pixel-fetching
			if fetcher.backgroundFifo.size == 0 {
				for i := 7; i >= 0; i-- {
					lowBit := (fetcher.fetchedTileDataLow >> i) & 1
					highBit := (fetcher.fetchedTileDataHigh >> i) & 1
					colorId := (highBit << 1) | lowBit
					pixel := FIFO{
						colorId:            colorId,
						palette:            0, // only applies to sprites
						backgroundPriority: 0, // only applies to sprites
					}
					fetcher.backgroundFifo.Push(pixel)
				}
				fetcher.state = StateGetTile
				fetcher.xPositionCounter++
				fetcher.isFetchingSprite = false
				pushedToFifo = true
			}
		}

		// We only reset the counter if we actually pushed to the FIFO.
		// By not resetting it, it allows us to retry pushing to the FIFO in the next tick.
		if pushedToFifo {
			fetcher.state = StateGetTile
			fetcher.counter = 0
		}
	}
}

func (fetcher *PixelFetcher) fetchTileData(offset uint16) uint8 {
	var address uint16
	if fetcher.isFetchingSprite {
		// sprites always use 8000 method
		address = 0x8000

		oamIndex := fetcher.spriteIndex
		spriteY := fetcher.ppu.oam[oamIndex*4]
		spriteFlags := fetcher.ppu.oam[oamIndex*4+3]
		spriteTileNumber := fetcher.fetchedTileNumber

		// determine which vertical row of the sprite we are on
		rowInSprite := (fetcher.ppu.ly + 16) - spriteY
		isTallSprite := (fetcher.ppu.lcdc>>2)&1 == 1
		flipY := (spriteFlags>>6)&1 == 1

		// See: https://ashiepaws.github.io/GBEDG/ppu/#lcdc2---sprite-size
		if isTallSprite {
			// handle y flipping
			if flipY {
				rowInSprite = 15 - rowInSprite
			}
			if rowInSprite < 8 {
				// the top tile, so force LSB to 0
				spriteTileNumber &= 0b1111_1110
			} else {
				// the bottom tile, so force LSB to 1
				spriteTileNumber |= 0x0000_0001
				rowInSprite -= 8
			}
		} else {
			// 8x8 Mode
			// Limit to 3 bits (0-7), so that if a tall sprite was picked during
			// Mode 2, we wrap correctly within the 8x8 tile.
			rowInSprite &= 0b0111
			// handle y flipping
			if flipY {
				rowInSprite = 7 - rowInSprite
			}
		}

		address += uint16(spriteTileNumber) * 16
		address += uint16(rowInSprite) * 2
	} else {
		// 8000 method
		if (fetcher.ppu.lcdc>>4)&1 == 1 {
			address = 0x8000
			address += uint16(fetcher.fetchedTileNumber) * 16
		} else
		// 8800 method
		{
			address = 0x9000
			address += uint16(int16(int8(fetcher.fetchedTileNumber))) * 16
		}
		// get the row offset
		if fetcher.isFetchingWindow {
			address += uint16(2 * (fetcher.windowLineCounter % 8))
		} else {
			address += uint16(2 * ((fetcher.ppu.ly + fetcher.ppu.scy) % 8))
		}
	}

	address += offset

	return fetcher.ppu.videoRam[address-0x8000]
}

func (fetcher *PixelFetcher) attemptToPushPixel() {
	if fetcher.currentX == 160 {
		fetcher.ppu.changeMode(HorizontalBlank)
		return
	}

	// do nothing if the FIFO is empty
	if fetcher.backgroundFifo.size == 0 {
		return
	}

	// handle background scrolling penalty: this occurs in chunks of 4 T-cycles (1 M-cycle)
	if fetcher.backgroundScrollingPenalty > 0 {
		fetcher.backgroundScrollingPenalty--

		if fetcher.pixelsToDiscard > 0 {
			fetcher.backgroundFifo.Pop()
			fetcher.pixelsToDiscard--
		}
		return
	}

	// otherwise add to the framebuffer
	backgroundPixel := fetcher.backgroundFifo.Pop()
	var color uint8

	// use the background pixel's color as the default
	colorId := backgroundPixel.colorId
	// if the background is disabled, we always use color ID 0
	bgEnabled := (fetcher.ppu.lcdc>>0)&1 == 1
	if !bgEnabled {
		colorId = 0
	}
	// color IDs are 2 bits, so we shift times 2, then mask 2 bits for the final color/shade
	color = (fetcher.ppu.bgp >> (colorId * 2)) & 0b11

	// See: https://ashiepaws.github.io/GBEDG/ppu/#pixel-mixing
	var spritePixel FIFO
	if fetcher.spriteFifo.size > 0 {
		spritePixel = fetcher.spriteFifo.Pop()
		spriteIsTransparent := spritePixel.colorId == 0
		backgroundHasPriority := spritePixel.backgroundPriority == 1 && backgroundPixel.colorId != 0
		objEnabled := (fetcher.ppu.lcdc>>1)&1 == 1

		if !spriteIsTransparent && !backgroundHasPriority && objEnabled {
			colorId := spritePixel.colorId
			if spritePixel.palette == 0 {
				color = (fetcher.ppu.obp0 >> (colorId * 2)) & 0b11
			} else if spritePixel.palette == 1 {
				color = (fetcher.ppu.obp1 >> (colorId * 2)) & 0b11
			}
		}
	}

	fetcher.ppu.frameBuffer[fetcher.ppu.ly][fetcher.currentX] = color
	fetcher.currentX++
}
