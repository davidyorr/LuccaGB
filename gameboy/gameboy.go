package gameboy

import (
	"github.com/davidyorr/EchoGB/cpu"
	"github.com/davidyorr/EchoGB/ppu"
)

type Gameboy struct {
	cpu *cpu.CPU
	ppu *ppu.PPU
}

func New() *Gameboy {
	return &Gameboy{
		cpu: cpu.New(),
		ppu: ppu.New(),
	}
}
