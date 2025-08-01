package dma

import (
	"fmt"

	"github.com/davidyorr/EchoGB/logger"
	"github.com/davidyorr/EchoGB/ppu"
)

type DMA struct {
	// 0xFF46 - DMA: OAM DMA source address & start
	dmaRegister            uint8
	state                  TransferState
	sourceAddress          uint16
	progress               uint8
	tCycleCounter          uint8
	currentTransferByte    uint8
	requestedSourceAddress uint8
	startingSourceAddress  uint8
	// true if the current transfer was started while another was in progress
	wasRestarted bool
	bus          MemoryBus
	ppu          *ppu.PPU
}

type TransferState uint8

const (
	StateIdle TransferState = iota
	StateRequested
	StateStarting
	StateActive
)

const (
	// 160 M-cycles: 640 dots (1.4 lines)
	transferDuration = 160
)

type MemoryBus interface {
	DirectRead(address uint16) uint8
	Write(address uint16, value uint8)
}

func New() *DMA {
	dma := &DMA{}

	dma.Reset()

	return dma
}

func (dma *DMA) Reset() {
	dma.state = StateIdle
	dma.sourceAddress = 0
	dma.progress = 0
}

func (dma *DMA) ConnectBus(bus MemoryBus) {
	dma.bus = bus
}

func (dma *DMA) ConnectPpu(ppu *ppu.PPU) {
	dma.ppu = ppu
}

// Perform 1 T-cycle of work
func (dma *DMA) Step() {
	dma.tCycleCounter++

	if dma.tCycleCounter == 4 {
		dma.tCycleCounter = 0
		dma.executeMachineCycle()
	}
}

// Perform 1 M-cycle of work
func (dma *DMA) executeMachineCycle() {
	switch dma.state {
	case StateActive:
		// Source:      $XX00 - $XX9F
		// Destination: $FE00 - $FE9F

		// if the source is VRAM and the PPU is in Mode 3, VRAM is locked
		sourceIsVram := dma.sourceAddress >= 0x8000 && dma.sourceAddress <= 0x9FFF
		if sourceIsVram && dma.ppu.Mode() == ppu.DrawingPixels {
			logger.Info("SOURCE IS VRAM, RETURNING")
			return
		}

		source := dma.sourceAddress + uint16(dma.progress)
		destination := 0xFE00 + uint16(dma.progress)

		dma.currentTransferByte = dma.bus.DirectRead(source)
		dma.ppu.WriteOam(destination, dma.currentTransferByte)

		dma.progress++
		logger.Info(
			"DMA WRITE",
			"PROGRESS", fmt.Sprintf("%d/%d", dma.progress, transferDuration),
			"ADDRESS", fmt.Sprintf("0x%04X", destination),
			"VALUE", fmt.Sprintf("0x%02X", dma.currentTransferByte),
		)

		if dma.progress == transferDuration {
			logger.Info("FINISHED DMA TRANSFER")
			dma.state = StateIdle
		}
	case StateStarting:
		logger.Info("DMA STATE MOVING FROM STARTING -> ACTIVE")
		dma.state = StateActive
		dma.progress = 0
		dma.sourceAddress = uint16(dma.startingSourceAddress) << 8
	case StateRequested:
		logger.Info("DMA STATE MOVING FROM REQUESTED -> STARTING")
		dma.startingSourceAddress = dma.requestedSourceAddress
		dma.state = StateStarting
	}
}

func (dma *DMA) StartTransfer(value uint8) {
	if dma.state == StateIdle {
		logger.Info("DMA STATE MOVING FROM IDLE -> REQUESTED")
		dma.wasRestarted = false
	}
	if dma.state == StateActive {
		logger.Info("DMA STATE MOVING FROM ACTIVE -> REQUESTED")
		dma.wasRestarted = true
	}
	dma.dmaRegister = value
	dma.requestedSourceAddress = value
	dma.state = StateRequested
}

func (dma *DMA) Active() bool {
	// for restarted DMA transfers
	if dma.wasRestarted {
		return dma.state != StateIdle
	}
	// for fresh DMA transfers
	return dma.state == StateActive
}

func (dma *DMA) CurrentTransferByte() uint8 {
	return dma.currentTransferByte
}

// SetDmaRegister has the side effect of starting a DMA transfer.
func (dma *DMA) SetDmaRegister(value uint8) {
	dma.dmaRegister = value
	dma.StartTransfer(value)
}

func (dma *DMA) DmaRegister() uint8 {
	return dma.dmaRegister
}
