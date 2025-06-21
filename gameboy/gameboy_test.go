package gameboy

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/davidyorr/EchoGB/logger"
)

func TestBlarggCpuInsructions(t *testing.T) {
	// 1. Redirect logs to a buffer
	var buf bytes.Buffer
	var logOutput io.Writer = &buf
	if testing.Verbose() {
		logOutput = os.Stdout
	} else {
		logOutput = &bytes.Buffer{}
	}
	testLogger := slog.New(slog.NewTextHandler(logOutput, &slog.HandlerOptions{
		Level: slog.LevelDebug,
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
			output := gb.mmu.SerialOutputBuffer()
			t.Log("unprefixed instructions remaining:", gb.cpu.GetNumberOfUnimplementedInstructions())
			t.Logf("++++++ test output: [%s] (hex: [% x])\n", string(output), output)
			t.Fatal(err)
		}
	}
	output := gb.mmu.SerialOutputBuffer()
	t.Log("unprefixed instructions remaining:", gb.cpu.GetNumberOfUnimplementedInstructions())
	t.Logf("++++++ test output: [%s] (hex: [% x])\n", string(output), output)
}
