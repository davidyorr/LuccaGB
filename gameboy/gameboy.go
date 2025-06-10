package gameboy

import (
	"fmt"

	"github.com/davidyorr/EchoGB/cpu"
	"github.com/davidyorr/EchoGB/mmu"
	"github.com/davidyorr/EchoGB/ppu"
)

type Gameboy struct {
	cpu *cpu.CPU
	ppu *ppu.PPU
	mmu *mmu.MMU
}

func New() *Gameboy {
	return &Gameboy{
		cpu: cpu.New(),
		ppu: ppu.New(),
		mmu: mmu.New(),
	}
}

func (gameboy *Gameboy) LoadRom(rom []uint8) {
	fmt.Println("Go: load ROM", len(rom))
}
