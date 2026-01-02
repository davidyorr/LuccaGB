//go:build !screenshots

package gameboy

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/davidyorr/LuccaGB/internal/logger"
)

var (
	testsRun    int32
	testsPassed int32
)

type TestType string

const (
	TestTypeBlargg  TestType = "blargg"
	TestTypeMooneye TestType = "mooneye"
)

func skipCi(t *testing.T, romName string) {
	if os.Getenv("CI") != "" {
		t.Skipf("SKIPPING TESTCASE: %s.gb", romName)
	}
}

func TestMain(m *testing.M) {
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

func TestBlargg__cpu_instrs(t *testing.T) {
	loadRomAndRunSteps(t, "blargg/cpu_instrs", 56_108_272, TestTypeBlargg)
}

func TestBlargg__instr_timing(t *testing.T) {
	loadRomAndRunSteps(t, "blargg/instr_timing", 676_091, TestTypeBlargg)
}

func TestBlargg__mem_timing(t *testing.T) {
	loadRomAndRunSteps(t, "blargg/mem_timing", 1_597_872, TestTypeBlargg)
}

func TestBlargg__halt_bug(t *testing.T) {
	loadRomAndRunSteps(t, "blargg/halt_bug", 2_000_000, TestTypeBlargg)
}

func TestMooneye__add_sp_e_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/add_sp_e_timing", 238_472, TestTypeMooneye)
}

func TestMooneye__call_cc_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/call_cc_timing", 256_184, TestTypeMooneye)
}

func TestMooneye__call_cc_timing2(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/call_cc_timing2", 256_184, TestTypeMooneye)
}

func TestMooneye__call_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/call_timing", 234_792, TestTypeMooneye)
}

func TestMooneye__call_timing2(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/call_timing2", 256_185, TestTypeMooneye)
}

// This tests DI instruction timing by setting up a vblank interrupt
// interrupt with a write to IE.
//
// This test is for DMG/MGB, so DI is expected to disable interrupts
// immediately
// On CGB/GBA DI has a delay and this test fails in round 2!!
func TestMooneye__di_timing_GS(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/di_timing-GS", 269_971, TestTypeMooneye)
}

func TestMooneye__div_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/div_timing", 185_726, TestTypeMooneye)
}

func TestMooneye__ei_sequence(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ei_sequence", 185_642, TestTypeMooneye)
}

// This tests EI instruction timing by forcing a serial
// interrupt with a write to IE/IF.
func TestMooneye__ei_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ei_timing", 185_635, TestTypeMooneye)
}

func TestMooneye__halt_ime0_ei(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/halt_ime0_ei", 182_199, TestTypeMooneye)
}

func TestMooneye__halt_ime0_nointr_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/halt_ime0_nointr_timing", 255_866, TestTypeMooneye)
}

func TestMooneye__halt_ime1_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/halt_ime1_timing", 185_568, TestTypeMooneye)
}

func TestMooneye__halt_ime1_timing2_GS(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/halt_ime1_timing2-GS", 326_246, TestTypeMooneye)
}

// This tests the behaviour of IE and IF flags by forcing a serial
// interrupt with a write to IF. The interrupt handler increments
// E, so we can track how many times the interrupt has been
// triggered
func TestMooneye__if_ie_registers(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/if_ie_registers", 185_802, TestTypeMooneye)
}

func TestMooneye__intr_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/intr_timing", 185_649, TestTypeMooneye)
}
func TestMooneye__jp_cc_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/jp_cc_timing", 234_788, TestTypeMooneye)
}

func TestMooneye__jp_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/jp_timing", 234_788, TestTypeMooneye)
}

func TestMooneye__ld_hl_sp_e_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ld_hl_sp_e_timing", 238_480, TestTypeMooneye)
}

// This tests starting another OAM DMA while one is already active
func TestMooneye__oam_dma_restart(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/oam_dma_restart", 237_556, TestTypeMooneye)
}

// This tests what happens in the first few cycles of OAM DMA.
// Also, when OAM DMA is restarted while a previous one is running, the previous one
// is not immediately stopped.
func TestMooneye__oam_dma_start(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/oam_dma_start", 258_931, TestTypeMooneye)
}

func TestMooneye__oam_dma_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/oam_dma_timing", 238_307, TestTypeMooneye)
}

func TestMooneye__pop_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/pop_timing", 185_122, TestTypeMooneye)
}

func TestMooneye__push_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/push_timing", 238_463, TestTypeMooneye)
}

func TestMooneye__rapid_di_ei(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/rapid_di_ei", 185_800, TestTypeMooneye)
}

func TestMooneye__ret_cc_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ret_cc_timing", 252_354, TestTypeMooneye)
}

func TestMooneye__ret_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ret_timing", 252_354, TestTypeMooneye)
}

func TestMooneye__reti_intr_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/reti_intr_timing", 185_724, TestTypeMooneye)
}

func TestMooneye__reti_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/reti_timing", 251_599, TestTypeMooneye)
}

func TestMooneye__rst_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/rst_timing", 238_481, TestTypeMooneye)
}

func TestMooneye__bits__mem_oam(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/bits/mem_oam", 186_664, TestTypeMooneye)
}

func TestMooneye__bits__reg_f(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/bits/reg_f", 185_645, TestTypeMooneye)
}

// This test checks all unused bits in working $FFxx IO,
// and all unused $FFxx IO. Unused bits and unused IO all return 1s.
func TestMooneye__bits__unused_hwio_GS(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/bits/unused_hwio-GS", 182_118, TestTypeMooneye)
}

// Tests the DAA instruction with all possible input combinations
func TestMooneye__instr__daa(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/instr/daa", 392_793, TestTypeMooneye)
}

// Tests what happens if the IE register is the target for one of the
// PC pushes during interrupt dispatch.
func TestMooneye__interrupts__ie_push(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/interrupts/ie_push", 181_550, TestTypeMooneye)
}

func TestMooneye__mbc1__bits_bank1(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc1/bits_bank1", 3_129_754, TestTypeMooneye)
}

func TestMooneye__mbc1__bits_bank2(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc1/bits_bank2", 3_077_085, TestTypeMooneye)
}

func TestMooneye__mbc1__bits_mode(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc1/bits_mode", 3_129_744, TestTypeMooneye)
}

func TestMooneye__mbc1__bits_ramg(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc1/bits_ramg", 6_114_264, TestTypeMooneye)
}

func TestMooneye__mbc1__multicart_rom_8Mb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc1/multicart_rom_8Mb", 5_500_000, TestTypeMooneye)
}

func TestMooneye__mbc1__ram_256kb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc1/ram_256kb", 970_358, TestTypeMooneye)
}

func TestMooneye__mbc1__ram_64kb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc1/ram_64kb", 970_368, TestTypeMooneye)
}

func TestMooneye__mbc1__rom_16Mb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc1/rom_16Mb", 233_004, TestTypeMooneye)
}

func TestMooneye__mbc1__rom_1Mb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc1/rom_1Mb", 233_004, TestTypeMooneye)
}

func TestMooneye__mbc1__rom_2Mb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc1/rom_2Mb", 233_004, TestTypeMooneye)
}

func TestMooneye__mbc1__rom_4Mb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc1/rom_4Mb", 233_004, TestTypeMooneye)
}

func TestMooneye__mbc1__rom_512kb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc1/rom_512kb", 233_004, TestTypeMooneye)
}

func TestMooneye__mbc1__rom_8Mb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc1/rom_8Mb", 233_004, TestTypeMooneye)
}

func TestMooneye__mbc2__bits_ramg(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc2/bits_ramg", 6_237_157, TestTypeMooneye)
}

func TestMooneye__mbc2__bits_romb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc2/bits_romb", 3_270_201, TestTypeMooneye)
}

func TestMooneye__mbc2__bits_unused(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc2/bits_unused", 3_059_522, TestTypeMooneye)
}

func TestMooneye__mbc2__ram(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc2/ram", 496_351, TestTypeMooneye)
}

func TestMooneye__mbc2__rom_1Mb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc2/rom_1Mb", 180_347, TestTypeMooneye)
}

func TestMooneye__mbc2__rom_2Mb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc2/rom_2Mb", 180_347, TestTypeMooneye)
}

func TestMooneye__mbc2__rom_512kb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc2/rom_512kb", 180_347, TestTypeMooneye)
}

func TestMooneye__mbc5__rom_16Mb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc5/rom_16Mb", 233_012, TestTypeMooneye)
}

func TestMooneye__mbc5__rom_1Mb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc5/rom_1Mb", 233_012, TestTypeMooneye)
}

func TestMooneye__mbc5__rom_2Mb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc5/rom_2Mb", 233_012, TestTypeMooneye)
}

func TestMooneye__mbc5__rom_32Mb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc5/rom_32Mb", 233_012, TestTypeMooneye)
}

func TestMooneye__mbc5__rom_4Mb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc5/rom_4Mb", 233_012, TestTypeMooneye)
}

func TestMooneye__mbc5__rom_512kb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc5/rom_512kb", 233_012, TestTypeMooneye)
}

func TestMooneye__mbc5__rom_64Mb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc5/rom_64Mb", 233_012, TestTypeMooneye)
}

func TestMooneye__mbc5__rom_8Mb(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/mbc5/rom_8Mb", 233_012, TestTypeMooneye)
}

// This test checks that OAM DMA copies all bytes correctly.
func TestMooneye__oam_dma__basic(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/oam_dma/basic", 186_754, TestTypeMooneye)
}

// This test checks what happens if you read the DMA register. Reads should
// always simply return the last written value, regardless of the state of the
// OAM DMA transfer or other things.
func TestMooneye__oam_dma__reg_read(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/oam_dma/reg_read", 183_038, TestTypeMooneye)
}

// This test checks that OAM DMA source memory areas work as expected,
// including the area past $DFFF.
func TestMooneye__oam_dma__sources_GS(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/oam_dma/sources-GS", 636_933, TestTypeMooneye)
}

// Tests how SCX affects the duration between STAT mode=0 interrupt and LY increment.
// No sprites or window.
//
// Expected behaviour:
//
//	(SCX mod 8) = 0   => LY increments 51 cycles after STAT interrupt
//	(SCX mod 8) = 1-4 => LY increments 50 cycles after STAT interrupt
//	(SCX mod 8) = 5-7 => LY increments 49 cycles after STAT interrupt
func TestMooneye__ppu__hblank_ly_scx_timing_GS(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ppu/hblank_ly_scx_timing-GS", 829_917, TestTypeMooneye)
}

// Tests how long does it take to get from STAT mode=1 interrupt to STAT mode=2 interrupt
// No sprites, scroll or window.
func TestMooneye__ppu__intr_1_2_timing_GS(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ppu/intr_1_2_timing-GS", 238_312, TestTypeMooneye)
}

// Tests how long does it take to get from STAT mode=2 interrupt to STAT mode=0 interrupt
// No sprites, scroll or window.
func TestMooneye__ppu__intr_2_0_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ppu/intr_2_0_timing", 300_000, TestTypeMooneye)
}

// Tests how long does it take to get from STAT=mode2 interrupt to mode0
// Includes sprites in various configurations
func TestMooneye__ppu__intr_2_mode0_timing_sprites(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ppu/intr_2_mode0_timing_sprites", 400_000, TestTypeMooneye)
}

// Tests how long does it take to get from STAT=mode2 interrupt to mode0
// No sprites, scroll, or window
func TestMooneye__ppu__intr_2_mode0_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ppu/intr_2_mode0_timing", 400_000, TestTypeMooneye)
}

// Tests how long does it take to get from STAT=mode2 interrupt to mode3
// No sprites, scroll, or window
func TestMooneye__ppu__intr_2_mode3_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ppu/intr_2_mode3_timing", 220_758, TestTypeMooneye)
}

func TestMooneye__ppu__intr_2_oam_ok_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ppu/intr_2_oam_ok_timing", 220_758, TestTypeMooneye)
}

// This tests how the internal STAT IRQ signal can block
// subsequent STAT interrupts if the signal is never cleared
func TestMooneye__ppu__stat_irq_blocking(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ppu/stat_irq_blocking", 197_935, TestTypeMooneye)
}

// If bit 5 (mode 2 OAM interrupt) is set, an interrupt is also triggered
// If bit 5 (mode 2 OAM interrupt) is set, an interrupt is also triggered
// If bit 5 (mode 2 OAM interrupt) is set, an interrupt is also triggered
// at line 144 when vblank starts.
// This test measures the cycles between vblank<->vblank and compares that to vblank<->stat_m2_144
// Expected behaviour: vblank and stat_m2_144 are triggered at the same time
func TestMooneye__ppu__vblank_stat_intr_GS(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ppu/vblank_stat_intr-GS", 324_462, TestTypeMooneye)
}

// This test verifies that the timer is affected by resetting the DIV register
// by writing to it. The timer uses the same internal counter as the DIV
// register, so resetting DIV also resets the timer.
// The basic idea of this test is very simple:
//  1. start the timer
//  2. keep resetting DIV in a loop by writing to it
//  3. run N iterations of the loop
//  4. if an interrupt happened, test failed
func TestMooneye__timer__div_write(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/timer/div_write", 901_159, TestTypeMooneye)
}

// This test rapidly starts and stops the timer.
// There are two behaviours that affect the test:
//  1. starting or stopping the timer does *not* reset its internal counter,
//     so repeated starting and stopping does not prevent timer increments
//  2. the timer circuit design causes some unexpected timer increases
func TestMooneye__timer__rapid_toggle(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/timer/rapid_toggle", 184_894, TestTypeMooneye)
}

func TestMooneye__timer__tim00_div_trigger(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/timer/tim00_div_trigger", 185_636, TestTypeMooneye)
}

func TestMooneye__timer__tim00(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/timer/tim00", 185_642, TestTypeMooneye)
}

func TestMooneye__timer__tim01_div_trigger(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/timer/tim01_div_trigger", 185_641, TestTypeMooneye)
}

func TestMooneye__timer__tim01(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/timer/tim01", 185_643, TestTypeMooneye)
}

func TestMooneye__timer__tim10_div_trigger(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/timer/tim10_div_trigger", 185_642, TestTypeMooneye)
}

func TestMooneye__timer__tim10(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/timer/tim10", 185_638, TestTypeMooneye)
}

func TestMooneye__timer__tim11_div_trigger(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/timer/tim11_div_trigger", 185_638, TestTypeMooneye)
}

func TestMooneye__timer__tim11(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/timer/tim11", 185_635, TestTypeMooneye)
}

// The test checks what values appear in the TIMA register when the
// timer overflows.
//
// Apparently the TIMA register contains 00 for 4 cycles before being
// reloaded with the value from the TMA register. The TIMA increments
// do still happen every 64 cycles, there is no additional 4 cycle
// delay.
func TestMooneye__timer__tima_reload(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/timer/tima_reload", 185_205, TestTypeMooneye)
}

// This test tests which write to the TIMA register is ignored when
// the timer is reloading.
func TestMooneye__timer__tima_write_reloading(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/timer/tima_write_reloading", 200_000, TestTypeMooneye)
}

// This test checks when writes to the TMA register get picked while
// the timer is reloading.
func TestMooneye__timer__tma_write_reloading(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/timer/tma_write_reloading", 200_000, TestTypeMooneye)
}

func loadRomAndRunSteps(t *testing.T, romName string, stepCount int, testType TestType) {
	atomic.AddInt32(&testsRun, 1)
	t.Logf("TESTCASE: %s.gb", romName)
	logBuffer, testLogger := initTestLogger()
	logger.Init(testLogger.Handler())
	defer logger.Init(slog.Default().Handler())

	romBytes, err := os.ReadFile(fmt.Sprintf("../../roms/test/%s.gb", romName))
	if err != nil {
		t.Fatal("Error reading file:", err)
	}

	gb := New()
	gb.LoadRom(romBytes)

	// track if the test emitted any pass/fail signal
	testComplete := false

	for i := range stepCount {
		gb.Step()
		output := string(gb.serial.SerialOutputBuffer())

		passed, failed := checkResult(output, testType)
		if failed {
			t.Logf("\n============ SERIAL OUTPUT ============\n%s\n=======================================\n", output)
			t.Logf("❌ TEST FAILED after %d steps", i+1)
			for _, line := range logBuffer.LastN(40) {
				fmt.Fprint(os.Stdout, line)
			}
			writeLogToFile(logBuffer)
			t.Fatal()
		}
		if passed {
			testComplete = true
			atomic.AddInt32(&testsPassed, 1)
			t.Logf("✅ TEST PASSED after %d steps", i+1)
			writeLogToFile(logBuffer)
			break
		}
	}

	if !testComplete {
		output := string(gb.serial.SerialOutputBuffer())
		t.Logf("\n============ SERIAL OUTPUT ============\n%s\n=======================================\n", output)
		for _, line := range logBuffer.LastN(40) {
			fmt.Fprint(os.Stdout, line)
		}
		writeLogToFile(logBuffer)
		t.Logf("❌ TIMED OUT after %d steps", stepCount)
		t.Fatal("Test timed out (no pass/fail signal detected)")
	}

	output := string(gb.serial.SerialOutputBuffer())
	t.Logf("\n============ SERIAL OUTPUT ============\n%s\n=======================================\n", output)
}

func checkResult(output string, testType TestType) (passed bool, failed bool) {
	switch testType {
	case TestTypeBlargg:
		passed = strings.Contains(output, "Passed\n") || strings.Contains(output, "Passed all tests\n")
		failed = strings.Contains(output, "Failed")
	case TestTypeMooneye:
		// check for Fibonacci sequence (3,5,8,13,21,34) in hex
		passed = strings.Contains(output, "\x03\x05\x08\x0D\x15\x22")
		failed = strings.Contains(output, "\x42\x42\x42\x42\x42\x42")
	}
	return
}

func writeLogToFile(logBuffer *logBuffer) {
	if os.Getenv("CI") == "" {
		if err := os.WriteFile("../testoutput.log", []byte(strings.Join(logBuffer.messages, "")), 0644); err != nil {
			fmt.Printf("Failed to write test output: %v\n", err)
		}
	}
}

type logBuffer struct {
	messages []string
	mu       sync.Mutex
}

func (lb *logBuffer) Write(p []byte) (n int, err error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.messages = append(lb.messages, string(p))
	return len(p), nil
}

func (lb *logBuffer) LastN(n int) []string {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	if len(lb.messages) <= n {
		return lb.messages
	}
	return lb.messages[len(lb.messages)-n:]
}

func initTestLogger() (*logBuffer, *slog.Logger) {
	// Redirect logs to a buffer
	logBuffer := &logBuffer{
		messages: make([]string, 0),
	}
	testLogger := slog.New(slog.NewTextHandler(logBuffer, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.String(a.Key, a.Value.Time().Format("15:04:05.000"))
			}
			return a
		},
	}))

	return logBuffer, testLogger
}
