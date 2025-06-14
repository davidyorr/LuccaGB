//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"

	"github.com/davidyorr/EchoGB/gameboy"
)

// 0xFFFF
var interruptRegister uint8

func main() {
	fmt.Println("Hello EchoGB!")

	js.Global().Set("loadRom", js.FuncOf(loadRom))

	<-make(chan struct{})
}

func initPostBootRomState() {
	interruptRegister = 0x00 // IE
}

func loadRom(this js.Value, args []js.Value) interface{} {
	jsRomData := args[0]
	cartridgeRom := make([]byte, jsRomData.Get("length").Int())
	js.CopyBytesToGo(cartridgeRom, jsRomData)

	gb := gameboy.New()
	gb.LoadRom(cartridgeRom)

	return nil
}
