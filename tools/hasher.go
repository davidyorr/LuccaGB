package tools

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
)

// StandardPalette: 0x00=Black, 0xFF=White
var StandardPalette = map[color.RGBA]byte{
	{0x00, 0x00, 0x00, 0xFF}: 0x00,
	{0x55, 0x55, 0x55, 0xFF}: 0x55,
	{0xAA, 0xAA, 0xAA, 0xFF}: 0xAA,
	{0xFF, 0xFF, 0xFF, 0xFF}: 0xFF,
}

// Map PPU values (0-3) to Hash Bytes.
var PpuMap = [4]byte{0xFF, 0xAA, 0x55, 0x00}

// HashFrameBuffer: Used by the Emulator Test
func HashFrameBuffer(buffer [144][160]uint8) string {
	hasher := sha256.New()
	for y := 0; y < 144; y++ {
		for x := 0; x < 160; x++ {
			val := buffer[y][x]
			if val > 3 {
				val = 3
			}
			hasher.Write([]byte{PpuMap[val]})
		}
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

// HashImage: Used by the CLI Tool (reads PNGs)
func HashImage(img image.Image, validColors map[color.RGBA]byte) (string, error) {
	if validColors == nil {
		validColors = StandardPalette
	}

	bounds := img.Bounds()
	if bounds.Dx() != 160 || bounds.Dy() != 144 {
		return "", fmt.Errorf("invalid dimensions: %dx%d", bounds.Dx(), bounds.Dy())
	}

	hasher := sha256.New()
	for y := 0; y < 144; y++ {
		for x := 0; x < 160; x++ {
			c := img.At(x, y)
			rgba := color.RGBAModel.Convert(c).(color.RGBA)

			canonicalByte, ok := validColors[rgba]
			if !ok {
				return "", fmt.Errorf("invalid pixel at (%d,%d): %v", x, y, rgba)
			}
			hasher.Write([]byte{canonicalByte})
		}
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
