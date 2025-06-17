package bus

import (
	"fmt"

	"github.com/davidyorr/EchoGB/mmu"
	"github.com/davidyorr/EchoGB/ppu"
	"github.com/davidyorr/EchoGB/timer"
)

type Bus struct {
	ppu   *ppu.PPU
	mmu   *mmu.MMU
	timer *timer.Timer
}

func New() *Bus {
	bus := &Bus{}

	return bus
}

func (bus *Bus) Connect(mmu *mmu.MMU, timer *timer.Timer, ppu *ppu.PPU) {
	bus.mmu = mmu
	bus.timer = timer
	bus.ppu = ppu
}

func (bus *Bus) Read(address uint16) uint8 {
	var value uint8 = 0

	switch {
	// ROM
	case address <= 0x7FFF:
		value = bus.mmu.Read(address)
	// working RAM
	case address >= 0xC000 && address <= 0xDFFF:
		value = bus.mmu.Read(address)
	// high RAM
	case address >= 0xFF80 && address <= 0xFFFE:
		value = bus.timer.Read(address)
	// timers
	case address >= 0xFF00 && address <= 0xFF07:
		value = bus.timer.Read(address)
	default:
		fmt.Println("unhandled address while reading ->")
		value = 0xFF
	}

	fmt.Printf("  [BUS READ] Address: 0x%04X, Value: 0x%02X\n", address, value)

	return value
}

func (bus *Bus) Write(address uint16, value uint8) {
	switch {
	// ROM
	case address <= 0x7FFF:
		bus.mmu.Write(address, value)
	// working RAM
	case address >= 0xC000 && address <= 0xDFFF:
		bus.mmu.Write(address, value)
	// high RAM
	case address >= 0xFF80 && address <= 0xFFFE:
		bus.timer.Write(address, value)
	// timers
	case address >= 0xFF04 && address <= 0xFF07:
		bus.timer.Write(address, value)
	default:
		fmt.Println("unhandled address while writing ->")
	}

	fmt.Printf("  [BUS WRITE] Address: 0x%04X, Value: 0x%02X\n", address, value)

	// SB is the Serial Data register at address 0xFF01
	// SC is the Serial Control register at address 0xFF02
	if address == 0xFF02 && value == 0x81 {
		bus.mmu.Write(address, value)
	}
}
