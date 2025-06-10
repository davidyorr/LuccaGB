package ppu

type PPU struct {
	// 0x8000 - 9FFF
	videoRam [8192]uint8
	// 0xFE00 - 0xFE9F
	oam [160]uint8
}

func New() *PPU {
	ppu := &PPU{}

	ppu.Reset()

	return ppu
}

func (ppu *PPU) Reset() {

}
