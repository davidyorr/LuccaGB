package ppu

import (
	"fmt"

	"github.com/davidyorr/EchoGB/interrupt"
	"github.com/davidyorr/EchoGB/logger"
)

type PPU struct {
	// 0x8000 - 0x9FFF
	videoRam [8192]uint8
	// 0xFE00 - 0xFE9F - Object Attribute Memory
	//	40 sprites (objects), each 4 bytes long
	oam [160]uint8
	// 0xFEA0 - 0xFEFF - Not usable
	//	Nintendo says use of this area is prohibited
	unusable [96]uint8
	// 10 sprites can be displayed per scanline
	visibleSprites []uint8
	// 0xFF40 - LCDC: LCD control
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
	obp1               uint8
	mode               Mode
	interruptRequester func(interruptType interrupt.Interrupt)
	counter            uint16
}

func New(interruptRequest func(interrupt.Interrupt)) *PPU {
	ppu := &PPU{}
	ppu.interruptRequester = interruptRequest

	ppu.Reset()

	return ppu
}

func (ppu *PPU) Reset() {
	ppu.lcdc = 0x91
	ppu.ly = 0x91
	ppu.lyc = 0x00
	ppu.stat = 0x81
	ppu.scy = 0x00
	ppu.scx = 0x00
	ppu.wy = 0x00
	ppu.wx = 0x00
	ppu.mode = OamScan
	ppu.counter = 0
}

// 1 dot = T-cycle
const dotsPerScanline = 456

// Perform 1 T-cycle of work
func (ppu *PPU) Step() {
	if !ppu.lcdEnabled() {
		return
	}

	ppu.counter++

	if ppu.counter == dotsPerScanline {
		ppu.counter = 0
		ppu.ly++

		if ppu.ly == 144 {
			ppu.changeMode(VerticalBlank)
			ppu.interruptRequester(interrupt.VBlankInterrupt)
		} else if ppu.ly >= 154 {
			ppu.ly = 0
		}

		ppu.compareLycLy()
	}

	if ppu.ly < 144 {
		switch {
		// Mode 2: OAM scan
		case ppu.counter < 80:
			if ppu.mode != OamScan {
				ppu.changeMode(OamScan)
				ppu.visibleSprites = nil
			}
			if ppu.counter%2 == 0 {
				if len(ppu.visibleSprites) < 10 {
					oamIndex := ppu.counter / 2
					spriteY := ppu.oam[oamIndex*4]
					var height uint8 = 8
					if ((ppu.lcdc & 0b0000_0100) >> 2) == 1 {
						height = 16
					}
					if ppu.ly+16 >= spriteY && ppu.ly+16 < spriteY+height {
						ppu.visibleSprites = append(ppu.visibleSprites, uint8(oamIndex))
					}
				}
			}
		// Mode 3: Drawing Pixels
		case ppu.counter < ppu.getMode3Duration():
			if ppu.mode != DrawingPixels {
				ppu.changeMode(DrawingPixels)
			}
		default:
			if ppu.mode != HorizontalBlank {
				ppu.changeMode(HorizontalBlank)
			}
		}
	}
}

func (ppu *PPU) getMode3Duration() uint16 {
	var baseDuration uint16 = 172
	var backgroundScrollingPenalty uint16 = uint16(ppu.scx) % 8
	var windowPenalty uint16 = 0
	var objectsPenalty uint16 = 0

	windowEnabled := ppu.lcdc&0b0010_0000>>5 == 1
	if windowEnabled && ppu.ly >= ppu.wy && ppu.wx < 167 && ppu.wx > 7 {
		windowPenalty = 6
	}

	type Position struct {
		x, y uint16
	}

	var processedBackgroundTiles map[Position]bool = make(map[Position]bool)
	for _, oamIndex := range ppu.visibleSprites {
		// each sprite is 4 bytes long
		baseAddress := oamIndex * 4
		spriteX := ppu.oam[baseAddress+1]

		if spriteX == 0 {
			objectsPenalty += 11
			continue
		}

		objectsPenalty += 6

		screenX := int(spriteX) - 8
		// ensure a positive result
		backgroundX := uint16(((int(screenX)+int(ppu.scx))%256 + 256) % 256)
		backgroundTileX := backgroundX / 8
		backgroundY := (uint16(ppu.ly) + uint16(ppu.scy)) % 256
		backgroundTileY := backgroundY / 8
		pos := Position{x: backgroundTileX, y: backgroundTileY}

		if processedBackgroundTiles[pos] {
			continue
		}

		processedBackgroundTiles[pos] = true
		fineX := backgroundX % 8
		penalty := max(0, 5-fineX)
		objectsPenalty += penalty
	}

	return min(baseDuration+backgroundScrollingPenalty+windowPenalty+objectsPenalty, 289)
}

func (ppu *PPU) Read(address uint16) uint8 {
	switch {
	case address == 0xFF40:
		return ppu.lcdc
	case address == 0xFF41:
		return ppu.stat
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
	logger.Debug(
		"PPU Write",
		"Address", fmt.Sprintf("0x%04X", address),
		"Value", fmt.Sprintf("0x%02X", value),
	)
	switch {
	case address == 0xFF40:
		lcdWasEnabled := ppu.lcdEnabled()
		ppu.lcdc = value
		lcdIsEnabled := ppu.lcdEnabled()

		// LCD ON -> LCD OFF
		if lcdWasEnabled && !lcdIsEnabled {
			ppu.ly = 0
			ppu.compareLycLy()
			ppu.counter = 0
			ppu.changeMode(VerticalBlank)
		}
		// LCD OFF -> LCD ON
		if !lcdWasEnabled && lcdIsEnabled {
			ppu.changeMode(OamScan)
			ppu.compareLycLy()
		}
	case address == 0xFF41:
		ppu.stat = (ppu.stat & 0b1000_0111) | (value & 0b0111_1000)
	case address == 0xFF42:
		ppu.scy = value
	case address == 0xFF43:
		ppu.scx = value
	case address == 0xFF45:
		ppu.lyc = value
		ppu.compareLycLy()
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

	// Mode 0 (Horizontal Blank)
	if mode == HorizontalBlank && ((ppu.stat & 0b0000_1000) != 0) {
		ppu.interruptRequester(interrupt.LcdInterrupt)
	}
	// Mode 1 (Vertical Blan)
	if mode == VerticalBlank && ((ppu.stat & 0b0001_0000) != 0) {
		ppu.interruptRequester(interrupt.LcdInterrupt)
	}
	// Mode 2 (OAM Scan)
	if mode == OamScan && ((ppu.stat & 0b0010_0000) != 0) {
		ppu.interruptRequester(interrupt.LcdInterrupt)
	}
}

func (ppu *PPU) compareLycLy() {
	if ppu.lyc == ppu.ly {
		ppu.stat |= 0b0000_0100
		// LYC int select
		if (ppu.stat & 0b0100_0000) != 0 {
			ppu.interruptRequester(interrupt.LcdInterrupt)
		}
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
