package bus

import (
	"fmt"

	"github.com/davidyorr/EchoGB/logger"
	"github.com/davidyorr/EchoGB/mmu"
	"github.com/davidyorr/EchoGB/ppu"
	"github.com/davidyorr/EchoGB/serial"
	"github.com/davidyorr/EchoGB/timer"
)

type Bus struct {
	ppu    *ppu.PPU
	mmu    *mmu.MMU
	timer  *timer.Timer
	serial *serial.Serial
}

func New() *Bus {
	bus := &Bus{}

	return bus
}

func (bus *Bus) Connect(mmu *mmu.MMU, timer *timer.Timer, serial *serial.Serial, ppu *ppu.PPU) {
	bus.mmu = mmu
	bus.timer = timer
	bus.serial = serial
	bus.ppu = ppu
}

func (bus *Bus) Read(address uint16) uint8 {
	var value uint8 = 0

	switch {
	// PPU
	case address >= 0xFF40 && address <= 0xFF4B:
		value = bus.ppu.Read(address)
	// timers
	case address >= 0xFF04 && address <= 0xFF07:
		value = bus.timer.Read(address)
	// serial data transfer
	case address == 0xFF01 || address == 0xFF02:
		value = bus.serial.Read(address)
	// handle everything else with the MMU
	default:
		value = bus.mmu.Read(address)
	}

	logger.Debug(
		"BUS READ",
		"Address", fmt.Sprintf("0x%04X", address),
		"Value", fmt.Sprintf("0x%02X", value),
	)

	return value
}

func (bus *Bus) Write(address uint16, value uint8) {
	switch {
	// timers
	case address >= 0xFF04 && address <= 0xFF07:
		bus.timer.Write(address, value)
	// serial data transfer
	case address == 0xFF01 || address == 0xFF02:
		bus.serial.Write(address, value)
	// handle everything else with the MMU
	default:
		bus.mmu.Write(address, value)
	}

	logger.Debug(
		"BUS WRITE",
		"Address", fmt.Sprintf("0x%04X", address),
		"Value", fmt.Sprintf("0x%02X", value),
	)
}
