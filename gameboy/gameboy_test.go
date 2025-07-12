package gameboy

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/davidyorr/EchoGB/logger"
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

func TestBlargg__cpu_instrs(t *testing.T) {
	skipCi(t, "blargg/cpu_instrs")
	loadRomAndRunSteps(t, "blargg/cpu_instrs", 25_000_000, TestTypeBlargg)
}

func TestBlargg__instr_timing(t *testing.T) {
	loadRomAndRunSteps(t, "blargg/instr_timing", 1_000_000, TestTypeBlargg)
}

func TestBlargg__mem_timing(t *testing.T) {
	loadRomAndRunSteps(t, "blargg/mem_timing", 2_000_000, TestTypeBlargg)
}

func TestBlargg__halt_bug(t *testing.T) {
	loadRomAndRunSteps(t, "blargg/halt_bug", 2_000_000, TestTypeBlargg)
}

func TestMooneye__add_sp_e_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/add_sp_e_timing", 140_335, TestTypeMooneye)
}

func TestMooneye__di_timing_GS(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/di_timing-GS", 269_971, TestTypeMooneye)
}

func TestMooneye__div_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/div_timing", 185_726, TestTypeMooneye)
}

func TestMooneye__ei_sequence(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ei_sequence", 185_642, TestTypeMooneye)
}

func TestMooneye__ei_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ei_timing", 185_635, TestTypeMooneye)
}

func TestMooneye__instr__daa(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/instr/daa", 392_793, TestTypeMooneye)
}

func TestMooneye__ppu__hblank_ly_scx_timing_GS(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ppu/hblank_ly_scx_timing-GS", 300_000, TestTypeMooneye)
}

func TestMooneye__ppu__intr_1_2_timing_GS(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ppu/intr_1_2_timing-GS", 238_312, TestTypeMooneye)
}

func TestMooneye__ppu__intr_2_0_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ppu/intr_2_0_timing", 300_000, TestTypeMooneye)
}
func TestMooneye__ppu__intr_2_mode0_timing_sprites(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ppu/intr_2_mode0_timing_sprites", 400_000, TestTypeMooneye)
}

func TestMooneye__ppu__intr_2_mode0_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ppu/intr_2_mode0_timing", 400_000, TestTypeMooneye)
}

func TestMooneye__ppu__intr_2_mode3_timing(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ppu/intr_2_mode3_timing", 220_758, TestTypeMooneye)
}

func TestMooneye__ppu__vblank_stat_intr_GS(t *testing.T) {
	loadRomAndRunSteps(t, "mooneye/ppu/vblank_stat_intr-GS", 400_000, TestTypeMooneye)
}

func loadRomAndRunSteps(t *testing.T, romName string, stepCount int, testType TestType) {
	t.Logf("TESTCASE: %s.gb", romName)
	logBuffer, testLogger := initTestLogger()
	logger.Init(testLogger.Handler())
	defer logger.Init(slog.Default().Handler())

	romBytes, err := os.ReadFile(fmt.Sprintf("../roms/test/%s.gb", romName))
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
			t.Logf("❌ TEST FAILED after %d steps", i+1)
			for _, line := range logBuffer.LastN(40) {
				fmt.Fprint(os.Stdout, line)
			}
			t.Fatal()
		}
		if passed {
			testComplete = true
			t.Logf("✅ TEST PASSED after %d steps", i+1)
			break
		}
	}

	if !testComplete {
		output := string(gb.serial.SerialOutputBuffer())
		t.Logf("\n============ SERIAL OUTPUT ============\n%s\n=======================================\n", output)
		for _, line := range logBuffer.LastN(40) {
			fmt.Fprint(os.Stdout, line)
		}
		if os.Getenv("CI") != "" {
			if err := os.WriteFile("../testoutput.log", []byte(strings.Join(logBuffer.messages, "")), 0644); err != nil {
				fmt.Printf("Failed to write test output: %v\n", err)
			}
		}
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
