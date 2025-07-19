package bus

import (
	"fmt"

	"github.com/davidyorr/EchoGB/dma"
	"github.com/davidyorr/EchoGB/logger"
	"github.com/davidyorr/EchoGB/mmu"
	"github.com/davidyorr/EchoGB/ppu"
	"github.com/davidyorr/EchoGB/serial"
	"github.com/davidyorr/EchoGB/timer"
)

type Bus struct {
	ppu    *ppu.PPU
	mmu    *mmu.MMU
	dma    *dma.DMA
	timer  *timer.Timer
	serial *serial.Serial
}

func New() *Bus {
	bus := &Bus{}

	return bus
}

func (bus *Bus) Connect(mmu *mmu.MMU, timer *timer.Timer, serial *serial.Serial, ppu *ppu.PPU, dma *dma.DMA) {
	bus.mmu = mmu
	bus.dma = dma
	bus.timer = timer
	bus.serial = serial
	bus.ppu = ppu
}

func (bus *Bus) Read(address uint16) (value uint8) {
	// handle DMA transfer
	if bus.dma.Active() {
		// DMA
		if address == 0xFF46 {
			return bus.DirectRead(address)
		} else if address >= 0xFF80 && address <= 0xFFFE {
			// HRAM
			logger.Info("DMA ACTIVE, READING FROM BUS FOR HRAM")
			return bus.DirectRead(address)
		} else if address >= 0xC000 && address <= 0xFDFF {
			// WRAM
			logger.Info("DMA ACTIVE, READING FROM BUS FOR WRAM")
			return bus.DirectRead(address)
		} else if address >= 0xFE00 && address <= 0xFE9F {
			// OAM
			logger.Info("DMA ACTIVE, 0xFF")
			return 0xFF
		} else {
			logger.Info(fmt.Sprintf("DMA ACTIVE, RETURNING CURRENT TRANSFER BYTE: 0x%0X2", bus.dma.CurrentTransferByte()))
			return bus.dma.CurrentTransferByte()
		}
	}

	// handle unusable area when OAM is blocked
	if bus.ppu.OamIsBlocked() && address >= 0xFEA0 && address <= 0xFEFF {
		return 0xFF
	}

	return bus.DirectRead(address)
}

// DirectRead performs a raw read from memory. It acts as the central
// dispatcher, routing the address to the appropriate hardware component. This
// is the lowest-level read operation on the bus. It deliberately contains no
// logic for bus contention.
func (bus *Bus) DirectRead(address uint16) (value uint8) {
	switch {
	// DMA
	case address == 0xFF46:
		return bus.dma.DmaRegister()
	// PPU LCD
	case address >= 0xFF40 && address <= 0xFF4B:
		return bus.ppu.Read(address)
	// PPU VRAM
	case address >= 0x8000 && address <= 0x9FFF:
		return bus.ppu.Read(address)
	// PPU OAM
	case address >= 0xFE00 && address <= 0xFE9F:
		return bus.ppu.Read(address)
	// Unusable
	case address >= 0xFEA0 && address <= 0xFEFF:
		return 0x00
	// timers
	case address >= 0xFF04 && address <= 0xFF07:
		return bus.timer.Read(address)
	// serial data transfer
	case address == 0xFF01 || address == 0xFF02:
		return bus.serial.Read(address)
	// handle everything else with the MMU
	default:
		return bus.mmu.Read(address)
	}
}

func (bus *Bus) Write(address uint16, value uint8) {
	// writes to the DMA register should succeed regardless of DMA transfer state
	if address == 0xFF46 {
		bus.dma.SetDmaRegister(value)
	}
	// during a transfer, only HRAM and WRAM can be accessed
	if bus.dma.Active() && (!(address >= 0xFF80 && address <= 0xFFFE) || !(address >= 0xC000 && address <= 0xFDFF)) {
		logger.Info("DMA ACTIVE, IGNORING WRITE")
		return
	}
	switch {
	// DMA
	case address == 0xFF46:
		bus.dma.StartTransfer(value)
	// PPU LCD
	case address >= 0xFF40 && address <= 0xFF4B:
		bus.ppu.Write(address, value)
	// PPU VRAM
	case address >= 0x8000 && address <= 0x9FFF:
		bus.ppu.Write(address, value)
	// PPU OAM
	case address >= 0xFE00 && address <= 0xFE9F:
		bus.ppu.Write(address, value)
	// Unusable
	case address >= 0xFEA0 && address <= 0xFEFF:
		logger.Info("UNUSABLE WRITE")
		return
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

func (bus *Bus) DmaActive() bool {
	return bus.dma.Active()
}
