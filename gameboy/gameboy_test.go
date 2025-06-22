package gameboy

import (
	"fmt"
	"log/slog"
	"os"
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

func TestBlarggCpuInsructions(t *testing.T) {
	// 1. Redirect logs to a buffer
	var logBuffer logBuffer
	testLogger := slog.New(slog.NewTextHandler(&logBuffer, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.String(a.Key, a.Value.Time().Format("15:04:05.000"))
			}
			return a
		},
	}))

	logger.Init(testLogger.Handler())
	defer logger.Init(slog.Default().Handler())

	romBytes, err := os.ReadFile("../roms/test/cpu_instrs.gb")
	if err != nil {
		t.Fatal("Error reading file:", err)
	}

	gb := New()
	gb.LoadRom(romBytes)

	for range 10_000_000 {
		err := gb.Step()
		if err != nil {
			output := gb.serial.SerialOutputBuffer()
			t.Log("unprefixed instructions remaining:", gb.cpu.GetNumberOfUnimplementedInstructions())
			t.Logf("++++++ test output: [%s]\n", string(output))
			if testing.Verbose() {
				for _, line := range logBuffer.LastN(20) {
					fmt.Fprint(os.Stdout, line)
				}
			}
			t.Fatal(err)
		}
	}
	output := gb.serial.SerialOutputBuffer()
	t.Log("unprefixed instructions remaining:", gb.cpu.GetNumberOfUnimplementedInstructions())
	t.Logf("++++++ test output: [%s]\n", string(output))

	if testing.Verbose() {
		for _, line := range logBuffer.LastN(20) {
			fmt.Fprint(os.Stdout, line)
		}
	}
}
