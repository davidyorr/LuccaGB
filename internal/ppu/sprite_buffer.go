package ppu

import "encoding/binary"

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

func (b *SpriteBuffer) Serialize(buf []byte) int {
	offset := 0

	n := copy(buf[offset:], b.data[:])
	offset += n

	binary.LittleEndian.PutUint64(buf[offset:], uint64(b.size))
	offset += 8

	return offset
}

func (b *SpriteBuffer) Deserialize(buf []byte) int {
	offset := 0

	n := copy(b.data[:], buf[offset:])
	offset += n

	b.size = int(binary.LittleEndian.Uint64(buf[offset:]))
	offset += 8

	return offset
}
