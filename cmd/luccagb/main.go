//go:build js && wasm

package main

import (
	"syscall/js"
	"unsafe"

	"github.com/davidyorr/LuccaGB/internal/gameboy"
	"github.com/davidyorr/LuccaGB/internal/joypad"
	"github.com/davidyorr/LuccaGB/internal/logger"
)

func main() {
	logger.Info("Hello LuccaGB!")

	js.Global().Set("loadRom", js.FuncOf(loadRom))
	js.Global().Set("getCartridgeRam", js.FuncOf(getCartridgeRam))
	js.Global().Set("setCartridgeRam", js.FuncOf(setCartridgeRam))
	js.Global().Set("processEmulatorCycles", js.FuncOf(processEmulatorCycles))
	js.Global().Set("pollFrame", js.FuncOf(pollFrame))
	js.Global().Set("pollAudioBuffer", js.FuncOf(pollAudioBuffer))
	js.Global().Set("handleJoypadButtonPressed", js.FuncOf(handleJoypadButtonPressed))
	js.Global().Set("handleJoypadButtonReleased", js.FuncOf(handleJoypadButtonReleased))
	js.Global().Set("enableTraceLogging", js.FuncOf(enableTraceLogging))
	js.Global().Set("disableTraceLogging", js.FuncOf(disableTraceLogging))
	js.Global().Set("getTraceLogs", js.FuncOf(getTraceLogs))
	js.Global().Set("getDebugInfo", js.FuncOf(getDebugInfo))

	jsImageData = js.Global().Get("Uint8Array").New(len(goImageData))

	<-make(chan struct{})
}

var gb *gameboy.Gameboy

func loadRom(this js.Value, args []js.Value) interface{} {
	jsRomData := args[0]
	cartridgeRom := make([]byte, jsRomData.Get("length").Int())
	js.CopyBytesToGo(cartridgeRom, jsRomData)

	gb = gameboy.New()
	cartridgeInfo := gb.LoadRom(cartridgeRom)

	return map[string]interface{}{
		"title":      cartridgeInfo.Title,
		"ramSize":    cartridgeInfo.RamSize,
		"hasBattery": cartridgeInfo.HasBattery,
	}
}

func getCartridgeRam(this js.Value, args []js.Value) interface{} {
	cartridgeRam := gb.CartridgeRam()
	jsBuffer := js.Global().Get("Uint8Array").New(len(cartridgeRam))
	js.CopyBytesToJS(jsBuffer, cartridgeRam)

	return jsBuffer
}

func setCartridgeRam(this js.Value, args []js.Value) interface{} {
	jsRamData := args[0]
	var cartridgeRam []byte

	if jsRamData.IsUndefined() || jsRamData.IsNull() {
		cartridgeRam = nil
	} else {
		cartridgeRam = make([]byte, jsRamData.Get("length").Int())
		js.CopyBytesToGo(cartridgeRam, jsRamData)
	}

	gb.SetCartridgeRam(cartridgeRam)

	return nil
}

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

var joypadInputFromString = map[string]joypad.JoypadInput{
	"RIGHT":  joypad.JoypadInputRight,
	"LEFT":   joypad.JoypadInputLeft,
	"UP":     joypad.JoypadInputUp,
	"DOWN":   joypad.JoypadInputDown,
	"A":      joypad.JoypadInputA,
	"B":      joypad.JoypadInputB,
	"SELECT": joypad.JoypadInputSelect,
	"START":  joypad.JoypadInputStart,
}

func handleJoypadButtonPressed(this js.Value, args []js.Value) interface{} {
	if gb == nil {
		return nil
	}

	s := args[0].String()
	input, ok := joypadInputFromString[s]
	if !ok {
		return nil
	}

	gb.PressJoypadInput(input)
	return nil
}

func handleJoypadButtonReleased(this js.Value, args []js.Value) interface{} {
	if gb == nil {
		return nil
	}

	s := args[0].String()
	input, ok := joypadInputFromString[s]
	if !ok {
		return nil
	}

	gb.ReleaseJoypadInput(input)
	return nil
}

const (
	displayWidth  = 160
	displayHeight = 144
)

var goImageData [displayWidth * displayHeight * 4]byte
var jsImageData js.Value
var frameReady bool = false

func presentFrame() {
	frameBuffer := gb.FrameBuffer()
	i := 0
	for screenY := 0; screenY < displayHeight; screenY++ {
		for screenX := 0; screenX < displayWidth; screenX++ {
			color := frameBuffer[screenY][screenX]
			switch color {
			case 0:
				goImageData[i], goImageData[i+1], goImageData[i+2], goImageData[i+3] = 208, 224, 64, 255
			case 1:
				goImageData[i], goImageData[i+1], goImageData[i+2], goImageData[i+3] = 160, 168, 48, 255
			case 2:
				goImageData[i], goImageData[i+1], goImageData[i+2], goImageData[i+3] = 96, 112, 40, 255
			case 3:
				goImageData[i], goImageData[i+1], goImageData[i+2], goImageData[i+3] = 56, 72, 40, 255
			}
			i += 4
		}
	}

	frameReady = true
}

// pollFrame returns a newly completed frame, if one is available.
// The frame is consumed exactly once.
func pollFrame(this js.Value, args []js.Value) interface{} {
	if !frameReady {
		return nil
	}

	frameReady = false
	js.CopyBytesToJS(jsImageData, goImageData[:])
	return jsImageData
}

func pollAudioBuffer(this js.Value, args []js.Value) interface{} {
	samples := gb.ConsumeAudioBuffer()

	if len(samples) == 0 {
		return nil
	}

	byteLen := len(samples) * 2
	bytes := unsafe.Slice((*byte)(unsafe.Pointer(&samples[0])), byteLen)

	jsUint8Array := js.Global().Get("Uint8Array").New(byteLen)
	js.CopyBytesToJS(jsUint8Array, bytes)

	jsInt16Array := js.Global().Get("Int16Array").New(
		jsUint8Array.Get("buffer"),
		jsUint8Array.Get("byteOffset"),
		len(samples),
	)

	return jsInt16Array
}

func enableTraceLogging(this js.Value, args []js.Value) interface{} {
	logger.GlobalTraceLogger.Enable()
	return nil
}

func disableTraceLogging(this js.Value, args []js.Value) interface{} {
	logger.GlobalTraceLogger.Disable()
	return nil
}

func getTraceLogs(this js.Value, args []js.Value) interface{} {
	buffer := logger.GlobalTraceLogger.GetBuffer()

	jsBuffer := js.Global().Get("Uint8Array").New(len(buffer))
	js.CopyBytesToJS(jsBuffer, buffer)

	return jsBuffer
}

func resetTraceLogs(this js.Value, args []js.Value) interface{} {
	logger.GlobalTraceLogger.Reset()
	return nil
}

// getDebugInfo returns a snapshot of the emulator's state.
func getDebugInfo(this js.Value, args []js.Value) interface{} {
	if gb == nil {
		return nil
	}

	return gb.Debug()
}
