package cpu

import "fmt"

// 0x00 No OPeration
func nop(cpu *CPU) uint8 {
	fmt.Println("Go: nop()")
	return 4
}

// 0xC3 Jump to address a16; effectively, copy a16 into PC
func jp_a16(cpu *CPU) uint8 {
	fmt.Println("Go: jp_a16()")
	fmt.Printf("  imm=[0x%04X]\n", cpu.immediateValue)
	cpu.pc = cpu.immediateValue
	return 16
}

// 0x31 Copy the value d16 into register SP
func ld_sp_d16(cpu *CPU) uint8 {
	cpu.sp = cpu.immediateValue
	return 12
}

// 0x3E
func ld_a_d8(cpu *CPU) uint8 {
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

// 0xF3 Disable Interrupts by clearing the IME flag
func di(cpu *CPU) uint8 {
	cpu.ime = false
	return 4
}
