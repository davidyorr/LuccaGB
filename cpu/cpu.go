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
	ime bool
	bus mmu.Bus
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

func (cpu *CPU) Step() {
	fmt.Println("Go: cpu.Step()")
	// fetch
	b := cpu.bus.Read(cpu.pc)
	fmt.Printf("Go: read byte 0x%0X\n", b)
	cpu.pc++

	instruction := instructions[b]
	if instruction.execute != nil {
		instructions[b].execute()
	} else {
		fmt.Printf("Go: unimplemented instruction 0x%02X\n", b)
	}
}
