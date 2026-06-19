package gameboy

import (
	"github.com/davidyorr/LuccaGB/internal/apu"
	"github.com/davidyorr/LuccaGB/internal/bus"
	"github.com/davidyorr/LuccaGB/internal/cartridge"
	"github.com/davidyorr/LuccaGB/internal/cpu"
	"github.com/davidyorr/LuccaGB/internal/dma"
	"github.com/davidyorr/LuccaGB/internal/interrupt"
	"github.com/davidyorr/LuccaGB/internal/joypad"
	"github.com/davidyorr/LuccaGB/internal/logger"
	"github.com/davidyorr/LuccaGB/internal/mmu"
	"github.com/davidyorr/LuccaGB/internal/ppu"
	"github.com/davidyorr/LuccaGB/internal/serial"
	"github.com/davidyorr/LuccaGB/internal/timer"
)

type Gameboy struct {
	cpu       *cpu.CPU
	ppu       *ppu.PPU
	apu       *apu.APU
	mmu       *mmu.MMU
	dma       *dma.DMA
	timer     *timer.Timer
	serial    *serial.Serial
	bus       *bus.Bus
	cartridge *cartridge.Cartridge
	joypad    *joypad.Joypad

	// non-hardware: circular rewind buffer
	rewindBuffer      [][]byte // Circular buffer of serialized states
	rewindCapacity    int      // Maximum number of states to hold
	rewindHead        int      // Index where the next state will be written
	rewindCount       int      // Current number of stored states
	pendingRewindSave bool     // Flag to save on the next safe cycle
	serializeBuf      []byte   // Reusable buffer
}

func New() *Gameboy {
	cartridge := cartridge.New()
	cpu := cpu.New()
	apu := apu.New()
	timer := timer.New()
	serial := serial.New()
	bus := bus.New()
	mmu := mmu.New(cartridge)
	dma := dma.New()
	ppu := ppu.New(mmu.RequestInterrupt)
	joypad := joypad.New(mmu.RequestInterrupt)

	mmu.ConnectJoypad(joypad)
	bus.Connect(mmu, timer, serial, ppu, apu, dma)
	cpu.ConnectBus(bus)
	dma.ConnectBus(bus)
	dma.ConnectPpu(ppu)

	return &Gameboy{
		cpu:          cpu,
		ppu:          ppu,
		apu:          apu,
		mmu:          mmu,
		dma:          dma,
		timer:        timer,
		serial:       serial,
		bus:          bus,
		cartridge:    cartridge,
		joypad:       joypad,
		serializeBuf: make([]byte, 1024*512), // 512KB
	}
}

func (gameboy *Gameboy) LoadRom(rom []uint8) cartridge.CartridgeInfo {
	logger.Info("GAMEBOY LOAD ROM", "SIZE", len(rom))

	// Reset rewind buffer so stale states from a previous ROM can't be loaded
	gameboy.ResetRewindBuffer()

	return gameboy.cartridge.LoadRom(rom)
}

func (gameboy *Gameboy) CartridgeRam() []uint8 {
	return gameboy.cartridge.Ram()
}

func (gameboy *Gameboy) SetCartridgeRam(ram []uint8) {
	gameboy.cartridge.SetRam(ram)
}

// Advance the entire system by 1 M-cycle (4 T-cycles)
func (gameboy *Gameboy) Step() (tCycles uint8, frameReady bool, err error) {
	for range 4 {
		gameboy.dma.Step()
		if gameboy.ppu.Step() {
			frameReady = true
			gameboy.pendingRewindSave = true
		}
		gameboy.apu.Step()
		gameboy.cpu.Step()
		requestTimerInterrupt := gameboy.timer.Step()
		if requestTimerInterrupt {
			gameboy.mmu.RequestInterrupt(interrupt.TimerInterrupt)
		}
		requestSerialInterrupt := gameboy.serial.Step()
		if requestSerialInterrupt {
			gameboy.mmu.RequestInterrupt(interrupt.SerialInterrupt)
		}
	}

	if gameboy.pendingRewindSave && gameboy.IsSafeToSerialize() {
		gameboy.saveRewindState()
		gameboy.pendingRewindSave = false
	}

	return 4, frameReady, nil
}

// StepFrames runs the emulator until exactly n frames are generated.
// It ignores real-time syncing and runs as fast as the CPU allows.
func (gameboy *Gameboy) StepFrames(frames int) {
	framesSeen := 0
	for {
		_, frameReady, _ := gameboy.Step()
		if frameReady {
			framesSeen++
			if framesSeen == frames {
				break
			}
		}
	}
}

// SetJoypadState sets the entire controller state in one go.
func (gameboy *Gameboy) SetJoypadState(state uint8) {
	mask := state

	// Helper to update button state based on bitmask
	setButton := func(input joypad.JoypadInput, pressed bool) {
		if pressed {
			gameboy.PressJoypadInput(input)
		} else {
			gameboy.ReleaseJoypadInput(input)
		}
	}

	// Bit 0: Right, 1: Left, 2: Up, 3: Down, 4: A, 5: B, 6: Select, 7: Start
	setButton(joypad.JoypadInputRight, (mask&0b0000_0001) != 0)
	setButton(joypad.JoypadInputLeft, (mask&0b0000_0010) != 0)
	setButton(joypad.JoypadInputUp, (mask&0b0000_0100) != 0)
	setButton(joypad.JoypadInputDown, (mask&0b0000_1000) != 0)
	setButton(joypad.JoypadInputA, (mask&0b0001_0000) != 0)
	setButton(joypad.JoypadInputB, (mask&0b0010_0000) != 0)
	setButton(joypad.JoypadInputSelect, (mask&0b0100_0000) != 0)
	setButton(joypad.JoypadInputStart, (mask&0b1000_0000) != 0)
}

func (gameboy *Gameboy) PressJoypadInput(input joypad.JoypadInput) {
	gameboy.joypad.Press(input)
}

func (gameboy *Gameboy) ReleaseJoypadInput(input joypad.JoypadInput) {
	gameboy.joypad.Release(input)
}

func (gameboy *Gameboy) FrameBuffer() [144][160]uint8 {
	return gameboy.ppu.FrameBuffer()
}

func (gameboy *Gameboy) FrameBufferDownsampled() [72][80]uint8 {
	return gameboy.ppu.FrameBufferDownsampled()
}

func (gameboy *Gameboy) ReadSamples(dst []int16) int {
	return gameboy.apu.ReadSamples(dst)
}

func (gameboy *Gameboy) SetAudioChannelEnabled(channel int, enabled bool) {
	gameboy.apu.SetChannelEnabled(channel, enabled)
}

func (gameboy *Gameboy) GetAudioChannelEnabled(channel int) bool {
	return gameboy.apu.GetChannelEnabled(channel)
}

func (gameboy *Gameboy) ReadMemory(address uint16) uint8 {
	return gameboy.bus.DirectRead(address)
}

// SetRewindBufferSize sets the maximum number of frames to store in the rewind buffer.
// If the new size is smaller than the current count, it preserves the most recent states.
func (gb *Gameboy) SetRewindBufferSize(size int) {
	if size < 0 {
		size = 0
	}
	if size == gb.rewindCapacity {
		return
	}

	newBuffer := make([][]byte, size)

	// Determine how many of the existing states we can keep
	keep := gb.rewindCount
	if keep > size {
		keep = size
	}

	if keep > 0 {
		oldStates := gb.GetRewindBuffer() // Gets them chronologically (oldest to newest)

		// Copy the "keep" newest states to the new buffer
		startIdx := len(oldStates) - keep
		for i := 0; i < keep; i++ {
			newBuffer[i] = oldStates[startIdx+i]
		}
	}

	gb.rewindBuffer = newBuffer
	gb.rewindCapacity = size
	gb.rewindCount = keep

	if size > 0 {
		gb.rewindHead = keep % size
	} else {
		gb.rewindHead = 0
	}
}

// GetRewindBuffer returns a copy of all currently saved states in chronological
// order (from oldest to newest). It does not modify or consume the buffer.
func (gb *Gameboy) GetRewindBuffer() [][]byte {
	if gb.rewindCount == 0 {
		return nil
	}

	states := make([][]byte, gb.rewindCount)

	// Find the index of the oldest element
	oldest := 0
	if gb.rewindCount == gb.rewindCapacity {
		oldest = gb.rewindHead
	}

	for i := 0; i < gb.rewindCount; i++ {
		idx := (oldest + i) % gb.rewindCapacity
		states[i] = gb.rewindBuffer[idx]
	}

	return states
}

// GetRewindCapacity gets the size of the rewind buffer.
func (gb *Gameboy) GetRewindCapacity() int {
	return gb.rewindCapacity
}

// Rewind pops the most recent state off the buffer and loads it.
// Returns true if a state was successfully loaded, false if the buffer is empty.
func (gb *Gameboy) Rewind() bool {
	if gb.rewindCount == 0 {
		return false
	}

	gb.rewindHead--
	if gb.rewindHead < 0 {
		gb.rewindHead = gb.rewindCapacity - 1
	}

	stateData := gb.rewindBuffer[gb.rewindHead]
	gb.DeserializeState(stateData)
	gb.rewindCount--

	return true
}

func (gb *Gameboy) ResetRewindBuffer() {
	gb.rewindHead = 0
	gb.rewindCount = 0
	gb.pendingRewindSave = false
}

// internal helper to save the state
func (gb *Gameboy) saveRewindState() {
	if gb.rewindCapacity <= 0 {
		return
	}

	// Serialize into our pre-allocated, reusable buffer
	data := gb.SerializeState(gb.serializeBuf)

	// Reuse existing buffer if possible to prevent allocations every frame
	snapshot := gb.rewindBuffer[gb.rewindHead]
	snapshot = append(snapshot[:0], data...)

	// Store in the ring buffer
	gb.rewindBuffer[gb.rewindHead] = snapshot

	// Advance head and count
	gb.rewindHead = (gb.rewindHead + 1) % gb.rewindCapacity
	if gb.rewindCount < gb.rewindCapacity {
		gb.rewindCount++
	}
}

// Debug gathers debug information from all components, acting as a single entry
// point for the frontend to get a snapshot of the machine state.
func (gb *Gameboy) Debug() map[string]interface{} {
	debugInfo := make(map[string]interface{})

	debugInfo["apu"] = gb.apu.Debug()
	debugInfo["cartridge"] = gb.cartridge.Debug()
	debugInfo["cpu"] = gb.cpu.Debug()
	// debugInfo["ppu"] = gb.ppu.Debug()

	return debugInfo
}

func (gb *Gameboy) IsSafeToSerialize() bool {
	return gb.cpu.IsSafeToSerialize()
}

func (gb *Gameboy) SerializeState(buf []byte) []byte {
	offset := 0

	offset += gb.cpu.Serialize(buf[offset:])
	offset += gb.apu.Serialize(buf[offset:])
	offset += gb.ppu.Serialize(buf[offset:])
	offset += gb.mmu.Serialize(buf[offset:])
	offset += gb.dma.Serialize(buf[offset:])
	offset += gb.timer.Serialize(buf[offset:])
	offset += gb.serial.Serialize(buf[offset:])
	offset += gb.cartridge.Serialize(buf[offset:])
	offset += gb.joypad.Serialize(buf[offset:])

	return buf[:offset]
}

func (gb *Gameboy) DeserializeState(data []byte) {
	offset := 0

	offset += gb.cpu.Deserialize(data[offset:])
	offset += gb.apu.Deserialize(data[offset:])
	offset += gb.ppu.Deserialize(data[offset:])
	offset += gb.mmu.Deserialize(data[offset:])
	offset += gb.dma.Deserialize(data[offset:])
	offset += gb.timer.Deserialize(data[offset:])
	offset += gb.serial.Deserialize(data[offset:])
	offset += gb.cartridge.Deserialize(data[offset:])
	offset += gb.joypad.Deserialize(data[offset:])
}
