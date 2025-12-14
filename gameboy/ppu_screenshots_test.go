//go:build screenshots

package gameboy

import (
	"fmt"
	"os"
	"testing"

	"github.com/davidyorr/LuccaGB/hasher"
)

func TestBoop__solid_color_0_background(t *testing.T) {
	runPpuTest(t, "boop/solid-color-0-background", 1, "beaca6d8b5aec6a02fe4db04662136db628ec2543d396d03601853067dd47eac")
}

func TestBoop__solid_color_0_window(t *testing.T) {
	runPpuTest(t, "boop/solid-color-0-window", 2, "beaca6d8b5aec6a02fe4db04662136db628ec2543d396d03601853067dd47eac")
}

func TestBoop__solid_color_1_background(t *testing.T) {
	runPpuTest(t, "boop/solid-color-1-background", 1, "fcbaf6ec8a002c189a1fa22a6c92b537d59ecb0eb54b833d614627b232f66f75")
}

func TestBoop__solid_color_1_window(t *testing.T) {
	runPpuTest(t, "boop/solid-color-1-window", 2, "fcbaf6ec8a002c189a1fa22a6c92b537d59ecb0eb54b833d614627b232f66f75")
}

func TestBoop__solid_color_2_background(t *testing.T) {
	runPpuTest(t, "boop/solid-color-2-background", 1, "8f0d54a23211730da42c276b2a528461963b8d474c617179fe79659d3c990b38")
}

func TestBoop__solid_color_2_window(t *testing.T) {
	runPpuTest(t, "boop/solid-color-2-window", 2, "8f0d54a23211730da42c276b2a528461963b8d474c617179fe79659d3c990b38")
}

func TestBoop__solid_color_3_background(t *testing.T) {
	runPpuTest(t, "boop/solid-color-3-background", 1, "46e2096b907947368d310929303a04005b39c4a278e3a7de2225c355b4522694")
}

func TestBoop__solid_color_3_window(t *testing.T) {
	runPpuTest(t, "boop/solid-color-3-window", 2, "46e2096b907947368d310929303a04005b39c4a278e3a7de2225c355b4522694")
}

func TestBoop__sprite_8x8(t *testing.T) {
	runPpuTest(t, "boop/sprite-8x8", 1, "c6f656c7c2a60a2359837ba585ebf16e72aade3fd471dce94b15b1922e5572bc")
}

func TestBoop__sprite_8x16(t *testing.T) {
	runPpuTest(t, "boop/sprite-8x16", 1, "8b47c129eea87cde106a185aaabd88dd606f26522f3ed1b4b5eece88eee82d0d")
}

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
