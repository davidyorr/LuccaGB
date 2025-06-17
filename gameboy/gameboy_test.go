package gameboy

import (
	"os"
	"testing"
)

func TestBlarggCpuInsructions(t *testing.T) {
	romBytes, err := os.ReadFile("../roms/test/cpu_instrs.gb")
	if err != nil {
		t.Fatal("Error reading file:", err)
	}

	gb := New()
	gb.LoadRom(romBytes)

	for i := range 100250 {
		// this logic will go in the actual code, not in the test code
		cycles, err := gb.cpu.Step()
		gb.mmu.Step(cycles)
		if err != nil {
			t.Log("unprefixed instructions remaining:", gb.cpu.GetNumberOfUnimplementedInstructions())
			t.Fatal(err)
		}
		// output := gb.mmu.SerialOutputBuffer()
		// t.Logf("++++++ test output: [%s] (hex: [% x])\n", string(output), output)
		i++
	}
	t.Log("unprefixed instructions remaining:", gb.cpu.GetNumberOfUnimplementedInstructions())
}
