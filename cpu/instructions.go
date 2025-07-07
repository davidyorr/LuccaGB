package cpu

type instruction struct {
	mnemonic      string
	operandLength uint8
	step          func(cpu *CPU) bool
}

// how to read opcode table
// ---------
// LD A, d8
//   2 8
// ---------
// 2 is length, 8 is duration
// 2 - 1 byte opcode 0x3E, 1 byte operand d8. set the operandLength to this value minus 1 (to exclude the opcode byte)
// 8 - duration in T cycles, what execute() should return
// ---------
// order of the registers: B, C, D, E, H, L, (HL), A

var instructions = [256]instruction{
	// load instructions
	0x06: {"LD B, n8", 1, ld_b_n8},
	0x0E: {"LD C, n8", 1, ld_c_n8},
	0x16: {"LD D, n8", 1, ld_d_n8},
	0x1E: {"LD E, n8", 1, ld_e_n8},
	0x26: {"LD H, n8", 1, ld_h_n8},
	0x2E: {"LD L, n8", 1, ld_l_n8},
	0x3E: {"LD A, n8", 1, ld_a_n8},

	0x40: {"LD B, B", 0, ld_b_b},
	0x41: {"LD B, C", 0, ld_b_c},
	0x42: {"LD B, D", 0, ld_b_d},
	0x43: {"LD B, E", 0, ld_b_e},
	0x44: {"LD B, H", 0, ld_b_h},
	0x45: {"LD B, L", 0, ld_b_l},
	0x47: {"LD B, A", 0, ld_b_a},

	0x48: {"LD C, B", 0, ld_c_b},
	0x49: {"LD C, C", 0, ld_c_c},
	0x4A: {"LD C, D", 0, ld_c_d},
	0x4B: {"LD C, E", 0, ld_c_e},
	0x4C: {"LD C, H", 0, ld_c_h},
	0x4D: {"LD C, L", 0, ld_c_l},
	0x4F: {"LD C, A", 0, ld_c_a},

	0x50: {"LD D, B", 0, ld_d_b},
	0x51: {"LD D, C", 0, ld_d_c},
	0x52: {"LD D, D", 0, ld_d_d},
	0x53: {"LD D, E", 0, ld_d_e},
	0x54: {"LD D, H", 0, ld_d_h},
	0x55: {"LD D, L", 0, ld_d_l},
	0x57: {"LD D, A", 0, ld_d_a},

	0x58: {"LD E, B", 0, ld_e_b},
	0x59: {"LD E, C", 0, ld_e_c},
	0x5A: {"LD E, D", 0, ld_e_d},
	0x5B: {"LD E, E", 0, ld_e_e},
	0x5C: {"LD E, H", 0, ld_e_h},
	0x5D: {"LD E, L", 0, ld_e_l},
	0x5F: {"LD E, A", 0, ld_e_a},

	0x60: {"LD H, B", 0, ld_h_b},
	0x61: {"LD H, C", 0, ld_h_c},
	0x62: {"LD H, D", 0, ld_h_d},
	0x63: {"LD H, E", 0, ld_h_e},
	0x64: {"LD H, H", 0, ld_h_h},
	0x65: {"LD H, L", 0, ld_h_l},
	0x67: {"LD H, A", 0, ld_h_a},

	0x68: {"LD L, B", 0, ld_l_b},
	0x69: {"LD L, C", 0, ld_l_c},
	0x6A: {"LD L, D", 0, ld_l_d},
	0x6B: {"LD L, E", 0, ld_l_e},
	0x6C: {"LD L, H", 0, ld_l_h},
	0x6D: {"LD L, L", 0, ld_l_l},
	0x6F: {"LD L, A", 0, ld_l_a},

	0x70: {"LD [HL], B", 0, ld_hl_b},
	0x71: {"LD [HL], C", 0, ld_hl_c},
	0x72: {"LD [HL], D", 0, ld_hl_d},
	0x73: {"LD [HL], E", 0, ld_hl_e},
	0x74: {"LD [HL], H", 0, ld_hl_h},
	0x75: {"LD [HL], L", 0, ld_hl_l},
	0x77: {"LD [HL], A", 0, ld_hl_a},
	0x36: {"LD [HL], n8", 1, ld_hl_n8},

	0x78: {"LD A, B", 0, ld_a_b},
	0x79: {"LD A, C", 0, ld_a_c},
	0x7A: {"LD A, D", 0, ld_a_d},
	0x7B: {"LD A, E", 0, ld_a_e},
	0x7C: {"LD A, H", 0, ld_a_h},
	0x7D: {"LD A, L", 0, ld_a_l},
	0x7F: {"LD A, A", 0, ld_a_a},

	0x46: {"LD B, [HL]", 0, ld_b_hl},
	0x4E: {"LD C, [HL]", 0, ld_c_hl},
	0x56: {"LD D, [HL]", 0, ld_d_hl},
	0x5E: {"LD E, [HL]", 0, ld_e_hl},
	0x66: {"LD H, [HL]", 0, ld_h_hl},
	0x6E: {"LD L, [HL]", 0, ld_l_hl},
	0x7E: {"LD A, [HL]", 0, ld_a_hl},

	0x01: {"LD BC, n16", 2, ld_bc_n16},
	0x11: {"LD DE, n16", 2, ld_de_n16},
	0x21: {"LD HL, n16", 2, ld_hl_n16},
	0x31: {"LD SP, n16", 2, ld_sp_n16},

	0x02: {"LD [BC], A", 0, ld_bc_a},
	0x12: {"LD [DE], A", 0, ld_de_a},
	0xEA: {"LD [a16], A", 2, ld_a16_a},

	0xE0: {"LDH [a8], A", 1, ldh_at_a8_a},
	0xE2: {"LDH [C], A", 0, ldh_at_c_a},
	0xF0: {"LDH A, [a8]", 1, ldh_a_a8},

	0x0A: {"LD A, [BC]", 0, ld_a_bc},
	0x1A: {"LD A, [DE]", 0, ld_a_de},
	0xFA: {"LD A, [a16]", 2, ld_a_a16},

	0xF2: {"LDH A, [C]", 0, ldh_a_at_c},

	0x22: {"LD [HL+], A", 0, ld_hli_a},
	0x32: {"LD [HL-], A", 0, ld_hld_a},
	0x2A: {"LD A, [HL+]", 0, ld_a_hli},
	0x3A: {"LD A, [HL-]", 0, ld_a_hld},

	// 8-bit arithmetic instructions
	0x88: {"ADC A, B", 0, adc_a_b},
	0x89: {"ADC A, C", 0, adc_a_c},
	0x8A: {"ADC A, D", 0, adc_a_d},
	0x8B: {"ADC A, E", 0, adc_a_e},
	0x8C: {"ADC A, H", 0, adc_a_h},
	0x8D: {"ADC A, L", 0, adc_a_l},
	0x8F: {"ADC A, A", 0, adc_a_a},
	0x8E: {"ADC A, [HL]", 0, adc_a_hl},
	0xCE: {"ADC A, n8", 1, adc_a_n8},

	0x80: {"ADD A, B", 0, add_a_b},
	0x81: {"ADD A, C", 0, add_a_c},
	0x82: {"ADD A, D", 0, add_a_d},
	0x83: {"ADD A, E", 0, add_a_e},
	0x84: {"ADD A, H", 0, add_a_h},
	0x85: {"ADD A, L", 0, add_a_l},
	0x87: {"ADD A, A", 0, add_a_a},
	0x86: {"ADD A, [HL]", 0, add_a_at_hl},
	0xC6: {"ADD A, n8", 1, add_a_n8},

	0xB8: {"CP A, B", 0, cp_a_b},
	0xB9: {"CP A, C", 0, cp_a_c},
	0xBA: {"CP A, D", 0, cp_a_d},
	0xBB: {"CP A, E", 0, cp_a_e},
	0xBC: {"CP A, H", 0, cp_a_h},
	0xBD: {"CP A, L", 0, cp_a_l},
	0xBF: {"CP A, A", 0, cp_a_a},
	0xBE: {"CP A, [HL]", 0, cp_a_hl},
	0xFE: {"CP A, n8", 1, cp_a_n8},

	0x05: {"DEC B", 0, dec_b},
	0x0D: {"DEC C", 0, dec_c},
	0x15: {"DEC D", 0, dec_d},
	0x1D: {"DEC E", 0, dec_e},
	0x25: {"DEC H", 0, dec_h},
	0x2D: {"DEC L", 0, dec_l},
	0x3D: {"DEC A", 0, dec_a},
	0x35: {"DEC [HL]", 0, dec_at_hl},

	0x04: {"INC B", 0, inc_b},
	0x0C: {"INC C", 0, inc_c},
	0x14: {"INC D", 0, inc_d},
	0x1C: {"INC E", 0, inc_e},
	0x24: {"INC H", 0, inc_h},
	0x2C: {"INC L", 0, inc_l},
	0x3C: {"INC A", 0, inc_a},
	0x34: {"INC [HL]", 0, inc_at_hl},

	0x98: {"SBC A, B", 0, sbc_a_b},
	0x99: {"SBC A, C", 0, sbc_a_c},
	0x9A: {"SBC A, D", 0, sbc_a_d},
	0x9B: {"SBC A, E", 0, sbc_a_e},
	0x9C: {"SBC A, H", 0, sbc_a_h},
	0x9D: {"SBC A, L", 0, sbc_a_l},
	0x9F: {"SBC A, A", 0, sbc_a_a},
	0x9E: {"SBC A, [HL]", 0, sbc_a_at_hl},
	0xDE: {"SBC A, n8", 1, sbc_a_n8},

	0x90: {"SUB A, B", 0, sub_a_b},
	0x91: {"SUB A, C", 0, sub_a_c},
	0x92: {"SUB A, D", 0, sub_a_d},
	0x93: {"SUB A, E", 0, sub_a_e},
	0x94: {"SUB A, H", 0, sub_a_h},
	0x95: {"SUB A, L", 0, sub_a_l},
	0x97: {"SUB A, A", 0, sub_a_a},
	0x96: {"SUB A, [HL]", 0, sub_a_hl},
	0xD6: {"SUB A, n8", 1, sub_a_n8},

	// 16-bit arithmetic instructions
	0x09: {"ADD HL, BC", 0, add_hl_bc},
	0x19: {"ADD HL, DE", 0, add_hl_de},
	0x29: {"ADD HL, HL", 0, add_hl_hl},
	0x0B: {"DEC BC", 0, dec_bc},
	0x1B: {"DEC DE", 0, dec_de},
	0x2B: {"DEC HL", 0, dec_hl},
	0x03: {"INC BC", 0, inc_bc},
	0x13: {"INC DE", 0, inc_de},
	0x23: {"INC HL", 0, inc_hl},

	// bitwise logic instructions
	0xA0: {"AND A, B", 0, and_a_b},
	0xA1: {"AND A, C", 0, and_a_c},
	0xA2: {"AND A, D", 0, and_a_d},
	0xA3: {"AND A, E", 0, and_a_e},
	0xA4: {"AND A, H", 0, and_a_h},
	0xA5: {"AND A, L", 0, and_a_l},
	0xA7: {"AND A, A", 0, and_a_a},
	0xA6: {"AND A, [HL]", 0, and_a_at_hl},
	0xE6: {"AND A, n8", 1, and_a_n8},

	0x2F: {"CPL", 0, cpl},

	0xA8: {"XOR A, B", 0, xor_a_b},
	0xA9: {"XOR A, C", 0, xor_a_c},
	0xAA: {"XOR A, D", 0, xor_a_d},
	0xAB: {"XOR A, E", 0, xor_a_e},
	0xAC: {"XOR A, H", 0, xor_a_h},
	0xAD: {"XOR A, L", 0, xor_a_l},
	0xAF: {"XOR A, A", 0, xor_a_a},
	0xAE: {"XOR A, [HL]", 0, xor_a_hl},
	0xEE: {"XOR A, n8", 1, xor_a_n8},

	0xB0: {"OR A, B", 0, or_a_b},
	0xB1: {"OR A, C", 0, or_a_c},
	0xB2: {"OR A, D", 0, or_a_d},
	0xB3: {"OR A, E", 0, or_a_e},
	0xB4: {"OR A, H", 0, or_a_h},
	0xB5: {"OR A, L", 0, or_a_l},
	0xB7: {"OR A, A", 0, or_a_a},
	0xB6: {"OR A, [HL]", 0, or_a_hl},
	0xF6: {"OR A, n8", 1, or_a_n8},

	// bit shift instructions
	0x17: {"RLA", 0, rla},
	0x07: {"RLCA", 0, rlca},
	0x1F: {"RRA", 0, rra},
	0x0F: {"RRCA", 0, rrca},

	// jumps and subroutine instructions
	0xCD: {"CALL a16", 2, call_a16},
	0xCC: {"CALL Z, a16", 2, call_z_a16},
	0xDC: {"CALL C, a16", 2, call_c_a16},
	0xC4: {"CALL NZ, a16", 2, call_nz_a16},
	0xD4: {"CALL NC, a16", 2, call_nc_a16},
	0xE9: {"JP HL", 0, jp_hl},
	0xC3: {"JP a16", 2, jp_a16},
	0xCA: {"JP Z, a16", 2, jp_z_a16},
	0xDA: {"JP C, a16", 2, jp_c_a16},
	0xC2: {"JP NZ, a16", 2, jp_nz_a16},
	0xD2: {"JP NC, a16", 2, jp_nc_a16},
	0x18: {"JR e8", 1, jr_e8},
	0x20: {"JR NZ, e8", 1, jr_nz_e8},
	0x28: {"JR Z, e8", 1, jr_z_e8},
	0x30: {"JR NC, e8", 1, jr_nc_e8},
	0x38: {"JR C, e8", 1, jr_c_e8},
	0xC8: {"RET Z", 0, ret_z},
	0xD8: {"RET C", 0, ret_c},
	0xC0: {"RET NZ", 0, ret_nz},
	0xD0: {"RET NC", 0, ret_nc},
	0xC9: {"RET", 0, ret},
	0xD9: {"RETI", 0, reti},

	0xC7: {"RST $00", 0, rst_00h},
	0xCF: {"RST $08", 0, rst_08h},
	0xD7: {"RST $10", 0, rst_10h},
	0xDF: {"RST $18", 0, rst_18h},
	0xE7: {"RST $20", 0, rst_20h},
	0xEF: {"RST $28", 0, rst_28h},
	0xF7: {"RST $30", 0, rst_30h},
	0xFF: {"RST $38", 0, rst_38h},

	// carry flag instructions
	0x3F: {"CCF", 0, ccf},
	0x37: {"SCF", 0, scf},

	// stack manipulation instructions
	0x39: {"ADD HL, SP", 0, add_hl_sp},
	0xE8: {"ADD SP, e8", 1, add_sp_e8},
	0x3B: {"DEC SP", 0, dec_sp},
	0x33: {"INC SP", 0, inc_sp},
	0xF9: {"LD SP, HL", 0, ld_sp_hl},
	0x08: {"LD [a16], SP", 2, ld_a16_sp},
	0xF8: {"LD HL, SP + e8", 1, ld_hl_sp_e8},
	0xC1: {"POP BC", 0, pop_bc},
	0xD1: {"POP DE", 0, pop_de},
	0xE1: {"POP HL", 0, pop_hl},
	0xF1: {"POP AF", 0, pop_af},
	0xC5: {"PUSH BC", 0, push_bc},
	0xD5: {"PUSH DE", 0, push_de},
	0xE5: {"PUSH HL", 0, push_hl},
	0xF5: {"PUSH AF", 0, push_af},

	// interrupt related instructions
	0xF3: {"DI", 0, di},
	0xFB: {"EI", 0, ei},
	0x76: {"HALT", 0, halt},

	// misc instructions
	0x00: {"NOP", 0, nop},
	0x27: {"DAA", 0, daa},
	0x10: {"STOP", 1, stop},

	// undefined opcodes
	0xD3: {"NOP", 0, nop},
	0xDB: {"NOP", 0, nop},
	0xDD: {"NOP", 0, nop},
	0xE3: {"NOP", 0, nop},
	0xE4: {"NOP", 0, nop},
	0xEB: {"NOP", 0, nop},
	0xEC: {"NOP", 0, nop},
	0xED: {"NOP", 0, nop},
	0xF4: {"NOP", 0, nop},
	0xFC: {"NOP", 0, nop},
	0xFD: {"NOP", 0, nop},
}

func (cpu *CPU) executeCbInstructionStep(opcode uint8) bool {
	operation := opcode >> 6
	u3 := (opcode & 0b0011_1000) >> 3
	r8 := (opcode & 0b0000_0111)

	switch operation {
	case 0b00:
		return cpu.shift_rotate_u3_r8(u3, r8)
	case 0b01:
		return cpu.bit_u3_r8(u3, r8)
	case 0b10:
		return cpu.res_u3_r8(u3, r8)
	case 0b11:
		return cpu.set_u3_r8(u3, r8)
	}

	return false
}

func (cpu *CPU) GetNumberOfUnimplementedInstructions() int {
	var count int = 0
	for i := range instructions {
		if instructions[i].step == nil {
			count++
		}
	}

	// subtract the prefix opcode 0xCB
	return count - 1
}
