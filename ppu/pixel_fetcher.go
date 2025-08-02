package ppu

import (
	"fmt"

	"github.com/davidyorr/EchoGB/logger"
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
	color uint8
	// only applies to objects (sprites)
	palette            uint8
	backgroundPriority uint8
}

func newPixelFetcher(ppu *PPU) *PixelFetcher {
	pixelFetcher := &PixelFetcher{
		ppu: ppu,
	}

	pixelFetcher.reset()

	return pixelFetcher
}

func (fetcher *PixelFetcher) reset() {
	fetcher.state = StateGetTile
	fetcher.counter = 0
	fetcher.xPositionCounter = 0
	fetcher.windowLineCounter = 0
	fetcher.isFirstFetchOfScanline = false
	fetcher.currentX = 0
}

func (fetcher *PixelFetcher) prepareForScanline() {
	fetcher.backgroundFifo = nil
	fetcher.pixelsToDiscard = fetcher.ppu.scx % 8
	fetcher.isFirstFetchOfScanline = true
}

func (fetcher *PixelFetcher) step() {
	// add to the framebuffer if the FIFO contains any data
	if len(fetcher.backgroundFifo) > 0 && fetcher.currentX < 160 {
		pixel := fetcher.backgroundFifo[0]
		fetcher.backgroundFifo = fetcher.backgroundFifo[1:]

		if fetcher.pixelsToDiscard > 0 {
			fetcher.pixelsToDiscard--
		} else {
			logger.Info(
				"PPU: ADDING TO FRAMEBUFFER",
				"COLOR", pixel.color,
				"LY", fetcher.ppu.ly,
				"X", fetcher.currentX,
			)
			fetcher.ppu.frameBuffer[fetcher.ppu.ly][fetcher.currentX] = pixel.color
			fetcher.currentX++
		}
	}

	windowEnabled := (fetcher.ppu.lcdc>>5)&1 == 1
	if !fetcher.isFetchingWindow && (windowEnabled) && (fetcher.wyEqualedLyDuringFrame) && (fetcher.currentX >= fetcher.ppu.wx-7) {
		fetcher.state = StateGetTile
		fetcher.backgroundFifo = nil
		fetcher.isFetchingWindow = true
		fetcher.scanlineHadWindowPixels = true
		fetcher.xPositionCounter = 0
		return
	}

	fetcher.counter++

	switch fetcher.state {
	case StateGetTile:
		if fetcher.counter == 2 {
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
			logger.Info(
				"e == GET TILE",
				"ADDRESS", fmt.Sprintf("0x%04X", address),
				"Y_TILE", fmt.Sprintf("0x%02X", yTile),
				"X_TILE", fmt.Sprintf("0x%02X", xTile),
				"TILE_NUMBER", fmt.Sprintf("0x%02X", fetcher.fetchedTileNumber),
			)

			fetcher.counter = 0
			fetcher.state = StateGetTileDataLow
		}
	case StateGetTileDataLow:
		if fetcher.counter == 2 {
			var address uint16
			// 8000 method
			if (fetcher.ppu.lcdc>>4)&1 == 1 {
				address = 0x8000
				address += uint16(fetcher.fetchedTileNumber * 16)
			} else
			// 8800 method
			{
				address = 0x9000
				address += uint16(int8(fetcher.fetchedTileNumber) * 16)
			}
			// get the row offset
			if fetcher.isFetchingWindow {
				address += uint16(2 * (fetcher.windowLineCounter % 8))
			} else {
				address += uint16(2 * ((fetcher.ppu.ly + fetcher.ppu.scy) % 8))
			}

			fetcher.fetchedTileDataLow = fetcher.ppu.videoRam[address-0x8000]

			fetcher.counter = 0
			fetcher.state = StateGetTileDataHigh
		}
	case StateGetTileDataHigh:
		if fetcher.counter == 2 {
			var address uint16
			// 8000 method
			if (fetcher.ppu.lcdc>>4)&1 == 1 {
				address = 0x8000
				address += uint16(fetcher.fetchedTileNumber * 16)
			} else
			// 8800 method
			{
				address = 0x9000
				address += uint16(int8(fetcher.fetchedTileNumber) * 16)
			}
			// get the row offset
			if fetcher.isFetchingWindow {
				address += uint16(2 * (fetcher.windowLineCounter % 8))
			} else {
				address += uint16(2 * ((fetcher.ppu.ly + fetcher.ppu.scy) % 8))
			}

			address += 1
			fetcher.fetchedTileDataHigh = fetcher.ppu.videoRam[address-0x8000]

			// Note: The first time the background fetcher completes this step on a scanline the status is fully reset and operation restarts at Step 1.
			// See: https://hacktix.github.io/GBEDG/Ground-pixel-fetching
			if fetcher.isFirstFetchOfScanline {
				fetcher.isFirstFetchOfScanline = false
				fetcher.state = StateGetTile
				fetcher.counter = 0
			} else {
				fetcher.counter = 0
				fetcher.state = StateSleep
			}
		}
	case StateSleep:
		if fetcher.counter == 2 {
			fetcher.counter = 0
			fetcher.state = StatePush
		}
	case StatePush:
		// Note: While fetching background pixels, this step is only executed if the background FIFO is fully empty.
		// See: https://hacktix.github.io/GBEDG/Ground-pixel-fetching
		if len(fetcher.backgroundFifo) == 0 {
			for i := 7; i >= 0; i-- {
				lowBit := (fetcher.fetchedTileDataLow >> i) & 1
				highBit := (fetcher.fetchedTileDataHigh >> i) & 1
				color := (highBit << 1) | lowBit
				pixel := FIFO{
					color:              color,
					palette:            0,
					backgroundPriority: 0,
				}
				fetcher.backgroundFifo = append(fetcher.backgroundFifo, pixel)
			}

			fetcher.state = StateGetTile
			fetcher.counter = 0
			fetcher.xPositionCounter++
		}
	}
}
