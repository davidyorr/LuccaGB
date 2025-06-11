package cartridge

type Cartridge struct {
	// 0x0000 - 0x7FFF, dynamically sized
	rom []uint8
}

func New() *Cartridge {
	cartridge := &Cartridge{}

	return cartridge
}

func (cartridge *Cartridge) SetRom(rom []uint8) {
	cartridge.rom = rom
}
