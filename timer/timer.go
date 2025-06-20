package timer

type Timer struct {
	// 0xFF04 divider register
	div uint8
	// 0xFF05 timer counter
	tima uint8
	// 0xFF06 timer modulo
	tma uint8
	// 0xFF07 timer control
	tac uint8

	divCounter  uint16
	timaCounter uint16
}

func New() *Timer {
	timer := &Timer{}

	return timer
}

func (timer *Timer) Reset() {
	timer.div = 0x18
	timer.divCounter = 0
	timer.tima = 0x00
	timer.tma = 0x00
	timer.tac = 0xF8
}

// 4194304 Hz / 16384 Hz
const threshold = 256

// return value of true represents a signal interrupt request
func (timer *Timer) Step(cycles uint8) bool {
	cycles16 := uint16(cycles)
	timer.divCounter += cycles16
	timer.timaCounter += cycles16

	if timer.divCounter >= threshold {
		timer.div++
		timer.divCounter -= threshold
	}

	// bit 2 -> timer enabled
	if (timer.tac & 0b0000_0100) != 0 {
		var frequency uint16 = 1024
		switch timer.tac & 0b11 {
		case 0b00:
			frequency = 1024
		case 0b01:
			frequency = 16
		case 0b10:
			frequency = 64
		case 0b11:
			frequency = 256
		}

		if timer.timaCounter >= frequency {
			timer.timaCounter -= frequency
			timer.tima++
			// check for overflow
			if timer.tima == 0 {
				timer.tima = timer.tma
				return true
			}
		}
	}

	return false
}

func (timer *Timer) Read(address uint16) uint8 {
	switch address & 0x000F {
	case 0x04:
		return timer.div
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
		timer.div = 0
		timer.divCounter = 0
	case 0x05:
		timer.tima = value
	case 0x06:
		timer.tma = value
	case 0x07:
		// upper 5 bits are unused
		timer.tac = (value & 0b0000_0111)
	}
}
