package cpu

import "fmt"

// 0x00 No OPeration
func nop(cpu *CPU) uint8 {
	fmt.Println("Go: nop()")
	return 4
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

// 0xE0
func ldh_a8_a(cpu *CPU) uint8 {
	cpu.bus.Write(0xFF00+cpu.immediateValue, cpu.a)
	return 2
}

// 0xEA Copy the value in register A into the byte at address a16
func ld_a16_a(cpu *CPU) uint8 {
	cpu.bus.Write(cpu.immediateValue, cpu.a)
	return 16
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

// 0xC3 Jump to address a16; effectively, copy a16 into PC
func jp_a16(cpu *CPU) uint8 {
	fmt.Println("Go: jp_a16()")
	fmt.Printf("  imm=[0x%04X]\n", cpu.immediateValue)
	cpu.pc = cpu.immediateValue
	return 16
}

// 0xFF Call address 0x38
func rst_38h(cpu *CPU) uint8 {
	cpu.pushToStack16(cpu.pc)
	cpu.pc = 0x38
	return 16
}

// 0x3C Increment the value in register SP by 1
func inc_sp(cpu *CPU) uint8 {
	cpu.sp++
	return 8
}

// 0xF3 Disable Interrupts by clearing the IME flag
func di(cpu *CPU) uint8 {
	cpu.ime = false
	return 4
}
