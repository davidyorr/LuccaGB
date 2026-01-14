package ppu

import (
	"fmt"

	"github.com/davidyorr/LuccaGB/internal/debug"
	"github.com/davidyorr/LuccaGB/internal/interrupt"
	"github.com/davidyorr/LuccaGB/internal/logger"
)

type PPU struct {
	pixelFetcher *PixelFetcher

	// 0x8000 - 0x9FFF
	//	Block 0 - 0x8000 - 0x87FF
	//	Block 1 - 0x8800 - 0x8FFF
	//	Block 2 - 0x9000 - 0x97FF
	//	Each block contains 384 tiles, each 16 bytes
	videoRam [8192]uint8
	// 0xFE00 - 0xFE9F - Object Attribute Memory
	//	40 sprites (objects), each 4 bytes long
	//	Byte 0 — Y Position
	//	Byte 1 — X Position
	//	Byte 2 — Tile Index
	//	Byte 3 — Attributes/Flags
	oam [160]uint8
	// 0xFEA0 - 0xFEFF - Not usable
	//	Nintendo says use of this area is prohibited
	unusable [96]uint8
	// 0xFF40 - LCDC: LCD control
	//	7 - LCD & PPU enable
	//	6 - Window tile map area: 0 = 9800–9BFF; 1 = 9C00–9FFF
	//	5 - Window enable
	//	4 - BG & Window tile data set: 0 = 8800–97FF; 1 = 8000–8FFF
	//	3 - BG tile map area: 0 = 9800–9BFF; 1 = 9C00–9FFF
	//	2 - OBJ size: 0 = 8×8; 1 = 8×16
	//	1 - OBJ enable
	//	0 - BG & Window enable / priority
	lcdc uint8
	// 0xFF44 - LY: LCD Y coordinate [read-only]
	ly uint8
	// 0xFF45 - LYC: LY compare
	lyc uint8
	// 0xFF41 - STAT: LCD status
	//	6 - LYC int select
	//	5 - Mode 2 int select
	//	4 - Mode 1 int select
	//	3 - Mode 0 int select
	//	2 - LYC == ly
	//	1 0 - PPU mode
	stat uint8
	// 0xFF42 - SCY: Background viewport Y position
	scy uint8
	// 0xFF43 - SCX: Background viewport X position
	scx uint8
	// 0xFF4A - WY: Window Y position
	wy uint8
	// 0xFF4B - WX: Window X position plus 7
	wx uint8
	// 0xFF47 - BGP: Background palette data
	bgp uint8
	// 0xFF48 - OBP0: Object palette 0 data
	obp0 uint8
	// 0xFF49 - OBP1: Object palette 1 data
	obp1                           uint8
	mode                           Mode
	previousStatInterruptLineState bool
	// 10 sprites can be displayed per scanline
	spriteBuffer       SpriteBuffer
	frameBuffer        [144][160]uint8
	interruptRequester func(interruptType interrupt.Interrupt)
	dot                uint16
}

func New(interruptRequest func(interrupt.Interrupt)) *PPU {
	ppu := &PPU{}
	ppu.interruptRequester = interruptRequest
	ppu.pixelFetcher = newPixelFetcher(ppu)
	ppu.Reset()

	return ppu
}

func (ppu *PPU) Reset() {
	ppu.lcdc = 0x91
	ppu.ly = 0x00
	ppu.lyc = 0x00
	ppu.stat = 0x85
	ppu.scy = 0x00
	ppu.scx = 0x00
	ppu.wy = 0x00
	ppu.wx = 0x00
	ppu.bgp = 0xFC
	ppu.mode = OamScan
	ppu.dot = 0
}

// 1 dot = T-cycle
const dotsPerScanline = 456

// Perform 1 T-cycle of work
func (ppu *PPU) Step() (frameReady bool) {
	if !ppu.lcdEnabled() {
		return
	}

	ppu.updateStatInterruptLine()

	if ppu.ly < 144 {
		if ppu.dot == 0 {
			ppu.changeMode(OamScan)
			ppu.spriteBuffer.Reset()
		} else if ppu.dot == 80 {
			ppu.changeMode(DrawingPixels)
			ppu.pixelFetcher.prepareForScanline()
		}
	}

	if ppu.ly < 144 {
		switch ppu.mode {
		// Mode 2
		case OamScan:
			if ppu.dot%2 == 0 {
				if ppu.spriteBuffer.size < MaxSpriteBufferSize {
					oamIndex := ppu.dot / 2
					spriteY := ppu.oam[oamIndex*4]
					spriteX := ppu.oam[oamIndex*4+1]
					var height uint8 = 8
					if ((ppu.lcdc & 0b0000_0100) >> 2) == 1 {
						height = 16
					}
					// See: https://ashiepaws.github.io/GBEDG/ppu/#oam-scan-mode-2
					if ppu.ly+16 >= spriteY && ppu.ly+16 < spriteY+height && spriteX > 0 {
						ppu.spriteBuffer.Push(uint8(oamIndex))
					}
				}
			}
		// Mode 3
		case DrawingPixels:
			ppu.pixelFetcher.step()
		// Mode 0
		case HorizontalBlank:
			// "pads" the duration of the scanline to a total of 456
		}
	}

	ppu.dot++

	// end of scanline
	if ppu.dot == dotsPerScanline {
		ppu.dot = 0
		ppu.ly++

		if ppu.pixelFetcher.scanlineHadWindowPixels {
			ppu.pixelFetcher.windowLineCounter++
		}
		ppu.pixelFetcher.scanlineHadWindowPixels = false

		if ppu.wy == ppu.ly {
			ppu.pixelFetcher.wyEqualedLyDuringFrame = true
		}

		if ppu.ly == 144 {
			frameReady = true
			ppu.changeMode(VerticalBlank)
			ppu.interruptRequester(interrupt.VBlankInterrupt)
			// If bit 5 (mode 2 OAM interrupt) is set, an LCD interrupt is also triggered.
			// See: https://github.com/Gekkio/mooneye-test-suite/blob/443f6e1f2a8d83ad9da051cbb960311c5aaaea66/acceptance/ppu/vblank_stat_intr-GS.s#L21
			if (ppu.stat & 0b0010_0000) != 0 {
				ppu.interruptRequester(interrupt.LcdInterrupt)
			}
		} else if ppu.ly == 154 {
			ppu.ly = 0
			ppu.pixelFetcher.windowLineCounter = 0
			ppu.pixelFetcher.wyEqualedLyDuringFrame = false

			// we need to check here as well, because LY was changed
			if ppu.wy == ppu.ly {
				ppu.pixelFetcher.wyEqualedLyDuringFrame = true
			}
		}

		ppu.updateLycCoincidenceFlag()
	}

	return frameReady
}

func (ppu *PPU) Read(address uint16) uint8 {
	switch {
	case address == 0xFF40:
		return ppu.lcdc
	case address == 0xFF41:
		return ppu.stat | 0b1000_0000
	case address == 0xFF42:
		return ppu.scy
	case address == 0xFF43:
		return ppu.scx
	case address == 0xFF44:
		return ppu.ly
	case address == 0xFF45:
		return ppu.lyc
	case address == 0xFF47:
		return ppu.bgp
	case address == 0xFF48:
		return ppu.obp0
	case address == 0xFF49:
		return ppu.obp1
	case address == 0xFF4A:
		return ppu.wy
	case address == 0xFF4B:
		return ppu.wx
	// VRAM
	case address >= 0x8000 && address <= 0x9FFF:
		if ppu.vramIsAccessible() {
			return ppu.videoRam[address-0x8000]
		}
		return 0xFF
	// OAM
	case address >= 0xFE00 && address <= 0xFE9F:
		if ppu.oamIsAccessible() {
			return ppu.oam[address-0xFE00]
		}
		return 0xFF
	}
	return 0
}

func (ppu *PPU) Write(address uint16, value uint8) {
	if debug.Enabled {
		logger.Debug(
			"PPU Write",
			"ADDRESS", fmt.Sprintf("0x%04X", address),
			"VALUE", fmt.Sprintf("0x%02X", value),
		)
	}
	switch {
	case address == 0xFF40:
		lcdWasEnabled := ppu.lcdEnabled()
		ppu.lcdc = value
		lcdIsEnabled := ppu.lcdEnabled()

		// LCD ON -> LCD OFF
		if lcdWasEnabled && !lcdIsEnabled {
			ppu.ly = 0
			ppu.updateLycCoincidenceFlag()
			ppu.dot = 0
			// When LCD is disabled, STAT mode reads as 0 (HBlank)
			// See: https://gbdev.io/pandocs/STAT.html#ff41--stat-lcd-status
			ppu.changeMode(HorizontalBlank)
		}
		// LCD OFF -> LCD ON
		if !lcdWasEnabled && lcdIsEnabled {
			ppu.changeMode(OamScan)
			ppu.updateLycCoincidenceFlag()
		}
	case address == 0xFF41:
		ppu.stat = (ppu.stat & 0b1000_0111) | (value & 0b0111_1000)
	case address == 0xFF42:
		ppu.scy = value
	case address == 0xFF43:
		ppu.scx = value
	case address == 0xFF45:
		ppu.lyc = value
		ppu.updateLycCoincidenceFlag()
	case address == 0xFF47:
		ppu.bgp = value
	case address == 0xFF48:
		ppu.obp0 = value
	case address == 0xFF49:
		ppu.obp1 = value
	case address == 0xFF4A:
		ppu.wy = value
	case address == 0xFF4B:
		ppu.wx = value
	// VRAM
	case address >= 0x8000 && address <= 0x9FFF:
		if ppu.vramIsAccessible() {
			ppu.videoRam[address-0x8000] = value
		}
	// OAM
	case address >= 0xFE00 && address <= 0xFE9F:
		if ppu.oamIsAccessible() {
			ppu.oam[address-0xFE00] = value
		}
	default:
		logger.Debug("unhandled address while writing <-")
	}
}

// to give direct access to the DMA
func (ppu *PPU) WriteOam(address uint16, value uint8) {
	ppu.oam[address-0xFE00] = value
}

func (ppu *PPU) Mode() Mode {
	return ppu.mode
}

func (ppu *PPU) OamIsBlocked() bool {
	return ppu.mode == OamScan
}

type Mode uint8

const (
	HorizontalBlank Mode = 0b00
	VerticalBlank   Mode = 0b01
	OamScan         Mode = 0b10
	DrawingPixels   Mode = 0b11
)

func (ppu *PPU) changeMode(mode Mode) {
	ppu.mode = mode
	ppu.stat = (ppu.stat & 0b1111_1100) | uint8(mode)
}

// INT $48 — STAT interrupt
//
// There are various sources which can trigger this interrupt to occur as
// described in STAT register ($FF41). The various STAT interrupt sources (modes
// 0-2 and LYC=LY) have their state (inactive=low and active=high) logically
// ORed into a shared “STAT interrupt line” if their respective enable bit is
// turned on. A STAT interrupt will be triggered by a rising edge (transition
// from low to high) on the STAT interrupt line.
// See: https://gbdev.io/pandocs/Interrupt_Sources.html#int-48--stat-interrupt
func (ppu *PPU) updateStatInterruptLine() {
	mode0 := (ppu.mode == HorizontalBlank) && ((ppu.stat & 0b0000_1000) != 0)
	mode1 := (ppu.mode == VerticalBlank) && ((ppu.stat & 0b0001_0000) != 0)
	mode2 := (ppu.mode == OamScan) && ((ppu.stat & 0b0010_0000) != 0)
	lycLyMatch := ((ppu.stat & 0b0100_0000) != 0) && ((ppu.stat & 0b0000_0100) != 0)

	currentStatInterruptLineState := mode0 || mode1 || mode2 || lycLyMatch

	// Check for a rising edge
	if !ppu.previousStatInterruptLineState && currentStatInterruptLineState {
		ppu.interruptRequester(interrupt.LcdInterrupt)
	}

	// Update the state for the next check
	ppu.previousStatInterruptLineState = currentStatInterruptLineState
}

// updateLycCoincidenceFlag sets or clears the STAT register's Coincidence Flag
// (bit 2) based on whether the LY and LYC registers are equal. This must be
// called whenever either LY or LYC is modified to ensure the flag's state is
// always accurate.
// See: https://ashiepaws.github.io/GBEDG/ppu/#stat2---coincidence-flag
func (ppu *PPU) updateLycCoincidenceFlag() {
	if ppu.lyc == ppu.ly {
		ppu.stat |= 0b0000_0100
	} else {
		ppu.stat &^= 0b0000_0100
	}
}

func (ppu *PPU) lcdEnabled() bool {
	return (ppu.lcdc & 0b1000_0000) != 0
}

func (ppu *PPU) vramIsAccessible() bool {
	return ppu.mode != DrawingPixels
}

func (ppu *PPU) oamIsAccessible() bool {
	return ppu.mode != OamScan && ppu.mode != DrawingPixels
}

func (ppu *PPU) FrameBuffer() [144][160]uint8 {
	return ppu.frameBuffer
}
