package apu

import "fmt"

type APU struct {
	// ======================================
	// ====== Global Control Registers ======
	// ======================================

	// 0xFF26 - NR52: Audio master control
	nr52 uint8
	// 0xFF25 - NR51: Sound panning
	nr51 uint8
	// 0xFF24 - NR50: Master volume & VIN panning
	nr50 uint8

	// =======================================================
	// ====== Sound Channel 1 — Pulse with period sweep ======
	// =======================================================

	// 0xFF10 - NR10: Channel 1 sweep
	nr10 uint8
	// 0xFF11 - NR11: Channel 1 length timer & duty cycle
	nr11 uint8
	// 0xFF12 - NR12: Channel 1 volume & envelope
	nr12 uint8
	// 0xFF13 - NR13: Channel 1 period low [write-only]
	nr13 uint8
	// 0xFF14 - NR14: Channel 1 period high & control
	nr14 uint8

	// =====================================
	// ====== Sound Channel 2 — Pulse ======
	// =====================================

	// 0xFF16 - NR21: Channel 2 length timer & duty cycle
	nr21 uint8
	// 0xFF17 - NR22: Channel 2 volume & envelope
	nr22 uint8
	// 0xFF18 - NR23: Channel 2 period low [write-only]
	nr23 uint8
	// 0xFF19 - NR24: Channel 2 period high & control
	nr24 uint8

	// ===========================================
	// ====== Sound Channel 3 — Wave output ======
	// ===========================================

	// 0xFF1A - NR30: Channel 3 DAC enable
	nr30 uint8
	// 0xFF1B - NR31: Channel 3 length timer [write-only]
	nr31 uint8
	// 0xFF1C - NR32: Channel 3 output level
	nr32 uint8
	// 0xFF1D - NR33: Channel 3 period low [write-only]
	nr33 uint8
	// 0xFF1E - NR34: Channel 3 period high & control
	nr34 uint8

	// 0xFF30-0xFF3F - Wave pattern RAM
	waveRam [16]uint8

	// =====================================
	// ====== Sound Channel 4 — Noise ======
	// =====================================

	// 0xFF20 - NR41: Channel 4 length timer [write-only]
	nr41 uint8
	// 0xFF21 - NR42: Channel 4 volume & envelope
	nr42 uint8
	// 0xFF22 - NR43: Channel 4 frequency & randomness
	nr43 uint8
	// 0xFF23 - NR44: Channel 4 control
	nr44 uint8

	// ===================
	// ====== State ======
	// ===================

	// wave position: 3 bits (0-7) because waveforms are 8 samples long

	internalTimer uint8
	sampleTimer   uint8
	outputBuffer  []int16

	// 0xFF04 - Timer DIV register
	// See: https://gbdev.io/pandocs/Audio_details.html#div-apu
	divApuCounter uint16

	// the step counter (0-7) to track which events to fire:
	// Event           Rate    Frequency
	// ----------------------------------
	// Envelope sweep  8       64 Hz
	// Sound length    2       256 Hz
	// CH1 freq sweep  4       128 Hz
	divApuStep uint8

	ch1Enabled       bool
	ch1PeriodDivider uint16
	ch1WavePosition  uint8
	ch1OutputBit     uint8
	ch1LengthTimer   uint8

	ch2Enabled       bool
	ch2PeriodDivider uint16
	ch2WavePosition  uint8
	ch2OutputBit     uint8
	ch2LengthTimer   uint8
}

func New() *APU {
	apu := &APU{}

	apu.Reset()

	return apu
}

func (apu *APU) Reset() {
	apu.nr10 = 0x80
	apu.nr11 = 0xBF
	apu.nr12 = 0xF3
	apu.nr13 = 0xFF
	apu.nr14 = 0xBF
	apu.nr21 = 0x3F
	apu.nr22 = 0x00
	apu.nr23 = 0xFF
	apu.nr24 = 0xBF
	apu.nr30 = 0x7F
	apu.nr31 = 0xFF
	apu.nr32 = 0x9F
	apu.nr33 = 0xFF
	apu.nr34 = 0xBF
	apu.nr41 = 0xFF
	apu.nr42 = 0x00
	apu.nr43 = 0x00
	apu.nr44 = 0xBF
	apu.nr50 = 0x77
	apu.nr51 = 0xF3
	apu.nr52 = 0xF1
}

// See: https://gbdev.gg8.se/wiki/articles/Gameboy_sound_hardware#Square_Wave
var dutyTable = [4][8]uint8{
	{0, 0, 0, 0, 0, 0, 0, 1}, // Duty 0 (12.5%)
	{1, 0, 0, 0, 0, 0, 0, 1}, // Duty 1 (25%)
	{1, 0, 0, 0, 0, 1, 1, 1}, // Duty 2 (50%)
	{0, 1, 1, 1, 1, 1, 1, 0}, // Duty 3 (75%)
}

// CPU clock speed / audio sampling rate
const downsampledRate = uint8(4194304 / 48000)

// See: https://gbdev.io/pandocs/Audio.html#length-timer
const MaxLengthTimer = uint8(64)

// Step performs 1 T-cycle of work
func (apu *APU) Step() {
	if !apu.poweredOn() {
		return
	}

	apu.divApuCounter++
	apu.internalTimer++
	apu.sampleTimer++

	// Runs at 512 Hz (4194304 / 512 = 8192)
	// See: https://gbdev.io/pandocs/Audio_details.html#div-apu
	if apu.divApuCounter == 8192 {
		apu.divApuCounter = 0

		apu.divApuStep++
		if apu.divApuStep > 7 {
			apu.divApuStep = 0
		}

		// Length runs at 256 Hz
		if apu.divApuStep%2 == 0 {
			lengthEnabled := (apu.nr24 & 0b0100_0000) != 0

			if lengthEnabled && apu.ch2LengthTimer < MaxLengthTimer {
				apu.ch2LengthTimer++

				if apu.ch2LengthTimer == MaxLengthTimer {
					apu.ch2Enabled = false
				}
			}
		}
	}

	if apu.internalTimer == 4 {
		apu.ch2PeriodDivider++
		apu.internalTimer = 0
	}

	// overflow check
	if apu.ch2PeriodDivider == 0b1_0000_0000_0000 {
		apu.ch2PeriodDivider = (uint16(apu.nr24&0b111) << 8) | uint16(apu.nr23)
		apu.ch2WavePosition++

		if apu.ch2WavePosition > 0b111 {
			apu.ch2WavePosition = 0
		}

		dutyType := (apu.nr21 & 0b1100_0000) >> 6
		apu.ch2OutputBit = dutyTable[dutyType][apu.ch2WavePosition&0b111]
	}

	if apu.sampleTimer == downsampledRate {
		apu.sampleTimer = 0

		var sample int16 = 0

		if apu.ch2Enabled {
			volume := (apu.nr22 & 0b1111_0000) >> 4
			// divide by 15.0 to normalize, then multiply by max int16 value
			// "The digital value produced by the generator, which ranges between $0 and $F (0 and 15)"
			// See: https://gbdev.io/pandocs/Audio_details.html#audio-details
			amplitude := int16((float32(volume) / 15.0) * 32767.0)
			sample = int16(apu.ch2OutputBit) * amplitude
		}

		apu.outputBuffer = append(apu.outputBuffer, sample)
	}

}

func (apu *APU) Read(address uint16) uint8 {
	switch {
	case address == 0xFF10:
		return apu.nr10 | 0b1000_0000
	case address == 0xFF11:
		return apu.nr11 | 0b0011_1111
	case address == 0xFF12:
		return apu.nr12
	case address == 0xFF13:
		return 0xFF
	case address == 0xFF14:
		return apu.nr14 | 0b1011_1111
	case address == 0xFF16:
		return apu.nr21 | 0b0011_1111
	case address == 0xFF17:
		return apu.nr22
	case address == 0xFF18:
		return 0xFF
	case address == 0xFF19:
		return apu.nr24 | 0b1011_1111
	case address == 0xFF1A:
		return apu.nr30 | 0b0111_1111
	case address == 0xFF1B:
		return 0xFF
	case address == 0xFF1C:
		return apu.nr32 | 0b1001_1111
	case address == 0xFF1D:
		return 0xFF
	case address == 0xFF1E:
		return apu.nr34 | 0b1011_1111
	case address == 0xFF20:
		return 0xFF
	case address == 0xFF21:
		return apu.nr42
	case address == 0xFF22:
		return apu.nr43
	case address == 0xFF23:
		return apu.nr44 | 0b1011_1111
	case address == 0xFF24:
		return apu.nr50
	case address == 0xFF25:
		return apu.nr51
	case address == 0xFF26:
		return apu.nr52 | 0b0111_0000
	case address >= 0xFF30 && address <= 0xFF3F:
		return apu.waveRam[address-0xFF30]
	}

	return 0xFF
}

func (apu *APU) Write(address uint16, value uint8) {
	// Only Wave RAM and NR52 are writable when APU is powered off
	if !apu.poweredOn() && address != 0xFF26 && !(address >= 0xFF30 && address <= 0xFF3F) {
		return
	}

	switch {
	case address == 0xFF10:
		apu.nr10 = value
	case address == 0xFF11:
		apu.nr11 = value

		apu.ch1LengthTimer = value & 0b0011_1111
		fmt.Println("APU: Channel 1 length timer set to", apu.ch1LengthTimer)
	case address == 0xFF12:
		apu.nr12 = value
	case address == 0xFF13:
		apu.nr13 = value
	case address == 0xFF14:
		apu.nr14 = value

		if (value & 0b1000_0000) != 0 {
			fmt.Println("APU: Channel 1 enabled")
			apu.ch1Enabled = true

			if apu.ch1LengthTimer == MaxLengthTimer {
				apu.ch1LengthTimer = 0
			}
		}
	case address == 0xFF16:
		apu.nr21 = value

		apu.ch2LengthTimer = value & 0b0011_1111
		fmt.Println("APU: Channel 2 length timer set to", apu.ch2LengthTimer)
	case address == 0xFF17:
		apu.nr22 = value
	case address == 0xFF18:
		apu.nr23 = value
	case address == 0xFF19:
		apu.nr24 = value

		if (value & 0b1000_0000) != 0 {
			fmt.Println("APU: Channel 2 enabled")
			apu.ch2Enabled = true
			apu.ch2PeriodDivider = (uint16(apu.nr24&0b111) << 8) | uint16(apu.nr23)

			if apu.ch2LengthTimer == MaxLengthTimer {
				apu.ch2LengthTimer = 0
			}
		}
	case address == 0xFF1A:
		apu.nr30 = value
	case address == 0xFF1B:
		apu.nr31 = value
	case address == 0xFF1C:
		apu.nr32 = value
	case address == 0xFF1D:
		apu.nr33 = value
	case address == 0xFF1E:
		apu.nr34 = value
	case address == 0xFF20:
		apu.nr41 = value
	case address == 0xFF21:
		apu.nr42 = value
	case address == 0xFF22:
		apu.nr43 = value
	case address == 0xFF23:
		apu.nr44 = value
	case address == 0xFF24:
		apu.nr50 = value
	case address == 0xFF25:
		apu.nr51 = value
	case address == 0xFF26:
		wasPoweredOn := apu.poweredOn()
		apu.nr52 = (apu.nr52 & 0b0111_1111) | (value & 0b1000_0000)
		isPoweredOn := apu.poweredOn()

		// APU ON -> APU OFF
		if wasPoweredOn && !isPoweredOn {
			powerBit := apu.nr52 & 0b1000_0000
			apu.nr52 = powerBit

			// Set all registers to 0x00
			apu.nr10 = 0x00
			apu.nr11 = 0x00
			apu.nr12 = 0x00
			apu.nr13 = 0x00
			apu.nr14 = 0x00
			apu.nr21 = 0x00
			apu.nr22 = 0x00
			apu.nr23 = 0x00
			apu.nr24 = 0x00
			apu.nr30 = 0x00
			apu.nr31 = 0x00
			apu.nr32 = 0x00
			apu.nr33 = 0x00
			apu.nr34 = 0x00
			apu.nr41 = 0x00
			apu.nr42 = 0x00
			apu.nr43 = 0x00
			apu.nr44 = 0x00
			apu.nr50 = 0x00
			apu.nr51 = 0x00
		}
		// APU OFF -> APU ON
		if !wasPoweredOn && isPoweredOn {
		}
	case address >= 0xFF30 && address <= 0xFF3F:
		apu.waveRam[address-0xFF30] = value
	}
}

func (apu *APU) poweredOn() bool {
	return (apu.nr52 & 0b1000_0000) != 0
}

// OnDivReset is called when the DIV register (0xFF04) is reset, to keep the
// DIV-APU in sync because they are physically the same counter.
func (apu *APU) OnDivReset() {
	apu.divApuCounter = 0
}

// ConsumeAudioBuffer returns the current audio buffer and clears it.
func (apu *APU) ConsumeAudioBuffer() []int16 {
	if (len(apu.outputBuffer)) == 0 {
		return nil
	}

	result := make([]int16, len(apu.outputBuffer))
	copy(result, apu.outputBuffer)

	apu.outputBuffer = apu.outputBuffer[:0]

	return result
}
