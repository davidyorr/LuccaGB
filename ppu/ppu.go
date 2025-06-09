package ppu

// 0x8000 - 9FFF
var videoRam [8192]uint8

// 0xFE00 - 0xFE9F
var oam [160]uint8
