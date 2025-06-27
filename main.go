//go:build js && wasm

package main

import (
	"syscall/js"
	"time"

	"github.com/davidyorr/EchoGB/gameboy"
	"github.com/davidyorr/EchoGB/logger"
)

// 0xFFFF
var interruptRegister uint8

func main() {
	logger.Info("Hello EchoGB!")

	js.Global().Set("loadRom", js.FuncOf(loadRom))

	<-make(chan struct{})
}

var gb *gameboy.Gameboy

func loadRom(this js.Value, args []js.Value) interface{} {
	jsRomData := args[0]
	cartridgeRom := make([]byte, jsRomData.Get("length").Int())
	js.CopyBytesToGo(cartridgeRom, jsRomData)

	gb = gameboy.New()
	gb.LoadRom(cartridgeRom)

	js.Global().Get("onRomLoaded").Invoke()

	lastFrameTime = time.Now()
	cycleAccumulator = 0

	return nil
}

const (
	// 4,194,304 T-cycles per second
	systemClockFrequency = 4.194304 // MHz
	framesPerSecond      = 59.73    // KHz
	cyclesPerFrame       = systemClockFrequency * 1_000_000 / framesPerSecond
)

var lastFrameTime time.Time
var cycleAccumulator float64

//go:wasmexport processEmulatorStep
func processEmulatorStep() {
	now := time.Now()
	delta := now.Sub(lastFrameTime)
	lastFrameTime = time.Now()
	cyclesToAdd := systemClockFrequency * 1_000_000 * delta.Seconds()
	cycleAccumulator += cyclesToAdd

	for cycleAccumulator >= 4 {
		tCycles, _ := gb.Step()
		cycleAccumulator -= float64(tCycles)
	}
}
