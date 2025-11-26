package gameboy

import (
	"github.com/davidyorr/LuccaGB/bus"
	"github.com/davidyorr/LuccaGB/cartridge"
	"github.com/davidyorr/LuccaGB/cpu"
	"github.com/davidyorr/LuccaGB/dma"
	"github.com/davidyorr/LuccaGB/interrupt"
	"github.com/davidyorr/LuccaGB/logger"
	"github.com/davidyorr/LuccaGB/mmu"
	"github.com/davidyorr/LuccaGB/ppu"
	"github.com/davidyorr/LuccaGB/serial"
	"github.com/davidyorr/LuccaGB/timer"
)

type Gameboy struct {
	cpu       *cpu.CPU
	ppu       *ppu.PPU
	mmu       *mmu.MMU
	dma       *dma.DMA
	timer     *timer.Timer
	serial    *serial.Serial
	cartridge *cartridge.Cartridge
}

func New() *Gameboy {
	cartridge := cartridge.New()
	cpu := cpu.New()
	timer := timer.New()
	serial := serial.New()
	bus := bus.New()
	mmu := mmu.New(cartridge)
	dma := dma.New()
	ppu := ppu.New(mmu.RequestInterrupt)

	bus.Connect(mmu, timer, serial, ppu, dma)
	cpu.ConnectBus(bus)
	dma.ConnectBus(bus)
	dma.ConnectPpu(ppu)

	return &Gameboy{
		cpu:       cpu,
		ppu:       ppu,
		mmu:       mmu,
		dma:       dma,
		timer:     timer,
		serial:    serial,
		cartridge: cartridge,
	}
}

func (gameboy *Gameboy) LoadRom(rom []uint8) {
	logger.Info("GAMEBOY LOAD ROM", "SIZE", len(rom))

	gameboy.cartridge.LoadRom(rom)
}

// Advance the entire system by 1 M-cycle (4 T-cycles)
func (gameboy *Gameboy) Step() (tCycles uint8, frameReady bool, err error) {
	for range 4 {
		gameboy.dma.Step()
		frameReady = gameboy.ppu.Step()
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

func (gameboy *Gameboy) FrameBuffer() [144][160]uint8 {
	return gameboy.ppu.FrameBuffer()
}

// Debug gathers debug information from all components, acting as a single entry
// point for the frontend to get a snapshot of the machine state.
func (gb *Gameboy) Debug() map[string]interface{} {
	debugInfo := make(map[string]interface{})

	debugInfo["cpu"] = gb.cpu.Debug()
	// debugInfo["ppu"] = gb.ppu.Debug()

	return debugInfo
}
