package gameboy

import (
	"fmt"

	"github.com/davidyorr/EchoGB/bus"
	"github.com/davidyorr/EchoGB/cartridge"
	"github.com/davidyorr/EchoGB/cpu"
	"github.com/davidyorr/EchoGB/dma"
	"github.com/davidyorr/EchoGB/interrupt"
	"github.com/davidyorr/EchoGB/logger"
	"github.com/davidyorr/EchoGB/mmu"
	"github.com/davidyorr/EchoGB/ppu"
	"github.com/davidyorr/EchoGB/serial"
	"github.com/davidyorr/EchoGB/timer"
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
	logger.Info("GAMEBOY LOAD ROM", "size", len(rom))

	gameboy.cartridge.LoadRom(rom)
}

// Advance the entire system by 1 M-cycle (4 T-cycles)
func (gameboy *Gameboy) Step() (uint8, error) {
	for range 4 {
		gameboy.cpu.Step()
		requestTimerInterrupt := gameboy.timer.Step()
		if requestTimerInterrupt {
			gameboy.mmu.RequestInterrupt(interrupt.TimerInterrupt)
		}
		requestSerialInterrupt := gameboy.serial.Step()
		if requestSerialInterrupt {
			gameboy.mmu.RequestInterrupt(interrupt.SerialInterrupt)
		}
		gameboy.ppu.Step()
		gameboy.dma.Step()
	}

	logger.Debug(
		"END OF GAMEBOY STEP",
		"IME", fmt.Sprintf("%t", gameboy.cpu.InterruptMasterEnable()),
		"IE", fmt.Sprintf("%0X", gameboy.mmu.InterruptEnable()),
		"IF", fmt.Sprintf("%0X", gameboy.mmu.InterruptFlag()),
	)

	return 4, nil
}
