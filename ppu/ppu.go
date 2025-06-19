package ppu

type PPU struct {
	// 0x8000 - 9FFF
	videoRam [8192]uint8
	// 0xFE00 - 0xFE9F
	oam [160]uint8
	// 0xFF44 - LY: LCD Y coordinate [read-only]
	ly uint8

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
	ppu.ly = 0x91
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
