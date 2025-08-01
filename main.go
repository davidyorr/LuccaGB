//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/davidyorr/EchoGB/gameboy"
	"github.com/davidyorr/EchoGB/logger"
)

func main() {
	logger.Info("Hello EchoGB!")

	js.Global().Set("loadRom", js.FuncOf(loadRom))
	js.Global().Set("processEmulatorCycles", js.FuncOf(processEmulatorCycles))

	jsImageData = js.Global().Get("Uint8Array").New(len(goImageData))

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

	return nil
}

const (
	displayWidth  = 160
	displayHeight = 144
)

func processEmulatorCycles(this js.Value, args []js.Value) interface{} {
	tCyclesToRun := args[0].Float()
	var tCyclesUsed float64

	for tCyclesToRun >= 4 {
		tCycles, frameReady, _ := gb.Step()
		tCyclesUsed += float64(tCycles)
		tCyclesToRun -= float64(tCycles)

		if frameReady {
			presentFrame()
		}
	}

	return js.ValueOf(map[string]interface{}{
		"tCyclesUsed": tCyclesUsed,
	})
}

var goImageData [displayWidth * displayHeight * 4]byte
var jsImageData js.Value

func presentFrame() {
	js.Global().Get("console").Call("log", "Go: presentFrame()")
	frameBuffer := gb.FrameBuffer()
	i := 0
	for screenY := 0; screenY < displayHeight; screenY++ {
		for screenX := 0; screenX < displayWidth; screenX++ {
			color := frameBuffer[screenY][screenX]
			switch color {
			case 0:
				goImageData[i], goImageData[i+1], goImageData[i+2], goImageData[i+3] = 255, 255, 255, 255
			case 1:
				goImageData[i], goImageData[i+1], goImageData[i+2], goImageData[i+3] = 170, 170, 170, 255
			case 2:
				goImageData[i], goImageData[i+1], goImageData[i+2], goImageData[i+3] = 85, 85, 85, 255
			case 3:
				goImageData[i], goImageData[i+1], goImageData[i+2], goImageData[i+3] = 0, 0, 0, 255
			}
			i += 4
		}
	}

	js.CopyBytesToJS(jsImageData, goImageData[:])
	js.Global().Get("updateCanvas").Invoke(jsImageData)
}
