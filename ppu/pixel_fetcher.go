package ppu

import (
	"github.com/davidyorr/LuccaGB/logger"
)

type PixelFetcher struct {
	ppu                     *PPU
	state                   FetcherState
	counter                 uint16
	fetchedTileNumber       uint8
	fetchedTileDataLow      uint8
	fetchedTileDataHigh     uint8
	xPositionCounter        uint8
	pixelsToDiscard         uint8
	isFetchingWindow        bool
	currentX                uint8
	windowLineCounter       uint8
	isFirstFetchOfScanline  bool
	scanlineHadWindowPixels bool
	wyEqualedLyDuringFrame  bool

	isFetchingSprite bool
	spriteIndex      uint8

	// background and window FIFO
	backgroundFifo []FIFO
	// sprite (object) FIFO
	spriteFifo []FIFO
}

type FetcherState int

const (
	StateGetTile FetcherState = iota
	StateGetTileDataLow
	StateGetTileDataHigh
	StateSleep
	StatePush
)

type FIFO struct {
	// 4 possible colors
	colorId uint8
	// only applies to objects (sprites)
	palette            uint8
	backgroundPriority uint8
}

func newPixelFetcher(ppu *PPU) *PixelFetcher {
	pixelFetcher := &PixelFetcher{
		ppu: ppu,
	}

	return pixelFetcher
}

func (fetcher *PixelFetcher) prepareForScanline() {
	fetcher.state = StateGetTile
	fetcher.backgroundFifo = nil
	fetcher.pixelsToDiscard = fetcher.ppu.scx % 8
	fetcher.isFirstFetchOfScanline = true
	fetcher.isFetchingSprite = false
	fetcher.isFetchingWindow = false
	fetcher.currentX = 0
	fetcher.xPositionCounter = 0
	fetcher.windowLineCounter = 0
}

func (fetcher *PixelFetcher) step() {
	fetcher.tick()
	fetcher.attemptToPushPixel()
}

func (fetcher *PixelFetcher) tick() {
	// If the X-Position of any sprite in the sprite buffer is less than or
	// equal to the current Pixel-X-Position + 8, a sprite fetch is initiated.
	// See: https://ashiepaws.github.io/GBEDG/ppu/#sprite-fetching
	for _, oamIndex := range fetcher.ppu.spriteBuffer {
		// each sprite is 4 bytes long
		baseAddress := oamIndex * 4
		spriteX := fetcher.ppu.oam[baseAddress+1]

		if spriteX <= fetcher.currentX+8 {
			logger.Info("SWITCHING TO SPRITE FETCHING MODE")
			fetcher.isFetchingSprite = true
			fetcher.state = StateGetTile
			fetcher.spriteIndex = oamIndex

			// reset background fetcher
			fetcher.counter = 0
			fetcher.fetchedTileNumber = 0
			fetcher.fetchedTileDataLow = 0
			fetcher.fetchedTileDataHigh = 0
			return
		}
	}

	windowEnabled := (fetcher.ppu.lcdc>>5)&1 == 1
	if !fetcher.isFetchingWindow && (windowEnabled) && (fetcher.wyEqualedLyDuringFrame) && (fetcher.currentX >= fetcher.ppu.wx-7) {
		fetcher.state = StateGetTile
		fetcher.backgroundFifo = nil
		fetcher.isFetchingWindow = true
		fetcher.scanlineHadWindowPixels = true
		fetcher.xPositionCounter = 0
		fetcher.counter = 0
		return
	}

	fetcher.counter++

	if fetcher.counter < 2 {
		return
	}
	fetcher.counter = 0

	switch fetcher.state {
	case StateGetTile:
		if fetcher.isFetchingSprite {
			oamIndex := fetcher.ppu.spriteBuffer[fetcher.spriteIndex]
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
		fetcher.fetchedTileDataLow = fetcher.fetchTileData(0)
		fetcher.state = StateGetTileDataHigh
	case StateGetTileDataHigh:
		fetcher.fetchedTileDataHigh = fetcher.fetchTileData(1)

		// Note: The first time the background fetcher completes this step on a scanline the status is fully reset and operation restarts at Step 1.
		// See: https://ashiepaws.github.io/GBEDG/ppu/#background-pixel-fetching
		if fetcher.isFirstFetchOfScanline {
			fetcher.isFirstFetchOfScanline = false
			fetcher.state = StateGetTile
			fetcher.counter = 0
		} else {
			fetcher.counter = 0
			fetcher.state = StateSleep
		}
	case StateSleep:
		fetcher.state = StatePush
	case StatePush:
		if fetcher.isFetchingSprite {
			oamIndex := fetcher.ppu.spriteBuffer[fetcher.spriteIndex]
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

				fifoIndex := i

				if fetcher.spriteFifo[fifoIndex].colorId == 0 && tempBuffer[i].colorId != 0 {
					fetcher.spriteFifo[i] = tempBuffer[i]
				}
			}
			fetcher.state = StateGetTile
			fetcher.isFetchingSprite = false
		} else {
			// Note: While fetching background pixels, this step is only executed if the background FIFO is fully empty.
			// See: https://ashiepaws.github.io/GBEDG/ppu/#background-pixel-fetching
			if len(fetcher.backgroundFifo) == 0 {
				for i := 7; i >= 0; i-- {
					lowBit := (fetcher.fetchedTileDataLow >> i) & 1
					highBit := (fetcher.fetchedTileDataHigh >> i) & 1
					colorId := (highBit << 1) | lowBit
					pixel := FIFO{
						colorId:            colorId,
						palette:            0, // only applies to sprites
						backgroundPriority: 0, // only applies to sprites
					}
					fetcher.backgroundFifo = append(fetcher.backgroundFifo, pixel)
				}
				fetcher.state = StateGetTile
				fetcher.xPositionCounter++
				fetcher.isFetchingSprite = false
			}
		}
	}
}

func (fetcher *PixelFetcher) fetchTileData(offset uint16) uint8 {
	var address uint16
	if fetcher.isFetchingSprite {
		// sprites always use 8000 method
		address = 0x8000

		oamIndex := fetcher.ppu.spriteBuffer[fetcher.spriteIndex]
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
			// handle y flipping
			if flipY {
				rowInSprite = 7 - rowInSprite
			}
		}

		address += uint16(spriteTileNumber * 16)
		address += uint16(rowInSprite * 2)
	} else {
		// 8000 method
		if (fetcher.ppu.lcdc>>4)&1 == 1 {
			address = 0x8000
			address += uint16(fetcher.fetchedTileNumber * 16)
		} else
		// 8800 method
		{
			address = 0x9000
			address += uint16(int16(int8(fetcher.fetchedTileNumber)) * 16)
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
	if len(fetcher.backgroundFifo) == 0 {
		return
	}

	// otherwise add to the framebuffer
	backgroundPixel := fetcher.backgroundFifo[0]
	fetcher.backgroundFifo = fetcher.backgroundFifo[1:]
	var color uint8

	// do nothing if we have pixels to discard
	// this must be done after popping from the FIFO
	if fetcher.pixelsToDiscard > 0 {
		fetcher.pixelsToDiscard--
		return
	}

	// use the background pixel's color as the default
	colorId := backgroundPixel.colorId
	// color IDs are 2 bits, so we shift times 2, then mask 2 bits for the final color/shade
	color = (fetcher.ppu.bgp >> (colorId * 2)) & 0b11

	// See: https://ashiepaws.github.io/GBEDG/ppu/#pixel-mixing
	if len(fetcher.spriteFifo) > 0 {
		spritePixel := fetcher.spriteFifo[0]
		spriteIsTransparent := spritePixel.colorId == 0
		backgroundHasPriority := spritePixel.backgroundPriority == 1 && backgroundPixel.colorId != 0

		if !spriteIsTransparent && !backgroundHasPriority {
			colorId := spritePixel.colorId
			if spritePixel.palette == 0 {
				color = (fetcher.ppu.obp0 >> (colorId * 2)) & 0b11
			} else if spritePixel.palette == 1 {
				color = (fetcher.ppu.obp1 >> (colorId * 2)) & 0b11
			}
		}
	}

	logger.Info(
		"PPU: ADDING TO FRAMEBUFFER",
		"COLOR", color,
		"LY", fetcher.ppu.ly,
		"X", fetcher.currentX,
	)
	fetcher.ppu.frameBuffer[fetcher.ppu.ly][fetcher.currentX] = color
	fetcher.currentX++
}
