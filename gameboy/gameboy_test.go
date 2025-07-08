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

func TestBlarggCpuInsructions(t *testing.T) {
	loadRomAndRunSteps(t, "blargg/cpu_instrs", 25_000_000)
}

func TestBlarggInstructionTiming(t *testing.T) {
	loadRomAndRunSteps(t, "blargg/instr_timing", 1_000_000)
}

func TestBlarggMemoryTiming(t *testing.T) {
	loadRomAndRunSteps(t, "blargg/mem_timing", 2_000_000)
}

func TestBlarggHaltBug(t *testing.T) {
	loadRomAndRunSteps(t, "blargg/halt_bug", 4_000_000)
}

func loadRomAndRunSteps(t *testing.T, romName string, stepCount int) {
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

		passed, failed := checkResult(output)
		if failed {
			t.Logf("❌ TEST FAILED after %d steps", i+1)
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
		t.Logf("❌ TIMED OUT after %d steps", stepCount)
		t.Fatal("Test timed out (no pass/fail signal detected)")
	}

	output := string(gb.serial.SerialOutputBuffer())
	t.Logf("\n============ SERIAL OUTPUT ============\n%s\n=======================================\n", output)
}

func checkResult(output string) (passed bool, failed bool) {
	passed = strings.Contains(output, "Passed\n") || strings.Contains(output, "Passed all tests\n")
	failed = strings.Contains(output, "Failed")
	return
}
