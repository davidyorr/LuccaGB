package cpu

import "fmt"

// 0x00 No OPeration
func nop(cpu *CPU) uint8 {
	fmt.Println("Go: nop()")
	return 4
}

// 0x01 Copy the value n16 into register r16
func ld_bc_n16(cpu *CPU) uint8 {
	cpu.setBC(cpu.immediateValue)
	return 12
}

// 0x21 Copy the value n16 into register HL.
func ld_hl_n16(cpu *CPU) uint8 {
	cpu.setHL(cpu.immediateValue)
	return 12
}

// 0x31 Copy the value n16 into register SP
func ld_sp_n16(cpu *CPU) uint8 {
	cpu.sp = cpu.immediateValue
	return 12
}

// 0x3E
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

// 0xE0
func ldh_a8_a(cpu *CPU) uint8 {
	cpu.bus.Write(0xFF00+cpu.immediateValue, cpu.a)
	return 12
}

// 0xEA Copy the value in register A into the byte at address a16
func ld_a16_a(cpu *CPU) uint8 {
	cpu.bus.Write(cpu.immediateValue, cpu.a)
	return 16
}

// 0x2A Copy the byte pointed to by HL into register A, and increment HL afterwards
func ld_a_hli(cpu *CPU) uint8 {
	value := cpu.getHL()
	cpu.a = cpu.bus.Read(value)
	cpu.setHL(value + 1)
	return 8
}

// 0x03 Increment the value in register r8 by 1
func inc_bc(cpu *CPU) uint8 {
	value := cpu.getBC()
	cpu.setBC(value + 1)
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

// 0x18 Relative Jump to address e8
func jr_e8(cpu *CPU) uint8 {
	operandLength := uint16(2)
	signedOffset := int8(cpu.immediateValue)
	destinationAddress := uint16(int(cpu.pc+operandLength) + int(signedOffset))
	cpu.pc = destinationAddress
	return 12
}

// 0xC3 Jump to address a16; effectively, copy a16 into PC
func jp_a16(cpu *CPU) uint8 {
	fmt.Println("Go: jp_a16()")
	fmt.Printf("  imm=[0x%04X]\n", cpu.immediateValue)
	cpu.pc = cpu.immediateValue
	return 16
}

// 0xCD Call address n16
func call_a16(cpu *CPU) uint8 {
	cpu.pushToStack16(cpu.pc + 3)
	cpu.pc = cpu.immediateValue
	return 24
}

// 0x30 Relative Jump to address e8 if condition nc is met
func jr_nc_e8(cpu *CPU) uint8 {
	if cpu.getFlag(FlagC) {
		return 8
	}

	operandLength := uint16(2)
	signedOffset := int8(cpu.immediateValue)
	destinationAddress := uint16(int(cpu.pc+operandLength) + int(signedOffset))
	cpu.pc = destinationAddress

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

// 0x3C Increment the value in register r8 by 1
func inc_a(cpu *CPU) uint8 {
	originalA := cpu.a
	cpu.a++

	cpu.setFlag(FlagZ, cpu.a == 0)
	cpu.setFlag(FlagN, false)
	// if overflow from bit 3
	cpu.setFlag(FlagH, (originalA&0x0F)+1 > 0x0F)

	return 4
}

// 0xF1 Pop register AF from the stack
func pop_af(cpu *CPU) uint8 {
	value := cpu.popFromStack16()
	cpu.setAF(value)
	return 12
}

// 0xE1 Pop register r16 from the stack
func pop_hl(cpu *CPU) uint8 {
	value := cpu.popFromStack16()
	cpu.setHL(value)
	return 12
}

// 0xC5 Push register r16 into the stack
func push_bc(cpu *CPU) uint8 {
	cpu.pushToStack16(cpu.getBC())
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
