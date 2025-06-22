package timer

type Timer struct {
	// 0xFF05 timer counter
	tima uint8
	// 0xFF06 timer modulo
	tma uint8
	// 0xFF07 timer control
	tac uint8

	counter               uint16
	previousTimerBitState bool
}

func New() *Timer {
	timer := &Timer{}

	return timer
}

func (timer *Timer) Reset() {
	timer.tima = 0x00
	timer.tma = 0x00
	timer.tac = 0xF8
	timer.counter = 0xABCC
	timer.previousTimerBitState = false
}

// return value of true represents a signal interrupt request
func (timer *Timer) Step(cycles uint8) bool {
	var interrupt bool = false
	cycles16 := uint16(cycles)

	for cycles16 > 0 {
		if timer.isTimerEnabled() {
			currentTimerBitState := timer.getTimerBitState()

			if timer.previousTimerBitState && !currentTimerBitState {
				timer.tima++
				// check for overflow
				if timer.tima == 0 {
					timer.tima = timer.tma
					interrupt = true
				}
			}

			timer.previousTimerBitState = currentTimerBitState
		}

		timer.counter++
		cycles16--
	}

	return interrupt

}

func (timer *Timer) Read(address uint16) uint8 {
	switch address & 0x000F {
	case 0x04:
		return uint8(timer.counter >> 8)
	case 0x05:
		return timer.tima
	case 0x06:
		return timer.tma
	case 0x07:
		// upper 5 bits are unused and should always be set
		return timer.tac | 0b1111_1000
	default:
		return 0
	}
}

func (timer *Timer) Write(address uint16, value uint8) {
	switch address & 0x000F {
	case 0x04:
		timer.counter = 0
	case 0x05:
		timer.tima = value
	case 0x06:
		timer.tma = value
	case 0x07:
		// upper 5 bits are unused
		timer.tac = (value & 0b0000_0111)
	}
}

func (timer *Timer) getTimerBitState() bool {
	// the bit index to check on the counter
	var bitIndex uint8 = 9

	switch timer.tac & 0b11 {
	case 0b00:
		bitIndex = 9
	case 0b01:
		bitIndex = 3
	case 0b10:
		bitIndex = 5
	case 0b11:
		bitIndex = 7
	}

	mask := uint16(1 << bitIndex)

	return (timer.counter & mask) != 0
}

func (timer *Timer) isTimerEnabled() bool {
	// bit 2 -> timer enabled
	return (timer.tac & 0b0000_0100) != 0
}
