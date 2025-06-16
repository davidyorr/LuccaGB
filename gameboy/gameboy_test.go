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

	for i := range 100 {
		_, err := gb.cpu.Step()
		if err != nil {
			t.Fatal(err)
		}
		output := gb.mmu.SerialOutputBuffer()
		t.Logf("++++++ test output: [%s] (hex: [% x])\n", string(output), output)
		i++
	}
}
