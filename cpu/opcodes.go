package cpu

import "fmt"

// No OPeration
func nop(cpu *CPU) uint8 {
	fmt.Println("Go: nop()")
	return 1
}

// Jump to address a16; effectively, copy a16 into PC
func jp_a16(cpu *CPU) uint8 {
	fmt.Println("Go: jp_a16()")
	fmt.Printf("  imm=[0x%04X]\n", cpu.immediateValue)
	return 4
}
