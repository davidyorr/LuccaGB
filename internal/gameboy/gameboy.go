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
	cartridge *cartridge.Cartridge
	joypad    *joypad.Joypad
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
		cpu:       cpu,
		ppu:       ppu,
		apu:       apu,
		mmu:       mmu,
		dma:       dma,
		timer:     timer,
		serial:    serial,
		cartridge: cartridge,
		joypad:    joypad,
	}
}

func (gameboy *Gameboy) LoadRom(rom []uint8) cartridge.CartridgeInfo {
	logger.Info("GAMEBOY LOAD ROM", "SIZE", len(rom))

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
		frameReady = gameboy.ppu.Step()
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
