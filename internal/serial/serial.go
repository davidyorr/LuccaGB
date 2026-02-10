package serial

import (
	"encoding/binary"
	"fmt"

	"github.com/davidyorr/LuccaGB/internal/logger"
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

// Perform 1 T-cycle of work
func (serial *Serial) Step() (requestInterrupt bool) {
	if !serial.transferInProgress {
		return false
	}

	// Bit 0 on SC is the Clock select bit
	// 0 = External clock, 1 = Internal clock
	// If it is 0, the Game Boy is waiting for a signal from another device.
	if (serial.sc & 0b0000_0001) == 0 {
		return false
	}

	serial.transferCyclesRemaining--
	requestInterrupt = false

	if serial.transferCyclesRemaining == 0 {
		serial.transferInProgress = false
		serial.transferCyclesRemaining = 0
		serial.sb = 0xFF
		// clear the transfer bit
		serial.sc &^= 0b1000_0000
		requestInterrupt = true
	}

	return requestInterrupt
}

func (serial *Serial) Read(address uint16) uint8 {
	// SB
	if address == 0xFF01 {
		return serial.sb
	}
	// SC
	if address == 0xFF02 {
		return serial.sc | 0b0111_1110
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
			// speed for DMG (4194304 / 8192)
			serial.transferCyclesRemaining = 512
			serial.serialOutputBuffer = append(serial.serialOutputBuffer, serial.sb)
		}
	}
}

func (serial *Serial) SerialOutputBuffer() []uint8 {
	return serial.serialOutputBuffer
}

func (serial *Serial) Serialize(buf []byte) int {
	offset := 0

	buf[offset] = serial.sb
	offset++
	buf[offset] = serial.sc
	offset++

	length := uint32(len(serial.serialOutputBuffer))
	binary.LittleEndian.PutUint32(buf[offset:], length)
	offset += 4

	n := copy(buf[offset:], serial.serialOutputBuffer)
	offset += n

	if serial.transferInProgress {
		buf[offset] = 1
	} else {
		buf[offset] = 0
	}
	offset++
	binary.LittleEndian.PutUint16(buf[offset:], serial.transferCyclesRemaining)
	offset += 2

	return offset
}

func (serial *Serial) Deserialize(buf []byte) int {
	offset := 0

	serial.sb = buf[offset]
	offset++
	serial.sc = buf[offset]
	offset++

	length := int(binary.LittleEndian.Uint32(buf[offset:]))
	offset += 4

	if cap(serial.serialOutputBuffer) < int(length) {
		serial.serialOutputBuffer = make([]uint8, length)
	} else {
		serial.serialOutputBuffer = serial.serialOutputBuffer[:length]
	}

	if length > 0 {
		n := copy(serial.serialOutputBuffer, buf[offset:])
		offset += n
	}

	if offset < len(buf) {
		serial.transferInProgress = buf[offset] == 1
		offset++
	}
	if offset+2 < len(buf) {
		serial.transferCyclesRemaining = binary.LittleEndian.Uint16(buf[offset:])
		offset += 2
	}

	return offset
}
