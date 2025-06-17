package mmu

import (
	"fmt"

	"github.com/davidyorr/EchoGB/cartridge"
	"github.com/davidyorr/EchoGB/timer"
)

type MMU struct {
	cartridge *cartridge.Cartridge
	timer     *timer.Timer
	// 0xC000 - 0xDFFF
	workingRam [8192]uint8
	// 0xFF00 - 0xFF7F
	ioRegisters [128]uint8
	// 0xFF80 - 0xFFFE
	highRam [127]uint8
	// for test output
	serialOutputBuffer []uint8
}

type Bus interface {
	Read(address uint16) uint8
	Write(address uint16, value uint8)
}

func New(cartridge *cartridge.Cartridge, timer *timer.Timer) *MMU {
	mmu := &MMU{}
	mmu.cartridge = cartridge
	mmu.timer = timer

	mmu.Reset()

	return mmu
}

func (mmu *MMU) Reset() {
	mmu.timer.Reset()
	mmu.ioRegisters[0xFF00-0xFF00] = 0xCF // P1
	mmu.ioRegisters[0xFF01-0xFF00] = 0x00 // SB
	mmu.ioRegisters[0xFF02-0xFF00] = 0x7E // SC
	// ioRegisters[0xFF03-0xFF00] = 0xFF // unused
	// mmu.ioRegisters[0xFF04-0xFF00] = 0x18 // DIV
	// mmu.ioRegisters[0xFF05-0xFF00] = 0x00 // TIMA
	// mmu.ioRegisters[0xFF06-0xFF00] = 0x00 // TMA
	// mmu.ioRegisters[0xFF07-0xFF00] = 0xF8 // TAC
	mmu.ioRegisters[0xFF0F-0xFF00] = 0xE1 // IF
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
	mmu.ioRegisters[0xFF40-0xFF00] = 0x91 // LCDC
	mmu.ioRegisters[0xFF41-0xFF00] = 0x81 // STAT
	mmu.ioRegisters[0xFF42-0xFF00] = 0x00 // SCY
	mmu.ioRegisters[0xFF43-0xFF00] = 0x00 // SCX
	mmu.ioRegisters[0xFF44-0xFF00] = 0x91 // LY
	mmu.ioRegisters[0xFF45-0xFF00] = 0x00 // LYC
	mmu.ioRegisters[0xFF46-0xFF00] = 0xFF // DMA
	mmu.ioRegisters[0xFF47-0xFF00] = 0xFC // BGP
	// left uninitialized
	// ioRegisters[0xFF48-0xFF00] = 0x00 // OBP0
	// ioRegisters[0xFF49-0xFF00] = 0x00 // OBP1
	mmu.ioRegisters[0xFF4A-0xFF00] = 0x00 // WY
	mmu.ioRegisters[0xFF4B-0xFF00] = 0x00 // WX
}

func (mmu *MMU) Step(cycles uint8) {
	mmu.timer.Step(cycles)
}

func (mmu *MMU) Read(address uint16) uint8 {
	var value uint8 = 0
	if address <= 0x7FFF {
		// ROM
		value = mmu.cartridge.Read(address)
	} else if address >= 0xC000 && address <= 0xDFFF {
		// working RAM
		value = mmu.workingRam[address-0xC000]
	} else if address >= 0xFF00 && address <= 0xFF7F {
		// IO registers
		// timers
		if address >= 0xFF00 && address <= 0xFF07 {
			value = mmu.timer.Read(address)
		} else {
			value = mmu.ioRegisters[address-0xFF00]
		}
	} else if address >= 0xFF80 && address <= 0xFFFE {
		// high RAM
		value = mmu.highRam[address-0xFF80]
	}
	fmt.Printf("  [MMU READ] Address: 0x%04X, Value: 0x%02X\n", address, value)

	return value
}

func (mmu *MMU) Write(address uint16, value uint8) {
	fmt.Printf("  [MMU WRITE] Address: 0x%04X, Value: 0x%02X\n", address, value)
	if address <= 0x7FFF {
		// ROM
	} else if address >= 0xC000 && address <= 0xDFFF {
		// working RAM
		mmu.workingRam[address-0xC000] = value
	} else if address >= 0xFF00 && address <= 0xFF7F {
		// IO registers
		// timers
		if address >= 0xFF00 && address <= 0xFF07 {
			mmu.timer.Write(address, value)
		} else {
			mmu.ioRegisters[address-0xFF00] = value
		}
	} else if address >= 0xFF80 && address <= 0xFFFE {
		// high RAM
		mmu.highRam[address-0xFF80] = value
	}
	// SB is the Serial Data register at address 0xFF01
	// SC is the Serial Control register at address 0xFF02
	if address == 0xFF02 && value == 0x81 {
		fmt.Println("========================== WRITING TO SERIAL OUTPUT BUFFER ================================")
		fmt.Println("                           ", mmu.Read(0xFF01))
		mmu.serialOutputBuffer = append(mmu.serialOutputBuffer, mmu.Read(0xFF01))
	}
}

func (mmu *MMU) SerialOutputBuffer() []byte {
	return mmu.serialOutputBuffer
}
