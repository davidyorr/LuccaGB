package apu

// At 48kHz - ~85ms of audio (4096 samples ÷ 48000 samples/second).
// Must be a power of 2 for the bitwise mask to work.
const AudioBufferSize = 4096
const bufferMask = AudioBufferSize - 1

type AudioBuffer struct {
	buffer [AudioBufferSize]int16
	head   int
	tail   int
	size   int
}

func (b *AudioBuffer) Write(sample int16) {
	// if full, advance head (overwrite oldest sample)
	if b.size == AudioBufferSize {
		b.head = (b.head + 1) & bufferMask
		b.size--
	}

	b.buffer[b.tail] = sample
	b.tail = (b.tail + 1) & bufferMask
	b.size++
}

// Read copies up to len(dst) samples into dst.
// Returns the actual number of samples copied.
func (b *AudioBuffer) Read(dst []int16) int {
	if b.size == 0 {
		return 0
	}

	// calculate how many we can read
	// either as much as requested, or as much as we have
	toRead := len(dst)
	if toRead > b.size {
		toRead = b.size
	}

	// if head < tail, data is contiguous in the buffer
	if b.head < b.tail {
		copy(dst[:toRead], b.buffer[b.head:b.head+toRead])
	} else {
		// if head >= tail, data wraps around the end of the buffer
		// copy from head to end of array
		chunk1 := AudioBufferSize - b.head
		// if the requested amount is smaller than chunk1, just copy what fits
		if toRead < chunk1 {
			copy(dst, b.buffer[b.head:b.head+toRead])
		} else {
			// copy all of chunk1
			copy(dst, b.buffer[b.head:])

			// copy remaining from start of array
			chunk2 := toRead - chunk1
			copy(dst[chunk1:], b.buffer[:chunk2])
		}
	}

	// advance pointers
	b.head = (b.head + toRead) & bufferMask
	b.size -= toRead

	return toRead
}
