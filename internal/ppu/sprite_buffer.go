package ppu

const MaxSpriteBufferSize = 10

type SpriteBuffer struct {
	data [MaxSpriteBufferSize]uint8
	size int
}

func (b *SpriteBuffer) Reset() {
	b.size = 0
}

func (b *SpriteBuffer) Push(v uint8) {
	if b.size < MaxSpriteBufferSize {
		b.data[b.size] = v
		b.size++
	}
}

func (b *SpriteBuffer) Remove(i int) {
	if i >= b.size {
		return
	}

	copy(b.data[i:], b.data[i+1:b.size])
	b.size--
}
