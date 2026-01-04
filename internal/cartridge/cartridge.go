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
	// memory bank controller
	mbc MBC
	// for persisting RAM
	hasBattery bool

	// 0x0134 - 0x0143 Title of the ROM in uppercase ASCII
	title []uint8
	// 0x0147 - Cartridge type
	cartridgeType uint8
	// 0x0148 ROM size
	romSizeCode uint8
	// 0x0149 RAM size, if any
	ramSizeCode uint8
}

var batteryBackedTypes = map[uint8]bool{
	0x03: true, // MBC1+RAM+BATTERY
	0x06: true, // MBC2+BATTERY
	0x09: true, // ROM+RAM+BATTERY
	0x0D: true, // MMM01+RAM+BATTERY
	0x0F: true, // MBC3+TIMER+BATTERY
	0x10: true, // MBC3+TIMER+RAM+BATTERY
	0x13: true, // MBC3+RAM+BATTERY
	0x1B: true, // MBC5+RAM+BATTERY
	0x1E: true, // MBC5+RUMBLE+RAM+BATTERY
	0x22: true, // MBC7+SENSOR+RUMBLE+RAM+BATTERY
	0xFF: true, // HuC1+RAM+BATTERY
}

func New() *Cartridge {
	cartridge := &Cartridge{}

	return cartridge
}

type CartridgeInfo struct {
	Title      string
	RamSize    int
	HasBattery bool
}

func (cartridge *Cartridge) LoadRom(rom []uint8) CartridgeInfo {
	cartridge.rom = rom

	cartridge.title = bytes.Trim(cartridge.rom[0x0134:0x0143], "\x00")
	cartridge.romSizeCode = cartridge.rom[0x148]
	cartridge.ramSizeCode = cartridge.rom[0x149]
	cartridge.cartridgeType = cartridge.rom[0x147]

	if batteryBackedTypes[cartridge.cartridgeType] {
		cartridge.hasBattery = true
	}

	switch cartridge.cartridgeType {
	// ROM only
	case 0x00:
		cartridge.mbc = nil
	// MBC1
	case 0x01, 0x02, 0x03:
		cartridge.mbc = newMbc1(cartridge)
	// MBC2
	case 0x05, 0x06:
		cartridge.mbc = newMbc2(cartridge)
	// MBC5
	case 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E:
		cartridge.mbc = newMbc5(cartridge)
	default:
		cartridge.mbc = nil
	}

	logger.Info(
		"CARTRIDGE LOAD ROM",
		"TITLE", string(cartridge.title),
		"TYPE", fmt.Sprintf("0x%02X", cartridge.cartridgeType),
		"ROM_SIZE_CODE", fmt.Sprintf("0x%02X", cartridge.romSizeCode),
		"RAM_SIZE_CODE", fmt.Sprintf("0x%02X", cartridge.ramSizeCode),
	)

	return CartridgeInfo{
		Title:      string(cartridge.title),
		RamSize:    len(cartridge.ram),
		HasBattery: cartridge.hasBattery,
	}
}

func (cartridge *Cartridge) SetRam(ram []uint8) {
	if cartridge.hasBattery && len(cartridge.ram) > 0 && len(ram) > 0 {
		n := min(len(cartridge.ram), len(ram))
		copy(cartridge.ram[:n], ram[:n])
	}

	if len(cartridge.ram) != len(ram) {
		logger.Warn(
			"SetRam() RAM size mismatch: cart=%d persisted=%d",
			len(cartridge.ram),
			len(ram),
		)
	}
}

func (cartridge *Cartridge) Ram() []uint8 {
	if cartridge.hasBattery && len(cartridge.ram) > 0 {
		return cartridge.ram
	}

	return nil
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
		"title":         string(cartridge.title),
		"cartridgeType": cartridge.cartridgeType,
		"romSizeCode":   cartridge.romSizeCode,
		"ramSizeCode":   cartridge.ramSizeCode,
	}
}
