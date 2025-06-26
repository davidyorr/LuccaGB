package ppu

type PPU struct {
	// 0x8000 - 0x97FF
	vramTileData [6144]uint8

	// 0x8000 - 9FFF
	videoRam [8192]uint8
	// 0xFE00 - 0xFE9F
	oam [160]uint8
	// 0xFF40 - LCDC: LCD control
	lcdc uint8
	// 0xFF44 - LY: LCD Y coordinate [read-only]
	ly uint8
	// 0xFF45 - LYC: LY compare
	lyc uint8
	// 0xFF41 - STAT: LCD status
	stat uint8
	// 0xFF42 - SCY: Background viewport Y position
	scy uint8
	// 0xFF43 - SCX: Background viewport X position
	scx uint8
	// 0xFF4A - WY: Window Y position
	wy uint8
	// 0xFF4B - WX: Window X position plus 7
	wx      uint8
	counter uint16
}

// 1 dot = T-cycle
const dotsPerScanline = 456

func New() *PPU {
	ppu := &PPU{}

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
	ppu.counter = 0
}

func (ppu *PPU) Step(cycles uint8) {
	ppu.counter += uint16(cycles)

	if ppu.counter > dotsPerScanline {
		ppu.ly++
		ppu.counter -= dotsPerScanline
		if ppu.ly == 154 {
			ppu.ly = 0
		}
	}
}

func (ppu *PPU) Read(address uint16) uint8 {
	if address == 0xFF44 {
		return ppu.ly
	}
	return 0
}
