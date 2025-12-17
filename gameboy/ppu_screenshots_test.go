//go:build screenshots

package gameboy

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/davidyorr/LuccaGB/hasher"
)

const screenshotOutDir = "../screenshots_out"

var ppuPalette = [4]color.RGBA{
	{0xFF, 0xFF, 0xFF, 0xFF}, // 0
	{0xAA, 0xAA, 0xAA, 0xFF}, // 1
	{0x55, 0x55, 0x55, 0xFF}, // 2
	{0x00, 0x00, 0x00, 0xFF}, // 3
}

func TestMain(m *testing.M) {
	// Clean the entire screenshots output directory before running tests
	if err := os.RemoveAll(screenshotOutDir); err != nil && !os.IsNotExist(err) {
		fmt.Printf("Warning: failed to clean screenshots output directory: %v\n", err)
	}

	// Run tests
	os.Exit(m.Run())
}

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

func TestDmg_acid2__dmg_acid2(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "dmg_acid2/dmg-acid2", 100, "cb26ef6174cd5b61662c2e74f2bdb2b0f8da0750234c75144c3a4d258fda8347")
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
		return
	}

	t.Errorf("❌ FAIL: %s", fmt.Sprintf("%s.gb", romName))

	outputName := strings.TrimPrefix(romName, "../")
	expectedPath := filepath.Join("../testdata", "screenshots", outputName+".png")
	actualPath := filepath.Join(screenshotOutDir, outputName, "actual.png")
	baselinePath := filepath.Join(screenshotOutDir, outputName, "expected.png")
	diffPath := filepath.Join(screenshotOutDir, outputName, "diff.png")

	// Write actual
	if err := WriteFrameBufferPng(actualPath, gb.FrameBuffer()); err != nil {
		t.Logf("Warning: failed to write actual screenshot: %v", err)
	}

	// Copy expected baseline
	if data, err := os.ReadFile(expectedPath); err == nil {
		_ = os.MkdirAll(filepath.Dir(baselinePath), 0o755)
		_ = os.WriteFile(baselinePath, data, 0o644)
	} else {
		t.Logf("Warning: Could not read expected baseline: %v", err)
	}

	// Load expected
	f, err := os.Open(expectedPath)
	if err != nil {
		t.Logf("Warning: missing baseline screenshot: %v", err)
		return
	}
	defer f.Close()

	expectedImg, _, err := image.Decode(f)
	if err != nil {
		t.Logf("Warning: failed to decode expected image: %v", err)
		return
	}
	expectedFrameBuffer, err := ImageToFrameBuffer(expectedImg)
	if err != nil {
		t.Logf("Warning: invalid expected image: %v", err)
		return
	}

	// Generate diff
	diffImg, changed, err := DiffFrameBuffers(expectedFrameBuffer, gb.FrameBuffer())
	if err != nil {
		t.Logf("Warning: failed to diff frame buffers: %v", err)
		return
	}
	percent := float64(changed) * 100 / float64(160*144)
	t.Logf("   %d pixels differ (~%.1f%%)", changed, percent)
	diffFile, err := os.Create(diffPath)
	if err != nil {
		t.Logf("Warning: Could not create diff image: %v", err)
		return
	}
	defer diffFile.Close()
	err = png.Encode(diffFile, diffImg)
	if err != nil {
		t.Logf("Warning: Could not encode diff image: %v", err)
		return
	}

	t.Logf("Expected: %s", expectedPath)
	t.Logf("Actual:   %s", actualPath)
	t.Logf("Diff:     %s", diffPath)
}

func WriteFrameBufferPng(path string, buffer [144][160]uint8) error {
	img := image.NewRGBA(image.Rect(0, 0, 160, 144))

	for y := 0; y < 144; y++ {
		for x := 0; x < 160; x++ {
			val := buffer[y][x]
			if val > 3 {
				val = 3
			}
			img.SetRGBA(x, y, ppuPalette[val])
		}
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func DiffFrameBuffers(expected, actual [144][160]uint8) (*image.RGBA, int, error) {
	diff := image.NewRGBA(image.Rect(0, 0, 160, 144))
	changed := 0

	for y := 0; y < 144; y++ {
		for x := 0; x < 160; x++ {
			e := expected[y][x]
			a := actual[y][x]

			if e > 3 || a > 3 {
				return diff, changed, fmt.Errorf("invalid palette index at (%d,%d): e=%d, a=%d", x, y, e, a)
			}

			delta := int(e) - int(a)
			switch {
			case delta == 0:
				diff.SetRGBA(x, y, ppuPalette[e])
			case delta > 0:
				// Expected was brighter - use blue
				intensity := uint8((delta * 85) % 256) // Scale to color
				diff.SetRGBA(x, y, color.RGBA{0x00, 0x00, 0x00 + intensity, 0xFF})
				changed++
			case delta < 0:
				// Actual was brighter - use red
				intensity := uint8((-delta * 85) % 256)
				diff.SetRGBA(x, y, color.RGBA{0x00 + intensity, 0x00, 0x00, 0xFF})
				changed++
			}
		}
	}

	return diff, changed, nil
}

func ImageToFrameBuffer(img image.Image) ([144][160]uint8, error) {
	var fb [144][160]uint8

	for y := 0; y < 144; y++ {
		for x := 0; x < 160; x++ {
			rgba := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)

			switch rgba {
			case color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}:
				fb[y][x] = 0
			case color.RGBA{0xAA, 0xAA, 0xAA, 0xFF}:
				fb[y][x] = 1
			case color.RGBA{0x55, 0x55, 0x55, 0xFF}:
				fb[y][x] = 2
			case color.RGBA{0x00, 0x00, 0x00, 0xFF}:
				fb[y][x] = 3
			default:
				return fb, fmt.Errorf("invalid pixel at (%d,%d): %v", x, y, rgba)
			}
		}
	}

	return fb, nil
}
