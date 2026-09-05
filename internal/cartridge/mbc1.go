package cartridge

import (
	"bytes"
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

	// 1 MiB Multi-Game Compilation Carts
	// See: https://gbdev.io/pandocs/MBC1.html#mbc1m-1-mib-multi-game-compilation-carts
	mbc1m bool

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

// See: https://gbdev.io/pandocs/The_Cartridge_Header.html#0104-0133--nintendo-logo
var nintendoLogo = []byte{
	0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B, 0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D,
	0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E, 0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99,
	0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC, 0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E,
}

func newMbc1(cartridge *Cartridge) *Mbc1 {
	mbc1 := &Mbc1{}

	mbc1.cartridge = cartridge
	mbc1.romAddressMask = addressMaskSizes[cartridge.romSizeCode]
	mbc1.ramAddressMask = ramAddressMaskSizes[cartridge.ramSizeCode]
	cartridge.ram = make([]uint8, ramSizes[cartridge.ramSizeCode])

	// check if MBC1M
	const bankSize = 0x4000
	const bank10Offset = 0x10 * bankSize
	const nintendoLogoStart = bank10Offset + 0x104
	const nintendoLogoEnd = bank10Offset + 0x134

	mbc1.mbc1m = false
	if cartridge.romSizeCode == 0x05 && len(cartridge.rom) >= nintendoLogoEnd {
		mbc1.mbc1m = bytes.Equal(cartridge.rom[nintendoLogoStart:nintendoLogoEnd], nintendoLogo)
	}

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
			var bank uint32
			if mbc.mbc1m {
				// Bank 2 selects Banks $00, $10, $20, or $30
				bank = uint32(mbc.bank2) << 4
			} else {
				// Bank 2 selects Banks $00, $20, $40, or $60
				bank = uint32(mbc.bank2) << 5
			}
			actualAddress := ((bank << 14) | uint32(address)) & mbc.romAddressMask
			return mbc.cartridge.rom[actualAddress]
		}

	// ROM BANK 01-7F
	case address >= 0x4000 && address <= 0x7FFF:
		var bank uint32
		if mbc.mbc1m {
			// ignore the top bit of the main ROM banking register
			bank = (uint32(mbc.bank2) << 4) | uint32(mbc.bank1&0x0F)
		} else {
			// lower 5 bits from Bank 1, upper 2 bits from Bank 2
			bank = (uint32(mbc.bank2) << 5) | uint32(mbc.bank1)
		}
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
