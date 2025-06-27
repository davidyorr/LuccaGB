package gameboy

import (
	"fmt"

	"github.com/davidyorr/EchoGB/bus"
	"github.com/davidyorr/EchoGB/cartridge"
	"github.com/davidyorr/EchoGB/cpu"
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
	timer     *timer.Timer
	serial    *serial.Serial
	cartridge *cartridge.Cartridge
}

func New() *Gameboy {
	cartridge := cartridge.New()
	cpu := cpu.New()
	ppu := ppu.New()
	timer := timer.New()
	serial := serial.New()
	bus := bus.New()
	mmu := mmu.New(cartridge)

	bus.Connect(mmu, timer, serial, ppu)
	cpu.ConnectBus(bus)

	return &Gameboy{
		cpu:       cpu,
		ppu:       ppu,
		mmu:       mmu,
		timer:     timer,
		serial:    serial,
		cartridge: cartridge,
	}
}

func (gameboy *Gameboy) LoadRom(rom []uint8) {
	logger.Info("GAMEBOY LOAD ROM", "size", len(rom))

	gameboy.cartridge.LoadRom(rom)
}

func (gameboy *Gameboy) Step() (uint8, error) {
	cycles, err := gameboy.cpu.Step()
	if err != nil {
		return 0, err
	}

	requestTimerInterrupt := gameboy.timer.Step(cycles)
	if requestTimerInterrupt {
		gameboy.mmu.RequestInterrupt(interrupt.TimerInterrupt)
	}
	requestSerialInterrupt := gameboy.serial.Step(cycles)
	if requestSerialInterrupt {
		gameboy.mmu.RequestInterrupt(interrupt.SerialInterrupt)
	}

	gameboy.ppu.Step(cycles)

	logger.Debug(
		"END OF GAMEBOY STEP",
		"IME", fmt.Sprintf("%t", gameboy.cpu.InterruptMasterEnable()),
		"IE", fmt.Sprintf("%0X", gameboy.mmu.InterruptEnable()),
		"IF", fmt.Sprintf("%0X", gameboy.mmu.InterruptFlag()),
	)

	pendingInterrupts := gameboy.mmu.InterruptEnable() & gameboy.mmu.InterruptFlag()

	if gameboy.cpu.Halted() && pendingInterrupts != 0 {
		gameboy.cpu.Unhalt()
	}

	if gameboy.cpu.InterruptMasterEnable() && (pendingInterrupts != 0) {
		gameboy.cpu.HandleInterrupts()
	}

	return cycles, nil
}
