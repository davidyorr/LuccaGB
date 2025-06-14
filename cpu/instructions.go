package cpu

type instruction struct {
	mnemonic      string
	operandLength uint8
	execute       func(cpu *CPU) uint8
}

// how to read opcode table
// ---------
// LD A, d8
//   2 8
// ---------
// 2 is length, 8 is duration
// 2 - 1 byte opcode 0x3E, 1 byte operand d8. set the operandLength to this value
// 8 - duration in T cycles, what execute() should return

var instructions = [256]instruction{
	// misc instructions
	0x00: {"NOP", 0, nop},

	// load instructions
	0x01: {"LD BC, ", 2, nil},
	0x31: {"LD SP, d16", 2, ld_sp_d16},
	0x3E: {"LD A, d8", 2, ld_a_d8},
	0xE0: {"LDH (a8), A", 1, ldh_a8_a},
	0xEA: {"LD (a16), A", 2, ld_a16_a},

	// bit shift instructions
	0x07: {"RLCA", 1, nil},

	// jumps and subroutine instructions
	0xC3: {"JP a16", 2, jp_a16},

	// interrupt related instructions
	0xF3: {"DI", 0, di},
}
