package cpu

import (
	"fmt"

	"github.com/davidyorr/EchoGB/bus"
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
	ime            bool
	immediateValue uint16
	bus            *bus.Bus
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
}

func (cpu *CPU) ConnectBus(bus *bus.Bus) {
	cpu.bus = bus
}

func (cpu *CPU) Step() (uint8, error) {
	fmt.Println("Go: cpu.Step() -------------")
	cpu.immediateValue = 0

	opcode := cpu.bus.Read(cpu.pc)
	instruction := instructions[opcode]
	fmt.Printf("  PC=0x%04X SP:0x%04X AF:0x%04X BC:0x%04X DE:0x%04X HL:0x%04X | (op:0x%02X, len:%d, imm:0x%04X) %s\n", cpu.pc, cpu.sp, cpu.getAF(), cpu.getBC(), cpu.getDE(), cpu.getHL(), opcode, instruction.operandLength, cpu.immediateValue, instruction.mnemonic)

	if instruction.execute == nil {
		fmt.Printf("Go: unimplemented instruction 0x%02X\n", opcode)
		return 0, fmt.Errorf("unimplemented instruction 0x%02X", opcode)
	}

	switch instruction.operandLength {
	case 1:
		cpu.immediateValue = uint16(cpu.bus.Read(cpu.pc + 1))
	case 2:
		lowByte := cpu.bus.Read(cpu.pc + 1)
		highByte := cpu.bus.Read(cpu.pc + 2)
		cpu.immediateValue = (uint16(highByte) << 8) | uint16(lowByte)
	}

	originalPc := cpu.pc
	cycles := instruction.execute(cpu)

	// don't update the PC if the opcode did
	if cpu.pc == originalPc {
		cpu.pc += 1 + uint16(instruction.operandLength)
	}

	return cycles, nil
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
