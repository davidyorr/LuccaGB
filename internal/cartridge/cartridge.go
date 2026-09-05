package cartridge

import (
	"bytes"
	"encoding/binary"
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

var ramTypes = map[uint8]bool{
	0x02: true, // MBC1+RAM
	0x03: true, // MBC1+RAM+BATTERY
	0x08: true, // ROM+RAM
	0x09: true, // ROM+RAM+BATTERY
	0x0C: true, // MMM01+RAM
	0x0D: true, // MMM01+RAM+BATTERY
	0x10: true, // MBC3+TIMER+RAM+BATTERY
	0x12: true, // MBC3+RAM
	0x13: true, // MBC3+RAM+BATTERY
	0x1A: true, // MBC5+RAM
	0x1B: true, // MBC5+RAM+BATTERY
	0x1D: true, // MBC5+RUMBLE+RAM
	0x1E: true, // MBC5+RUMBLE+RAM+BATTERY
	0x22: true, // MBC7+SENSOR+RUMBLE+RAM+BATTERY
	0xFF: true, // HuC1+RAM+BATTERY
}

// The largest RAM size code supported by each cartridge type. Used as a
// fallback when the header declares RAM size 0 but the cartridge type requires
// RAM. No commercial games should have this bug, but some homebrew ROMs do,
// including Blargg's halt_bug.gb.
var maxRamSizeCodeForType = map[uint8]uint8{
	// MBC1: 32KiB
	// See: https://gbdev.io/pandocs/MBC1.html
	0x02: 0x03, // MBC1+RAM
	0x03: 0x03, // MBC1+RAM+BATTERY

	// No MBC: no banking, optionally up to 8KiB
	// See: https://gbdev.io/pandocs/nombc.html
	0x08: 0x02, // ROM+RAM
	0x09: 0x02, // ROM+RAM+BATTERY

	// MMM01: emulates an MBC1, 32KiB
	// See: https://gbdev.io/pandocs/MMM01.html
	0x0C: 0x03,
	0x0D: 0x03,

	// MBC3: 32KiB
	// See: https://gbdev.io/pandocs/MBC3.html
	0x10: 0x03,
	0x12: 0x03,
	0x13: 0x03,

	// MBC5: 128KiB
	// See: https://gbdev.io/pandocs/MBC5.html
	0x1A: 0x04,
	0x1B: 0x04,
	0x1D: 0x04,
	0x1E: 0x04,

	// MBC7: 8KiB
	// See: https://gbdev.io/pandocs/MBC7.html
	0x22: 0x02,

	// HuC1: banks like MBC1, 32KiB
	// See: https://gbdev.io/pandocs/HuC1.html
	0xFF: 0x03,
}

func New() *Cartridge {
	cartridge := &Cartridge{}

	return cartridge
}

type CartridgeInfo struct {
	Title       string
	RomSizeCode int
	RamSizeCode int
	RamSize     int
	RomSize     int
	HasBattery  bool
	Type        int
}

func (cartridge *Cartridge) LoadRom(rom []uint8) CartridgeInfo {
	cartridge.rom = rom

	cartridge.title = bytes.Trim(cartridge.rom[0x0134:0x0143], "\x00")
	cartridge.romSizeCode = cartridge.rom[0x148]
	cartridge.ramSizeCode = cartridge.rom[0x149]
	cartridge.cartridgeType = cartridge.rom[0x147]

	if cartridge.ramSizeCode == 0x00 && ramTypes[cartridge.cartridgeType] {
		fallbackCode, ok := maxRamSizeCodeForType[cartridge.cartridgeType]
		if !ok {
			fallbackCode = 0x02
		}
		logger.Warn(
			"cartridge type requires RAM but header declares size 0, using mapper's max addressable RAM as fallback",
			"TYPE", fmt.Sprintf("0x%02X", cartridge.cartridgeType),
			"FALLBACK_RAM_SIZE_CODE", fmt.Sprintf("0x%02X", fallbackCode),
		)
		cartridge.ramSizeCode = fallbackCode
	}

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
		Title:       string(cartridge.title),
		RomSizeCode: int(cartridge.romSizeCode),
		RamSizeCode: int(cartridge.ramSizeCode),
		RamSize:     len(cartridge.ram),
		RomSize:     len(cartridge.rom),
		HasBattery:  cartridge.hasBattery,
		Type:        int(cartridge.cartridgeType),
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
		if int(address) >= len(cartridge.rom) {
			return 0xFF
		}
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

func (cart *Cartridge) Serialize(buf []byte) int {
	offset := 0

	ramLen := uint32(len(cart.ram))
	binary.LittleEndian.PutUint32(buf[offset:], ramLen)
	offset += 4

	if ramLen > 0 {
		n := copy(buf[offset:], cart.ram)
		offset += n
	}

	if cart.hasBattery {
		buf[offset] = 1
	} else {
		buf[offset] = 0
	}
	offset++

	// Skip Title, CartType, RomSize, RamSize because they are
	// read-only properties of the ROM file itself.

	if cart.mbc != nil {
		offset += cart.mbc.Serialize(buf[offset:])
	}

	return offset
}

func (cart *Cartridge) Deserialize(buf []byte) int {
	offset := 0

	savedRamLen := int(binary.LittleEndian.Uint32(buf[offset:]))
	offset += 4

	if len(cart.ram) != savedRamLen {
		if savedRamLen > 0 {
			if cap(cart.ram) < savedRamLen {
				cart.ram = make([]uint8, savedRamLen)
			} else {
				cart.ram = cart.ram[:savedRamLen]
			}
		}
	}

	if savedRamLen > 0 {
		n := copy(cart.ram, buf[offset:])
		offset += n
	}

	cart.hasBattery = buf[offset] == 1
	offset++

	if cart.mbc != nil {
		offset += cart.mbc.Deserialize(buf[offset:])
	}

	return offset
}
