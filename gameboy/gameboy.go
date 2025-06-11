package gameboy

import (
	"fmt"

	"github.com/davidyorr/EchoGB/cartridge"
	"github.com/davidyorr/EchoGB/cpu"
	"github.com/davidyorr/EchoGB/mmu"
	"github.com/davidyorr/EchoGB/ppu"
)

type Gameboy struct {
	cpu       *cpu.CPU
	ppu       *ppu.PPU
	mmu       *mmu.MMU
	cartridge *cartridge.Cartridge
}

func New() *Gameboy {
	cartridge := cartridge.New()
	mmu := mmu.New(cartridge)
	cpu := cpu.New()
	cpu.ConnectBus(mmu)
	ppu := ppu.New()

	return &Gameboy{
		cpu:       cpu,
		ppu:       ppu,
		mmu:       mmu,
		cartridge: cartridge,
	}
}

func (gameboy *Gameboy) LoadRom(rom []uint8) {
	fmt.Println("Go: load ROM", len(rom))

	gameboy.cartridge.SetRom(rom)
	gameboy.cpu.Step()
}
