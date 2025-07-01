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
	logBuffer, testLogger := initTestLogger()
	logger.Init(testLogger.Handler())
	defer logger.Init(slog.Default().Handler())

	romBytes, err := os.ReadFile("../roms/test/cpu_instrs.gb")
	if err != nil {
		t.Fatal("Error reading file:", err)
	}

	gb := New()
	gb.LoadRom(romBytes)

	for range 20_000_000 {
		_, err := gb.Step()
		if err != nil {
			output := string(gb.serial.SerialOutputBuffer())
			t.Logf("\n============\n%s\n============\n", output)
			if testing.Verbose() {
				for _, line := range logBuffer.LastN(40) {
					fmt.Fprint(os.Stdout, line)
				}
			}
			t.Fatal(err)
		}
	}
	output := string(gb.serial.SerialOutputBuffer())
	t.Logf("\n============\n%s\n============\n", output)

	if testing.Verbose() {
		for _, line := range logBuffer.LastN(40) {
			fmt.Fprint(os.Stdout, line)
		}
	}
}

func TestBlarggInstructionTiming(t *testing.T) {
	loadRomAndRunSteps(t, "instr_timing", 1_000_000)
}

func TestBlarggMemoryTiming(t *testing.T) {
	loadRomAndRunSteps(t, "mem_timing", 5_000_000)
}

func TestBlarggHaltBug(t *testing.T) {
	loadRomAndRunSteps(t, "halt_bug", 40_000_000)
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

	for range stepCount {
		_, err := gb.Step()
		output := string(gb.serial.SerialOutputBuffer())
		if err != nil || strings.Contains(output, "Failed") {
			t.Logf("\n============\n%s\n============\n", output)
			for _, line := range logBuffer.LastN(10) {
				fmt.Fprint(os.Stdout, line)
			}
			t.Fatal(err)
		}
		if strings.Contains(output, "Passed") {
			break
		}
	}
	output := string(gb.serial.SerialOutputBuffer())
	t.Logf("\n============\n%s\n============\n", output)

	if testing.Verbose() {
		for _, line := range logBuffer.LastN(10) {
			fmt.Fprint(os.Stdout, line)
		}
	}
}
