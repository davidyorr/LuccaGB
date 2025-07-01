package cpu

import (
	"fmt"

	"github.com/davidyorr/EchoGB/bus"
	"github.com/davidyorr/EchoGB/interrupt"
	"github.com/davidyorr/EchoGB/logger"
)

type CPU struct {
	// program counter
	pc uint16
	// stack pointer
	sp uint16
	// accumulator
	a uint8
	// flags register, lower 4 bits are unused and always 0
	f uint8
	b uint8
	c uint8
	d uint8
	e uint8
	h uint8
	l uint8
	// interrupt master enable flag
	ime          bool
	imeScheduled bool
	// the current instruction
	instruction *instruction
	// the current opcode
	opcode uint8
	// the current CB opcode
	cbOpcode *uint8
	// the current immediate value
	immediateValue uint16
	// the current M-cycle, 1 indexed
	mCycle uint8
	// Memory Data Register
	//	for storing temporary values used across different M-cycles
	//	the Game Boy hardware does not have a MDR, but instead uses latches
	mdr           uint8
	halted        bool
	haltBugActive bool
	bus           *bus.Bus
}

func New() *CPU {
	cpu := &CPU{}

	cpu.Reset()

	return cpu
}

func (cpu *CPU) Reset() {
	cpu.a = 0x01
	cpu.f = 0x80
	cpu.b = 0x00
	cpu.c = 0x13
	cpu.d = 0x00
	cpu.e = 0xD8
	cpu.h = 0x01
	cpu.l = 0x4D
	cpu.pc = 0x0100
	cpu.sp = 0xFFFE
	cpu.halted = false
	cpu.haltBugActive = false
}

func (cpu *CPU) ConnectBus(bus *bus.Bus) {
	cpu.bus = bus
}

// Return the number of T-Cycles taken
func (cpu *CPU) Step() (uint8, error) {
	if cpu.halted {
		return 4, nil
	}
	haltBugWasActive := cpu.haltBugActive
	if cpu.haltBugActive {
		cpu.haltBugActive = false
	}
	if cpu.imeScheduled {
		cpu.ime = true
		cpu.imeScheduled = false
	}

	// start new instruction
	if cpu.instruction == nil {
		cpu.opcode = cpu.bus.Read(cpu.pc)
		cpu.instruction = &instructions[cpu.opcode]
		cpu.mCycle = 1
		cpu.immediateValue = 0
	}

	// handle cb prefixed instructions
	if cpu.opcode == 0xCB {
		if cpu.cbOpcode == nil {
			cbOpcode := cpu.bus.Read(cpu.pc + 1)
			cpu.cbOpcode = &cbOpcode
		}
		tCycles, done := cpu.executeCbInstructionStep(*cpu.cbOpcode)

		if done {
			cpu.pc += 2
			cpu.instruction = nil
			cpu.cbOpcode = nil
		}

		cpu.mCycle++

		return tCycles, nil
	}

	var cbo uint8 = 0
	if cpu.cbOpcode != nil {
		cbo = *cpu.cbOpcode
	}

	originalPc := cpu.pc
	tCycles, done := cpu.instruction.step(cpu)

	logger.Info(
		"CPU STEP",
		"PC", fmt.Sprintf("0x%04X", originalPc),
		"SP", fmt.Sprintf("0x%04X", cpu.sp),
		"AF", fmt.Sprintf("0x%04X", cpu.getAF()),
		"BC", fmt.Sprintf("0x%04X", cpu.getBC()),
		"DE", fmt.Sprintf("0x%04X", cpu.getDE()),
		"HL", fmt.Sprintf("0x%04X", cpu.getHL()),
		"op", fmt.Sprintf("(op:0x%02X, len:%d, imm:0x%04X)", cpu.opcode, cpu.instruction.operandLength, cpu.immediateValue),
		"cb", fmt.Sprintf("0x%02X", cbo),
		"instruction", cpu.instruction.mnemonic,
	)

	if done {
		// don't update the PC if the opcode did
		if cpu.pc == originalPc {
			if haltBugWasActive {
				// don't update the PC if the halt bug was active
			} else {
				cpu.pc += 1 + uint16(cpu.instruction.operandLength)
			}
		}
		// reset state
		cpu.instruction = nil
	}

	cpu.mCycle++

	return tCycles, nil
}

func (cpu *CPU) HandleInterrupts() {
	ifRegister := cpu.bus.Read(0xFF0F)
	ieRegister := cpu.bus.Read(0xFFFF)
	var interruptType interrupt.Interrupt
	var vectorAddress uint16

	if (ifRegister&ieRegister)&uint8(interrupt.VBlankInterrupt) != 0 {
		interruptType = interrupt.VBlankInterrupt
		vectorAddress = 0x0040
	} else if (ifRegister&ieRegister)&uint8(interrupt.LcdInterrupt) != 0 {
		interruptType = interrupt.LcdInterrupt
		vectorAddress = 0x0048
	} else if (ifRegister&ieRegister)&uint8(interrupt.TimerInterrupt) != 0 {
		interruptType = interrupt.TimerInterrupt
		vectorAddress = 0x0050
	} else if (ifRegister&ieRegister)&uint8(interrupt.SerialInterrupt) != 0 {
		interruptType = interrupt.SerialInterrupt
		vectorAddress = 0x0058
	} else if (ifRegister&ieRegister)&uint8(interrupt.JoypadInterrupt) != 0 {
		interruptType = interrupt.JoypadInterrupt
		vectorAddress = 0x0060
	}

	cpu.ime = false
	clearedFlag := ifRegister & ^uint8(interruptType)
	cpu.bus.Write(0xFF0F, clearedFlag)
	cpu.pushToStack16(cpu.pc)
	cpu.pc = vectorAddress
}

func (cpu *CPU) pushToStack16(returnAddress uint16) {
	highByte := uint8(returnAddress >> 8)
	lowByte := uint8(returnAddress & 0x00FF)
	cpu.sp--
	cpu.bus.Write(cpu.sp, highByte)
	cpu.sp--
	cpu.bus.Write(cpu.sp, lowByte)
}

func (cpu *CPU) popFromStack16() uint16 {
	lowByte := cpu.bus.Read(cpu.sp)
	cpu.sp++
	highByte := cpu.bus.Read(cpu.sp)
	cpu.sp++

	return (uint16(highByte) << 8) | uint16(lowByte)
}

func (cpu *CPU) getAF() uint16 {
	return (uint16(cpu.a) << 8) | uint16(cpu.f)
}

func (cpu *CPU) setAF(value uint16) {
	cpu.a = uint8(value >> 8)
	lowByte := uint8(value & 0x00FF)
	// clear the lower 4 bits because the lower 4 bits of the F register must always be 0
	cpu.f = lowByte & 0xF0
}

func (cpu *CPU) getBC() uint16 {
	return (uint16(cpu.b) << 8) | uint16(cpu.c)
}

func (cpu *CPU) setBC(value uint16) {
	cpu.b = uint8(value >> 8)
	cpu.c = uint8(value & 0x00FF)
}

func (cpu *CPU) getDE() uint16 {
	return (uint16(cpu.d) << 8) | uint16(cpu.e)
}

func (cpu *CPU) setDE(value uint16) {
	cpu.d = uint8(value >> 8)
	cpu.e = uint8(value & 0x00FF)
}

func (cpu *CPU) getHL() uint16 {
	return (uint16(cpu.h) << 8) | uint16(cpu.l)
}

func (cpu *CPU) setHL(value uint16) {
	cpu.h = uint8(value >> 8)
	cpu.l = uint8(value & 0x00FF)
}

// return the value and the number of cycles it took
func (cpu *CPU) get_r8(r8 uint8) uint8 {
	var value uint8 = 0

	switch r8 {
	case 0b000:
		value = cpu.b
	case 0b001:
		value = cpu.c
	case 0b010:
		value = cpu.d
	case 0b011:
		value = cpu.e
	case 0b100:
		value = cpu.h
	case 0b101:
		value = cpu.l
	case 0b111:
		value = cpu.a
	case 0b110:
		// special case
		value = cpu.bus.Read(cpu.getHL())
	}

	return value
}

// return the number of cycles it took
func (cpu *CPU) set_r8(r8 uint8, value uint8) {
	switch r8 {
	case 0b000:
		cpu.b = value
	case 0b001:
		cpu.c = value
	case 0b010:
		cpu.d = value
	case 0b011:
		cpu.e = value
	case 0b100:
		cpu.h = value
	case 0b101:
		cpu.l = value
	case 0b111:
		cpu.a = value
	case 0b110:
		// special case
		cpu.bus.Write(cpu.getHL(), value)
	}
}

type Flag uint8

const (
	// zero flag
	FlagZ Flag = 7
	// subtraction flag (BCD)
	FlagN Flag = 6
	// half carry flag (BCD)
	FlagH Flag = 5
	// carry flag
	FlagC Flag = 4
)

func (cpu *CPU) getFlag(flag Flag) bool {
	return (cpu.f & (1 << flag)) != 0
}

func (cpu *CPU) setFlag(flag Flag, value bool) {
	if value {
		// set the bit
		cpu.f |= (1 << flag)
	} else {
		// clear the bit
		cpu.f &= ^(1 << flag)
	}
}

func (cpu *CPU) fetchImmLowByte() {
	lowByte := cpu.bus.Read(cpu.pc + 1)
	cpu.immediateValue |= uint16(lowByte)
}

func (cpu *CPU) fetchImmHighByte() {
	highByte := cpu.bus.Read(cpu.pc + 2)
	cpu.immediateValue |= (uint16(highByte) << 8)
}

func (cpu *CPU) InterruptMasterEnable() bool {
	return cpu.ime
}

func (cpu *CPU) ScheduleIme() {
	cpu.imeScheduled = true
}

func (cpu *CPU) Halt() {
	cpu.halted = true
}

func (cpu *CPU) Unhalt() {
	cpu.halted = false
}

func (cpu *CPU) Halted() bool {
	return cpu.halted
}
