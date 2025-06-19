package gameboy

import (
	"fmt"

	"github.com/davidyorr/EchoGB/bus"
	"github.com/davidyorr/EchoGB/cartridge"
	"github.com/davidyorr/EchoGB/cpu"
	"github.com/davidyorr/EchoGB/mmu"
	"github.com/davidyorr/EchoGB/ppu"
	"github.com/davidyorr/EchoGB/timer"
)

type Gameboy struct {
	cpu       *cpu.CPU
	ppu       *ppu.PPU
	mmu       *mmu.MMU
	timer     *timer.Timer
	cartridge *cartridge.Cartridge
}

func New() *Gameboy {
	cartridge := cartridge.New()
	cpu := cpu.New()
	ppu := ppu.New()
	timer := timer.New()
	bus := bus.New()
	mmu := mmu.New(cartridge)

	bus.Connect(mmu, timer, ppu)
	cpu.ConnectBus(bus)

	return &Gameboy{
		cpu:       cpu,
		ppu:       ppu,
		mmu:       mmu,
		timer:     timer,
		cartridge: cartridge,
	}
}

func (gameboy *Gameboy) LoadRom(rom []uint8) {
	fmt.Println("Go: load ROM", len(rom))

	gameboy.cartridge.LoadRom(rom)
}

func (gameboy *Gameboy) Step() error {
	cycles, err := gameboy.cpu.Step()
	if err != nil {
		return err
	}

	gameboy.timer.Step(cycles)
	gameboy.ppu.Step(cycles)

	return nil
}
