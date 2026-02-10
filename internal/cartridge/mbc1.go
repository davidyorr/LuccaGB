package cartridge

import (
	"fmt"

	"github.com/davidyorr/LuccaGB/internal/logger"
)

type Mbc1 struct {
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

	// 2000–3FFF — ROM Bank Number (Write Only)
	bank1 uint8

	// 4000–5FFF — RAM Bank Number — or — Upper Bits of ROM Bank Number (Write Only)
	bank2 uint8

	// 6000–7FFF — Banking Mode Select (Write Only)
	mode uint8
}

func newMbc1(cartridge *Cartridge) *Mbc1 {
	mbc1 := &Mbc1{}

	mbc1.cartridge = cartridge
	mbc1.romAddressMask = addressMaskSizes[cartridge.romSizeCode]
	mbc1.ramAddressMask = ramAddressMaskSizes[cartridge.ramSizeCode]
	cartridge.ram = make([]uint8, ramSizes[cartridge.ramSizeCode])
	mbc1.Reset()

	return mbc1
}

func (mbc *Mbc1) Reset() {
	mbc.ramg = 0x00
	mbc.bank1 = 0x01
	mbc.bank2 = 0x00
	mbc.mode = 0x00
}

func (mbc *Mbc1) Read(address uint16) uint8 {
	switch {

	// ROM Bank 00
	case address >= 0x000 && address <= 0x3FFF:
		if mbc.mode == 0b00 {
			// Mode 0: Simple banking mode - read from bank 0
			actualAddress := uint32(address) & mbc.romAddressMask
			return mbc.cartridge.rom[actualAddress]
		} else if mbc.mode == 0b01 {
			// Mode 1: Advanced banking mode
			bank := uint32(mbc.bank2) << 5
			actualAddress := ((bank << 14) | uint32(address)) & mbc.romAddressMask
			return mbc.cartridge.rom[actualAddress]
		}

	// ROM BANK 01-7F
	case address >= 0x4000 && address <= 0x7FFF:
		// lower 5 bits from Bank 1, upper 2 bits from Bank 2
		bank := (uint32(mbc.bank2) << 5) | uint32(mbc.bank1)
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

		// Default to Bank 0
		bank := uint32(0)
		offset := uint32(address - 0xA000)

		// RAM banking enabled, use Bank 2 to select RAM bank 0-3
		if mbc.mode == 0b01 {
			bank = uint32(mbc.bank2)
		}

		actualAddress := ((bank << 13) | offset) & uint32(mbc.ramAddressMask)
		return mbc.cartridge.ram[actualAddress]
	}

	logger.Error(
		"MBC1 returning 0xFF",
		"ADDRESS", fmt.Sprintf("0x%04X", address),
		"BANK1", fmt.Sprintf("0x%08b", mbc.bank1),
		"BANK2", fmt.Sprintf("0x%08b", mbc.bank2),
		"MODE", fmt.Sprintf("0x%08b", mbc.mode),
	)
	return 0xFF
}

func (mbc *Mbc1) Write(address uint16, value uint8) {
	switch {

	// RAM Enable
	case address >= 0x0000 && address <= 0x1FFF:
		mbc.ramg = value & 0b0000_1111

	// Bank 1
	case address >= 0x2000 && address <= 0x3FFF:
		if (value & 0b1_1111) == 0x00 {
			mbc.bank1 = 0x01
		} else {
			mbc.bank1 = value & 0b1_1111
		}

	// Bank 2
	case address >= 0x4000 && address <= 0x5FFF:
		mbc.bank2 = value & 0b11

	// Banking Mode Select
	case address >= 0x6000 && address <= 0x7FFF:
		mbc.mode = value & 0b01

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

		// Default to Bank 0
		bank := uint32(0)
		offset := uint32(address - 0xA000)

		// Mode 1: Write to specific Bank
		if mbc.mode == 0b01 {
			bank = uint32(mbc.bank2)
		}

		actualAddress := ((bank << 13) | offset) & mbc.ramAddressMask
		mbc.cartridge.ram[actualAddress] = value
	}
}

func (mbc *Mbc1) Serialize(buf []byte) int {
	offset := 0

	buf[offset] = mbc.ramg
	offset++
	buf[offset] = mbc.bank1
	offset++
	buf[offset] = mbc.bank2
	offset++
	buf[offset] = mbc.mode
	offset++

	return offset
}

func (mbc *Mbc1) Deserialize(buf []byte) int {
	offset := 0

	mbc.ramg = buf[offset]
	offset++
	mbc.bank1 = buf[offset]
	offset++
	mbc.bank2 = buf[offset]
	offset++
	mbc.mode = buf[offset]
	offset++

	return offset
}
