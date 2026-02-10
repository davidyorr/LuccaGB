package cartridge

type Mbc2 struct {
	cartridge *Cartridge
	// bitmask to wrap addresses to the physical ROM capacity,
	// derived from the ROM size code
	romAddressMask uint32
	// bitmask to wrap addresses to the physical RAM capacity,
	// derived from the RAM size code
	ramAddressMask uint32

	// =======================
	// ====== Registers ======
	// =======================

	// 0000–3FFF — ROM Bank 0 [read-only]
	ramg uint8

	// 4000–7FFF — ROM Bank $01-0F [read-only]
	romb uint8

	// A000–A1FF — Built-in RAM
	// A200–BFFF — 15 “echoes” of A000–A1FF
}

func newMbc2(cartridge *Cartridge) *Mbc2 {
	mbc2 := &Mbc2{}

	mbc2.cartridge = cartridge
	mbc2.romAddressMask = addressMaskSizes[cartridge.romSizeCode]
	// MBC2 RAM is fixed at 512 bytes.
	// The mask to wrap 0-511 is 0x1FF (511).
	mbc2.ramAddressMask = 0x1FF
	cartridge.ram = make([]uint8, 512)
	mbc2.Reset()

	return mbc2
}

func (mbc *Mbc2) Reset() {
	mbc.ramg = 0x00
	mbc.romb = 0x01
}

func (mbc *Mbc2) Read(address uint16) uint8 {
	switch {

	// ROM Bank 00
	case address >= 0x000 && address <= 0x3FFF:
		return mbc.cartridge.rom[address]

	// ROM BANK 01-0F
	case address >= 0x4000 && address <= 0x7FFF:
		bank := uint32(mbc.romb) & 0b1111

		// map 0x4000-0x7FFF down to 0x0000-0x3FFF
		offset := uint32(address) & 0b11_1111_1111_1111
		actualAddress := ((bank << 14) | offset) & mbc.romAddressMask
		return mbc.cartridge.rom[actualAddress]

	// Built-in RAM
	case address >= 0xA000 && address <= 0xBFFF:
		// RAM disabled
		if (mbc.ramg & 0b1111) != 0b1010 {
			return 0xFF
		}

		offset := uint32(address - 0xA000)
		actualAddress := offset & mbc.ramAddressMask
		return mbc.cartridge.ram[actualAddress] | 0b1111_0000
	}

	return 0xFF
}

func (mbc *Mbc2) Write(address uint16, value uint8) {
	switch {

	// RAM Enable, ROM Bank Number
	case address >= 0x0000 && address <= 0x3FFF:
		// controls whether the RAM is enabled
		if uint16(address)&0b1_0000_0000 == 0 {
			mbc.ramg = value & 0b0000_1111
		} else
		// controls the selected ROM bank
		{
			if (value & 0b1111) == 0x00 {
				mbc.romb = 0x01
			} else {
				mbc.romb = value & 0b1111
			}
		}

	// Write to RAM
	case address >= 0xA000 && address <= 0xBFFF:
		// RAM disabled
		if (mbc.ramg & 0b1111) != 0b1010 {
			return
		}

		offset := uint32(address - 0xA000)
		actualAddress := offset & mbc.ramAddressMask
		mbc.cartridge.ram[actualAddress] = value & 0b0000_1111
	}
}

func (mbc *Mbc2) Serialize(buf []byte) int {
	offset := 0

	buf[offset] = mbc.ramg
	offset++
	buf[offset] = mbc.romb
	offset++

	return offset
}

func (mbc *Mbc2) Deserialize(buf []byte) int {
	offset := 0

	mbc.ramg = buf[offset]
	offset++
	mbc.romb = buf[offset]
	offset++

	return offset
}
