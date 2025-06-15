package cpu

import (
	"fmt"

	"github.com/davidyorr/EchoGB/mmu"
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
	bus            mmu.Bus
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

func (cpu *CPU) ConnectBus(bus *mmu.MMU) {
	cpu.bus = bus
}

func (cpu *CPU) pushToStack16(returnAddress uint16) {
	highByte := uint8((returnAddress >> 8) & 0xFF)
	lowByte := uint8(returnAddress & 0xFF)
	cpu.sp--
	cpu.bus.Write(cpu.sp, highByte)
	cpu.sp--
	cpu.bus.Write(cpu.sp, lowByte)
}

func (cpu *CPU) Step() (uint8, error) {
	fmt.Println("Go: cpu.Step() -------------")
	fmt.Printf("  pc=0x%02X\n", cpu.pc)

	opcode := cpu.bus.Read(cpu.pc)
	fmt.Printf("Go: opcode 0x%0X\n", opcode)
	instruction := instructions[opcode]

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
	} else {
		fmt.Println("~~~~~~~~~~~~ not updating PC ~~~~~~~~~~~~~~~~~")
	}

	return cycles, nil
}

type Flag uint8

const (
	FlagZ Flag = 7
	FlagN Flag = 6
	FlagH Flag = 5
	FlagC Flag = 4
)

func (cpu *CPU) setFlag(flag Flag, value bool) {
	if value {
		// set the bit
		cpu.f |= (1 << flag)
	} else {
		// clear the bit
		cpu.f &= ^(1 << flag)
	}
}
