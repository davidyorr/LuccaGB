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

	for i := range 10 {
		gb.cpu.Step()
		output := gb.mmu.SerialOutputBuffer()
		t.Logf("++++++ test output=[%s]\n", output)
		i++
	}
}
