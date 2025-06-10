package cpu

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
