package ppu

const fifoSize = 8
const bufferMask = fifoSize - 1

type FIFO struct {
	// 4 possible colors
	colorId uint8
	// only applies to objects (sprites)
	palette            uint8
	backgroundPriority uint8
}

type PixelFifo struct {
	buffer [fifoSize]FIFO
	head   int
	tail   int
	size   int
}

func (f *PixelFifo) Reset() {
	f.head = 0
	f.tail = 0
	f.size = 0
}

func (f *PixelFifo) Push(p FIFO) {
	if f.size >= fifoSize {
		panic("FIFO push overflow")
	}

	f.buffer[f.tail] = p
	f.tail = (f.tail + 1) & bufferMask // mod 8, keep lower 3 bits
	f.size++
}

func (f *PixelFifo) Pop() FIFO {
	if f.size == 0 {
		panic("FIFO pop underflow")
	}

	p := f.buffer[f.head]
	f.head = (f.head + 1) & bufferMask // mod 8, keep lower 3 bits
	f.size--
	return p
}

func (f *PixelFifo) Peek(i int) *FIFO {
	if i >= f.size {
		panic("FIFO peek out of bounds")
	}

	index := (f.head + i) & bufferMask
	return &f.buffer[index]
}
