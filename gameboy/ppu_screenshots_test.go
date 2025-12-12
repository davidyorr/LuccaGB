//go:build screenshots

package gameboy

import (
	"fmt"
	"os"
	"testing"

	"github.com/davidyorr/LuccaGB/hasher"
)

func Test__lucca(t *testing.T) {
	runPpuTest(t, "../lucca", 2, "b77a59fe8c635f5db714d0b5eea19b23cfab3fbe7001a541c8056bbc6834a3e5")
}

func runPpuTest(t *testing.T, romName string, framesToRun int, expectedHash string) {
	t.Helper()
	t.Logf("TESTCASE: %s.gb", romName)

	// 1. Setup
	romBytes, err := os.ReadFile(fmt.Sprintf("../roms/test/%s.gb", romName))
	if err != nil {
		t.Fatalf("❌ SETUP FAIL: %v", err)
	}

	gb := New()
	gb.LoadRom(romBytes)

	// 2. Run
	for i := 0; i < framesToRun; i++ {
		for {
			_, ready, _ := gb.Step()
			if ready {
				break
			}
		}
	}

	// 3. Check
	actualHash := hasher.HashFrameBuffer(gb.FrameBuffer())

	if actualHash == expectedHash {
		t.Logf("✅ PASS: %s", fmt.Sprintf("%s.gb", romName))
	} else {
		t.Errorf("❌ FAIL: %s", fmt.Sprintf("%s.gb", romName))
	}
}
