package serial

import (
	"fmt"

	"github.com/davidyorr/EchoGB/logger"
)

type Serial struct {
	// 0xFF01 — SB: Serial transfer data
	sb uint8
	// 0xFF02 — SC: Serial transfer control
	sc uint8
	// output buffer
	serialOutputBuffer      []uint8
	transferInProgress      bool
	transferCyclesRemaining uint16
}

func New() *Serial {
	serial := &Serial{}

	serial.Reset()

	return serial
}

func (serial *Serial) Reset() {
	serial.sb = 0x00
	serial.sc = 0x7E
}

// return value of true represents a serial interrupt request
func (serial *Serial) Step(cycles uint8) bool {
	if !serial.transferInProgress {
		return false
	}

	serial.transferCyclesRemaining -= uint16(cycles)

	if serial.transferCyclesRemaining <= 0 {
		serial.transferInProgress = false
		serial.transferCyclesRemaining = 0
		serial.sb = 0xFF
		// clear the transfer bit
		serial.sc &^= 0b1000_0000
		return true
	}

	return false
}

func (serial *Serial) Read(address uint16) uint8 {
	if address == 0xFF01 {
		return serial.sb
	}
	if address == 0xFF02 {
		return serial.sc
	}

	logger.Error("SERIAL", fmt.Sprintf("invalid address on read: 0x%0X", address))
	return 0xFF
}

func (serial *Serial) Write(address uint16, value uint8) {
	if address == 0xFF01 {
		serial.sb = value
	} else if address == 0xFF02 {
		serial.sc = value

		// bit 7 is the transfer bit
		if (value & 0b1000_0000) != 0 {
			serial.transferInProgress = true
			// speed for DMG
			serial.transferCyclesRemaining = 8192
			serial.serialOutputBuffer = append(serial.serialOutputBuffer, serial.sb)
		}
	}
}

func (serial *Serial) SerialOutputBuffer() []uint8 {
	return serial.serialOutputBuffer
}
