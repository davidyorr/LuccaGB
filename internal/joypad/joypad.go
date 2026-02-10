package joypad

import (
	"github.com/davidyorr/LuccaGB/internal/interrupt"
	"github.com/davidyorr/LuccaGB/internal/logger"
)

type Joypad struct {
	// P1/JOYP: stores the value written to 0xFF00
	p1Register uint8

	// Internal state: 1 = Pressed, 0 = Released
	buttons uint8
	dpad    uint8

	interruptRequester func(interruptType interrupt.Interrupt)
}

type JoypadInput uint8

const (
	JoypadInputStart JoypadInput = iota
	JoypadInputSelect
	JoypadInputB
	JoypadInputA
	JoypadInputDown
	JoypadInputUp
	JoypadInputLeft
	JoypadInputRight
)

var joypadInputMask = map[JoypadInput]uint8{
	JoypadInputStart:  0b1000,
	JoypadInputSelect: 0b0100,
	JoypadInputB:      0b0010,
	JoypadInputA:      0b0001,
	JoypadInputDown:   0b1000,
	JoypadInputUp:     0b0100,
	JoypadInputLeft:   0b0010,
	JoypadInputRight:  0b0001,
}

func New(interruptRequest func(interrupt.Interrupt)) *Joypad {
	joypad := &Joypad{}
	joypad.interruptRequester = interruptRequest
	joypad.Reset()

	return joypad
}

func (joypad *Joypad) Reset() {
	joypad.buttons = 0
	joypad.dpad = 0
	joypad.p1Register = 0xCF
}

func (joypad *Joypad) Write(value uint8) {
	logger.GlobalTraceLogger.LogMemWrite(0xFF00, value)

	oldState := joypad.calculateP1Register()

	// only update the "Select" bits (4-5)
	joypad.p1Register = value & 0b0011_0000

	newState := joypad.calculateP1Register()

	joypad.checkInterrupt(oldState, newState)
}

func (joypad *Joypad) Read() uint8 {
	value := joypad.calculateP1Register()
	logger.GlobalTraceLogger.LogMemRead(0xFF00, value)

	return value
}

func (joypad *Joypad) Press(input JoypadInput) {
	oldState := joypad.calculateP1Register()

	mask := joypadInputMask[input]
	if input.isDpad() {
		joypad.dpad |= mask
	} else {
		joypad.buttons |= mask
	}

	newState := joypad.calculateP1Register()

	joypad.checkInterrupt(oldState, newState)
}

func (joypad *Joypad) Release(input JoypadInput) {
	oldState := joypad.calculateP1Register()

	mask := joypadInputMask[input]
	if input.isDpad() {
		joypad.dpad &^= mask
	} else {
		joypad.buttons &^= mask
	}

	newState := joypad.calculateP1Register()

	joypad.checkInterrupt(oldState, newState)
}

// calculateP1Register calculates the values of the lower input bits (0-3) based
// on the state of the internal variables `buttons` and `dpad`
func (joypad *Joypad) calculateP1Register() uint8 {
	// Force bits 6-7 to 1.
	value := uint8(0b1100_0000)
	// Keep the "Select" bits (4-5) from the register.
	value |= (joypad.p1Register & 0b0011_0000)
	// Start with bits 0-3 as 1 (Released).
	value |= 0b0000_1111

	// Reading Start, Select, B, A
	if (joypad.p1Register & 0b0010_0000) == 0 {
		value &= (^joypad.buttons & 0b0000_1111) | 0b1111_0000
	}

	// Reading Down, Up, Left, Right
	if (joypad.p1Register & 0b0001_0000) == 0 {
		dpadState := ^joypad.dpad & 0b0000_1111

		// == impossible input sanitization ==
		// prevent Left+Right
		if (dpadState & 0b0011) == 0 {
			dpadState |= 0b0011 // release both
		}
		// prevent Up+Down
		if (dpadState & 0b1100) == 0 {
			dpadState |= 0b1100 // release both
		}

		value &= (dpadState | 0b1111_0000)
	}

	return value
}

// The Joypad interrupt is requested when any of P1 bits 0-3 change from High to
// Low. This happens when a button is pressed (provided that the
// action/direction buttons are enabled by bit 5/4, respectively), however, due
// to switch bounce, one or more High to Low transitions are usually produced
// when pressing a button.
// See: https://gbdev.io/pandocs/Interrupt_Sources.html#int-60--joypad-interrupt
func (joypad *Joypad) checkInterrupt(oldState, newState uint8) {
	// (Old was High) AND (New is Low)
	if ((oldState & 0b0000_1111) & ^(newState & 0b0000_1111)) != 0 {
		joypad.interruptRequester(interrupt.JoypadInterrupt)
	}
}

func (input JoypadInput) isDpad() bool {
	return input == JoypadInputDown || input == JoypadInputUp || input == JoypadInputLeft || input == JoypadInputRight
}

func (joypad *Joypad) Serialize(buf []byte) int {
	offset := 0

	buf[offset] = joypad.p1Register
	offset++
	buf[offset] = joypad.buttons
	offset++
	buf[offset] = joypad.dpad
	offset++

	return offset
}

func (joypad *Joypad) Deserialize(buf []byte) int {
	offset := 0

	joypad.p1Register = buf[offset]
	offset++
	joypad.buttons = buf[offset]
	offset++
	joypad.dpad = buf[offset]
	offset++

	return offset
}
