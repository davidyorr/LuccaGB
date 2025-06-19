package cpu

import "fmt"

// 0x00 No OPeration
func nop(cpu *CPU) uint8 {
	return 4
}

// 0x06 Copy the value n8 into register r8
func ld_b_n8(cpu *CPU) uint8 {
	cpu.b = uint8(cpu.immediateValue)
	return 8
}

// 0x0E Copy the value n8 into register r8
func ld_c_n8(cpu *CPU) uint8 {
	cpu.c = uint8(cpu.immediateValue)
	return 8
}

// 0x16 Copy the value n8 into register r8
func ld_d_n8(cpu *CPU) uint8 {
	cpu.d = uint8(cpu.immediateValue)
	return 8
}

// 0x1E Copy the value n8 into register r8
func ld_e_n8(cpu *CPU) uint8 {
	cpu.e = uint8(cpu.immediateValue)
	return 8
}

// 0x26 Copy the value n8 into register r8
func ld_h_n8(cpu *CPU) uint8 {
	cpu.h = uint8(cpu.immediateValue)
	return 8
}

// 0x2E Copy the value n8 into register r8
func ld_l_n8(cpu *CPU) uint8 {
	cpu.l = uint8(cpu.immediateValue)
	return 8
}

// 0x3E Copy the value n8 into register r8
func ld_a_n8(cpu *CPU) uint8 {
	cpu.a = uint8(cpu.immediateValue)
	return 8
}

// 0x40 Copy (aka Load) the value in register on the right into the register on the left
func ld_b_b(cpu *CPU) uint8 {
	cpu.b = cpu.b
	return 4
}

// 0x41 Copy (aka Load) the value in register on the right into the register on the left
func ld_b_c(cpu *CPU) uint8 {
	cpu.b = cpu.c
	return 4
}

// 0x42 Copy (aka Load) the value in register on the right into the register on the left
func ld_b_d(cpu *CPU) uint8 {
	cpu.b = cpu.d
	return 4
}

// 0x43 Copy (aka Load) the value in register on the right into the register on the left
func ld_b_e(cpu *CPU) uint8 {
	cpu.b = cpu.e
	return 4
}

// 0x44 Copy (aka Load) the value in register on the right into the register on the left
func ld_b_h(cpu *CPU) uint8 {
	cpu.b = cpu.h
	return 4
}

// 0x45 Copy (aka Load) the value in register on the right into the register on the left
func ld_b_l(cpu *CPU) uint8 {
	cpu.b = cpu.l
	return 4
}

// 0x47 Copy (aka Load) the value in register on the right into the register on the left
func ld_b_a(cpu *CPU) uint8 {
	cpu.b = cpu.a
	return 4
}

// 0x48 Copy (aka Load) the value in register on the right into the register on the left
func ld_c_b(cpu *CPU) uint8 {
	cpu.c = cpu.b
	return 4
}

// 0x49 Copy (aka Load) the value in register on the right into the register on the left
func ld_c_c(cpu *CPU) uint8 {
	cpu.c = cpu.c
	return 4
}

// 0x4A Copy (aka Load) the value in register on the right into the register on the left
func ld_c_d(cpu *CPU) uint8 {
	cpu.c = cpu.d
	return 4
}

// 0x4B Copy (aka Load) the value in register on the right into the register on the left
func ld_c_e(cpu *CPU) uint8 {
	cpu.c = cpu.e
	return 4
}

// 0x4C Copy (aka Load) the value in register on the right into the register on the left
func ld_c_h(cpu *CPU) uint8 {
	cpu.c = cpu.h
	return 4
}

// 0x4L Copy (aka Load) the value in register on the right into the register on the left
func ld_c_l(cpu *CPU) uint8 {
	cpu.c = cpu.l
	return 4
}

// 0x4F Copy (aka Load) the value in register on the right into the register on the left
func ld_c_a(cpu *CPU) uint8 {
	cpu.c = cpu.a
	return 4
}

// 0x50 Copy (aka Load) the value in register on the right into the register on the left
func ld_d_b(cpu *CPU) uint8 {
	cpu.d = cpu.b
	return 4
}

// 0x51 Copy (aka Load) the value in register on the right into the register on the left
func ld_d_c(cpu *CPU) uint8 {
	cpu.d = cpu.c
	return 4
}

// 0x52 Copy (aka Load) the value in register on the right into the register on the left
func ld_d_d(cpu *CPU) uint8 {
	cpu.d = cpu.d
	return 4
}

// 0x53 Copy (aka Load) the value in register on the right into the register on the left
func ld_d_e(cpu *CPU) uint8 {
	cpu.d = cpu.e
	return 4
}

// 0x54 Copy (aka Load) the value in register on the right into the register on the left
func ld_d_h(cpu *CPU) uint8 {
	cpu.d = cpu.h
	return 4
}

// 0x55 Copy (aka Load) the value in register on the right into the register on the left
func ld_d_l(cpu *CPU) uint8 {
	cpu.d = cpu.l
	return 4
}

// 0x57 Copy (aka Load) the value in register on the right into the register on the left
func ld_d_a(cpu *CPU) uint8 {
	cpu.d = cpu.a
	return 4
}

// 0x58 Copy (aka Load) the value in register on the right into the register on the left
func ld_e_b(cpu *CPU) uint8 {
	cpu.e = cpu.b
	return 4
}

// 0x59 Copy (aka Load) the value in register on the right into the register on the left
func ld_e_c(cpu *CPU) uint8 {
	cpu.e = cpu.c
	return 4
}

// 0x5A Copy (aka Load) the value in register on the right into the register on the left
func ld_e_d(cpu *CPU) uint8 {
	cpu.e = cpu.d
	return 4
}

// 0x5B Copy (aka Load) the value in register on the right into the register on the left
func ld_e_e(cpu *CPU) uint8 {
	cpu.e = cpu.e
	return 4
}

// 0x5C Copy (aka Load) the value in register on the right into the register on the left
func ld_e_h(cpu *CPU) uint8 {
	cpu.e = cpu.h
	return 4
}

// 0x5D Copy (aka Load) the value in register on the right into the register on the left
func ld_e_l(cpu *CPU) uint8 {
	cpu.e = cpu.l
	return 4
}

// 0x5F Copy (aka Load) the value in register on the right into the register on the left
func ld_e_a(cpu *CPU) uint8 {
	cpu.e = cpu.a
	return 4
}

// 0x60 Copy (aka Load) the value in register on the right into the register on the left
func ld_h_b(cpu *CPU) uint8 {
	cpu.h = cpu.b
	return 4
}

// 0x61 Copy (aka Load) the value in register on the right into the register on the left
func ld_h_c(cpu *CPU) uint8 {
	cpu.h = cpu.c
	return 4
}

// 0x62 Copy (aka Load) the value in register on the right into the register on the left
func ld_h_d(cpu *CPU) uint8 {
	cpu.h = cpu.d
	return 4
}

// 0x63 Copy (aka Load) the value in register on the right into the register on the left
func ld_h_e(cpu *CPU) uint8 {
	cpu.h = cpu.e
	return 4
}

// 0x64 Copy (aka Load) the value in register on the right into the register on the left
func ld_h_h(cpu *CPU) uint8 {
	cpu.h = cpu.h
	return 4
}

// 0x65 Copy (aka Load) the value in register on the right into the register on the left
func ld_h_l(cpu *CPU) uint8 {
	cpu.h = cpu.l
	return 4
}

// 0x67 Copy (aka Load) the value in register on the right into the register on the left
func ld_h_a(cpu *CPU) uint8 {
	cpu.h = cpu.a
	return 4
}

// 0x68 Copy (aka Load) the value in register on the right into the register on the left
func ld_l_b(cpu *CPU) uint8 {
	cpu.l = cpu.b
	return 4
}

// 0x69 Copy (aka Load) the value in register on the right into the register on the left
func ld_l_c(cpu *CPU) uint8 {
	cpu.l = cpu.c
	return 4
}

// 0x6A Copy (aka Load) the value in register on the right into the register on the left
func ld_l_d(cpu *CPU) uint8 {
	cpu.l = cpu.d
	return 4
}

// 0x6B Copy (aka Load) the value in register on the right into the register on the left
func ld_l_e(cpu *CPU) uint8 {
	cpu.l = cpu.e
	return 4
}

// 0x6C Copy (aka Load) the value in register on the right into the register on the left
func ld_l_h(cpu *CPU) uint8 {
	cpu.l = cpu.h
	return 4
}

// 0x6D Copy (aka Load) the value in register on the right into the register on the left
func ld_l_l(cpu *CPU) uint8 {
	cpu.l = cpu.l
	return 4
}

// 0x6F Copy (aka Load) the value in register on the right into the register on the left
func ld_l_a(cpu *CPU) uint8 {
	cpu.l = cpu.a
	return 4
}

// 0x78 Copy (aka Load) the value in register on the right into the register on the left
func ld_a_b(cpu *CPU) uint8 {
	cpu.a = cpu.b
	return 4
}

// 0x79 Copy (aka Load) the value in register on the right into the register on the left
func ld_a_c(cpu *CPU) uint8 {
	cpu.a = cpu.c
	return 4
}

// 0x7A Copy (aka Load) the value in register on the right into the register on the left
func ld_a_d(cpu *CPU) uint8 {
	cpu.a = cpu.d
	return 4
}

// 0x7B Copy (aka Load) the value in register on the right into the register on the left
func ld_a_e(cpu *CPU) uint8 {
	cpu.a = cpu.e
	return 4
}

// 0x7F Copy (aka Load) the value in register on the right into the register on the left
func ld_a_a(cpu *CPU) uint8 {
	cpu.a = cpu.a
	return 4
}

// 0x46 Copy the value pointed to by HL into register r8
func ld_b_hl(cpu *CPU) uint8 {
	cpu.b = cpu.bus.Read(cpu.getHL())
	return 8
}

// 0x4E Copy the value pointed to by HL into register r8
func ld_c_hl(cpu *CPU) uint8 {
	cpu.c = cpu.bus.Read(cpu.getHL())
	return 8
}

// 0x56 Copy the value pointed to by HL into register r8
func ld_d_hl(cpu *CPU) uint8 {
	cpu.d = cpu.bus.Read(cpu.getHL())
	return 8
}

// 0x5E Copy the value pointed to by HL into register r8
func ld_e_hl(cpu *CPU) uint8 {
	cpu.e = cpu.bus.Read(cpu.getHL())
	return 8
}

// 0x66 Copy the value pointed to by HL into register r8
func ld_h_hl(cpu *CPU) uint8 {
	cpu.h = cpu.bus.Read(cpu.getHL())
	return 8
}

// 0x6E Copy the value pointed to by HL into register r8
func ld_l_hl(cpu *CPU) uint8 {
	cpu.l = cpu.bus.Read(cpu.getHL())
	return 8
}

// 0x01 Copy the value n16 into register r16
func ld_bc_n16(cpu *CPU) uint8 {
	cpu.setBC(cpu.immediateValue)
	return 12
}

// 0x11 Copy the value n16 into register r16
func ld_de_n16(cpu *CPU) uint8 {
	cpu.setDE(cpu.immediateValue)
	return 12
}

// 0x21 Copy the value n16 into register r16
func ld_hl_n16(cpu *CPU) uint8 {
	cpu.setHL(cpu.immediateValue)
	return 12
}

// 0x31 Copy the value n16 into register r16
func ld_sp_n16(cpu *CPU) uint8 {
	cpu.sp = cpu.immediateValue
	return 12
}

// 0x70 Copy the value in register r8 into the byte pointed to by HL
func ld_hl_b(cpu *CPU) uint8 {
	address := cpu.getHL()
	cpu.bus.Write(uint16(address), cpu.b)
	return 8
}

// 0x71 Copy the value in register r8 into the byte pointed to by HL
func ld_hl_c(cpu *CPU) uint8 {
	address := cpu.getHL()
	cpu.bus.Write(uint16(address), cpu.c)
	return 8
}

// 0x72 Copy the value in register r8 into the byte pointed to by HL
func ld_hl_d(cpu *CPU) uint8 {
	address := cpu.getHL()
	cpu.bus.Write(uint16(address), cpu.d)
	return 8
}

// 0x73 Copy the value in register r8 into the byte pointed to by HL
func ld_hl_e(cpu *CPU) uint8 {
	address := cpu.getHL()
	cpu.bus.Write(uint16(address), cpu.e)
	return 8
}

// 0x74 Copy the value in register r8 into the byte pointed to by HL
func ld_hl_h(cpu *CPU) uint8 {
	address := cpu.getHL()
	cpu.bus.Write(uint16(address), cpu.h)
	return 8
}

// 0x75 Copy the value in register r8 into the byte pointed to by HL
func ld_hl_l(cpu *CPU) uint8 {
	address := cpu.getHL()
	cpu.bus.Write(uint16(address), cpu.l)
	return 8
}

// 0x77 Copy the value in register r8 into the byte pointed to by HL
func ld_hl_a(cpu *CPU) uint8 {
	address := cpu.getHL()
	cpu.bus.Write(uint16(address), cpu.a)
	return 8
}

// 0x7C Copy (aka Load) the value in register on the right into the register on the left
func ld_a_h(cpu *CPU) uint8 {
	cpu.a = cpu.h
	return 4
}

// 0x7D Copy (aka Load) the value in register on the right into the register on the left
func ld_a_l(cpu *CPU) uint8 {
	cpu.a = cpu.l
	return 4
}

// 0xE0 Copy the value in register A into the byte at address n16, provided the address is between $FF00 and $FFFF
func ldh_a8_a(cpu *CPU) uint8 {
	cpu.bus.Write(0xFF00+cpu.immediateValue, cpu.a)
	return 12
}

// 0xF0 Copy the byte at address n16 into register A, provided the address is between $FF00 and $FFFF
func ldh_a_a8(cpu *CPU) uint8 {
	cpu.a = cpu.bus.Read(0xFF00 + cpu.immediateValue)
	return 12
}

// 0xEA Copy the value in register A into the byte at address n16
func ld_a16_a(cpu *CPU) uint8 {
	cpu.bus.Write(cpu.immediateValue, cpu.a)
	return 16
}

// 0x0A Copy the byte pointed to by r16 into register A
func ld_a_bc(cpu *CPU) uint8 {
	cpu.a = cpu.bus.Read(cpu.getBC())
	return 8
}

// 0x1A Copy the byte pointed to by r16 into register A
func ld_a_de(cpu *CPU) uint8 {
	cpu.a = cpu.bus.Read(cpu.getDE())
	return 8
}

// 0x7E Copy the byte pointed to by r16 into register A
func ld_a_hl(cpu *CPU) uint8 {
	cpu.a = cpu.bus.Read(cpu.getHL())
	return 8
}

// 0xFA Copy the byte at address n16 into register A
func ld_a_a16(cpu *CPU) uint8 {
	cpu.a = cpu.bus.Read(cpu.immediateValue)
	return 16
}

// 0x22 Copy the value in register A into the byte pointed by HL and increment HL afterwards
func ld_hli_a(cpu *CPU) uint8 {
	address := cpu.getHL()
	cpu.bus.Write(address, cpu.a)
	cpu.setHL(address + 1)
	return 8
}

// 0x32 Copy the value in register A into the byte pointed by HL and decrement HL afterwards
func ld_hld_a(cpu *CPU) uint8 {
	address := cpu.getHL()
	cpu.bus.Write(address, cpu.a)
	cpu.setHL(address - 1)
	return 8
}

// 0x2A Copy the byte pointed to by HL into register A, and increment HL afterwards
func ld_a_hli(cpu *CPU) uint8 {
	value := cpu.getHL()
	cpu.a = cpu.bus.Read(value)
	cpu.setHL(value + 1)
	return 8
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
func add_a_b(cpu *CPU) uint8 {
	cpu.add_a_r8(cpu.b)
	return 4
}

// 0x81 Add the value in r8 to A
func add_a_c(cpu *CPU) uint8 {
	cpu.add_a_r8(cpu.c)
	return 4
}

// 0x82 Add the value in r8 to A
func add_a_d(cpu *CPU) uint8 {
	cpu.add_a_r8(cpu.d)
	return 4
}

// 0x83 Add the value in r8 to A
func add_a_e(cpu *CPU) uint8 {
	cpu.add_a_r8(cpu.e)
	return 4
}

// 0x84 Add the value in r8 to A
func add_a_h(cpu *CPU) uint8 {
	cpu.add_a_r8(cpu.h)
	return 4
}

// 0x85 Add the value in r8 to A
func add_a_l(cpu *CPU) uint8 {
	cpu.add_a_r8(cpu.l)
	return 4
}

// 0x87 Add the value in r8 to A
func add_a_a(cpu *CPU) uint8 {
	cpu.add_a_r8(cpu.a)
	return 4
}

// 0xC6 Add the value n8 to A
func add_a_n8(cpu *CPU) uint8 {
	cpu.add_a_r8(uint8(cpu.immediateValue))
	return 8
}

// ComPare the value in A with the value in r8, then return the difference.
func (cpu *CPU) cp_a_r8(r8 uint8) uint8 {
	originalA := cpu.a
	difference := originalA - r8
	cpu.setFlag(FlagZ, difference == 0)
	cpu.setFlag(FlagN, true)
	cpu.setFlag(FlagC, difference > originalA)
	cpu.setFlag(FlagH, (originalA&0x0F) < (r8&0x0F))

	return difference
}

// 0xB8 ComPare the value in A with the value in r8
func cp_a_b(cpu *CPU) uint8 {
	cpu.cp_a_r8(cpu.b)
	return 4
}

// 0xB9 ComPare the value in A with the value in r8
func cp_a_c(cpu *CPU) uint8 {
	cpu.cp_a_r8(cpu.c)
	return 4
}

// 0xBA ComPare the value in A with the value in r8
func cp_a_d(cpu *CPU) uint8 {
	cpu.cp_a_r8(cpu.d)
	return 4
}

// 0xBB ComPare the value in A with the value in r8
func cp_a_e(cpu *CPU) uint8 {
	cpu.cp_a_r8(cpu.e)
	return 4
}

// 0xBC ComPare the value in A with the value in r8
func cp_a_h(cpu *CPU) uint8 {
	cpu.cp_a_r8(cpu.h)
	return 4
}

// 0xBD ComPare the value in A with the value in r8
func cp_a_l(cpu *CPU) uint8 {
	cpu.cp_a_r8(cpu.l)
	return 4
}

// 0xBF ComPare the value in A with the value in r8
func cp_a_a(cpu *CPU) uint8 {
	cpu.cp_a_r8(cpu.a)
	return 4
}

// 0xFE ComPare the value in A with the value n8
func cp_a_n8(cpu *CPU) uint8 {
	cpu.cp_a_r8(uint8(cpu.immediateValue))
	return 4
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
func dec_b(cpu *CPU) uint8 {
	cpu.b = cpu.dec_r8(cpu.b)
	return 4
}

// 0x0D Decrement the value in register r8 by 1
func dec_c(cpu *CPU) uint8 {
	cpu.c = cpu.dec_r8(cpu.c)
	return 4
}

// 0x15 Decrement the value in register r8 by 1
func dec_d(cpu *CPU) uint8 {
	cpu.d = cpu.dec_r8(cpu.d)
	return 4
}

// 0x1D Decrement the value in register r8 by 1
func dec_e(cpu *CPU) uint8 {
	cpu.e = cpu.dec_r8(cpu.e)
	return 4
}

// 0x25 Decrement the value in register r8 by 1
func dec_h(cpu *CPU) uint8 {
	cpu.h = cpu.dec_r8(cpu.h)
	return 4
}

// 0x2D Decrement the value in register r8 by 1
func dec_l(cpu *CPU) uint8 {
	cpu.l = cpu.dec_r8(cpu.l)
	return 4
}

// 0x3D Decrement the value in register r8 by 1
func dec_a(cpu *CPU) uint8 {
	cpu.a = cpu.dec_r8(cpu.a)
	return 4
}

// Increment the value in register r8 by 1
func (cpu *CPU) inc_r8(r8 uint8) uint8 {
	result := r8 + 1
	cpu.setFlag(FlagZ, result == 0)
	cpu.setFlag(FlagN, false)
	// if overflow from bit 3
	cpu.setFlag(FlagH, (r8&0x0F)+1 > 0x0F)
	return result
}

// 0x04 Increment the value in register r8 by 1
func inc_b(cpu *CPU) uint8 {
	cpu.b = cpu.inc_r8(cpu.b)
	return 4
}

// 0x0C Increment the value in register r8 by 1
func inc_c(cpu *CPU) uint8 {
	cpu.c = cpu.inc_r8(cpu.c)
	return 4
}

// 0x14 Increment the value in register r8 by 1
func inc_d(cpu *CPU) uint8 {
	cpu.d = cpu.inc_r8(cpu.d)
	return 4
}

// 0x1C Increment the value in register r8 by 1
func inc_e(cpu *CPU) uint8 {
	cpu.e = cpu.inc_r8(cpu.e)
	return 4
}

// 0x24 Increment the value in register r8 by 1
func inc_h(cpu *CPU) uint8 {
	cpu.h = cpu.inc_r8(cpu.h)
	return 4
}

// 0x2C Increment the value in register r8 by 1
func inc_l(cpu *CPU) uint8 {
	cpu.l = cpu.inc_r8(cpu.l)
	return 4
}

// 0x3C Increment the value in register r8 by 1
func inc_a(cpu *CPU) uint8 {
	cpu.a = cpu.inc_r8(cpu.a)
	return 4
}

// 0x90 Subtract the value in r8 from A
func sub_a_b(cpu *CPU) uint8 {
	cpu.a = cpu.cp_a_r8(cpu.b)
	return 4
}

// 0x91 Subtract the value in r8 from A
func sub_a_c(cpu *CPU) uint8 {
	cpu.a = cpu.cp_a_r8(cpu.c)
	return 4
}

// 0x92 Subtract the value in r8 from A
func sub_a_d(cpu *CPU) uint8 {
	cpu.a = cpu.cp_a_r8(cpu.d)
	return 4
}

// 0x93 Subtract the value in r8 from A
func sub_a_e(cpu *CPU) uint8 {
	cpu.a = cpu.cp_a_r8(cpu.e)
	return 4
}

// 0x94 Subtract the value in r8 from A
func sub_a_h(cpu *CPU) uint8 {
	cpu.a = cpu.cp_a_r8(cpu.h)
	return 4
}

// 0x95 Subtract the value in r8 from A
func sub_a_l(cpu *CPU) uint8 {
	cpu.a = cpu.cp_a_r8(cpu.l)
	return 4
}

// 0x97 Subtract the value in r8 from A
func sub_a_a(cpu *CPU) uint8 {
	cpu.a = cpu.cp_a_r8(cpu.a)
	return 4
}

// 0x03 Increment the value in register r8 by 1
func inc_bc(cpu *CPU) uint8 {
	value := cpu.getBC()
	cpu.setBC(value + 1)
	return 8
}

// 0x13 Increment the value in register r8 by 1
func inc_de(cpu *CPU) uint8 {
	value := cpu.getDE()
	cpu.setDE(value + 1)
	return 8
}

// 0x23 Increment the value in register r8 by 1
func inc_hl(cpu *CPU) uint8 {
	value := cpu.getHL()
	cpu.setHL(value + 1)
	return 8
}

// 0xD6 Subtract the value n8 from A
func sub_a_n8(cpu *CPU) uint8 {
	cpu.a -= uint8(cpu.immediateValue)
	return 8
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
func and_a_b(cpu *CPU) uint8 {
	cpu.and_a_r8(cpu.b)
	return 4
}

// 0xA1 Set A to the bitwise AND between the value in r8 and A
func and_a_c(cpu *CPU) uint8 {
	cpu.and_a_r8(cpu.c)
	return 4
}

// 0xA2 Set A to the bitwise AND between the value in r8 and A
func and_a_d(cpu *CPU) uint8 {
	cpu.and_a_r8(cpu.d)
	return 4
}

// 0xA3 Set A to the bitwise AND between the value in r8 and A
func and_a_e(cpu *CPU) uint8 {
	cpu.and_a_r8(cpu.e)
	return 4
}

// 0xA4 Set A to the bitwise AND between the value in r8 and A
func and_a_h(cpu *CPU) uint8 {
	cpu.and_a_r8(cpu.h)
	return 4
}

// 0xA5 Set A to the bitwise AND between the value in r8 and A
func and_a_l(cpu *CPU) uint8 {
	cpu.and_a_r8(cpu.l)
	return 4
}

// 0xA7 Set A to the bitwise AND between the value in r8 and A
func and_a_a(cpu *CPU) uint8 {
	cpu.and_a_r8(cpu.a)
	return 4
}

// 0xE6 Set A to the bitwise AND between the value n8 and A
func and_a_n8(cpu *CPU) uint8 {
	cpu.and_a_r8(uint8(cpu.immediateValue))
	return 8
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
func xor_a_b(cpu *CPU) uint8 {
	cpu.xor_a_n8(cpu.b)
	return 4
}

// 0xA9 Set A to the bitwise XOR between the value in r8 and A
func xor_a_c(cpu *CPU) uint8 {
	fmt.Println("--- EXECUTING NEW XOR HANDLER ---")
	cpu.xor_a_n8(cpu.c)
	return 4
}

// 0xAA Set A to the bitwise XOR between the value in r8 and A
func xor_a_d(cpu *CPU) uint8 {
	cpu.xor_a_n8(cpu.d)
	return 4
}

// 0xAB Set A to the bitwise XOR between the value in r8 and A
func xor_a_e(cpu *CPU) uint8 {
	cpu.xor_a_n8(cpu.e)
	return 4
}

// 0xAC Set A to the bitwise XOR between the value in r8 and A
func xor_a_h(cpu *CPU) uint8 {
	cpu.xor_a_n8(cpu.h)
	return 4
}

// 0xAD Set A to the bitwise XOR between the value in r8 and A
func xor_a_l(cpu *CPU) uint8 {
	cpu.xor_a_n8(cpu.l)
	return 4
}

// 0xAF Set A to the bitwise XOR between the value in r8 and A
func xor_a_a(cpu *CPU) uint8 {
	cpu.xor_a_n8(cpu.a)
	return 4
}

// 0xAE Set A to the bitwise XOR between the byte pointed to by HL and A
func xor_a_hl(cpu *CPU) uint8 {
	cpu.xor_a_n8(cpu.bus.Read(cpu.getHL()))
	return 8
}

// 0xEE Set A to the bitwise XOR between the value n8 and A
func xor_a_n8(cpu *CPU) uint8 {
	cpu.xor_a_n8(uint8(cpu.immediateValue))
	return 8
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
func or_a_b(cpu *CPU) uint8 {
	cpu.or_a_r8(cpu.b)
	return 4
}

// 0xB1 Set A to the bitwise OR between the value in r8 and A
func or_a_c(cpu *CPU) uint8 {
	cpu.or_a_r8(cpu.c)
	return 4
}

// 0xB2 Set A to the bitwise OR between the value in r8 and A
func or_a_d(cpu *CPU) uint8 {
	cpu.or_a_r8(cpu.d)
	return 4
}

// 0xB3 Set A to the bitwise OR between the value in r8 and A
func or_a_e(cpu *CPU) uint8 {
	cpu.or_a_r8(cpu.e)
	return 4
}

// 0xB4 Set A to the bitwise OR between the value in r8 and A
func or_a_h(cpu *CPU) uint8 {
	cpu.or_a_r8(cpu.h)
	return 4
}

// 0xB5 Set A to the bitwise OR between the value in r8 and A
func or_a_l(cpu *CPU) uint8 {
	cpu.or_a_r8(cpu.l)
	return 4
}

// 0xB7 Set A to the bitwise OR between the value in r8 and A
func or_a_a(cpu *CPU) uint8 {
	cpu.or_a_r8(cpu.a)
	return 4
}

// 0xF6 Set A to the bitwise OR between the value in r8 and A
func or_a_n8(cpu *CPU) uint8 {
	cpu.or_a_r8(uint8(cpu.immediateValue))
	return 8
}

// 0x07 Rotate register A left
func rlca(cpu *CPU) uint8 {
	carry := (cpu.a >> 7) & 1
	cpu.a = cpu.a << 1
	cpu.a |= carry

	cpu.setFlag(FlagZ, false)
	cpu.setFlag(FlagN, false)
	cpu.setFlag(FlagH, false)
	cpu.setFlag(FlagC, carry == 1)

	fmt.Printf("RCLA: carry=[0x%0X] a=[0x%0X]\n", carry, cpu.a)

	return 4
}

// 0xCD Call address n16
func call_a16(cpu *CPU) uint8 {
	cpu.pushToStack16(cpu.pc + 3)
	cpu.pc = cpu.immediateValue
	return 24
}

// 0xC4 Call address n16 if condition cc is not met
func call_nz_a16(cpu *CPU) uint8 {
	if cpu.getFlag(FlagZ) {
		return 12
	}

	cpu.pushToStack16(cpu.pc + 3)
	cpu.pc = cpu.immediateValue
	return 24
}

// 0xC3 Jump to address a16; effectively, copy a16 into PC
func jp_a16(cpu *CPU) uint8 {
	cpu.pc = cpu.immediateValue
	return 16
}

// Relative Jump
func (cpu *CPU) jr() {
	operandLength := uint16(2)
	signedOffset := int8(cpu.immediateValue)
	destinationAddress := uint16(int(cpu.pc+operandLength) + int(signedOffset))
	cpu.pc = destinationAddress
}

// 0x18 Relative Jump to address e8
func jr_e8(cpu *CPU) uint8 {
	cpu.jr()
	return 12
}

// 0x20 Relative Jump to address e8 if condition z is not met
func jr_nz_e8(cpu *CPU) uint8 {
	if cpu.getFlag(FlagZ) {
		return 8
	}

	cpu.jr()

	return 12
}

// 0x28 Relative Jump to address e8 if condition z is met
func jr_z_e8(cpu *CPU) uint8 {
	if !cpu.getFlag(FlagZ) {
		return 8
	}

	cpu.jr()

	return 12
}

// 0x30 Relative Jump to address e8 if condition c is not met
func jr_nc_e8(cpu *CPU) uint8 {
	if cpu.getFlag(FlagC) {
		return 8
	}

	cpu.jr()

	return 12
}

// 0x38 Relative Jump to address e8 if condition c is met
func jr_c_e8(cpu *CPU) uint8 {
	if !cpu.getFlag(FlagC) {
		return 8
	}

	cpu.jr()

	return 12
}

// 0xC9 Return from subroutine
func ret(cpu *CPU) uint8 {
	returnAddress := cpu.popFromStack16()
	cpu.pc = returnAddress
	return 16
}

// 0xFF Call address 0x38
func rst_38h(cpu *CPU) uint8 {
	cpu.pushToStack16(cpu.pc + 1)
	cpu.pc = 0x38
	return 16
}

// 0x33 Increment the value in register SP by 1
func inc_sp(cpu *CPU) uint8 {
	cpu.sp++
	return 8
}

// 0xC1 Pop register r16 from the stack
func pop_bc(cpu *CPU) uint8 {
	value := cpu.popFromStack16()
	cpu.setBC(value)
	return 12
}

// 0xD1 Pop register r16 from the stack
func pop_de(cpu *CPU) uint8 {
	value := cpu.popFromStack16()
	cpu.setDE(value)
	return 12
}

// 0xE1 Pop register r16 from the stack
func pop_hl(cpu *CPU) uint8 {
	value := cpu.popFromStack16()
	cpu.setHL(value)
	return 12
}

// 0xF1 Pop register AF from the stack
func pop_af(cpu *CPU) uint8 {
	value := cpu.popFromStack16()
	cpu.setAF(value)
	return 12
}

// 0xC5 Push register r16 into the stack
func push_bc(cpu *CPU) uint8 {
	cpu.pushToStack16(cpu.getBC())
	return 16
}

// 0xD5 Push register r16 into the stack
func push_de(cpu *CPU) uint8 {
	cpu.pushToStack16(cpu.getDE())
	return 16
}

// 0xE5 Push register r16 into the stack
func push_hl(cpu *CPU) uint8 {
	cpu.pushToStack16(cpu.getHL())
	return 16
}

// 0xF5 Push register r16 into the stack
func push_af(cpu *CPU) uint8 {
	cpu.pushToStack16(cpu.getAF())
	return 16
}

// 0xF3 Disable Interrupts by clearing the IME flag
func di(cpu *CPU) uint8 {
	cpu.ime = false
	return 4
}

// CB Prefixed Instructions

// BIT u3,r8 - Test bit u3 in register r8, set the zero flag if bit not set.
//
// BIT u3,[HL] - Test bit u3 in the byte pointed by HL, set the zero flag if bit not set.
func (cpu *CPU) bit_u3_r8(u3 uint8, r8 uint8) uint8 {
	testValue, cycles := cpu.get_r8(r8)
	bit := (testValue >> u3) & 1

	cpu.setFlag(FlagZ, bit == 0)
	cpu.setFlag(FlagN, false)
	cpu.setFlag(FlagH, true)

	return 8 + cycles
}

// RES u3,r8 - Set bit u3 in register r8 to 0.
//
// RES u3,[HL] - Set bit u3 in the byte pointed by HL to 0.
func (cpu *CPU) res_u3_r8(u3 uint8, r8 uint8) uint8 {
	value, fetchCycles := cpu.get_r8(r8)
	writeCycles := cpu.set_r8(r8, value & ^(1<<u3))

	return 8 + fetchCycles + writeCycles
}

// SET u3,r8 - Set bit u3 in register r8 to 1.
//
// SET u3,[HL] - Set bit u3 in the byte pointed by HL to 1.
func (cpu *CPU) set_u3_r8(u3 uint8, r8 uint8) uint8 {
	value, fetchCycles := cpu.get_r8(r8)
	writeCycles := cpu.set_r8(r8, value|(1<<u3))

	return 8 + fetchCycles + writeCycles
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
func (cpu *CPU) shift_rotate_u3_r8(u3 uint8, r8 uint8) uint8 {
	value, fetchCycles := cpu.get_r8(r8)
	var writeCycles uint8 = 0

	// parse the middle 3 bits
	switch u3 {
	case 0b000:
		// RLC
		bit7 := (value & 0b1000_0000)
		result := (value << 1) | (bit7 >> 7)
		writeCycles = cpu.set_r8(r8, result)
		cpu.setFlag(FlagZ, result == 0)
		cpu.setFlag(FlagN, false)
		cpu.setFlag(FlagH, false)
		cpu.setFlag(FlagC, (bit7>>7) == 1)
	case 0b001:
		// RRC
		bit0 := value & 1
		result := (value >> 1) | (bit0 << 7)
		writeCycles = cpu.set_r8(r8, result)
		cpu.setFlag(FlagZ, result == 0)
		cpu.setFlag(FlagN, false)
		cpu.setFlag(FlagH, false)
		cpu.setFlag(FlagC, bit0 == 1)
	case 0b010:
		// RL
		bit7 := (value & 0b1000_0000)
		var mask uint8 = 0
		if cpu.getFlag(FlagC) {
			mask = 1
		}
		result := (value << 1) | (mask)
		writeCycles = cpu.set_r8(r8, result)
		cpu.setFlag(FlagZ, result == 0)
		cpu.setFlag(FlagN, false)
		cpu.setFlag(FlagH, false)
		cpu.setFlag(FlagC, (bit7>>7) == 1)
	case 0b011:
		// RR
		bit0 := value & 1
		var mask uint8 = 0
		if cpu.getFlag(FlagC) {
			mask = 1
		}
		result := (value >> 1) | (mask << 7)
		writeCycles = cpu.set_r8(r8, result)
		cpu.setFlag(FlagZ, result == 0)
		cpu.setFlag(FlagN, false)
		cpu.setFlag(FlagH, false)
		cpu.setFlag(FlagC, bit0 == 1)
	case 0b100:
		// SLA
		bit7 := (value & 0b1000_0000)
		result := value << 1
		writeCycles = cpu.set_r8(r8, result)
		cpu.setFlag(FlagZ, result == 0)
		cpu.setFlag(FlagN, false)
		cpu.setFlag(FlagH, false)
		cpu.setFlag(FlagC, (bit7>>7) == 1)
	case 0b101:
		// SRA
		bit0 := value & 1
		bit7 := (value & 0b1000_0000)
		result := (value >> 1) | bit7
		writeCycles = cpu.set_r8(r8, result)
		cpu.setFlag(FlagZ, result == 0)
		cpu.setFlag(FlagN, false)
		cpu.setFlag(FlagH, false)
		cpu.setFlag(FlagC, bit0 == 1)
	case 0b111:
		// SRL
		bit0 := value & 1
		result := value >> 1
		writeCycles = cpu.set_r8(r8, result)
		cpu.setFlag(FlagZ, result == 0)
		cpu.setFlag(FlagN, false)
		cpu.setFlag(FlagH, false)
		cpu.setFlag(FlagC, bit0 == 1)
	case 0b110:
		// SWAP
		upperNibble := value & 0xF0
		lowerNibble := value & 0x0F
		result := (lowerNibble << 4) | (upperNibble >> 4)
		writeCycles = cpu.set_r8(r8, result)
		cpu.setFlag(FlagZ, result == 0)
		cpu.setFlag(FlagN, false)
		cpu.setFlag(FlagH, false)
		cpu.setFlag(FlagC, false)
	}

	return 8 + fetchCycles + writeCycles
}
