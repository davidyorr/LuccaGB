package cpu

import "fmt"

type instruction struct {
	mnemonic string
	step     func(cpu *CPU) bool
}

// how to read opcode table
// ---------
// LD A, d8
//   2 8
// ---------
// 2 is length, 8 is duration
// 2 - 1 byte opcode 0x3E, 1 byte operand d8
// 8 - duration in T cycles
// ---------
// order of the registers: B, C, D, E, H, L, (HL), A

var instructions = [256]instruction{
	// load instructions
	0x06: {"LD B, n8", ld_b_n8},
	0x0E: {"LD C, n8", ld_c_n8},
	0x16: {"LD D, n8", ld_d_n8},
	0x1E: {"LD E, n8", ld_e_n8},
	0x26: {"LD H, n8", ld_h_n8},
	0x2E: {"LD L, n8", ld_l_n8},
	0x3E: {"LD A, n8", ld_a_n8},

	0x40: {"LD B, B", ld_b_b},
	0x41: {"LD B, C", ld_b_c},
	0x42: {"LD B, D", ld_b_d},
	0x43: {"LD B, E", ld_b_e},
	0x44: {"LD B, H", ld_b_h},
	0x45: {"LD B, L", ld_b_l},
	0x47: {"LD B, A", ld_b_a},

	0x48: {"LD C, B", ld_c_b},
	0x49: {"LD C, C", ld_c_c},
	0x4A: {"LD C, D", ld_c_d},
	0x4B: {"LD C, E", ld_c_e},
	0x4C: {"LD C, H", ld_c_h},
	0x4D: {"LD C, L", ld_c_l},
	0x4F: {"LD C, A", ld_c_a},

	0x50: {"LD D, B", ld_d_b},
	0x51: {"LD D, C", ld_d_c},
	0x52: {"LD D, D", ld_d_d},
	0x53: {"LD D, E", ld_d_e},
	0x54: {"LD D, H", ld_d_h},
	0x55: {"LD D, L", ld_d_l},
	0x57: {"LD D, A", ld_d_a},

	0x58: {"LD E, B", ld_e_b},
	0x59: {"LD E, C", ld_e_c},
	0x5A: {"LD E, D", ld_e_d},
	0x5B: {"LD E, E", ld_e_e},
	0x5C: {"LD E, H", ld_e_h},
	0x5D: {"LD E, L", ld_e_l},
	0x5F: {"LD E, A", ld_e_a},

	0x60: {"LD H, B", ld_h_b},
	0x61: {"LD H, C", ld_h_c},
	0x62: {"LD H, D", ld_h_d},
	0x63: {"LD H, E", ld_h_e},
	0x64: {"LD H, H", ld_h_h},
	0x65: {"LD H, L", ld_h_l},
	0x67: {"LD H, A", ld_h_a},

	0x68: {"LD L, B", ld_l_b},
	0x69: {"LD L, C", ld_l_c},
	0x6A: {"LD L, D", ld_l_d},
	0x6B: {"LD L, E", ld_l_e},
	0x6C: {"LD L, H", ld_l_h},
	0x6D: {"LD L, L", ld_l_l},
	0x6F: {"LD L, A", ld_l_a},

	0x70: {"LD [HL], B", ld_hl_b},
	0x71: {"LD [HL], C", ld_hl_c},
	0x72: {"LD [HL], D", ld_hl_d},
	0x73: {"LD [HL], E", ld_hl_e},
	0x74: {"LD [HL], H", ld_hl_h},
	0x75: {"LD [HL], L", ld_hl_l},
	0x77: {"LD [HL], A", ld_hl_a},
	0x36: {"LD [HL], n8", ld_hl_n8},

	0x78: {"LD A, B", ld_a_b},
	0x79: {"LD A, C", ld_a_c},
	0x7A: {"LD A, D", ld_a_d},
	0x7B: {"LD A, E", ld_a_e},
	0x7C: {"LD A, H", ld_a_h},
	0x7D: {"LD A, L", ld_a_l},
	0x7F: {"LD A, A", ld_a_a},

	0x46: {"LD B, [HL]", ld_b_hl},
	0x4E: {"LD C, [HL]", ld_c_hl},
	0x56: {"LD D, [HL]", ld_d_hl},
	0x5E: {"LD E, [HL]", ld_e_hl},
	0x66: {"LD H, [HL]", ld_h_hl},
	0x6E: {"LD L, [HL]", ld_l_hl},
	0x7E: {"LD A, [HL]", ld_a_hl},

	0x01: {"LD BC, n16", ld_bc_n16},
	0x11: {"LD DE, n16", ld_de_n16},
	0x21: {"LD HL, n16", ld_hl_n16},
	0x31: {"LD SP, n16", ld_sp_n16},

	0x02: {"LD [BC], A", ld_bc_a},
	0x12: {"LD [DE], A", ld_de_a},
	0xEA: {"LD [a16], A", ld_a16_a},

	0xE0: {"LDH [a8], A", ldh_at_a8_a},
	0xE2: {"LDH [C], A", ldh_at_c_a},
	0xF0: {"LDH A, [a8]", ldh_a_a8},

	0x0A: {"LD A, [BC]", ld_a_bc},
	0x1A: {"LD A, [DE]", ld_a_de},
	0xFA: {"LD A, [a16]", ld_a_a16},

	0xF2: {"LDH A, [C]", ldh_a_at_c},

	0x22: {"LD [HL+], A", ld_hli_a},
	0x32: {"LD [HL-], A", ld_hld_a},
	0x2A: {"LD A, [HL+]", ld_a_hli},
	0x3A: {"LD A, [HL-]", ld_a_hld},

	// 8-bit arithmetic instructions
	0x88: {"ADC A, B", adc_a_b},
	0x89: {"ADC A, C", adc_a_c},
	0x8A: {"ADC A, D", adc_a_d},
	0x8B: {"ADC A, E", adc_a_e},
	0x8C: {"ADC A, H", adc_a_h},
	0x8D: {"ADC A, L", adc_a_l},
	0x8F: {"ADC A, A", adc_a_a},
	0x8E: {"ADC A, [HL]", adc_a_hl},
	0xCE: {"ADC A, n8", adc_a_n8},

	0x80: {"ADD A, B", add_a_b},
	0x81: {"ADD A, C", add_a_c},
	0x82: {"ADD A, D", add_a_d},
	0x83: {"ADD A, E", add_a_e},
	0x84: {"ADD A, H", add_a_h},
	0x85: {"ADD A, L", add_a_l},
	0x87: {"ADD A, A", add_a_a},
	0x86: {"ADD A, [HL]", add_a_at_hl},
	0xC6: {"ADD A, n8", add_a_n8},

	0xB8: {"CP A, B", cp_a_b},
	0xB9: {"CP A, C", cp_a_c},
	0xBA: {"CP A, D", cp_a_d},
	0xBB: {"CP A, E", cp_a_e},
	0xBC: {"CP A, H", cp_a_h},
	0xBD: {"CP A, L", cp_a_l},
	0xBF: {"CP A, A", cp_a_a},
	0xBE: {"CP A, [HL]", cp_a_hl},
	0xFE: {"CP A, n8", cp_a_n8},

	0x05: {"DEC B", dec_b},
	0x0D: {"DEC C", dec_c},
	0x15: {"DEC D", dec_d},
	0x1D: {"DEC E", dec_e},
	0x25: {"DEC H", dec_h},
	0x2D: {"DEC L", dec_l},
	0x3D: {"DEC A", dec_a},
	0x35: {"DEC [HL]", dec_at_hl},

	0x04: {"INC B", inc_b},
	0x0C: {"INC C", inc_c},
	0x14: {"INC D", inc_d},
	0x1C: {"INC E", inc_e},
	0x24: {"INC H", inc_h},
	0x2C: {"INC L", inc_l},
	0x3C: {"INC A", inc_a},
	0x34: {"INC [HL]", inc_at_hl},

	0x98: {"SBC A, B", sbc_a_b},
	0x99: {"SBC A, C", sbc_a_c},
	0x9A: {"SBC A, D", sbc_a_d},
	0x9B: {"SBC A, E", sbc_a_e},
	0x9C: {"SBC A, H", sbc_a_h},
	0x9D: {"SBC A, L", sbc_a_l},
	0x9F: {"SBC A, A", sbc_a_a},
	0x9E: {"SBC A, [HL]", sbc_a_at_hl},
	0xDE: {"SBC A, n8", sbc_a_n8},

	0x90: {"SUB A, B", sub_a_b},
	0x91: {"SUB A, C", sub_a_c},
	0x92: {"SUB A, D", sub_a_d},
	0x93: {"SUB A, E", sub_a_e},
	0x94: {"SUB A, H", sub_a_h},
	0x95: {"SUB A, L", sub_a_l},
	0x97: {"SUB A, A", sub_a_a},
	0x96: {"SUB A, [HL]", sub_a_hl},
	0xD6: {"SUB A, n8", sub_a_n8},

	// 16-bit arithmetic instructions
	0x09: {"ADD HL, BC", add_hl_bc},
	0x19: {"ADD HL, DE", add_hl_de},
	0x29: {"ADD HL, HL", add_hl_hl},
	0x0B: {"DEC BC", dec_bc},
	0x1B: {"DEC DE", dec_de},
	0x2B: {"DEC HL", dec_hl},
	0x03: {"INC BC", inc_bc},
	0x13: {"INC DE", inc_de},
	0x23: {"INC HL", inc_hl},

	// bitwise logic instructions
	0xA0: {"AND A, B", and_a_b},
	0xA1: {"AND A, C", and_a_c},
	0xA2: {"AND A, D", and_a_d},
	0xA3: {"AND A, E", and_a_e},
	0xA4: {"AND A, H", and_a_h},
	0xA5: {"AND A, L", and_a_l},
	0xA7: {"AND A, A", and_a_a},
	0xA6: {"AND A, [HL]", and_a_at_hl},
	0xE6: {"AND A, n8", and_a_n8},

	0x2F: {"CPL", cpl},

	0xA8: {"XOR A, B", xor_a_b},
	0xA9: {"XOR A, C", xor_a_c},
	0xAA: {"XOR A, D", xor_a_d},
	0xAB: {"XOR A, E", xor_a_e},
	0xAC: {"XOR A, H", xor_a_h},
	0xAD: {"XOR A, L", xor_a_l},
	0xAF: {"XOR A, A", xor_a_a},
	0xAE: {"XOR A, [HL]", xor_a_hl},
	0xEE: {"XOR A, n8", xor_a_n8},

	0xB0: {"OR A, B", or_a_b},
	0xB1: {"OR A, C", or_a_c},
	0xB2: {"OR A, D", or_a_d},
	0xB3: {"OR A, E", or_a_e},
	0xB4: {"OR A, H", or_a_h},
	0xB5: {"OR A, L", or_a_l},
	0xB7: {"OR A, A", or_a_a},
	0xB6: {"OR A, [HL]", or_a_hl},
	0xF6: {"OR A, n8", or_a_n8},

	// bit shift instructions
	0x17: {"RLA", rla},
	0x07: {"RLCA", rlca},
	0x1F: {"RRA", rra},
	0x0F: {"RRCA", rrca},

	// jumps and subroutine instructions
	0xCD: {"CALL a16", call_a16},
	0xCC: {"CALL Z, a16", call_z_a16},
	0xDC: {"CALL C, a16", call_c_a16},
	0xC4: {"CALL NZ, a16", call_nz_a16},
	0xD4: {"CALL NC, a16", call_nc_a16},
	0xE9: {"JP HL", jp_hl},
	0xC3: {"JP a16", jp_a16},
	0xCA: {"JP Z, a16", jp_z_a16},
	0xDA: {"JP C, a16", jp_c_a16},
	0xC2: {"JP NZ, a16", jp_nz_a16},
	0xD2: {"JP NC, a16", jp_nc_a16},
	0x18: {"JR e8", jr_e8},
	0x20: {"JR NZ, e8", jr_nz_e8},
	0x28: {"JR Z, e8", jr_z_e8},
	0x30: {"JR NC, e8", jr_nc_e8},
	0x38: {"JR C, e8", jr_c_e8},
	0xC8: {"RET Z", ret_z},
	0xD8: {"RET C", ret_c},
	0xC0: {"RET NZ", ret_nz},
	0xD0: {"RET NC", ret_nc},
	0xC9: {"RET", ret},
	0xD9: {"RETI", reti},

	0xC7: {"RST $00", rst_00h},
	0xCF: {"RST $08", rst_08h},
	0xD7: {"RST $10", rst_10h},
	0xDF: {"RST $18", rst_18h},
	0xE7: {"RST $20", rst_20h},
	0xEF: {"RST $28", rst_28h},
	0xF7: {"RST $30", rst_30h},
	0xFF: {"RST $38", rst_38h},

	// carry flag instructions
	0x3F: {"CCF", ccf},
	0x37: {"SCF", scf},

	// stack manipulation instructions
	0x39: {"ADD HL, SP", add_hl_sp},
	0xE8: {"ADD SP, e8", add_sp_e8},
	0x3B: {"DEC SP", dec_sp},
	0x33: {"INC SP", inc_sp},
	0xF9: {"LD SP, HL", ld_sp_hl},
	0x08: {"LD [a16], SP", ld_a16_sp},
	0xF8: {"LD HL, SP + e8", ld_hl_sp_e8},
	0xC1: {"POP BC", pop_bc},
	0xD1: {"POP DE", pop_de},
	0xE1: {"POP HL", pop_hl},
	0xF1: {"POP AF", pop_af},
	0xC5: {"PUSH BC", push_bc},
	0xD5: {"PUSH DE", push_de},
	0xE5: {"PUSH HL", push_hl},
	0xF5: {"PUSH AF", push_af},

	// interrupt related instructions
	0xF3: {"DI", di},
	0xFB: {"EI", ei},
	0x76: {"HALT", halt},

	// misc instructions
	0x00: {"NOP", nop},
	0x27: {"DAA", daa},
	0x10: {"STOP", stop},

	// undefined opcodes
	0xD3: {"NOP", nop},
	0xDB: {"NOP", nop},
	0xDD: {"NOP", nop},
	0xE3: {"NOP", nop},
	0xE4: {"NOP", nop},
	0xEB: {"NOP", nop},
	0xEC: {"NOP", nop},
	0xED: {"NOP", nop},
	0xF4: {"NOP", nop},
	0xFC: {"NOP", nop},
	0xFD: {"NOP", nop},

	// the mnemonic is generated in the handler function
	0xCB: {"", executeCbInstructionStep},
}

// Helper for CB-prefixed instruction mnemonics
var cbRegisters = []string{"B", "C", "D", "E", "H", "L", "[HL]", "A"}

// Corresponds to u3 when operation is 0b00
var cbShiftRotates = []string{"RLC", "RRC", "RL", "RR", "SLA", "SRA", "SWAP", "SRL"}

func executeCbInstructionStep(cpu *CPU) bool {
	// the CB opcode does not get fetched until M-cycle 2
	if cpu.mCycle < 2 {
		return false
	}

	operation := *cpu.cbOpcode >> 6
	u3 := (*cpu.cbOpcode & 0b0011_1000) >> 3
	r8 := (*cpu.cbOpcode & 0b0000_0111)

	switch operation {
	case 0b00:
		cpu.instruction.mnemonic = fmt.Sprintf("%s %s", cbShiftRotates[u3], cbRegisters[r8])
		return cpu.shift_rotate_u3_r8(u3, r8)
	case 0b01:
		cpu.instruction.mnemonic = fmt.Sprintf("BIT %d, %s", u3, cbRegisters[r8])
		return cpu.bit_u3_r8(u3, r8)
	case 0b10:
		cpu.instruction.mnemonic = fmt.Sprintf("RES %d, %s", u3, cbRegisters[r8])
		return cpu.res_u3_r8(u3, r8)
	case 0b11:
		cpu.instruction.mnemonic = fmt.Sprintf("SET %d, %s", u3, cbRegisters[r8])
		return cpu.set_u3_r8(u3, r8)
	}

	return false
}
