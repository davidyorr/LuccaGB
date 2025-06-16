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
// 2 - 1 byte opcode 0x3E, 1 byte operand d8. set the operandLength to this value minus 1 (to exclude the opcode byte)
// 8 - duration in T cycles, what execute() should return

var instructions = [256]instruction{
	// misc instructions
	0x00: {"NOP", 0, nop},

	// load instructions
	0x01: {"LD BC, n16", 2, ld_bc_n16},
	0x21: {"LD HL, n16", 2, ld_hl_n16},
	0x31: {"LD SP, n16", 2, ld_sp_n16},
	0x3E: {"LD A, n8", 1, ld_a_n8},
	0x78: {"LD A, B", 0, ld_a_b},
	0x7C: {"LD A, H", 0, ld_a_h},
	0x7D: {"LD A, L", 0, ld_a_l},
	0xE0: {"LDH (a8), A", 1, ldh_a8_a},
	0xEA: {"LD (a16), A", 2, ld_a16_a},
	0x2A: {"LD A, [HL+]", 0, ld_a_hli},

	// 8-bit arithmetic instructions
	0x03: {"INC BC", 0, inc_bc},
	0x23: {"INC HL", 0, inc_hl},
	0xD6: {"SUB A, n8", 1, sub_a_n8},

	// bit shift instructions
	0x07: {"RLCA", 0, rlca},

	// jumps and subroutine instructions
	0x18: {"JR e8", 1, jr_e8},
	0xC3: {"JP a16", 2, jp_a16},
	0xCD: {"CALL a16", 2, call_a16},
	0x30: {"JR NC, e8", 1, jr_nc_e8},
	0xC9: {"RET", 0, ret},
	0xFF: {"RST 38h", 0, rst_38h},

	// stack manipulation instructions
	0x33: {"INC SP", 0, inc_sp},
	0x3C: {"INC A", 0, inc_a},
	0xF1: {"POP AF", 0, pop_af},
	0xE1: {"POP HL", 0, pop_hl},
	0xF5: {"PUSH AF", 0, push_af},
	0xC5: {"PUSH BC", 0, push_bc},
	0xE5: {"PUSH HL", 0, push_hl},

	// interrupt related instructions
	0xF3: {"DI", 0, di},
}
