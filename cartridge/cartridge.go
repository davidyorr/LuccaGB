package cartridge

import (
	"bytes"
	"fmt"

	"github.com/davidyorr/EchoGB/logger"
)

type Cartridge struct {
	// 0x0000 - 0x7FFF, dynamically sized
	rom []uint8
	// title of the ROM in uppercase ASCII
	title []uint8
	// memory bank controller type
	mbc uint8
	// ROM size
	romSizeCode uint8
	// RAM size, if any
	ramSizeCode uint8
}

func New() *Cartridge {
	cartridge := &Cartridge{}

	return cartridge
}

func (cartridge *Cartridge) LoadRom(rom []uint8) {
	cartridge.rom = rom

	cartridge.title = bytes.Trim(cartridge.rom[0x0134:0x0143], "\x00")
	logger.Debug(
		"CARTRIDGE LOAD ROM",
		"TITLE", string(cartridge.title),
		"MBC", fmt.Sprintf("0x%02X", cartridge.mbc),
		"ROM size code", fmt.Sprintf("0x%02X", cartridge.romSizeCode),
		"RAM size code", fmt.Sprintf("0x%02X", cartridge.ramSizeCode),
	)

	cartridge.mbc = cartridge.rom[0x147]

	cartridge.romSizeCode = cartridge.rom[0x148]

	cartridge.ramSizeCode = cartridge.rom[0x149]
}

func (cartridge *Cartridge) Read(address uint16) uint8 {
	return cartridge.rom[address]
}

func (cartridge *Cartridge) Write(address uint16, value uint8) {

}
