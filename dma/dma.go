package dma

import "github.com/davidyorr/EchoGB/logger"

type DMA struct {
	active        bool
	sourceAddress uint16
	progress      uint8
	startDelay    uint8
	tCycleCounter uint8
	bus           MemoryBus
}

const (
	// 160 M-cycles: 640 dots (1.4 lines)
	transferDuration = 160
)

type MemoryBus interface {
	Read(address uint16) uint8
	Write(address uint16, value uint8)
}

func New() *DMA {
	dma := &DMA{}

	dma.Reset()

	return dma
}

func (dma *DMA) Reset() {
	dma.active = false
	dma.sourceAddress = 0
	dma.progress = 0
}

func (dma *DMA) ConnectBus(bus MemoryBus) {
	dma.bus = bus
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
	if !dma.active {
		return
	}

	if dma.startDelay > 0 {
		dma.startDelay--
		return
	}

	// Source:      $XX00 - $XX9F
	// Destination: $FE00 - $FE9F

	source := dma.sourceAddress + uint16(dma.progress)
	destination := 0xFE00 + uint16(dma.progress)

	sourceByte := dma.bus.Read(source)
	dma.bus.Write(destination, sourceByte)

	dma.progress++

	if dma.progress == transferDuration {
		logger.Info("FINISHED DMA TRANSFER")
		dma.active = false
	}
}

func (dma *DMA) StartTransfer(value uint8) {
	dma.active = true
	dma.progress = 0
	// there's a 4 T-cycle delay before the transfer begins
	dma.startDelay = 4
	dma.sourceAddress = uint16(value) << 8
}

func (dma *DMA) Active() bool {
	return dma.active
}
