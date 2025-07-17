package dma

import (
	"fmt"

	"github.com/davidyorr/EchoGB/logger"
	"github.com/davidyorr/EchoGB/ppu"
)

type DMA struct {
	state                  TransferState
	sourceAddress          uint16
	progress               uint8
	tCycleCounter          uint8
	currentTransferByte    uint8
	requestedSourceAddress uint8
	startingSourceAddress  uint8
	bus                    MemoryBus
	ppu                    *ppu.PPU
}

type TransferState uint8

const (
	Idle TransferState = iota
	Requested
	Starting
	Active
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
	dma.state = Idle
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
	case Active:
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
		)

		if dma.progress == transferDuration {
			logger.Info("FINISHED DMA TRANSFER")
			dma.state = Idle
		}
	case Starting:
		logger.Info("DMA STATE MOVING FROM STARTING -> ACTIVE")
		dma.state = Active
		dma.progress = 0
		dma.sourceAddress = uint16(dma.startingSourceAddress) << 8
	case Requested:
		logger.Info("DMA STATE MOVING FROM REQUESTED -> STARTING")
		dma.startingSourceAddress = dma.requestedSourceAddress
		dma.state = Starting
	}
}

func (dma *DMA) StartTransfer(value uint8) {
	logger.Info("DMA STATE MOVING FROM IDLE -> REQUESTED")
	dma.requestedSourceAddress = value
	dma.state = Requested
}

func (dma *DMA) Active() bool {
	return dma.state == Active
}

func (dma *DMA) CurrentTransferByte() uint8 {
	return dma.currentTransferByte
}
