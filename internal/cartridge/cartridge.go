package cartridge

import (
	"bytes"
	"fmt"

	"github.com/davidyorr/LuccaGB/internal/logger"
)

type Cartridge struct {
	// 0x0000 - 0x7FFF, dynamically sized
	// holds the raw bytes from the loaded ROM
	rom []uint8
	// holds the External RAM (SRAM)
	ram []uint8
	// title of the ROM in uppercase ASCII
	title []uint8
	// memory bank controller
	mbc MBC
	// memory bank controller type
	mbcType uint8
	// ROM size
	romSizeCode uint8
	// RAM size, if any
	ramSizeCode uint8
}

func New() *Cartridge {
	cartridge := &Cartridge{}

	return cartridge
}

func (cartridge *Cartridge) LoadRom(rom []uint8) {
	cartridge.rom = rom

	cartridge.title = bytes.Trim(cartridge.rom[0x0134:0x0143], "\x00")
	cartridge.romSizeCode = cartridge.rom[0x148]
	cartridge.ramSizeCode = cartridge.rom[0x149]

	cartridge.mbcType = cartridge.rom[0x147]
	switch cartridge.mbcType {
	// ROM only
	case 0x00:
		cartridge.mbc = nil
	// MBC1
	case 0x01:
		cartridge.mbc = newMbc1(cartridge)
	// MBC1+RAM
	case 0x02:
		cartridge.mbc = newMbc1(cartridge)
	// MBC1+RAM+BATTERY
	case 0x03:
		cartridge.mbc = newMbc1(cartridge)
	// MBC2
	case 0x05:
		cartridge.mbc = newMbc2(cartridge)
	// MBC2+BATTERY
	case 0x06:
		cartridge.mbc = newMbc2(cartridge)
	// MBC5
	case 0x19:
		cartridge.mbc = newMbc5(cartridge)
	// MBC5+RAM
	case 0x1A:
		cartridge.mbc = newMbc5(cartridge)
	// MBC5+RAM+BATTERY
	case 0x1B:
		cartridge.mbc = newMbc5(cartridge)
	// MBC5+RUMBLE
	case 0x1C:
		cartridge.mbc = newMbc5(cartridge)
	// MBC5+RUMBLE+RAM
	case 0x1D:
		cartridge.mbc = newMbc5(cartridge)
	// MBC5+RUMBLE+RAM+BATTERY
	case 0x1E:
		cartridge.mbc = newMbc5(cartridge)
	default:
		cartridge.mbc = nil
	}

	logger.Info(
		"CARTRIDGE LOAD ROM",
		"TITLE", string(cartridge.title),
		"MBC", fmt.Sprintf("0x%02X", cartridge.mbcType),
		"ROM_SIZE_CODE", fmt.Sprintf("0x%02X", cartridge.romSizeCode),
		"RAM_SIZE_CODE", fmt.Sprintf("0x%02X", cartridge.ramSizeCode),
	)
}

func (cartridge *Cartridge) Read(address uint16) uint8 {
	if cartridge.mbc == nil {
		return cartridge.rom[address]
	}

	return cartridge.mbc.Read(address)
}

func (cartridge *Cartridge) Write(address uint16, value uint8) {
	if cartridge.mbc == nil {
		return
	}

	cartridge.mbc.Write(address, value)
}

// Debug gathers the current state of the Cartridge into a structured map.
func (cartridge *Cartridge) Debug() map[string]interface{} {
	return map[string]interface{}{
		"title":       string(cartridge.title),
		"mbcType":     cartridge.mbcType,
		"romSizeCode": cartridge.romSizeCode,
		"ramSizeCode": cartridge.ramSizeCode,
	}
}
