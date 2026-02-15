package main

/*
#include <stdlib.h>
#include <stdint.h>
*/
import "C"
import (
	"unsafe"

	"github.com/davidyorr/LuccaGB/internal/gameboy"
)

var gb *gameboy.Gameboy

//export Init
func Init() {
	gb = gameboy.New()
}

//export Step
func Step(cycles C.int) {
	gb.StepFrames(int(cycles))
}

//export LoadRom
func LoadRom(data *C.uint8_t, length C.int) {
	gb = gameboy.New()

	rom := C.GoBytes(unsafe.Pointer(data), length)
	gb.LoadRom(rom)
}

//export SetJoypad
func SetJoypad(state C.uint8_t) {
	gb.SetJoypadState(uint8(state))
}

var frameCache [72 * 80]uint8

//export GetFrame
func GetFrame() *C.uint8_t {
	downsampled := gb.FrameBufferDownsampled() // Returns [72][80]uint8
	i := 0
	for y := 0; y < 72; y++ {
		for x := 0; x < 80; x++ {
			frameCache[i] = downsampled[y][x]
			i++
		}
	}
	return (*C.uint8_t)(unsafe.Pointer(&frameCache[0]))
}

//export ReadMemory
func ReadMemory(address C.uint16_t) C.uint8_t {
	return C.uint8_t(gb.ReadMemory(uint16(address)))
}

type SerializedData struct {
	data   *C.uint8_t
	length C.int
}

//export GetSerializedState
func GetSerializedState(outLength *C.int) *C.uint8_t {
	// Fast-forward to a safe boundary.
	safetyLimit := 20
	for !gb.IsSafeToSerialize() && safetyLimit > 0 {
		gb.Step()
		safetyLimit--
	}

	const BufferSize = 256 * 1024 // 256KB
	buf := make([]byte, BufferSize)
	stateData := gb.SerializeState(buf)

	// Allocate C memory and copy the data
	cBuf := C.malloc(C.size_t(len(stateData)))
	copy((*[1 << 30]byte)(cBuf)[:len(stateData)], stateData)

	// Set the output length
	*outLength = C.int(len(stateData))

	return (*C.uint8_t)(cBuf)
}

//export FreeSerializedState
func FreeSerializedState(ptr *C.uint8_t) {
	C.free(unsafe.Pointer(ptr))
}

//export LoadSerializedState
func LoadSerializedState(data *C.uint8_t, length C.int) {
	stateData := C.GoBytes(unsafe.Pointer(data), length)
	gb.DeserializeState(stateData)
}

func main() {}
