package mmu

import (
	"fmt"

	"github.com/davidyorr/LuccaGB/internal/cartridge"
	"github.com/davidyorr/LuccaGB/internal/debug"
	"github.com/davidyorr/LuccaGB/internal/interrupt"
	"github.com/davidyorr/LuccaGB/internal/joypad"
	"github.com/davidyorr/LuccaGB/internal/logger"
)

type MMU struct {
	cartridge *cartridge.Cartridge
	joypad    *joypad.Joypad
	// 0xC000 - 0xDFFF
	workingRam [8192]uint8
	// 0xFF00 - 0xFF7F
	ioRegisters [128]uint8
	// 0xFF80 - 0xFFFE
	highRam [127]uint8
	// 0xFFFF - Interrupt enable
	ieRegister uint8
	// 0xFF0F - Interrupt flag
	ifRegister uint8
}

func New(cartridge *cartridge.Cartridge) *MMU {
	mmu := &MMU{}
	mmu.cartridge = cartridge

	mmu.Reset()

	return mmu
}

func (mmu *MMU) Reset() {
}

func (mmu *MMU) ConnectJoypad(joypad *joypad.Joypad) {
	mmu.joypad = joypad
}

func (mmu *MMU) Read(address uint16) (value uint8) {
	switch {
	// ROM
	case address <= 0x7FFF:
		value = mmu.cartridge.Read(address)
	// External RAM
	case address >= 0xA000 && address <= 0xBFFF:
		value = mmu.cartridge.Read(address)
	// working RAM
	case address >= 0xC000 && address <= 0xDFFF:
		value = mmu.workingRam[address-0xC000]
	// echo RAM
	case address >= 0xE000 && address <= 0xFDFF:
		value = mmu.workingRam[address-0xE000]
	// IF
	case address == 0xFF0F:
		value = mmu.ifRegister | 0b1110_0000
	// IE
	case address == 0xFFFF:
		value = mmu.ieRegister
	// JOYPAD
	case address == 0xFF00:
		value = mmu.joypad.Read()
	// IO registers
	case address >= 0xFF01 && address <= 0xFF7F:
		value = mmu.ioRegisters[address-0xFF00]

		// unused bits in IO registers should return 1
		value |= 0b1111_1111
	// high RAM
	case address >= 0xFF80 && address <= 0xFFFE:
		value = mmu.highRam[address-0xFF80]
	default:
		logger.Info("unhandled address while reading ->", "ADDRESS", fmt.Sprintf("%04X", address))
		value = 0xFF
	}

	if debug.Enabled {
		logger.Debug(
			"MMU READ",
			"Address", fmt.Sprintf("0x%04X", address),
			"Value", fmt.Sprintf("0x%02X", value),
		)
	}

	return value
}

func (mmu *MMU) Write(address uint16, value uint8) {
	if debug.Enabled {
		logger.Debug(
			"MMU Write",
			"Address", fmt.Sprintf("0x%04X", address),
			"Value", fmt.Sprintf("0x%02X", value),
		)
	}
	switch {
	// ROM
	case address <= 0x7FFF:
		mmu.cartridge.Write(address, value)
	// External RAM
	case address >= 0xA000 && address <= 0xBFFF:
		mmu.cartridge.Write(address, value)
	// working RAM
	case address >= 0xC000 && address <= 0xDFFF:
		mmu.workingRam[address-0xC000] = value
	// echo RAM
	case address >= 0xE000 && address <= 0xFDFF:
		mmu.workingRam[address-0xE000] = value
	// IF
	case address == 0xFF0F:
		mmu.ifRegister = value & 0b0001_1111
	// IE
	case address == 0xFFFF:
		mmu.ieRegister = value
	// JOYPAD
	case address == 0xFF00:
		mmu.joypad.Write(value)
	// IO registers
	case address >= 0xFF01 && address <= 0xFF7F:
		mmu.ioRegisters[address-0xFF00] = value
	// high RAM
	case address >= 0xFF80 && address <= 0xFFFE:
		mmu.highRam[address-0xFF80] = value
	default:
		logger.Info("unhandled address while writing <-", "ADDRESS", fmt.Sprintf("%04X", address))
	}
}

func (mmu *MMU) RequestInterrupt(interrupt interrupt.Interrupt) {
	mmu.ifRegister |= uint8(interrupt)
}

func (mmu *MMU) ClearInterrupt(interrupt interrupt.Interrupt) {
	mmu.ifRegister &= ^uint8(interrupt)
}

func (mmu *MMU) InterruptEnable() uint8 {
	return mmu.ieRegister | 0b1110_0000
}

func (mmu *MMU) InterruptFlag() uint8 {
	return mmu.ifRegister | 0b1110_0000
}
