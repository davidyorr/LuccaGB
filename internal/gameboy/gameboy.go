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

func (gameboy *Gameboy) PressJoypadInput(input joypad.JoypadInput) {
	gameboy.joypad.Press(input)
}

func (gameboy *Gameboy) ReleaseJoypadInput(input joypad.JoypadInput) {
	gameboy.joypad.Release(input)
}

func (gameboy *Gameboy) FrameBuffer() [144][160]uint8 {
	return gameboy.ppu.FrameBuffer()
}

// Debug gathers debug information from all components, acting as a single entry
// point for the frontend to get a snapshot of the machine state.
func (gb *Gameboy) Debug() map[string]interface{} {
	debugInfo := make(map[string]interface{})

	debugInfo["cartridge"] = gb.cartridge.Debug()
	debugInfo["cpu"] = gb.cpu.Debug()
	// debugInfo["ppu"] = gb.ppu.Debug()

	return debugInfo
}
