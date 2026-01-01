package cartridge

import (
	"fmt"

	"github.com/davidyorr/LuccaGB/internal/logger"
)

type Mbc5 struct {
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

	// 0000–1FFF — RAM Enable (Write Only)
	ramg uint8

	// 2000-2FFF - 8 least significant bits of ROM bank number (Write Only)
	romb0 uint8

	// 3000-3FFF - 9th bit of ROM bank number (Write Only)
	romb1 uint8

	// 4000-5FFF - RAM bank number (Write Only)
	ramb uint8
}

func newMbc5(cartridge *Cartridge) *Mbc5 {
	mbc5 := &Mbc5{}

	mbc5.cartridge = cartridge
	mbc5.romAddressMask = addressMaskSizes[cartridge.romSizeCode]
	mbc5.ramAddressMask = ramAddressMaskSizes[cartridge.ramSizeCode]
	cartridge.ram = make([]uint8, ramSizes[cartridge.ramSizeCode])
	mbc5.Reset()

	return mbc5
}

func (mbc *Mbc5) Reset() {
	mbc.ramg = 0x00
	mbc.romb0 = 0x01
	mbc.romb1 = 0x00
	mbc.ramb = 0x00
}

func (mbc *Mbc5) Read(address uint16) uint8 {
	switch {

	// ROM Bank 00
	case address >= 0x000 && address <= 0x3FFF:
		return mbc.cartridge.rom[address]

	// ROM BANK 00-1FF
	case address >= 0x4000 && address <= 0x7FFF:
		// lower 8 bits from ROMB0, bit 8 from ROMB1
		bank := (uint32(mbc.romb1) << 8) | uint32(mbc.romb0)
		// map 0x4000-0x7FFF down to 0x0000-0x3FFF
		offset := uint32(address) & 0b11_1111_1111_1111
		actualAddress := ((bank << 14) | offset) & mbc.romAddressMask
		return mbc.cartridge.rom[actualAddress]

	// External RAM
	case address >= 0xA000 && address <= 0xBFFF:
		// RAM disabled
		if (mbc.ramg & 0b1111) != 0b1010 {
			return 0xFF
		}

		// no RAM hardware
		if len(mbc.cartridge.ram) == 0 {
			return 0xFF
		}

		bank := uint32(mbc.ramb)
		offset := uint32(address - 0xA000)

		actualAddress := ((bank << 13) | offset) & uint32(mbc.ramAddressMask)
		return mbc.cartridge.ram[actualAddress]
	}

	logger.Error(
		"MBC5 returning 0xFF",
		"ADDRESS", fmt.Sprintf("0x%04X", address),
		"ROMB0", fmt.Sprintf("0x%08b", mbc.romb0),
		"ROMB1", fmt.Sprintf("0x%08b", mbc.romb1),
		"RAMB", fmt.Sprintf("0x%08b", mbc.ramb),
	)
	return 0xFF
}

func (mbc *Mbc5) Write(address uint16, value uint8) {
	switch {

	// RAM Enable
	case address >= 0x0000 && address <= 0x1FFF:
		mbc.ramg = value

	// ROMB0 - lower ROM bank register (bits 0-7)
	case address >= 0x2000 && address <= 0x2FFF:
		mbc.romb0 = value

	// ROMB1 - upper ROM bank register (bit 8)
	case address >= 0x3000 && address <= 0x3FFF:
		mbc.romb1 = value & 0b0001

	// RAM Bank Select
	case address >= 0x4000 && address <= 0x5FFF:
		mbc.ramb = value & 0b0000_1111

	// Write to RAM
	case address >= 0xA000 && address <= 0xBFFF:
		// RAM disabled
		if (mbc.ramg & 0b1111) != 0b1010 {
			return
		}

		// no RAM hardware
		if len(mbc.cartridge.ram) == 0 {
			return
		}

		bank := uint32(mbc.ramb)
		offset := uint32(address - 0xA000)

		actualAddress := ((bank << 13) | offset) & uint32(mbc.ramAddressMask)
		mbc.cartridge.ram[actualAddress] = value
	}
}
