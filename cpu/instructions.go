package cpu

type instruction struct {
	mnemonic      string
	operandLength uint8
	execute       func(cpu *CPU) uint8
}

var instructions = [256]instruction{
	0x00: {"NOP", 0, nop},
	0x01: {"LD BC, ", 2, nil},
	0xC3: {"JP a16", 2, jp_a16},
}
