package cpu

// 0x06 Copy the value n8 into register r8
func ld_b_n8(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		cpu.b = uint8(cpu.immediateValue)
		return 4, true
	}
	return 4, true
}

// 0x0E Copy the value n8 into register r8
func ld_c_n8(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		cpu.c = uint8(cpu.immediateValue)
		return 4, true
	}
	return 4, true
}

// 0x16 Copy the value n8 into register r8
func ld_d_n8(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		cpu.d = uint8(cpu.immediateValue)
		return 4, true
	}
	return 4, true
}

// 0x1E Copy the value n8 into register r8
func ld_e_n8(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		cpu.e = uint8(cpu.immediateValue)
		return 4, true
	}
	return 4, true
}

// 0x26 Copy the value n8 into register r8
func ld_h_n8(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		cpu.h = uint8(cpu.immediateValue)
		return 4, true
	}
	return 4, true
}

// 0x2E Copy the value n8 into register r8
func ld_l_n8(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		cpu.l = uint8(cpu.immediateValue)
		return 4, true
	}
	return 4, true
}

// 0x3E Copy the value n8 into register r8
func ld_a_n8(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		cpu.a = uint8(cpu.immediateValue)
		return 4, true
	}
	return 4, true
}

// 0x40 Copy (aka Load) the value in register on the right into the register on the left
func ld_b_b(cpu *CPU) (uint8, bool) {
	cpu.b = cpu.b
	return 4, true
}

// 0x41 Copy (aka Load) the value in register on the right into the register on the left
func ld_b_c(cpu *CPU) (uint8, bool) {
	cpu.b = cpu.c
	return 4, true
}

// 0x42 Copy (aka Load) the value in register on the right into the register on the left
func ld_b_d(cpu *CPU) (uint8, bool) {
	cpu.b = cpu.d
	return 4, true
}

// 0x43 Copy (aka Load) the value in register on the right into the register on the left
func ld_b_e(cpu *CPU) (uint8, bool) {
	cpu.b = cpu.e
	return 4, true
}

// 0x44 Copy (aka Load) the value in register on the right into the register on the left
func ld_b_h(cpu *CPU) (uint8, bool) {
	cpu.b = cpu.h
	return 4, true
}

// 0x45 Copy (aka Load) the value in register on the right into the register on the left
func ld_b_l(cpu *CPU) (uint8, bool) {
	cpu.b = cpu.l
	return 4, true
}

// 0x47 Copy (aka Load) the value in register on the right into the register on the left
func ld_b_a(cpu *CPU) (uint8, bool) {
	cpu.b = cpu.a
	return 4, true
}

// 0x48 Copy (aka Load) the value in register on the right into the register on the left
func ld_c_b(cpu *CPU) (uint8, bool) {
	cpu.c = cpu.b
	return 4, true
}

// 0x49 Copy (aka Load) the value in register on the right into the register on the left
func ld_c_c(cpu *CPU) (uint8, bool) {
	cpu.c = cpu.c
	return 4, true
}

// 0x4A Copy (aka Load) the value in register on the right into the register on the left
func ld_c_d(cpu *CPU) (uint8, bool) {
	cpu.c = cpu.d
	return 4, true
}

// 0x4B Copy (aka Load) the value in register on the right into the register on the left
func ld_c_e(cpu *CPU) (uint8, bool) {
	cpu.c = cpu.e
	return 4, true
}

// 0x4C Copy (aka Load) the value in register on the right into the register on the left
func ld_c_h(cpu *CPU) (uint8, bool) {
	cpu.c = cpu.h
	return 4, true
}

// 0x4L Copy (aka Load) the value in register on the right into the register on the left
func ld_c_l(cpu *CPU) (uint8, bool) {
	cpu.c = cpu.l
	return 4, true
}

// 0x4F Copy (aka Load) the value in register on the right into the register on the left
func ld_c_a(cpu *CPU) (uint8, bool) {
	cpu.c = cpu.a
	return 4, true
}

// 0x50 Copy (aka Load) the value in register on the right into the register on the left
func ld_d_b(cpu *CPU) (uint8, bool) {
	cpu.d = cpu.b
	return 4, true
}

// 0x51 Copy (aka Load) the value in register on the right into the register on the left
func ld_d_c(cpu *CPU) (uint8, bool) {
	cpu.d = cpu.c
	return 4, true
}

// 0x52 Copy (aka Load) the value in register on the right into the register on the left
func ld_d_d(cpu *CPU) (uint8, bool) {
	cpu.d = cpu.d
	return 4, true
}

// 0x53 Copy (aka Load) the value in register on the right into the register on the left
func ld_d_e(cpu *CPU) (uint8, bool) {
	cpu.d = cpu.e
	return 4, true
}

// 0x54 Copy (aka Load) the value in register on the right into the register on the left
func ld_d_h(cpu *CPU) (uint8, bool) {
	cpu.d = cpu.h
	return 4, true
}

// 0x55 Copy (aka Load) the value in register on the right into the register on the left
func ld_d_l(cpu *CPU) (uint8, bool) {
	cpu.d = cpu.l
	return 4, true
}

// 0x57 Copy (aka Load) the value in register on the right into the register on the left
func ld_d_a(cpu *CPU) (uint8, bool) {
	cpu.d = cpu.a
	return 4, true
}

// 0x58 Copy (aka Load) the value in register on the right into the register on the left
func ld_e_b(cpu *CPU) (uint8, bool) {
	cpu.e = cpu.b
	return 4, true
}

// 0x59 Copy (aka Load) the value in register on the right into the register on the left
func ld_e_c(cpu *CPU) (uint8, bool) {
	cpu.e = cpu.c
	return 4, true
}

// 0x5A Copy (aka Load) the value in register on the right into the register on the left
func ld_e_d(cpu *CPU) (uint8, bool) {
	cpu.e = cpu.d
	return 4, true
}

// 0x5B Copy (aka Load) the value in register on the right into the register on the left
func ld_e_e(cpu *CPU) (uint8, bool) {
	cpu.e = cpu.e
	return 4, true
}

// 0x5C Copy (aka Load) the value in register on the right into the register on the left
func ld_e_h(cpu *CPU) (uint8, bool) {
	cpu.e = cpu.h
	return 4, true
}

// 0x5D Copy (aka Load) the value in register on the right into the register on the left
func ld_e_l(cpu *CPU) (uint8, bool) {
	cpu.e = cpu.l
	return 4, true
}

// 0x5F Copy (aka Load) the value in register on the right into the register on the left
func ld_e_a(cpu *CPU) (uint8, bool) {
	cpu.e = cpu.a
	return 4, true
}

// 0x60 Copy (aka Load) the value in register on the right into the register on the left
func ld_h_b(cpu *CPU) (uint8, bool) {
	cpu.h = cpu.b
	return 4, true
}

// 0x61 Copy (aka Load) the value in register on the right into the register on the left
func ld_h_c(cpu *CPU) (uint8, bool) {
	cpu.h = cpu.c
	return 4, true
}

// 0x62 Copy (aka Load) the value in register on the right into the register on the left
func ld_h_d(cpu *CPU) (uint8, bool) {
	cpu.h = cpu.d
	return 4, true
}

// 0x63 Copy (aka Load) the value in register on the right into the register on the left
func ld_h_e(cpu *CPU) (uint8, bool) {
	cpu.h = cpu.e
	return 4, true
}

// 0x64 Copy (aka Load) the value in register on the right into the register on the left
func ld_h_h(cpu *CPU) (uint8, bool) {
	cpu.h = cpu.h
	return 4, true
}

// 0x65 Copy (aka Load) the value in register on the right into the register on the left
func ld_h_l(cpu *CPU) (uint8, bool) {
	cpu.h = cpu.l
	return 4, true
}

// 0x67 Copy (aka Load) the value in register on the right into the register on the left
func ld_h_a(cpu *CPU) (uint8, bool) {
	cpu.h = cpu.a
	return 4, true
}

// 0x68 Copy (aka Load) the value in register on the right into the register on the left
func ld_l_b(cpu *CPU) (uint8, bool) {
	cpu.l = cpu.b
	return 4, true
}

// 0x69 Copy (aka Load) the value in register on the right into the register on the left
func ld_l_c(cpu *CPU) (uint8, bool) {
	cpu.l = cpu.c
	return 4, true
}

// 0x6A Copy (aka Load) the value in register on the right into the register on the left
func ld_l_d(cpu *CPU) (uint8, bool) {
	cpu.l = cpu.d
	return 4, true
}

// 0x6B Copy (aka Load) the value in register on the right into the register on the left
func ld_l_e(cpu *CPU) (uint8, bool) {
	cpu.l = cpu.e
	return 4, true
}

// 0x6C Copy (aka Load) the value in register on the right into the register on the left
func ld_l_h(cpu *CPU) (uint8, bool) {
	cpu.l = cpu.h
	return 4, true
}

// 0x6D Copy (aka Load) the value in register on the right into the register on the left
func ld_l_l(cpu *CPU) (uint8, bool) {
	cpu.l = cpu.l
	return 4, true
}

// 0x6F Copy (aka Load) the value in register on the right into the register on the left
func ld_l_a(cpu *CPU) (uint8, bool) {
	cpu.l = cpu.a
	return 4, true
}

// 0x70 Copy the value in register r8 into the byte pointed to by HL
func ld_hl_b(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		address := cpu.getHL()
		cpu.bus.Write(uint16(address), cpu.b)
		return 4, true
	default:
		return 4, true
	}
}

// 0x71 Copy the value in register r8 into the byte pointed to by HL
func ld_hl_c(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		address := cpu.getHL()
		cpu.bus.Write(address, cpu.c)
		return 4, true
	default:
		return 4, true
	}
}

// 0x72 Copy the value in register r8 into the byte pointed to by HL
func ld_hl_d(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		address := cpu.getHL()
		cpu.bus.Write(address, cpu.d)
		return 4, true
	default:
		return 4, true
	}
}

// 0x73 Copy the value in register r8 into the byte pointed to by HL
func ld_hl_e(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		address := cpu.getHL()
		cpu.bus.Write(address, cpu.e)
		return 4, true
	default:
		return 4, true
	}
}

// 0x74 Copy the value in register r8 into the byte pointed to by HL
func ld_hl_h(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		address := cpu.getHL()
		cpu.bus.Write(address, cpu.h)
		return 4, true
	default:
		return 4, true
	}
}

// 0x75 Copy the value in register r8 into the byte pointed to by HL
func ld_hl_l(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		address := cpu.getHL()
		cpu.bus.Write(address, cpu.l)
		return 4, true
	default:
		return 4, true
	}
}

// 0x77 Copy the value in register r8 into the byte pointed to by HL
func ld_hl_a(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		address := cpu.getHL()
		cpu.bus.Write(address, cpu.a)
		return 4, true
	default:
		return 4, true
	}
}

// 0x36 Copy the value n8 into the byte pointed to by HL
func ld_hl_n8(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		return 4, false
	case 3:
		address := cpu.getHL()
		cpu.bus.Write(address, uint8(cpu.immediateValue))
		return 4, true
	default:
		return 4, true
	}
}

// 0x78 Copy (aka Load) the value in register on the right into the register on the left
func ld_a_b(cpu *CPU) (uint8, bool) {
	cpu.a = cpu.b
	return 4, true
}

// 0x79 Copy (aka Load) the value in register on the right into the register on the left
func ld_a_c(cpu *CPU) (uint8, bool) {
	cpu.a = cpu.c
	return 4, true
}

// 0x7A Copy (aka Load) the value in register on the right into the register on the left
func ld_a_d(cpu *CPU) (uint8, bool) {
	cpu.a = cpu.d
	return 4, true
}

// 0x7B Copy (aka Load) the value in register on the right into the register on the left
func ld_a_e(cpu *CPU) (uint8, bool) {
	cpu.a = cpu.e
	return 4, true
}

// 0x7F Copy (aka Load) the value in register on the right into the register on the left
func ld_a_a(cpu *CPU) (uint8, bool) {
	cpu.a = cpu.a
	return 4, true
}

// 0x46 Copy the value pointed to by HL into register r8
func ld_b_hl(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.b = cpu.bus.Read(cpu.getHL())
		return 4, true
	default:
		return 4, true
	}
}

// 0x4E Copy the value pointed to by HL into register r8
func ld_c_hl(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.c = cpu.bus.Read(cpu.getHL())
		return 4, true
	default:
		return 4, true
	}
}

// 0x56 Copy the value pointed to by HL into register r8
func ld_d_hl(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.d = cpu.bus.Read(cpu.getHL())
		return 4, true
	default:
		return 4, true
	}
}

// 0x5E Copy the value pointed to by HL into register r8
func ld_e_hl(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.e = cpu.bus.Read(cpu.getHL())
		return 4, true
	default:
		return 4, true
	}
}

// 0x66 Copy the value pointed to by HL into register r8
func ld_h_hl(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.h = cpu.bus.Read(cpu.getHL())
		return 4, true
	default:
		return 4, true
	}
}

// 0x6E Copy the value pointed to by HL into register r8
func ld_l_hl(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.l = cpu.bus.Read(cpu.getHL())
		return 4, true
	default:
		return 4, true
	}
}

// 0x01 Copy the value n16 into register r16
func ld_bc_n16(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		return 4, false
	case 3:
		cpu.fetchImmHighByte()
		cpu.setBC(cpu.immediateValue)
		return 4, true
	default:
		return 4, true
	}
}

// 0x11 Copy the value n16 into register r16
func ld_de_n16(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		return 4, false
	case 3:
		cpu.fetchImmHighByte()
		cpu.setDE(cpu.immediateValue)
		return 4, true
	default:
		return 4, true
	}
}

// 0x21 Copy the value n16 into register r16
func ld_hl_n16(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		return 4, false
	case 3:
		cpu.fetchImmHighByte()
		cpu.setHL(cpu.immediateValue)
		return 4, true
	default:
		return 4, true
	}
}

// 0x31 Copy the value n16 into register r16
func ld_sp_n16(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		return 4, false
	case 3:
		cpu.fetchImmHighByte()
		cpu.sp = cpu.immediateValue
		return 4, true
	default:
		return 4, true
	}
}

// 0x02 Copy the value in register A into the byte pointed to by r16
func ld_bc_a(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.bus.Write(cpu.getBC(), cpu.a)
		return 4, true
	default:
		return 4, true
	}
}

// 0x12 Copy the value in register A into the byte pointed to by r16
func ld_de_a(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.bus.Write(cpu.getDE(), cpu.a)
		return 4, true
	default:
		return 4, true
	}
}

// 0x7C Copy (aka Load) the value in register on the right into the register on the left
func ld_a_h(cpu *CPU) (uint8, bool) {
	cpu.a = cpu.h
	return 4, true
}

// 0x7D Copy (aka Load) the value in register on the right into the register on the left
func ld_a_l(cpu *CPU) (uint8, bool) {
	cpu.a = cpu.l
	return 4, true
}

// 0xE0 Copy the value in register A into the byte at address n16, provided the address is between $FF00 and $FFFF
func ldh_at_a8_a(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		return 4, false
	case 3:
		cpu.bus.Write(0xFF00+cpu.immediateValue, cpu.a)
		return 4, true
	default:
		return 4, true
	}
}

// 0xE2 Copy the value in register A into the byte at address $FF00+C
func ldh_at_c_a(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		cpu.bus.Write(0xFF00+uint16(cpu.c), cpu.a)
		return 4, true
	default:
		return 4, true
	}
}

// 0xF0 Copy the byte at address n16 into register A, provided the address is between $FF00 and $FFFF
func ldh_a_a8(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		return 4, false
	case 3:
		cpu.a = cpu.bus.Read(0xFF00 + cpu.immediateValue)
		return 4, true
	}
	return 4, true
}

// 0xEA Copy the value in register A into the byte at address n16
func ld_a16_a(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		return 4, false
	case 3:
		cpu.fetchImmHighByte()
		return 4, false
	case 4:
		cpu.bus.Write(cpu.immediateValue, cpu.a)
		return 4, true
	default:
		return 4, true
	}
}

// 0x0A Copy the byte pointed to by r16 into register A
func ld_a_bc(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.a = cpu.bus.Read(cpu.getBC())
		return 4, true
	default:
		return 4, true
	}
}

// 0x1A Copy the byte pointed to by r16 into register A
func ld_a_de(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.a = cpu.bus.Read(cpu.getDE())
		return 4, true
	default:
		return 4, true
	}
}

// 0x7E Copy the byte pointed to by r16 into register A
func ld_a_hl(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.a = cpu.bus.Read(cpu.getHL())
		return 4, true
	default:
		return 4, true
	}
}

// 0xFA Copy the byte at address n16 into register A
func ld_a_a16(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		return 4, false
	case 3:
		cpu.fetchImmHighByte()
		return 4, false
	case 4:
		cpu.a = cpu.bus.Read(cpu.immediateValue)
		return 4, true
	default:
		return 4, true
	}
}

// 0xF2 Copy the byte at address $FF00+C into register A
func ldh_a_at_c(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.a = cpu.bus.Read(0xFF00 + uint16(cpu.c))
		return 4, true
	default:
		return 4, true
	}
}

// 0x22 Copy the value in register A into the byte pointed by HL and increment HL afterwards
func ld_hli_a(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		address := cpu.getHL()
		cpu.bus.Write(address, cpu.a)
		cpu.setHL(address + 1)
		return 4, true
	default:
		return 4, true
	}
}

// 0x32 Copy the value in register A into the byte pointed by HL and decrement HL afterwards
func ld_hld_a(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		address := cpu.getHL()
		cpu.bus.Write(address, cpu.a)
		cpu.setHL(address - 1)
		return 4, true
	default:
		return 4, true
	}
}

// 0x2A Copy the byte pointed to by HL into register A, and increment HL afterwards
func ld_a_hli(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		value := cpu.getHL()
		cpu.a = cpu.bus.Read(value)
		cpu.setHL(value + 1)
		return 4, true
	default:
		return 4, true
	}
}

// 0x3A Copy the byte pointed to by HL into register A, and decrement HL afterwards
func ld_a_hld(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		value := cpu.getHL()
		cpu.a = cpu.bus.Read(value)
		cpu.setHL(value - 1)
		return 4, true
	default:
		return 4, true
	}
}

// Add the value in r8 plus the carry flag to A
func (cpu *CPU) adc_a_r8(r8 uint8) {
	originalA := cpu.a
	var carryFlag uint8 = 0
	if cpu.getFlag(FlagC) {
		carryFlag = 1
	}
	sum := uint16(originalA) + uint16(r8) + uint16(carryFlag)
	cpu.a = uint8(sum)

	cpu.setFlag(FlagZ, cpu.a == 0)
	cpu.setFlag(FlagN, false)
	cpu.setFlag(FlagC, sum > 0xFF)
	cpu.setFlag(FlagH, ((originalA&0x0F)+(r8&0x0F)+carryFlag) > 0x0F)
}

// 0x88 Add the value in r8 plus the carry flag to A
func adc_a_b(cpu *CPU) (uint8, bool) {
	cpu.adc_a_r8(cpu.b)
	return 4, true
}

// 0x89 Add the value in r8 plus the carry flag to A
func adc_a_c(cpu *CPU) (uint8, bool) {
	cpu.adc_a_r8(cpu.c)
	return 4, true
}

// 0x8A Add the value in r8 plus the carry flag to A
func adc_a_d(cpu *CPU) (uint8, bool) {
	cpu.adc_a_r8(cpu.d)
	return 4, true
}

// 0x8B Add the value in r8 plus the carry flag to A
func adc_a_e(cpu *CPU) (uint8, bool) {
	cpu.adc_a_r8(cpu.e)
	return 4, true
}

// 0x8C Add the value in r8 plus the carry flag to A
func adc_a_h(cpu *CPU) (uint8, bool) {
	cpu.adc_a_r8(cpu.h)
	return 4, true
}

// 0x8D Add the value in r8 plus the carry flag to A
func adc_a_l(cpu *CPU) (uint8, bool) {
	cpu.adc_a_r8(cpu.l)
	return 4, true
}

// 0x8F Add the value in r8 plus the carry flag to A
func adc_a_a(cpu *CPU) (uint8, bool) {
	cpu.adc_a_r8(cpu.a)
	return 4, true
}

// 0x8E Add the byte pointed to by HL plus the carry flag to A
func adc_a_hl(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.adc_a_r8(cpu.bus.Read(cpu.getHL()))
		return 4, true
	}
	return 4, true
}

// 0xCE Add the value n8 plus the carry flag to A
func adc_a_n8(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		cpu.adc_a_r8(uint8(cpu.immediateValue))
		return 4, true
	}
	return 4, true
}

// Add the value in r8 to A
func (cpu *CPU) add_a_r8(r8 uint8) {
	originalA := cpu.a
	sum := originalA + r8
	cpu.a = sum
	cpu.setFlag(FlagZ, sum == 0)
	cpu.setFlag(FlagN, false)
	cpu.setFlag(FlagC, sum < originalA)
	cpu.setFlag(FlagH, ((originalA&0x0F)+(r8&0x0F)) > 0x0F)
}

// 0x80 Add the value in r8 to A
func add_a_b(cpu *CPU) (uint8, bool) {
	cpu.add_a_r8(cpu.b)
	return 4, true
}

// 0x81 Add the value in r8 to A
func add_a_c(cpu *CPU) (uint8, bool) {
	cpu.add_a_r8(cpu.c)
	return 4, true
}

// 0x82 Add the value in r8 to A
func add_a_d(cpu *CPU) (uint8, bool) {
	cpu.add_a_r8(cpu.d)
	return 4, true
}

// 0x83 Add the value in r8 to A
func add_a_e(cpu *CPU) (uint8, bool) {
	cpu.add_a_r8(cpu.e)
	return 4, true
}

// 0x84 Add the value in r8 to A
func add_a_h(cpu *CPU) (uint8, bool) {
	cpu.add_a_r8(cpu.h)
	return 4, true
}

// 0x85 Add the value in r8 to A
func add_a_l(cpu *CPU) (uint8, bool) {
	cpu.add_a_r8(cpu.l)
	return 4, true
}

// 0x87 Add the value in r8 to A
func add_a_a(cpu *CPU) (uint8, bool) {
	cpu.add_a_r8(cpu.a)
	return 4, true
}

// 0x86 Add the byte pointed to by HL to A
func add_a_at_hl(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.add_a_r8(cpu.bus.Read(cpu.getHL()))
		return 4, true
	}
	return 4, true
}

// 0xC6 Add the value n8 to A
func add_a_n8(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		cpu.add_a_r8(uint8(cpu.immediateValue))
		return 4, true
	}
	return 4, true
}

// Helper function for opcodes:
//
//	CP A,r8
//	CP A,[HL]
//	CP A,n8
//	SUB A,r8
//	SUB A,[HL]
//	SUB A,n8
//
// Compare the value in A with the given value, then return the difference.
// This function does not set the value in A because the "CP" opcodes do not.
func (cpu *CPU) sub(value uint8) uint8 {
	originalA := cpu.a
	difference := originalA - value
	cpu.setFlag(FlagZ, difference == 0)
	cpu.setFlag(FlagN, true)
	cpu.setFlag(FlagC, value > originalA)
	cpu.setFlag(FlagH, (originalA&0x0F) < (value&0x0F))

	return difference
}

// 0xB8 ComPare the value in A with the value in r8
func cp_a_b(cpu *CPU) (uint8, bool) {
	cpu.sub(cpu.b)
	return 4, true
}

// 0xB9 ComPare the value in A with the value in r8
func cp_a_c(cpu *CPU) (uint8, bool) {
	cpu.sub(cpu.c)
	return 4, true
}

// 0xBA ComPare the value in A with the value in r8
func cp_a_d(cpu *CPU) (uint8, bool) {
	cpu.sub(cpu.d)
	return 4, true
}

// 0xBB ComPare the value in A with the value in r8
func cp_a_e(cpu *CPU) (uint8, bool) {
	cpu.sub(cpu.e)
	return 4, true
}

// 0xBC ComPare the value in A with the value in r8
func cp_a_h(cpu *CPU) (uint8, bool) {
	cpu.sub(cpu.h)
	return 4, true
}

// 0xBD ComPare the value in A with the value in r8
func cp_a_l(cpu *CPU) (uint8, bool) {
	cpu.sub(cpu.l)
	return 4, true
}

// 0xBF ComPare the value in A with the value in r8
func cp_a_a(cpu *CPU) (uint8, bool) {
	cpu.sub(cpu.a)
	return 4, true
}

// 0xBE ComPare the value in A with the byte pointed to by HL
func cp_a_hl(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.sub(cpu.bus.Read(cpu.getHL()))
		return 4, true
	}
	return 4, true
}

// 0xFE ComPare the value in A with the value n8
func cp_a_n8(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		cpu.sub(uint8(cpu.immediateValue))
		return 4, true
	}
	return 4, true
}

// Decrement the value in register r8 by 1
func (cpu *CPU) dec_r8(r8 uint8) uint8 {
	result := r8 - 1
	cpu.setFlag(FlagZ, result == 0)
	cpu.setFlag(FlagN, true)
	// if borrow from bit 4
	cpu.setFlag(FlagH, (r8&0x0F) == 0)
	return result
}

// 0x05 Decrement the value in register r8 by 1
func dec_b(cpu *CPU) (uint8, bool) {
	cpu.b = cpu.dec_r8(cpu.b)
	return 4, true
}

// 0x0D Decrement the value in register r8 by 1
func dec_c(cpu *CPU) (uint8, bool) {
	cpu.c = cpu.dec_r8(cpu.c)
	return 4, true
}

// 0x15 Decrement the value in register r8 by 1
func dec_d(cpu *CPU) (uint8, bool) {
	cpu.d = cpu.dec_r8(cpu.d)
	return 4, true
}

// 0x1D Decrement the value in register r8 by 1
func dec_e(cpu *CPU) (uint8, bool) {
	cpu.e = cpu.dec_r8(cpu.e)
	return 4, true
}

// 0x25 Decrement the value in register r8 by 1
func dec_h(cpu *CPU) (uint8, bool) {
	cpu.h = cpu.dec_r8(cpu.h)
	return 4, true
}

// 0x2D Decrement the value in register r8 by 1
func dec_l(cpu *CPU) (uint8, bool) {
	cpu.l = cpu.dec_r8(cpu.l)
	return 4, true
}

// 0x3D Decrement the value in register r8 by 1
func dec_a(cpu *CPU) (uint8, bool) {
	cpu.a = cpu.dec_r8(cpu.a)
	return 4, true
}

// 0x35 Decrement the byte pointed to by HL by 1
func dec_at_hl(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.mdr = cpu.bus.Read(cpu.getHL())
		return 4, false
	case 3:
		result := cpu.dec_r8(cpu.mdr)
		cpu.bus.Write(cpu.getHL(), result)
		return 4, true
	}
	return 12, true
}

// Increment the value in register r8 by 1
func (cpu *CPU) inc_r8(r8 uint8) uint8 {
	result := r8 + 1
	cpu.setFlag(FlagZ, result == 0)
	cpu.setFlag(FlagN, false)
	// if overflow from bit 3
	cpu.setFlag(FlagH, (r8&0x0F) == 0x0F)
	return result
}

// 0x04 Increment the value in register r8 by 1
func inc_b(cpu *CPU) (uint8, bool) {
	cpu.b = cpu.inc_r8(cpu.b)
	return 4, true
}

// 0x0C Increment the value in register r8 by 1
func inc_c(cpu *CPU) (uint8, bool) {
	cpu.c = cpu.inc_r8(cpu.c)
	return 4, true
}

// 0x14 Increment the value in register r8 by 1
func inc_d(cpu *CPU) (uint8, bool) {
	cpu.d = cpu.inc_r8(cpu.d)
	return 4, true
}

// 0x1C Increment the value in register r8 by 1
func inc_e(cpu *CPU) (uint8, bool) {
	cpu.e = cpu.inc_r8(cpu.e)
	return 4, true
}

// 0x24 Increment the value in register r8 by 1
func inc_h(cpu *CPU) (uint8, bool) {
	cpu.h = cpu.inc_r8(cpu.h)
	return 4, true
}

// 0x2C Increment the value in register r8 by 1
func inc_l(cpu *CPU) (uint8, bool) {
	cpu.l = cpu.inc_r8(cpu.l)
	return 4, true
}

// 0x3C Increment the value in register r8 by 1
func inc_a(cpu *CPU) (uint8, bool) {
	cpu.a = cpu.inc_r8(cpu.a)
	return 4, true
}

// 0x34 Increment the byte pointed to by HL by 1
func inc_at_hl(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.mdr = cpu.bus.Read(cpu.getHL())
		return 4, false
	case 3:
		result := cpu.mdr + 1
		cpu.bus.Write(cpu.getHL(), result)
		cpu.setFlag(FlagZ, result == 0)
		cpu.setFlag(FlagN, false)
		// if overflow from bit 3
		cpu.setFlag(FlagH, (cpu.mdr&0x0F)+1 > 0x0F)
		return 4, true
	}
	return 12, true
}

// Subtract the value in r8 and the carry flag from A
func (cpu *CPU) sbc_a_r8(r8 uint8) {
	var carryValue uint8 = 0
	if cpu.getFlag(FlagC) {
		carryValue = 1
	}
	originalA := cpu.a
	subtrahend := uint16(r8) + uint16(carryValue)
	difference := originalA - uint8(subtrahend)
	cpu.a = difference
	cpu.setFlag(FlagZ, difference == 0)
	cpu.setFlag(FlagN, true)
	cpu.setFlag(FlagH, (originalA&0x0F) < (r8&0x0F)+carryValue)
	cpu.setFlag(FlagC, uint16(originalA) < subtrahend)
}

// 0x98 Subtract the value in r8 and the carry flag from A
func sbc_a_b(cpu *CPU) (uint8, bool) {
	cpu.sbc_a_r8(cpu.b)
	return 4, true
}

// 0x99 Subtract the value in r8 and the carry flag from A
func sbc_a_c(cpu *CPU) (uint8, bool) {
	cpu.sbc_a_r8(cpu.c)
	return 4, true
}

// 0x9A Subtract the value in r8 and the carry flag from A
func sbc_a_d(cpu *CPU) (uint8, bool) {
	cpu.sbc_a_r8(cpu.d)
	return 4, true
}

// 0x9B Subtract the value in r8 and the carry flag from A
func sbc_a_e(cpu *CPU) (uint8, bool) {
	cpu.sbc_a_r8(cpu.e)
	return 4, true
}

// 0x9C Subtract the value in r8 and the carry flag from A
func sbc_a_h(cpu *CPU) (uint8, bool) {
	cpu.sbc_a_r8(cpu.h)
	return 4, true
}

// 0x9D Subtract the value in r8 and the carry flag from A
func sbc_a_l(cpu *CPU) (uint8, bool) {
	cpu.sbc_a_r8(cpu.l)
	return 4, true
}

// 0x9F Subtract the value in r8 and the carry flag from A
func sbc_a_a(cpu *CPU) (uint8, bool) {
	cpu.sbc_a_r8(cpu.a)
	return 4, true
}

// 0x9E Subtract the byte pointed to by HL and the carry flag from A
func sbc_a_at_hl(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.sbc_a_r8(cpu.bus.Read(cpu.getHL()))
		return 4, true
	}
	return 4, true
}

// 0xDE Subtract the value n8 and the carry flag from A
func sbc_a_n8(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		cpu.sbc_a_r8(uint8(cpu.immediateValue))
		return 4, true
	}
	return 4, true
}

// 0x90 Subtract the value in r8 from A
func sub_a_b(cpu *CPU) (uint8, bool) {
	cpu.a = cpu.sub(cpu.b)
	return 4, true
}

// 0x91 Subtract the value in r8 from A
func sub_a_c(cpu *CPU) (uint8, bool) {
	cpu.a = cpu.sub(cpu.c)
	return 4, true
}

// 0x92 Subtract the value in r8 from A
func sub_a_d(cpu *CPU) (uint8, bool) {
	cpu.a = cpu.sub(cpu.d)
	return 4, true
}

// 0x93 Subtract the value in r8 from A
func sub_a_e(cpu *CPU) (uint8, bool) {
	cpu.a = cpu.sub(cpu.e)
	return 4, true
}

// 0x94 Subtract the value in r8 from A
func sub_a_h(cpu *CPU) (uint8, bool) {
	cpu.a = cpu.sub(cpu.h)
	return 4, true
}

// 0x95 Subtract the value in r8 from A
func sub_a_l(cpu *CPU) (uint8, bool) {
	cpu.a = cpu.sub(cpu.l)
	return 4, true
}

// 0x97 Subtract the value in r8 from A
func sub_a_a(cpu *CPU) (uint8, bool) {
	cpu.a = cpu.sub(cpu.a)
	return 4, true
}

// 0x96 Subtract the byte pointed to by HL from A
func sub_a_hl(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.a = cpu.sub(cpu.bus.Read(cpu.getHL()))
		return 4, true
	}
	return 4, true
}

// Add the value in r16 to HL
func (cpu *CPU) add_hl_r16(r16 uint16) {
	originalHL := cpu.getHL()
	sum := cpu.getHL() + r16
	cpu.setHL(sum)
	cpu.setFlag(FlagN, false)
	cpu.setFlag(FlagH, ((originalHL&0x0FFF)+(r16&0x0FFF)) > 0x0FFF)
	cpu.setFlag(FlagC, sum < originalHL)
}

// 0x09 Add the value in r16 to HL
func add_hl_bc(cpu *CPU) (uint8, bool) {
	cpu.add_hl_r16(cpu.getBC())
	return 8, true
}

// 0x19 Add the value in r16 to HL
func add_hl_de(cpu *CPU) (uint8, bool) {
	cpu.add_hl_r16(cpu.getDE())
	return 8, true
}

// 0x29 Add the value in r16 to HL
func add_hl_hl(cpu *CPU) (uint8, bool) {
	cpu.add_hl_r16(cpu.getHL())
	return 8, true
}

// 0x0B Decrement the value in register r16 by 1
func dec_bc(cpu *CPU) (uint8, bool) {
	cpu.setBC(cpu.getBC() - 1)
	return 8, true
}

// 0x1B Decrement the value in register r16 by 1
func dec_de(cpu *CPU) (uint8, bool) {
	cpu.setDE(cpu.getDE() - 1)
	return 8, true
}

// 0x2B Decrement the value in register r16 by 1
func dec_hl(cpu *CPU) (uint8, bool) {
	cpu.setHL(cpu.getHL() - 1)
	return 8, true
}

// 0x03 Increment the value in register r8 by 1
func inc_bc(cpu *CPU) (uint8, bool) {
	cpu.setBC(cpu.getBC() + 1)
	return 8, true
}

// 0x13 Increment the value in register r8 by 1
func inc_de(cpu *CPU) (uint8, bool) {
	cpu.setDE(cpu.getDE() + 1)
	return 8, true
}

// 0x23 Increment the value in register r8 by 1
func inc_hl(cpu *CPU) (uint8, bool) {
	cpu.setHL(cpu.getHL() + 1)
	return 8, true
}

// 0xD6 Subtract the value n8 from A
func sub_a_n8(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		cpu.a = cpu.sub(uint8(cpu.immediateValue))
		return 4, true
	}
	return 4, true
}

// Set A to the bitwise AND between the value in r8 and A
func (cpu *CPU) and_a_r8(r8 uint8) {
	cpu.a &= r8
	cpu.setFlag(FlagZ, cpu.a == 0)
	cpu.setFlag(FlagN, false)
	cpu.setFlag(FlagH, true)
	cpu.setFlag(FlagC, false)
}

// 0xA0 Set A to the bitwise AND between the value in r8 and A
func and_a_b(cpu *CPU) (uint8, bool) {
	cpu.and_a_r8(cpu.b)
	return 4, true
}

// 0xA1 Set A to the bitwise AND between the value in r8 and A
func and_a_c(cpu *CPU) (uint8, bool) {
	cpu.and_a_r8(cpu.c)
	return 4, true
}

// 0xA2 Set A to the bitwise AND between the value in r8 and A
func and_a_d(cpu *CPU) (uint8, bool) {
	cpu.and_a_r8(cpu.d)
	return 4, true
}

// 0xA3 Set A to the bitwise AND between the value in r8 and A
func and_a_e(cpu *CPU) (uint8, bool) {
	cpu.and_a_r8(cpu.e)
	return 4, true
}

// 0xA4 Set A to the bitwise AND between the value in r8 and A
func and_a_h(cpu *CPU) (uint8, bool) {
	cpu.and_a_r8(cpu.h)
	return 4, true
}

// 0xA5 Set A to the bitwise AND between the value in r8 and A
func and_a_l(cpu *CPU) (uint8, bool) {
	cpu.and_a_r8(cpu.l)
	return 4, true
}

// 0xA7 Set A to the bitwise AND between the value in r8 and A
func and_a_a(cpu *CPU) (uint8, bool) {
	cpu.and_a_r8(cpu.a)
	return 4, true
}

// 0xA6 Set A to the bitwise AND between the byte pointed to by HL and A
func and_a_at_hl(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.and_a_r8(cpu.bus.Read(cpu.getHL()))
		return 4, true
	}
	return 4, true
}

// 0xE6 Set A to the bitwise AND between the value n8 and A
func and_a_n8(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		cpu.and_a_r8(uint8(cpu.immediateValue))
		return 4, true
	}
	return 4, true
}

// 0x2F ComPLement accumulator (A = ~A); also called bitwise NOT
func cpl(cpu *CPU) (uint8, bool) {
	cpu.a = ^cpu.a
	cpu.setFlag(FlagN, true)
	cpu.setFlag(FlagH, true)
	return 4, true
}

// Set A to the bitwise XOR between n8 and A
func (cpu *CPU) xor_a_n8(r8 uint8) {
	cpu.a ^= r8
	cpu.setFlag(FlagZ, cpu.a == 0)
	cpu.setFlag(FlagN, false)
	cpu.setFlag(FlagH, false)
	cpu.setFlag(FlagC, false)
}

// 0xA8 Set A to the bitwise XOR between the value in r8 and A
func xor_a_b(cpu *CPU) (uint8, bool) {
	cpu.xor_a_n8(cpu.b)
	return 4, true
}

// 0xA9 Set A to the bitwise XOR between the value in r8 and A
func xor_a_c(cpu *CPU) (uint8, bool) {
	cpu.xor_a_n8(cpu.c)
	return 4, true
}

// 0xAA Set A to the bitwise XOR between the value in r8 and A
func xor_a_d(cpu *CPU) (uint8, bool) {
	cpu.xor_a_n8(cpu.d)
	return 4, true
}

// 0xAB Set A to the bitwise XOR between the value in r8 and A
func xor_a_e(cpu *CPU) (uint8, bool) {
	cpu.xor_a_n8(cpu.e)
	return 4, true
}

// 0xAC Set A to the bitwise XOR between the value in r8 and A
func xor_a_h(cpu *CPU) (uint8, bool) {
	cpu.xor_a_n8(cpu.h)
	return 4, true
}

// 0xAD Set A to the bitwise XOR between the value in r8 and A
func xor_a_l(cpu *CPU) (uint8, bool) {
	cpu.xor_a_n8(cpu.l)
	return 4, true
}

// 0xAF Set A to the bitwise XOR between the value in r8 and A
func xor_a_a(cpu *CPU) (uint8, bool) {
	cpu.xor_a_n8(cpu.a)
	return 4, true
}

// 0xAE Set A to the bitwise XOR between the byte pointed to by HL and A
func xor_a_hl(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.xor_a_n8(cpu.bus.Read(cpu.getHL()))
		return 4, true
	}
	return 4, true
}

// 0xEE Set A to the bitwise XOR between the value n8 and A
func xor_a_n8(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		cpu.xor_a_n8(uint8(cpu.immediateValue))
		return 4, true
	}
	return 4, true
}

// Set A to the bitwise OR between the value in r8 and A
func (cpu *CPU) or_a_r8(r8 uint8) {
	cpu.a |= r8
	cpu.setFlag(FlagZ, cpu.a == 0)
	cpu.setFlag(FlagN, false)
	cpu.setFlag(FlagH, false)
	cpu.setFlag(FlagC, false)
}

// 0xB0 Set A to the bitwise OR between the value in r8 and A
func or_a_b(cpu *CPU) (uint8, bool) {
	cpu.or_a_r8(cpu.b)
	return 4, true
}

// 0xB1 Set A to the bitwise OR between the value in r8 and A
func or_a_c(cpu *CPU) (uint8, bool) {
	cpu.or_a_r8(cpu.c)
	return 4, true
}

// 0xB2 Set A to the bitwise OR between the value in r8 and A
func or_a_d(cpu *CPU) (uint8, bool) {
	cpu.or_a_r8(cpu.d)
	return 4, true
}

// 0xB3 Set A to the bitwise OR between the value in r8 and A
func or_a_e(cpu *CPU) (uint8, bool) {
	cpu.or_a_r8(cpu.e)
	return 4, true
}

// 0xB4 Set A to the bitwise OR between the value in r8 and A
func or_a_h(cpu *CPU) (uint8, bool) {
	cpu.or_a_r8(cpu.h)
	return 4, true
}

// 0xB5 Set A to the bitwise OR between the value in r8 and A
func or_a_l(cpu *CPU) (uint8, bool) {
	cpu.or_a_r8(cpu.l)
	return 4, true
}

// 0xB7 Set A to the bitwise OR between the value in r8 and A
func or_a_a(cpu *CPU) (uint8, bool) {
	cpu.or_a_r8(cpu.a)
	return 4, true
}

// 0xB6 Set A to the bitwise OR between the value in r8 and A
func or_a_hl(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.or_a_r8(cpu.bus.Read(cpu.getHL()))
		return 4, true
	}
	return 4, true
}

// 0xF6 Set A to the bitwise OR between the value in r8 and A
func or_a_n8(cpu *CPU) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		cpu.fetchImmLowByte()
		cpu.or_a_r8(uint8(cpu.immediateValue))
		return 4, true
	}
	return 4, true
}

// 0x17 Rotate register A left, through the carry flag
func rla(cpu *CPU) (uint8, bool) {
	bit7 := (cpu.a & 0b1000_0000)
	var mask uint8 = 0
	if cpu.getFlag(FlagC) {
		mask = 1
	}
	cpu.a = (cpu.a << 1) | (mask)
	cpu.setFlag(FlagZ, false)
	cpu.setFlag(FlagN, false)
	cpu.setFlag(FlagH, false)
	cpu.setFlag(FlagC, (bit7>>7) == 1)
	return 4, true
}

// 0x07 Rotate register A left
func rlca(cpu *CPU) (uint8, bool) {
	carry := (cpu.a >> 7) & 1
	cpu.a = cpu.a << 1
	cpu.a |= carry

	cpu.setFlag(FlagZ, false)
	cpu.setFlag(FlagN, false)
	cpu.setFlag(FlagH, false)
	cpu.setFlag(FlagC, carry == 1)

	return 4, true
}

// 0x1F Rotate register A right, through the carry flag
func rra(cpu *CPU) (uint8, bool) {
	bit0 := cpu.a & 1
	var mask uint8 = 0
	if cpu.getFlag(FlagC) {
		mask = 1
	}
	result := (cpu.a >> 1) | (mask << 7)
	cpu.a = result
	cpu.setFlag(FlagZ, false)
	cpu.setFlag(FlagN, false)
	cpu.setFlag(FlagH, false)
	cpu.setFlag(FlagC, bit0 == 1)
	return 4, true
}

// 0x0F Rotate register A right
func rrca(cpu *CPU) (uint8, bool) {
	bit0 := cpu.a & 1
	cpu.a = (cpu.a >> 1) | (bit0 << 7)
	cpu.setFlag(FlagZ, false)
	cpu.setFlag(FlagN, false)
	cpu.setFlag(FlagH, false)
	cpu.setFlag(FlagC, bit0 == 1)
	return 4, true
}

// 0xCD Call address n16
func call_a16(cpu *CPU) (uint8, bool) {
	cpu.pushToStack16(cpu.pc + 3)
	cpu.fetchImmLowByte()
	cpu.fetchImmHighByte()
	cpu.pc = cpu.immediateValue
	return 24, true
}

// 0xCC Call address n16 if condition cc is met
func call_z_a16(cpu *CPU) (uint8, bool) {
	if cpu.getFlag(FlagZ) {
		cpu.pushToStack16(cpu.pc + 3)
		cpu.fetchImmLowByte()
		cpu.fetchImmHighByte()
		cpu.pc = cpu.immediateValue
		return 24, true
	}

	return 12, true
}

// 0xDC Call address n16 if condition cc is met
func call_c_a16(cpu *CPU) (uint8, bool) {
	if cpu.getFlag(FlagC) {
		cpu.pushToStack16(cpu.pc + 3)
		cpu.fetchImmLowByte()
		cpu.fetchImmHighByte()
		cpu.pc = cpu.immediateValue
		return 24, true
	}

	return 12, true
}

// 0xC4 Call address n16 if condition cc is met
func call_nz_a16(cpu *CPU) (uint8, bool) {
	if cpu.getFlag(FlagZ) {
		return 12, true
	}

	cpu.pushToStack16(cpu.pc + 3)
	cpu.fetchImmLowByte()
	cpu.fetchImmHighByte()
	cpu.pc = cpu.immediateValue
	return 24, true
}

// 0xD4 Call address n16 if condition cc is met
func call_nc_a16(cpu *CPU) (uint8, bool) {
	if cpu.getFlag(FlagC) {
		return 12, true
	}

	cpu.pushToStack16(cpu.pc + 3)
	cpu.fetchImmLowByte()
	cpu.fetchImmHighByte()
	cpu.pc = cpu.immediateValue
	return 24, true
}

// 0xE9 Jump to address in HL; effectively, copy the value in register HL into PC
func jp_hl(cpu *CPU) (uint8, bool) {
	cpu.pc = cpu.getHL()
	return 4, true
}

// 0xC3 Jump to address a16; effectively, copy a16 into PC
func jp_a16(cpu *CPU) (uint8, bool) {
	cpu.fetchImmLowByte()
	cpu.fetchImmHighByte()
	cpu.pc = cpu.immediateValue
	return 16, true
}

// 0xCA Jump to address n16 if condition cc is met
func jp_z_a16(cpu *CPU) (uint8, bool) {
	if !cpu.getFlag(FlagZ) {
		return 12, true
	}

	cpu.fetchImmLowByte()
	cpu.fetchImmHighByte()
	cpu.pc = cpu.immediateValue
	return 16, true
}

// 0xDA Jump to address n16 if condition cc is met
func jp_c_a16(cpu *CPU) (uint8, bool) {
	if !cpu.getFlag(FlagC) {
		return 12, true
	}

	cpu.fetchImmLowByte()
	cpu.fetchImmHighByte()
	cpu.pc = cpu.immediateValue
	return 16, true
}

// 0xC2 Jump to address n16 if condition cc is met
func jp_nz_a16(cpu *CPU) (uint8, bool) {
	if cpu.getFlag(FlagZ) {
		return 12, true
	}

	cpu.fetchImmLowByte()
	cpu.fetchImmHighByte()
	cpu.pc = cpu.immediateValue
	return 16, true
}

// 0xD2 Jump to address n16 if condition cc is met
func jp_nc_a16(cpu *CPU) (uint8, bool) {
	if cpu.getFlag(FlagC) {
		return 12, true
	}

	cpu.fetchImmLowByte()
	cpu.fetchImmHighByte()
	cpu.pc = cpu.immediateValue
	return 16, true
}

// Relative Jump
func (cpu *CPU) jr() {
	operandLength := uint16(2)
	cpu.fetchImmLowByte()
	signedOffset := int8(cpu.immediateValue)
	destinationAddress := uint16(int(cpu.pc+operandLength) + int(signedOffset))
	cpu.pc = destinationAddress
}

// 0x18 Relative Jump to address e8
func jr_e8(cpu *CPU) (uint8, bool) {
	cpu.fetchImmLowByte()
	cpu.jr()
	return 12, true
}

// 0x20 Relative Jump to address e8 if condition z is met
func jr_nz_e8(cpu *CPU) (uint8, bool) {
	if cpu.getFlag(FlagZ) {
		return 8, true
	}

	cpu.fetchImmLowByte()
	cpu.jr()

	return 12, true
}

// 0x28 Relative Jump to address e8 if condition z is met
func jr_z_e8(cpu *CPU) (uint8, bool) {
	if !cpu.getFlag(FlagZ) {
		return 8, true
	}

	cpu.fetchImmLowByte()
	cpu.jr()

	return 12, true
}

// 0x30 Relative Jump to address e8 if condition c is met
func jr_nc_e8(cpu *CPU) (uint8, bool) {
	if cpu.getFlag(FlagC) {
		return 8, true
	}

	cpu.fetchImmLowByte()
	cpu.jr()

	return 12, true
}

// 0x38 Relative Jump to address e8 if condition c is met
func jr_c_e8(cpu *CPU) (uint8, bool) {
	if !cpu.getFlag(FlagC) {
		return 8, true
	}

	cpu.fetchImmLowByte()
	cpu.jr()

	return 12, true
}

// helper function for opcodes:
//
//	RET
//	REC cc
//	RETI
//
// Return from subroutine.
func (cpu *CPU) return_from_subroutine() {
	returnAddress := cpu.popFromStack16()
	cpu.pc = returnAddress
}

// 0xC8 Return from subroutine if condition cc is met
func ret_z(cpu *CPU) (uint8, bool) {
	if cpu.getFlag(FlagZ) {
		cpu.return_from_subroutine()
		return 20, true
	}

	return 8, true
}

// 0xD8 Return from subroutine if condition cc is met
func ret_c(cpu *CPU) (uint8, bool) {
	if cpu.getFlag(FlagC) {
		cpu.return_from_subroutine()
		return 20, true
	}

	return 8, true
}

// 0xC0 Return from subroutine if condition cc is met
func ret_nz(cpu *CPU) (uint8, bool) {
	if !cpu.getFlag(FlagZ) {
		cpu.return_from_subroutine()
		return 20, true
	}

	return 8, true
}

// 0xD0 Return from subroutine if condition cc is met
func ret_nc(cpu *CPU) (uint8, bool) {
	if !cpu.getFlag(FlagC) {
		cpu.return_from_subroutine()
		return 20, true
	}

	return 8, true
}

// 0xC9 Return from subroutine
func ret(cpu *CPU) (uint8, bool) {
	cpu.return_from_subroutine()
	return 16, true
}

// 0xD9 Return from subroutine and enable interrupts
func reti(cpu *CPU) (uint8, bool) {
	cpu.ScheduleIme()
	cpu.return_from_subroutine()
	return 16, true
}

// 0xC7 Call address vec
func rst_00h(cpu *CPU) (uint8, bool) {
	cpu.pushToStack16(cpu.pc + 1)
	cpu.pc = 0x00
	return 16, true
}

// 0xCF Call address vec
func rst_08h(cpu *CPU) (uint8, bool) {
	cpu.pushToStack16(cpu.pc + 1)
	cpu.pc = 0x08
	return 16, true
}

// 0xD7 Call address vec
func rst_10h(cpu *CPU) (uint8, bool) {
	cpu.pushToStack16(cpu.pc + 1)
	cpu.pc = 0x10
	return 16, true
}

// 0xDF Call address vec
func rst_18h(cpu *CPU) (uint8, bool) {
	cpu.pushToStack16(cpu.pc + 1)
	cpu.pc = 0x18
	return 16, true
}

// 0xE7 Call address vec
func rst_20h(cpu *CPU) (uint8, bool) {
	cpu.pushToStack16(cpu.pc + 1)
	cpu.pc = 0x20
	return 16, true
}

// 0xEF Call address vec
func rst_28h(cpu *CPU) (uint8, bool) {
	cpu.pushToStack16(cpu.pc + 1)
	cpu.pc = 0x28
	return 16, true
}

// 0xF7 Call address vec
func rst_30h(cpu *CPU) (uint8, bool) {
	cpu.pushToStack16(cpu.pc + 1)
	cpu.pc = 0x30
	return 16, true
}

// 0xFF Call address vec
func rst_38h(cpu *CPU) (uint8, bool) {
	cpu.pushToStack16(cpu.pc + 1)
	cpu.pc = 0x38
	return 16, true
}

// 0x3F Complement Carry Flag
func ccf(cpu *CPU) (uint8, bool) {
	cpu.setFlag(FlagN, false)
	cpu.setFlag(FlagH, false)
	cpu.setFlag(FlagC, !cpu.getFlag(FlagC))
	return 4, true
}

// 0x37 Set Carry Flag
func scf(cpu *CPU) (uint8, bool) {
	cpu.setFlag(FlagN, false)
	cpu.setFlag(FlagH, false)
	cpu.setFlag(FlagC, true)
	return 4, true
}

// 0x39 Add the value in SP to HL
func add_hl_sp(cpu *CPU) (uint8, bool) {
	cpu.add_hl_r16(cpu.sp)
	return 8, true
}

// helper function for opcodes:
//
//	ADD SP,e8
//	ADD SP,e8
//	LD HL,SP+e8
func (cpu *CPU) set_e8_carry_flags(original uint16, e8Unsigned uint16) {
	cpu.setFlag(FlagH, ((original&0x0F)+(e8Unsigned&0x0F)) > 0x0F)
	cpu.setFlag(FlagC, (original&0xFF)+(e8Unsigned&0xFF) > 0xFF)
}

// 0xE8 Add the signed value e8 to SP
func add_sp_e8(cpu *CPU) (uint8, bool) {
	originalSp := cpu.sp
	cpu.fetchImmLowByte()
	e8Signed := int8(cpu.immediateValue)
	result := int32(cpu.sp) + int32(e8Signed)
	cpu.sp = uint16(result)

	// flags are based on unsigned addition
	e8Unsigned := uint16(cpu.immediateValue)

	cpu.setFlag(FlagZ, false)
	cpu.setFlag(FlagN, false)
	cpu.set_e8_carry_flags(originalSp, e8Unsigned)
	return 16, true
}

// 0x3B Decrement the value in register SP by 1
func dec_sp(cpu *CPU) (uint8, bool) {
	cpu.sp--
	return 8, true
}

// 0x33 Increment the value in register SP by 1
func inc_sp(cpu *CPU) (uint8, bool) {
	cpu.sp++
	return 8, true
}

// 0xF9 Copy register HL into register SP
func ld_sp_hl(cpu *CPU) (uint8, bool) {
	cpu.sp = cpu.getHL()
	return 8, true
}

// 0x08 Copy SP & $FF at address n16 and SP >> 8 at address n16 + 1
func ld_a16_sp(cpu *CPU) (uint8, bool) {
	cpu.fetchImmLowByte()
	cpu.fetchImmHighByte()
	cpu.bus.Write(cpu.immediateValue, uint8(cpu.sp&0xFF))
	cpu.bus.Write(cpu.immediateValue+1, uint8(cpu.sp>>8))
	return 20, true
}

// 0xF8 Add the signed value e8 to SP and copy the result in HL
func ld_hl_sp_e8(cpu *CPU) (uint8, bool) {
	originalSp := cpu.sp
	cpu.fetchImmLowByte()
	e8Signed := int8(cpu.immediateValue)
	result := uint16(int32(cpu.sp) + int32(e8Signed))
	cpu.setHL(result)

	e8Unsigned := uint16(cpu.immediateValue)
	cpu.setFlag(FlagZ, false)
	cpu.setFlag(FlagN, false)
	cpu.set_e8_carry_flags(originalSp, e8Unsigned)
	return 12, true
}

// 0xC1 Pop register r16 from the stack
func pop_bc(cpu *CPU) (uint8, bool) {
	value := cpu.popFromStack16()
	cpu.setBC(value)
	return 12, true
}

// 0xD1 Pop register r16 from the stack
func pop_de(cpu *CPU) (uint8, bool) {
	value := cpu.popFromStack16()
	cpu.setDE(value)
	return 12, true
}

// 0xE1 Pop register r16 from the stack
func pop_hl(cpu *CPU) (uint8, bool) {
	value := cpu.popFromStack16()
	cpu.setHL(value)
	return 12, true
}

// 0xF1 Pop register AF from the stack
func pop_af(cpu *CPU) (uint8, bool) {
	value := cpu.popFromStack16()
	cpu.setAF(value)
	return 12, true
}

// 0xC5 Push register r16 into the stack
func push_bc(cpu *CPU) (uint8, bool) {
	cpu.pushToStack16(cpu.getBC())
	return 16, true
}

// 0xD5 Push register r16 into the stack
func push_de(cpu *CPU) (uint8, bool) {
	cpu.pushToStack16(cpu.getDE())
	return 16, true
}

// 0xE5 Push register r16 into the stack
func push_hl(cpu *CPU) (uint8, bool) {
	cpu.pushToStack16(cpu.getHL())
	return 16, true
}

// 0xF5 Push register r16 into the stack
func push_af(cpu *CPU) (uint8, bool) {
	cpu.pushToStack16(cpu.getAF())
	return 16, true
}

// 0xF3 Disable Interrupts by clearing the IME flag
func di(cpu *CPU) (uint8, bool) {
	cpu.ime = false
	return 4, true
}

// 0xFB Enable Interrupts by setting the IME flag
func ei(cpu *CPU) (uint8, bool) {
	cpu.ScheduleIme()
	return 4, true
}

// 0x76 Enter CPU low-power consumption mode until an interrupt occurs
func halt(cpu *CPU) (uint8, bool) {
	interruptEnable := cpu.bus.Read(0xFFFF)
	interruptFlag := cpu.bus.Read(0xFF0F)

	if cpu.InterruptMasterEnable() {
		cpu.Halt()
	} else {
		if (interruptEnable & interruptFlag) != 0 {
			cpu.haltBugActive = false
		} else {
			cpu.Halt()
		}
	}
	cpu.Halt()
	return 4, true
}

// 0x00 No OPeration
func nop(cpu *CPU) (uint8, bool) {
	return 4, true
}

// 0x27 Decimal Adjust Accumulator
func daa(cpu *CPU) (uint8, bool) {
	var flagC bool = cpu.getFlag(FlagC)

	if cpu.getFlag(FlagN) {
		// subtraction
		var adjustment uint8 = 0
		if cpu.getFlag(FlagH) {
			adjustment += 0x6
		}
		if cpu.getFlag(FlagC) {
			adjustment += 0x60
		}
		cpu.a -= adjustment
	} else {
		// addition
		var adjustment uint8 = 0
		if cpu.getFlag(FlagH) || (cpu.a&0xF) > 0x9 {
			adjustment += 0x6
		}
		if cpu.getFlag(FlagC) || (cpu.a > 0x99) {
			adjustment += 0x60
			flagC = true
		} else {
			flagC = false
		}
		cpu.a += adjustment
	}

	cpu.setFlag(FlagZ, cpu.a == 0)
	cpu.setFlag(FlagH, false)
	cpu.setFlag(FlagC, flagC)

	return 4, true
}

// 0x10 Enter CPU very low power mode
//
// Not implementing this for now because:
//
//	No licensed rom makes use of STOP outside of CGB speed switching.
func stop(cpu *CPU) (uint8, bool) {
	return 4, true
}

// CB Prefixed Instructions

// BIT u3,r8 - Test bit u3 in register r8, set the zero flag if bit not set.
//
// BIT u3,[HL] - Test bit u3 in the byte pointed by HL, set the zero flag if bit not set.
func (cpu *CPU) bit_u3_r8(u3 uint8, r8 uint8) (uint8, bool) {
	hlOperation := r8 == 0b110

	setFlags := func(bit uint8) {
		cpu.setFlag(FlagZ, bit == 0)
		cpu.setFlag(FlagN, false)
		cpu.setFlag(FlagH, true)
	}

	switch cpu.mCycle {
	case 1:
		// fetch CB prefix
		return 4, false
	case 2:
		// fetch opcode
		if !hlOperation {
			testValue := cpu.get_r8(r8)
			cpu.immediateValue = uint16(testValue)
			bit := (uint8(cpu.immediateValue) >> u3) & 1
			setFlags(bit)
			return 4, true
		}
		return 4, false
	case 3:
		// at this point, it should be the special HL case
		testValue := cpu.get_r8(r8)
		cpu.immediateValue = uint16(testValue)
		bit := (uint8(cpu.immediateValue) >> u3) & 1
		setFlags(bit)
		return 4, true
	}
	return 4, true
}

// RES u3,r8 - Set bit u3 in register r8 to 0.
//
// RES u3,[HL] - Set bit u3 in the byte pointed by HL to 0.
func (cpu *CPU) res_u3_r8(u3 uint8, r8 uint8) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		if r8 != 0b110 {
			value := cpu.get_r8(r8)
			cpu.set_r8(r8, value & ^(1<<u3))
			return 4, true
		}
		return 4, false
	case 3:
		cpu.mdr = cpu.get_r8(r8)
		return 4, false
	case 4:
		cpu.set_r8(r8, cpu.mdr & ^(1<<u3))
		return 4, true
	}

	return 4, true
}

// SET u3,r8 - Set bit u3 in register r8 to 1.
//
// SET u3,[HL] - Set bit u3 in the byte pointed by HL to 1.
func (cpu *CPU) set_u3_r8(u3 uint8, r8 uint8) (uint8, bool) {
	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		if r8 != 0b110 {
			value := cpu.get_r8(r8)
			cpu.set_r8(r8, value|(1<<u3))
			return 4, true
		}
		return 4, false
	case 3:
		cpu.mdr = cpu.get_r8(r8)
		return 4, false
	case 4:
		cpu.set_r8(r8, cpu.mdr|(1<<u3))
		return 4, true
	}

	return 4, true
}

// SLA r8 - Shift Left Arithmetically register r8.
//
// SLA [HL] - Shift Left Arithmetically the byte pointed to by HL.
//
// SRA r8 - Shift Right Arithmetically register r8 (bit 7 of r8 is unchanged).
//
// SRA [HL] - Shift Right Arithmetically the byte pointed to by HL (bit 7 of the byte pointed to by HL is unchanged).
//
// SRL r8 - Shift Right Logically register r8.
//
// SRL [HL] - Shift Right Logically the byte pointed to by HL.
//
// RLC r8 - Rotate register r8 left.
//
// RLC [HL] - Rotate the byte pointed to by HL left.
//
// RRC r8 - Rotate register r8 right.
//
// RRC [HL] - Rotate the byte pointed to by HL right.
//
// RL r8 - Rotate bits in register r8 left, through the carry flag.
//
// RL [HL] - Rotate the byte pointed to by HL left, through the carry flag.
//
// RR r8 - Rotate register r8 right, through the carry flag.
//
// RR [HL] - Rotate the byte pointed to by HL right, through the carry flag.
func (cpu *CPU) shift_rotate_u3_r8(u3 uint8, r8 uint8) (uint8, bool) {
	var fetchCycles uint8 = 4
	var writeCycles uint8 = 0

	// parse the middle 3 bits
	switch u3 {
	case 0b000:
		// RLC
		return cpu.rlc(r8)
	case 0b001:
		// RRC
		return cpu.rrc(r8)
	case 0b010:
		// RL
		return cpu.rl(r8)
	case 0b011:
		// RR
		return cpu.rr(r8)
	case 0b100:
		// SLA
		return cpu.sla(r8)
	case 0b101:
		// SRA
		return cpu.sra(r8)
	case 0b111:
		// SRL
		return cpu.srl(r8)
	case 0b110:
		// SWAP
		return cpu.swap(r8)
	}

	return 8 + fetchCycles + writeCycles, true
}

func (cpu *CPU) rlc(r8 uint8) (uint8, bool) {
	var setFlags = func(result uint8, bit7 uint8) {
		cpu.setFlag(FlagZ, result == 0)
		cpu.setFlag(FlagN, false)
		cpu.setFlag(FlagH, false)
		cpu.setFlag(FlagC, (bit7>>7) == 1)
	}

	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		if r8 != 0b110 {
			value := cpu.get_r8(r8)
			bit7 := (value & 0b1000_0000)
			result := (value << 1) | (bit7 >> 7)
			cpu.set_r8(r8, result)
			setFlags(result, bit7)
			return 4, true
		}
		// fetch the CB opcode
		return 4, false
	case 3:
		cpu.mdr = cpu.get_r8(r8)
		return 4, false
	case 4:
		bit7 := (cpu.mdr & 0b1000_0000)
		result := (cpu.mdr << 1) | (bit7 >> 7)
		cpu.set_r8(r8, result)
		setFlags(result, bit7)
		return 4, true
	}

	return 4, true
}

func (cpu *CPU) rrc(r8 uint8) (uint8, bool) {
	var setFlags = func(result uint8, bit0 uint8) {
		cpu.setFlag(FlagZ, result == 0)
		cpu.setFlag(FlagN, false)
		cpu.setFlag(FlagH, false)
		cpu.setFlag(FlagC, bit0 == 1)
	}

	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		if r8 != 0b110 {
			value := cpu.get_r8(r8)
			bit0 := value & 1
			result := (value >> 1) | (bit0 << 7)
			cpu.set_r8(r8, result)
			cpu.set_r8(r8, result)
			setFlags(result, bit0)
			return 4, true
		}
		// fetch the CB opcode
		return 4, false
	case 3:
		cpu.mdr = cpu.get_r8(r8)
		return 4, false
	case 4:
		bit0 := cpu.mdr & 1
		result := (cpu.mdr >> 1) | (bit0 << 7)
		cpu.set_r8(r8, result)
		setFlags(result, bit0)
		return 4, true
	}

	return 4, true
}

func (cpu *CPU) rl(r8 uint8) (uint8, bool) {
	var setFlags = func(result uint8, bit7 uint8) {
		cpu.setFlag(FlagZ, result == 0)
		cpu.setFlag(FlagN, false)
		cpu.setFlag(FlagH, false)
		cpu.setFlag(FlagC, (bit7>>7) == 1)
	}

	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		if r8 != 0b110 {
			value := cpu.get_r8(r8)
			bit7 := (value & 0b1000_0000)
			var mask uint8 = 0
			if cpu.getFlag(FlagC) {
				mask = 1
			}
			result := (value << 1) | (mask)
			cpu.set_r8(r8, result)
			setFlags(result, bit7)
			return 4, true
		}
		// fetch the CB opcode
		return 4, false
	case 3:
		cpu.mdr = cpu.get_r8(r8)
		return 4, false
	case 4:
		bit7 := (cpu.mdr & 0b1000_0000)
		var mask uint8 = 0
		if cpu.getFlag(FlagC) {
			mask = 1
		}
		result := (cpu.mdr << 1) | (mask)
		cpu.set_r8(r8, result)
		setFlags(result, bit7)
		return 4, true
	}

	return 4, true
}

func (cpu *CPU) rr(r8 uint8) (uint8, bool) {
	var setFlags = func(result uint8, bit0 uint8) {
		cpu.setFlag(FlagZ, result == 0)
		cpu.setFlag(FlagN, false)
		cpu.setFlag(FlagH, false)
		cpu.setFlag(FlagC, bit0 == 1)
	}

	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		if r8 != 0b110 {
			value := cpu.get_r8(r8)
			bit0 := value & 1
			var mask uint8 = 0
			if cpu.getFlag(FlagC) {
				mask = 1
			}
			result := (value >> 1) | (mask << 7)
			cpu.set_r8(r8, result)
			setFlags(result, bit0)
			return 4, true
		}
		// fetch the CB opcode
		return 4, false
	case 3:
		cpu.mdr = cpu.get_r8(r8)
		return 4, false
	case 4:
		bit0 := cpu.mdr & 1
		var mask uint8 = 0
		if cpu.getFlag(FlagC) {
			mask = 1
		}
		result := (cpu.mdr >> 1) | (mask << 7)
		cpu.set_r8(r8, result)
		setFlags(result, bit0)
		return 4, true
	}

	return 4, true
}

func (cpu *CPU) sla(r8 uint8) (uint8, bool) {
	var setFlags = func(result uint8, bit7 uint8) {
		cpu.setFlag(FlagZ, result == 0)
		cpu.setFlag(FlagN, false)
		cpu.setFlag(FlagH, false)
		cpu.setFlag(FlagC, (bit7>>7) == 1)
	}

	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		if r8 != 0b110 {
			value := cpu.get_r8(r8)
			bit7 := (value & 0b1000_0000)
			result := value << 1
			cpu.set_r8(r8, result)
			setFlags(result, bit7)
			return 4, true
		}
		// fetch the CB opcode
		return 4, false
	case 3:
		cpu.mdr = cpu.get_r8(r8)
		return 4, false
	case 4:
		bit7 := (cpu.mdr & 0b1000_0000)
		result := cpu.mdr << 1
		cpu.set_r8(r8, result)
		setFlags(result, bit7)
		return 4, true
	}

	return 4, true
}

func (cpu *CPU) sra(r8 uint8) (uint8, bool) {
	var setFlags = func(result uint8, bit0 uint8) {
		cpu.setFlag(FlagZ, result == 0)
		cpu.setFlag(FlagN, false)
		cpu.setFlag(FlagH, false)
		cpu.setFlag(FlagC, bit0 == 1)
	}

	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		if r8 != 0b110 {
			value := cpu.get_r8(r8)
			bit0 := value & 1
			bit7 := (value & 0b1000_0000)
			result := (value >> 1) | bit7
			cpu.set_r8(r8, result)
			setFlags(result, bit0)
			return 4, true
		}
		// fetch the CB opcode
		return 4, false
	case 3:
		cpu.mdr = cpu.get_r8(r8)
		return 4, false
	case 4:
		bit0 := cpu.mdr & 1
		bit7 := (cpu.mdr & 0b1000_0000)
		result := (cpu.mdr >> 1) | bit7
		cpu.set_r8(r8, result)
		setFlags(result, bit0)
		return 4, true
	}

	return 4, true
}

func (cpu *CPU) srl(r8 uint8) (uint8, bool) {
	var setFlags = func(result uint8, bit0 uint8) {
		cpu.setFlag(FlagZ, result == 0)
		cpu.setFlag(FlagN, false)
		cpu.setFlag(FlagH, false)
		cpu.setFlag(FlagC, bit0 == 1)
	}

	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		if r8 != 0b110 {
			value := cpu.get_r8(r8)
			bit0 := value & 1
			result := value >> 1
			cpu.set_r8(r8, result)
			setFlags(result, bit0)
			return 4, true
		}
		// fetch the CB opcode
		return 4, false
	case 3:
		cpu.mdr = cpu.get_r8(r8)
		return 4, false
	case 4:
		bit0 := cpu.mdr & 1
		result := cpu.mdr >> 1
		cpu.set_r8(r8, result)
		setFlags(result, bit0)
		return 4, true
	}

	return 4, true
}

func (cpu *CPU) swap(r8 uint8) (uint8, bool) {
	var setFlags = func(result uint8) {
		cpu.setFlag(FlagZ, result == 0)
		cpu.setFlag(FlagN, false)
		cpu.setFlag(FlagH, false)
		cpu.setFlag(FlagC, false)
	}

	switch cpu.mCycle {
	case 1:
		return 4, false
	case 2:
		if r8 != 0b110 {
			value := cpu.get_r8(r8)
			upperNibble := value & 0xF0
			lowerNibble := value & 0x0F
			result := (lowerNibble << 4) | (upperNibble >> 4)
			cpu.set_r8(r8, result)
			setFlags(result)
			return 4, true
		}
		// fetch the CB opcode
		return 4, false
	case 3:
		cpu.mdr = cpu.get_r8(r8)
		return 4, false
	case 4:
		upperNibble := cpu.mdr & 0xF0
		lowerNibble := cpu.mdr & 0x0F
		result := (lowerNibble << 4) | (upperNibble >> 4)
		cpu.set_r8(r8, result)
		setFlags(result)
		return 4, true
	}

	return 4, true
}
