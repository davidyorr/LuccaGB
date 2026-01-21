package apu

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
	sampleTimer   int
	outputBuffer  AudioBuffer

	// 0xFF04 - Timer DIV register
	// See: https://gbdev.io/pandocs/Audio_details.html#div-apu
	divApuCounter uint16

	cyclesSinceWaveRamFetch uint16

	// the step counter (0-7) to track which events to fire:
	//
	// Event           Rate    Frequency
	// ----------------------------------
	// Envelope sweep  8       64 Hz
	// Sound length    2       256 Hz
	// CH1 freq sweep  4       128 Hz
	//
	// 	Step   Length Ctr  Vol Env     Sweep
	// ---------------------------------------
	// 0      Clock       -           -
	// 1      -           -           -
	// 2      Clock       -           Clock
	// 3      -           -           -
	// 4      Clock       -           -
	// 5      -           -           -
	// 6      Clock       -           Clock
	// 7      -           Clock       -
	// ---------------------------------------
	// Rate   256 Hz      64 Hz       128 Hz
	//
	// See: https://gbdev.gg8.se/wiki/articles/Gameboy_sound_hardware#Frame_Sequencer
	divApuStep uint8

	ch1 channel
	ch2 channel
	ch3 channel
	ch4 channel

	// state controlled by the UI. 1-indexed to match the channel name.
	channelsEnabled [5]bool
}

type channel struct {
	// ====== Common Fields ======
	register         *uint8
	lengthTimer      uint16
	enabled          bool
	envelopeEnabled  bool
	currentVolume    uint8
	envelopeTimer    uint8
	envelopeRegister *uint8
	maxLength        uint16
	nr52BitMask      uint8
	dacRegister      *uint8
	dacRegisterMask  uint8
	outputBit        uint8
	periodDivider    uint16

	// ====== Not Common Fields ======
	wavePosition uint8

	// Ch 1 Sweep Unit
	sweepTimer          uint16
	sweepEnabled        bool
	sweepShadowRegister uint16 // output period
	sweepNegateModeUsed bool

	// Ch 3
	sampleIndex  uint8
	sampleBuffer uint8

	// Ch 4
	lfsr uint16
}

func New() *APU {
	apu := &APU{}
	apu.ch1 = channel{
		register:         &apu.nr14,
		maxLength:        MaxLengthTimer_Ch1Ch2Ch4,
		nr52BitMask:      0b0001,
		dacRegister:      &apu.nr12,
		dacRegisterMask:  0b1111_1000,
		envelopeRegister: &apu.nr12,
	}
	apu.ch2 = channel{
		register:         &apu.nr24,
		maxLength:        MaxLengthTimer_Ch1Ch2Ch4,
		nr52BitMask:      0b0010,
		dacRegister:      &apu.nr22,
		dacRegisterMask:  0b1111_1000,
		envelopeRegister: &apu.nr22,
	}
	apu.ch3 = channel{
		register:        &apu.nr34,
		maxLength:       MaxLengthTimer_Ch3,
		nr52BitMask:     0b0100,
		dacRegister:     &apu.nr30,
		dacRegisterMask: 0b1000_0000,
	}
	apu.ch4 = channel{
		register:         &apu.nr44,
		maxLength:        MaxLengthTimer_Ch1Ch2Ch4,
		nr52BitMask:      0b1000,
		dacRegister:      &apu.nr42,
		dacRegisterMask:  0b1111_1000,
		envelopeRegister: &apu.nr42,
	}
	apu.channelsEnabled = [5]bool{false, true, true, true, true}

	apu.cyclesSinceWaveRamFetch = 255

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

// how many bits to shift
// See: https://gbdev.gg8.se/wiki/articles/Gameboy_sound_hardware#Wave_Channel
var ch3OutputLevelMap = [4]uint8{
	4, // 0% (silent)
	0, // 100%
	1, // 50%
	2, // 25%
}

const CpuClockSpeed = 4_194_304
const TargetSampleRate = 48_000

// See: https://gbdev.io/pandocs/Audio.html#length-timer
const MaxLengthTimer_Ch1Ch2Ch4 = uint16(64)
const MaxLengthTimer_Ch3 = uint16(256)

// Step performs 1 T-cycle of work
func (apu *APU) Step() {
	if !apu.poweredOn() {
		return
	}

	apu.divApuCounter++
	apu.internalTimer++
	apu.sampleTimer += TargetSampleRate
	apu.cyclesSinceWaveRamFetch++

	ch1 := &apu.ch1
	ch2 := &apu.ch2
	ch3 := &apu.ch3
	ch4 := &apu.ch4

	// Runs at 512 Hz (4194304 / 512 = 8192)
	// See: https://gbdev.io/pandocs/Audio_details.html#div-apu
	if apu.divApuCounter == 8192 {
		apu.divApuCounter = 0

		apu.divApuStep++
		if apu.divApuStep > 7 {
			apu.divApuStep = 0
		}

		// Ch1 Freq Sweep runs at 128 Hz (2, 6)
		if apu.divApuStep&0b11 == 0b10 {
			if ch1.sweepTimer > 0 {
				ch1.sweepTimer--
			}

			if ch1.sweepTimer == 0 {
				apu.reloadCh1SweepTimer()

				pace := (apu.nr10 & 0b0111_0000) >> 4
				if ch1.sweepEnabled && pace != 0 {
					frequency := apu.calculateCh1Frequency()
					individualStep := apu.nr10 & 0b0000_0111

					if frequency <= 2047 && individualStep != 0 {
						ch1.sweepShadowRegister = uint16(frequency)
						apu.nr13 = uint8(frequency & 0xFF)
						apu.nr14 = (apu.nr14 & 0b1111_1000) | uint8((frequency>>8)&0b0111)
						ch1.periodDivider = uint16(frequency)
						apu.calculateCh1Frequency()
					}
				}
			}
		}

		// Length runs at 256 Hz (0, 2, 4, 6)
		if apu.divApuStep&1 == 0 {
			apu.clockLength(ch1)
			apu.clockLength(ch2)
			apu.clockLength(ch3)
			apu.clockLength(ch4)
		}

		// Ch 1, 2, and 4 Envelope Sweep runs at 64 Hz (7)
		if apu.divApuStep == 7 {
			apu.clockEnvelope(ch1)
			apu.clockEnvelope(ch2)
			apu.clockEnvelope(ch4)
		}
	}

	// Channels 1 & 2 (Pulse Channels) period dividers are clocked at 1048576 Hz, once per four dots
	// See: https://gbdev.io/pandocs/Audio_Registers.html#ff13--nr13-channel-1-period-low-write-only
	if (apu.internalTimer&0b11) == 0 && ch1.enabled {
		ch1.periodDivider++
	}

	if (apu.internalTimer&0b11) == 0 && ch2.enabled {
		ch2.periodDivider++
	}

	// Channel 3 (Wave Channel) period divider is clocked at 2097152 Hz, once per two dots
	// See: https://gbdev.io/pandocs/Audio_Registers.html#ff1d--nr33-channel-3-period-low-write-only
	if (apu.internalTimer&0b1) == 0 && ch3.enabled {
		ch3.periodDivider++
	}

	// Channel 4 (Noise) is clocked at 262144 Hz, once per 16 dots.
	// Shift being equal to 14 or 15 stops the channel from being clocked entirely.
	// See: https://gbdev.io/pandocs/Audio_Registers.html#ff22--nr43-channel-4-frequency--randomness
	clockShift := (apu.nr43 & 0b1111_0000) >> 4
	if (apu.internalTimer&0b1111) == 0 && (clockShift != 14 && clockShift != 15) && ch4.enabled {
		ch4.periodDivider--
	}

	// ch 1 overflow check
	if ch1.periodDivider == 0b1000_0000_0000 {
		ch1.periodDivider = (uint16(apu.nr14&0b111) << 8) | uint16(apu.nr13)
		ch1.wavePosition++

		if ch1.wavePosition > 0b111 {
			ch1.wavePosition = 0
		}

		dutyType := (apu.nr11 & 0b1100_0000) >> 6
		ch1.outputBit = dutyTable[dutyType][ch1.wavePosition&0b111]
	}

	// ch 2 overflow check
	if ch2.periodDivider == 0b1000_0000_0000 {
		ch2.periodDivider = (uint16(apu.nr24&0b111) << 8) | uint16(apu.nr23)
		ch2.wavePosition++

		if ch2.wavePosition > 0b111 {
			ch2.wavePosition = 0
		}

		dutyType := (apu.nr21 & 0b1100_0000) >> 6
		ch2.outputBit = dutyTable[dutyType][ch2.wavePosition&0b111]
	}

	// ch 3 overflow check
	if ch3.periodDivider == 0b1000_0000_0000 {
		apu.cyclesSinceWaveRamFetch = 0
		ch3.periodDivider = (uint16(apu.nr34&0b111) << 8) | uint16(apu.nr33)

		// wave RAM is 16 bytes long, we fetch a nibble at a time (half a byte),
		// so we limit to 31
		ch3.sampleIndex = (ch3.sampleIndex + 1) & 31

		byteIndex := ch3.sampleIndex >> 1
		sampleByte := apu.waveRam[byteIndex]

		if ch3.sampleIndex&1 == 0 {
			ch3.sampleBuffer = sampleByte >> 4
		} else {
			ch3.sampleBuffer = sampleByte & 0b0000_1111
		}
	}

	// ch 4 period divider counts down to 0
	if ch4.periodDivider == 0 {
		// See: https://gbdev.io/pandocs/Audio_details.html#noise-channel-ch4

		// reset timer
		ch4.periodDivider = apu.calculateCh4Period()

		// XNOR ("1 if bit 0 and bit 1 are identical")
		bit0 := ch4.lfsr & 0b01
		bit1 := (ch4.lfsr & 0b10) >> 1
		result := ^(bit0 ^ bit1) & 1

		// apply result to bit 15
		ch4.lfsr &= 0b0111_1111_1111_1111 // clear bit
		ch4.lfsr |= (result << 15)        // OR in the result

		// short mode (if width mode bit is 1, apply to bit 6)
		if (apu.nr43 & 0b0000_1000) != 0 {
			ch4.lfsr &= 0b1111_1111_0111_1111 // clear bit
			ch4.lfsr |= (result << 7)         // OR in the result
		}

		// shift LFSR right
		ch4.lfsr >>= 1

		// set output bit
		if (ch4.lfsr & 1) == 0 {
			ch4.outputBit = 1
		} else {
			ch4.outputBit = 0
		}
	}

	if apu.sampleTimer >= CpuClockSpeed {
		apu.sampleTimer -= CpuClockSpeed

		var leftAccumulator int32 = 0
		var rightAccumulator int32 = 0

		if ch1.enabled && apu.channelsEnabled[1] {
			volume := ch1.currentVolume
			sample := int32(ch1.outputBit) * int32(volume)
			sample -= int32(volume) / 2 // center to [-volume/2, +volume/2]

			if apu.nr51&0b0001_0000 != 0 {
				leftAccumulator += sample
			}
			if apu.nr51&0b0001 != 0 {
				rightAccumulator += sample
			}
		}

		if ch2.enabled && apu.channelsEnabled[2] {
			volume := ch2.currentVolume
			sample := int32(ch2.outputBit) * int32(volume)
			sample -= int32(volume) / 2 // center to [-volume/2, +volume/2]

			if apu.nr51&0b0010_0000 != 0 {
				leftAccumulator += sample
			}
			if apu.nr51&0b0010 != 0 {
				rightAccumulator += sample
			}
		}

		if ch3.enabled && apu.channelsEnabled[3] {
			shift := ch3OutputLevelMap[ch3.currentVolume]
			sample := int32(ch3.sampleBuffer >> shift)
			sample -= (15 >> shift) / 2 // center to [-max/2, +max/2]

			if apu.nr51&0b0100_0000 != 0 {
				leftAccumulator += sample
			}
			if apu.nr51&0b0100 != 0 {
				rightAccumulator += sample
			}
		}

		if ch4.enabled && apu.channelsEnabled[4] {
			volume := ch4.currentVolume
			sample := int32(ch4.outputBit) * int32(volume)
			sample -= int32(volume) / 2 // center to [-volume/2, +volume/2]

			if apu.nr51&0b1000_0000 != 0 {
				leftAccumulator += sample
			}
			if apu.nr51&0b1000 != 0 {
				rightAccumulator += sample
			}
		}

		// "The digital value produced by the generator, which ranges between $0 and $F (0 and 15)"
		// See: https://gbdev.io/pandocs/Audio_details.html#audio-details

		// Normalize the result
		// We have 4 channels, each capable of outputting 0-15.
		// The max accumulator = 15 * 4 = 60.
		// So we want to map the range [0, 60] to [-32768, 32767] (for int16).
		leftSample := int16((float32(leftAccumulator) / 60.0) * 32767.0)
		rightSample := int16((float32(rightAccumulator) / 60.0) * 32767.0)

		apu.outputBuffer.Write(leftSample)
		apu.outputBuffer.Write(rightSample)
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
		// If the wave channel is enabled, accessing any byte from $FF30-$FF3F
		// is equivalent to accessing the current byte selected by the waveform
		// position. Further, on the DMG accesses will only work in this manner
		// if made within a couple of clocks of the wave channel accessing wave
		// RAM; if made at any other time, reads return $FF and writes have no
		// effect.
		// See: https://gbdev.gg8.se/wiki/articles/Gameboy_sound_hardware#Obscure_Behavior
		if apu.ch3.enabled {
			if apu.cyclesSinceWaveRamFetch < 2 {
				return apu.waveRam[apu.ch3.sampleIndex>>1]
			}
			return 0xFF
		}
		return apu.waveRam[address-0xFF30]
	}

	return 0xFF
}

func (apu *APU) Write(address uint16, value uint8) {
	isPoweredOn := apu.poweredOn()
	isNR52 := address == 0xFF26
	isWaveRam := address >= 0xFF30 && address <= 0xFF3F
	// On DMG the length timer bits can still be written to while powered off
	// See: https://gbdev.gg8.se/wiki/articles/Gameboy_sound_hardware#Differences
	isLengthTimerRegister := address == 0xFF11 || address == 0xFF16 || address == 0xFF1B || address == 0xFF20

	// Only Wave RAM, NR52, and length timer bits are writable when APU is powered off.
	if !isPoweredOn && !isNR52 && !isWaveRam && !isLengthTimerRegister {
		return
	}

	switch {
	case address == 0xFF10:
		// Clearing the sweep negate mode bit in NR10 after at least one sweep
		// calculation has been made using the negate mode since the last
		// trigger causes the channel to be immediately disabled. This prevents
		// you from having the sweep lower the frequency then raise the
		// frequency without a trigger inbetween.
		// See: https://gbdev.gg8.se/wiki/articles/Gameboy_sound_hardware#Obscure_Behavior
		prevDirection := (apu.nr10 & 0b000_1000) >> 3
		newDirection := (value & 0b000_1000) >> 3

		// if we have used negation since the last trigger
		if apu.ch1.sweepNegateModeUsed {
			// and we are switching from subtraction to addition
			if prevDirection == 1 && newDirection == 0 {
				// then disable the channel
				apu.ch1.enabled = false
				apu.nr52 &^= 0b0001
			}
		}

		apu.nr10 = value
	case address == 0xFF11:
		if isPoweredOn {
			apu.nr11 = value
		} else {
			// On DMG we can modify the length timer bits, but must preserve the wave duty bits
			// See: https://gbdev.gg8.se/wiki/articles/Gameboy_sound_hardware#Differences
			apu.nr11 = value & 0b0011_1111
		}

		apu.ch1.lengthTimer = uint16(value & 0b0011_1111)
	case address == 0xFF12:
		apu.nr12 = value

		// Ch 1 DAC disabled
		if (value & 0b1111_1000) == 0 {
			apu.ch1.enabled = false
			apu.nr52 &^= 0b0001
		}
	case address == 0xFF13:
		apu.nr13 = value
	case address == 0xFF14:
		apu.writeNRx4(address, value)

		// Trigger bit is set
		if (value & 0b1000_0000) != 0 {
			apu.ch1.periodDivider = (uint16(apu.nr14&0b111) << 8) | uint16(apu.nr13)

			// sweep initialization
			apu.ch1.sweepNegateModeUsed = false
			pace := (apu.nr10 & 0b0111_0000) >> 4
			individualStep := apu.nr10 & 0b0000_0111

			apu.ch1.sweepShadowRegister = (uint16(apu.nr14&0b111) << 8) | uint16(apu.nr13)
			apu.ch1.sweepEnabled = pace != 0 || individualStep != 0

			apu.reloadCh1SweepTimer()

			if individualStep != 0 {
				apu.calculateCh1Frequency()
			}
		}
	case address == 0xFF16:
		if isPoweredOn {
			apu.nr21 = value
		} else {
			// On DMG we can modify the length timer bits, but must preserve the wave duty bits
			// See: https://gbdev.gg8.se/wiki/articles/Gameboy_sound_hardware#Differences
			apu.nr21 = value & 0b0011_1111
		}

		apu.ch2.lengthTimer = uint16(value & 0b0011_1111)
	case address == 0xFF17:
		apu.nr22 = value

		// Ch 2 DAC disabled
		if (value & 0b1111_1000) == 0 {
			apu.ch2.enabled = false
			apu.nr52 &^= 0b0010
		}
	case address == 0xFF18:
		apu.nr23 = value
	case address == 0xFF19:
		apu.writeNRx4(address, value)

		// Trigger bit is set
		if (value & 0b1000_0000) != 0 {
			apu.ch2.periodDivider = (uint16(apu.nr24&0b111) << 8) | uint16(apu.nr23)
		}
	case address == 0xFF1A:
		apu.nr30 = value

		// Ch 3 DAC disabled
		if (value & 0b1000_0000) == 0 {
			apu.ch3.enabled = false
			apu.nr52 &^= 0b0000_0100
		}
	case address == 0xFF1B:
		apu.nr31 = value

		apu.ch3.lengthTimer = uint16(value)
	case address == 0xFF1C:
		apu.nr32 = value

		apu.ch3.currentVolume = (value & 0b0110_0000) >> 6
	case address == 0xFF1D:
		apu.nr33 = value
	case address == 0xFF1E:
		apu.writeNRx4(address, value)

		// Trigger bit is set
		if (value & 0b1000_0000) != 0 {
			// Triggering the wave channel on the DMG while it reads a sample
			// byte will alter the first four bytes of wave RAM.
			// See: https://gbdev.gg8.se/wiki/articles/Gameboy_sound_hardware#Obscure_Behavior
			// Each wave RAM fetch takes 4 cycles, so when cyclesSinceWaveRamFetch == 2
			// we are in the middle of a fetch. 2 cycles must be when the
			// actual data is being latched.
			if apu.ch3.enabled && apu.cyclesSinceWaveRamFetch == 2 {
				// + 1 to use the sample index that is about to be read, because
				// we haven't actually read the nibble yet
				byteIndex := ((apu.ch3.sampleIndex + 1) >> 1) & 0xF

				if byteIndex < 4 {
					apu.waveRam[0] = apu.waveRam[byteIndex]
				} else {
					// Copy the aligned block of 4 bytes to the start of wave RAM
					srcStart := byteIndex & 0b1100
					copy(apu.waveRam[0:4], apu.waveRam[srcStart:srcStart+4])
				}
			}

			// Delay by 3 extra clocks.
			// I can't find any docs that mention this, but the Blargg
			// 09-wave_read_while_on.gb test fails without this delay.
			apu.ch3.periodDivider = (uint16(apu.nr34&0b111) << 8) | uint16(apu.nr33) - 3
			apu.ch3.currentVolume = (apu.nr32 & 0b0110_0000) >> 6

			apu.ch3.sampleIndex = 0
		}
	case address == 0xFF20:
		apu.nr41 = (value & 0b0011_1111)

		apu.ch4.lengthTimer = uint16(value & 0b0011_1111)
	case address == 0xFF21:
		apu.nr42 = value

		// Ch 4 DAC disabled
		if (value & 0b1111_1000) == 0 {
			apu.ch4.enabled = false
			apu.nr52 &^= 0b1000
		}
	case address == 0xFF22:
		apu.nr43 = value
	case address == 0xFF23:
		apu.writeNRx4(address, value)

		// Trigger bit is set
		if (value & 0b1000_0000) != 0 {
			apu.ch4.lfsr = 0
			apu.ch4.periodDivider = apu.calculateCh4Period()
		}
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
			apu.divApuCounter = 0
			// set to 7 so that the first tick wraps it to 0
			apu.divApuStep = 7
		}
	case address >= 0xFF30 && address <= 0xFF3F:
		// If the wave channel is enabled, accessing any byte from $FF30-$FF3F
		// is equivalent to accessing the current byte selected by the waveform
		// position. Further, on the DMG accesses will only work in this manner
		// if made within a couple of clocks of the wave channel accessing wave
		// RAM; if made at any other time, reads return $FF and writes have no
		// effect.
		// See: https://gbdev.gg8.se/wiki/articles/Gameboy_sound_hardware#Obscure_Behavior
		if apu.ch3.enabled {
			if apu.cyclesSinceWaveRamFetch < 2 {
				apu.waveRam[apu.ch3.sampleIndex>>1] = value
			}
			return
		}
		apu.waveRam[address-0xFF30] = value
	}
}

func (apu *APU) poweredOn() bool {
	return (apu.nr52 & 0b1000_0000) != 0
}

func (apu *APU) writeNRx4(address uint16, value uint8) {
	var ch *channel
	switch address {
	case 0xFF14:
		ch = &apu.ch1
	case 0xFF19:
		ch = &apu.ch2
	case 0xFF1E:
		ch = &apu.ch3
	case 0xFF23:
		ch = &apu.ch4
	}

	// See: https://gbdev.io/pandocs/Audio_details.html#obscure-behavior
	triggerBit := (value & 0b1000_0000) != 0
	lengthTimerWasDisabled := (*ch.register & 0b0100_0000) == 0
	lengthTimerIsEnabled := (value & 0b0100_0000) != 0

	// do the actual write
	*ch.register = value

	nextStepDoesNotClock := (apu.divApuStep & 1) == 0

	if nextStepDoesNotClock && lengthTimerWasDisabled && lengthTimerIsEnabled {
		if ch.lengthTimer < ch.maxLength {
			ch.lengthTimer++

			if ch.lengthTimer == ch.maxLength {
				if !triggerBit {
					ch.enabled = false
					apu.nr52 &^= ch.nr52BitMask
				}
			}
		}
	}

	if triggerBit {
		if ch.lengthTimer == ch.maxLength {
			ch.lengthTimer = 0

			if nextStepDoesNotClock && lengthTimerIsEnabled {
				ch.lengthTimer++
			}
		}

		// Ch 1, 2, and 4 Envelope initialization
		if ch == &apu.ch1 || ch == &apu.ch2 || ch == &apu.ch4 {
			ch.envelopeEnabled = true
			ch.currentVolume = (*ch.envelopeRegister & 0b1111_0000) >> 4
			ch.envelopeTimer = (*ch.envelopeRegister & 0b0000_0111)

			if ch.envelopeTimer == 0 {
				ch.envelopeTimer = 8 // to reload
			}
		}

		dacEnabled := (*ch.dacRegister & ch.dacRegisterMask) != 0
		if dacEnabled {
			ch.enabled = true
			apu.nr52 |= ch.nr52BitMask
		}
	}
}

// calculateCh1Frequency consists of taking the value in the frequency “shadow
// register”, shifting it right by the individual step, optionally negating the
// value (depending on the direction) and summing this with the frequency
// “shadow register” to produce a new frequency.
// See: https://gbdev.io/pandocs/Audio_details.html#pulse-channel-with-sweep-ch1
func (apu *APU) calculateCh1Frequency() int {
	direction := (apu.nr10 & 0b0000_1000) >> 3
	individualStep := apu.nr10 & 0b0000_0111

	frequency := int(apu.ch1.sweepShadowRegister)
	delta := int(apu.ch1.sweepShadowRegister >> individualStep)
	if direction == 1 {
		apu.ch1.sweepNegateModeUsed = true
		frequency -= delta
	} else {
		frequency += delta
	}

	if frequency > 2047 {
		apu.ch1.enabled = false
		apu.nr52 &^= 0b0001
	}

	return frequency
}

func (apu *APU) reloadCh1SweepTimer() {
	pace := (apu.nr10 & 0b0111_0000) >> 4
	apu.ch1.sweepTimer = uint16(pace)
	if pace == 0 {
		apu.ch1.sweepTimer = 8
	}
}

// See: https://gbdev.gg8.se/wiki/articles/Gameboy_sound_hardware#Noise_Channel
var noiseChannelLookupTable [8]uint16 = [8]uint16{
	8, 16, 32, 48, 64, 80, 96, 112,
}

func (apu *APU) calculateCh4Period() uint16 {
	clockShift := (apu.nr43 & 0b1111_0000) >> 4
	clockDivider := (apu.nr43 & 0b0000_0111)

	baseDivisor := noiseChannelLookupTable[clockDivider]

	// Calculate the number of clock ticks (at 262144 Hz) to wait.
	// Period = Base Divisor * 2^Shift
	// Left shift is the same as multiplying by power of 2
	period := baseDivisor << clockShift

	return period
}

func (apu *APU) clockLength(ch *channel) {
	lengthEnabled := (*ch.register & 0b0100_0000) != 0
	if lengthEnabled && ch.lengthTimer < ch.maxLength {
		ch.lengthTimer++

		if ch.lengthTimer == ch.maxLength {
			ch.enabled = false
			apu.nr52 &^= ch.nr52BitMask
		}
	}
}

func (apu *APU) clockEnvelope(ch *channel) {
	if !ch.envelopeEnabled {
		return
	}

	envelopePeriod := *ch.envelopeRegister & 0b111
	if envelopePeriod == 0 {
		return
	}

	ch.envelopeTimer--

	if ch.envelopeTimer == 0 {
		ch.envelopeTimer = envelopePeriod

		direction := (*ch.envelopeRegister & 0b0000_1000) >> 3
		if direction == 1 {
			if ch.currentVolume < 15 {
				ch.currentVolume++
			} else {
				ch.envelopeEnabled = false
			}
		} else {
			if ch.currentVolume > 0 {
				ch.currentVolume--
			} else {
				ch.envelopeEnabled = false
			}
		}
	}
}

// OnDivReset is called when the DIV register (0xFF04) is reset, to keep the
// DIV-APU in sync because they are physically the same counter.
func (apu *APU) OnDivReset() {
	apu.divApuCounter = 0
}

func (apu *APU) SetChannelEnabled(channel int, enabled bool) {
	if channel >= 1 && channel <= 4 {
		apu.channelsEnabled[channel] = enabled
	}
}

func (apu *APU) GetChannelEnabled(channel int) bool {
	if channel >= 1 && channel <= 4 {
		return apu.channelsEnabled[channel]
	}

	return false
}

// ReadSamples copies up to len(dst) samples into dst.
// Returns number of samples written.
func (apu *APU) ReadSamples(dst []int16) int {
	return apu.outputBuffer.Read(dst)
}
