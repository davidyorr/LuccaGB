package cartridge

type MBC interface {
	Read(address uint16) uint8
	Write(address uint16, value uint8)
}

var addressMaskSizes = map[uint8]uint32{
	0x00: 0x7FFF,   // 15 bits - 32KiB (2 banks, no banking)
	0x01: 0xFFFF,   // 16 bits - 64KiB (4 banks)
	0x02: 0x1FFFF,  // 17 bits - 128KiB (8 banks)
	0x03: 0x3FFFF,  // 18 bits - 256KiB (16 banks)
	0x04: 0x7FFFF,  // 19 bits - 512KiB (32 banks)
	0x05: 0xFFFFF,  // 20 bits - 1MiB (64 banks)
	0x06: 0x1FFFFF, // 21 bits - 2MiB (128 banks)
	0x07: 0x3FFFFF, // 22 bits - 4MiB (256 banks)
	0x08: 0x7FFFFF, // 23 bits - 8MiB (512 banks)
}

var ramAddressMaskSizes = map[uint8]uint32{
	0x02: 0x1FFF,  // 13 bits
	0x03: 0x7FFF,  // 15 bits
	0x04: 0x1FFFF, // 17 bits
	0x05: 0xFFFF,  // 16 bits
}

var ramSizes = map[uint8]int{
	0x00: 0, // No RAM
	// 0x01: - Unused
	0x02: 8192,   // 8 KiB (1 bank)
	0x03: 32768,  // 32 KiB (4 banks)
	0x04: 131072, // 128 KiB (16 banks)
	0x05: 65536,  // 64 KiB (8 banks)
}
