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
	// ioRegisters[0xFF03-0xFF00] = 0xFF // unused
	mmu.ioRegisters[0xFF10-0xFF00] = 0x80 // NR10
	mmu.ioRegisters[0xFF11-0xFF00] = 0xBF // NR11
	mmu.ioRegisters[0xFF12-0xFF00] = 0xF3 // NR12
	mmu.ioRegisters[0xFF13-0xFF00] = 0xFF // NR13
	mmu.ioRegisters[0xFF14-0xFF00] = 0xBF // NR14
	mmu.ioRegisters[0xFF16-0xFF00] = 0x3F // NR21
	mmu.ioRegisters[0xFF17-0xFF00] = 0x00 // NR22
	mmu.ioRegisters[0xFF18-0xFF00] = 0xFF // NR23
	mmu.ioRegisters[0xFF19-0xFF00] = 0xBF // NR24
	mmu.ioRegisters[0xFF1A-0xFF00] = 0x7F // NR30
	mmu.ioRegisters[0xFF1B-0xFF00] = 0xFF // NR31
	mmu.ioRegisters[0xFF1C-0xFF00] = 0x9F // NR32
	mmu.ioRegisters[0xFF1D-0xFF00] = 0xFF // NR33
	mmu.ioRegisters[0xFF1E-0xFF00] = 0xBF // NR34
	mmu.ioRegisters[0xFF20-0xFF00] = 0xFF // NR41
	mmu.ioRegisters[0xFF21-0xFF00] = 0x00 // NR42
	mmu.ioRegisters[0xFF22-0xFF00] = 0x00 // NR43
	mmu.ioRegisters[0xFF23-0xFF00] = 0xBF // NR44
	mmu.ioRegisters[0xFF24-0xFF00] = 0x77 // NR50
	mmu.ioRegisters[0xFF25-0xFF00] = 0xF3 // NR51
	mmu.ioRegisters[0xFF26-0xFF00] = 0xF1 // NR52
	mmu.ioRegisters[0xFF46-0xFF00] = 0xFF // DMA
	mmu.ioRegisters[0xFF47-0xFF00] = 0xFC // BGP
	// left uninitialized
	// ioRegisters[0xFF48-0xFF00] = 0x00 // OBP0
	// ioRegisters[0xFF49-0xFF00] = 0x00 // OBP1
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
		switch {
		// NR10
		case address == 0xFF10:
			value |= 0b1000_0000
		// NR30
		case address == 0xFF1A:
			value |= 0b0111_1111
		// NR31
		case address == 0xFF1B:
			value |= 0b1000_0001
		// NR32
		case address == 0xFF1C:
			value |= 0b1001_1111
		// NR41
		case address == 0xFF20:
			value |= 0b1100_0000
		// NR44
		case address == 0xFF23:
			value |= 0b0011_1111
		// NR52
		case address == 0xFF26:
			value |= 0b0111_0000
		// unmapped
		default:
			value |= 0b1111_1111
		}
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
