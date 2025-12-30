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
	"sync/atomic"
	"testing"

	"github.com/davidyorr/LuccaGB/tools"
)

var (
	testsRun    int32
	testsPassed int32
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
	code := m.Run()

	total := atomic.LoadInt32(&testsRun)
	passed := atomic.LoadInt32(&testsPassed)

	fmt.Printf(
		"\n==================== TEST SUMMARY ===================\n"+
			"  %d / %d tests passing\n"+
			"=====================================================\n",
		passed,
		total,
	)

	os.Exit(code)
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

func TestMealybug_tearoom__m2_win_en_toggle(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m2_win_en_toggle", 100, "0f10f209a9435e1168b519dd433b09394bedf4c9fddb3b98656d608cab045949")
}

func TestMealybug_tearoom__m3_bgp_change(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_bgp_change", 100, "ebba92b2babcdc02536c05806efbb492addc4f948d0aacf87fd199752cceeb2c")
}

func TestMealybug_tearoom__m3_bgp_change_sprites(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_bgp_change_sprites", 100, "3da61d2ebb0a19610fd7637838db42e843b4214ae0ed3f40fce080205dbf08b4")
}

func TestMealybug_tearoom__m3_lcdc_bg_en_change(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_lcdc_bg_en_change", 100, "89deadb700b8a4c4296822c9f9ff352fa074dddabb865699f858d3e9caa07ad1")
}

func TestMealybug_tearoom__m3_lcdc_bg_map_change(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_lcdc_bg_map_change", 100, "a5ffbe52e64ca1e5d92bcbe5ba7cc7d945eb87ec53e7c11ec72ba15fa64b8fc5")
}

func TestMealybug_tearoom__m3_lcdc_obj_en_change(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_lcdc_obj_en_change", 100, "45cefce7e3f7f311105eabf24f6dcb8c12e5b5d06d0da7c60185fe350ca472cb")
}

func TestMealybug_tearoom__m3_lcdc_obj_en_change_variant(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_lcdc_obj_en_change_variant", 100, "32ad4bf8f14579a63dc4fa23bb96a97f929708e8627d30ff28677034a4a29858")
}

func TestMealybug_tearoom__m3_lcdc_obj_size_change(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_lcdc_obj_size_change", 100, "b390224ed7efda72f9dbacf84bbecd18fbae42314c8988f4c0dfa097fad0b28d")
}

func TestMealybug_tearoom__m3_lcdc_obj_size_change_scx(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_lcdc_obj_size_change_scx", 100, "198d1c79aff25c46262fc0f2b6a800ffe9112014498d27417302a6e06ce916d8")
}

func TestMealybug_tearoom__m3_lcdc_tile_sel_change(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_lcdc_tile_sel_change", 100, "d1536b0be6acb0c8b08586050a79693e23e450f2b17072f47139f0776585d286")
}

func TestMealybug_tearoom__m3_lcdc_tile_sel_win_change(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_lcdc_tile_sel_win_change", 100, "a950c37cdcc203bb0b8403a053581d8c52fa3a8bcca5ab0973eebd6cb502e7de")
}

func TestMealybug_tearoom__m3_lcdc_win_en_change_multiple(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_lcdc_win_en_change_multiple", 100, "115c91d1dd6cf5b2a2c137940d2fb1d4fd36ae6033e1a281de3150559fda486c")
}

func TestMealybug_tearoom__m3_lcdc_win_en_change_multiple_wx(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_lcdc_win_en_change_multiple_wx", 100, "60ce199f7ae3915f8ccf274703960ee8616868d68916e4eea97e239e2e4add74")
}

func TestMealybug_tearoom__m3_lcdc_win_map_change(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_lcdc_win_map_change", 100, "5eb6270ceab73e59e17275c629dd46273dbc2d32d5b90d6dfecb3d51a445728b")
}

func TestMealybug_tearoom__m3_obp0_change(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_obp0_change", 100, "93f44509f7cb7abab96d895c59e828c9302be79cb83c06312e98fcca4ff6e6b5")
}

func TestMealybug_tearoom__m3_scx_high_5_bits(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_scx_high_5_bits", 100, "6b0c2f51914064fa6f682a352fe205d42a10eabd971df7ef4c116d25b30c6528")
}

func TestMealybug_tearoom__m3_scx_low_3_bits(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_scx_low_3_bits", 100, "c30f8a5cd8926d426f4da3f196aada599c8fd0aecf731de1437be357abc25a0f")
}

func TestMealybug_tearoom__m3_scy_change(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_scy_change", 100, "b2daaac1a178651c295d1bee0fe3b42bc0ab62fbb2f3b5aecd2c53e927095da0")
}

func TestMealybug_tearoom__m3_window_timing(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_window_timing", 100, "d0d47a58f4977f7b107820c994241a49dfe384e261aa4a84231a01c93397f81b")
}

func TestMealybug_tearoom__m3_window_timing_wx_0(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_window_timing_wx_0", 100, "2896371077fda2b52ced293b26330d891afcd91235696347c303542d21d46ea7")
}

func TestMealybug_tearoom__m3_wx_4_change(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_wx_4_change", 100, "f9d77c00a28f85eddce4e253c30d08afc8b90758f17889249c138b01af01ec9d")
}

func TestMealybug_tearoom__m3_wx_4_change_sprites(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_wx_4_change_sprites", 100, "468410845f9e85400e7f9fda831157f8ebcf9b0ad76a477e0281db07c45d17b1")
}

func TestMealybug_tearoom__m3_wx_5_change(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_wx_5_change", 100, "fb6364ee948631dc2d731ddec053cd526704307aae9bae25e5e1aa17a7239004")
}

func TestMealybug_tearoom__m3_wx_6_change(t *testing.T) {
	// adjust the frameToRun after this tests passes
	runPpuTest(t, "mealybug_tearoom/m3_wx_6_change", 100, "e35b12a8ca4c5e216bd39c652e7c44a2dcf9c92b5b67764553262aeed698b15b")
}

func Test__lucca(t *testing.T) {
	runPpuTest(t, "../lucca", 2, "b77a59fe8c635f5db714d0b5eea19b23cfab3fbe7001a541c8056bbc6834a3e5")
}

func runPpuTest(t *testing.T, romName string, framesToRun int, expectedHash string) {
	atomic.AddInt32(&testsRun, 1)
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
	actualHash := tools.HashFrameBuffer(gb.FrameBuffer())

	if actualHash == expectedHash {
		atomic.AddInt32(&testsPassed, 1)
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
