export class AudioController {
	private audioContext: AudioContext;
	private nextStartTime = 0;

	constructor() {
		this.audioContext = new AudioContext();
	}

	public scheduleAudioSamples(samples: number[] | null) {
		if (samples === null || samples.length === 0) {
			return;
		}

		// sample are interleaved stereo [L, R, L, R, ...]
		const frameCount = samples.length / 2;
		const buffer = this.audioContext.createBuffer(2, frameCount, 48000);

		const leftChannel = buffer.getChannelData(0);
		const rightChannel = buffer.getChannelData(1);

		// deinterleave the samples
		for (let i = 0; i < frameCount; i++) {
			// Normalize sample to range [-1.0, 1.0]
			// Divide by 32768.0 to convert int16 range to float range
			// 32768.0 = max positive int16 value (2^15)
			leftChannel[i] = samples[i * 2] / 32768.0;
			rightChannel[i] = samples[i * 2 + 1] / 32768.0;
		}

		// Handle underrun
		// If nextStartTime is in the past, reset it to "now" so we don't delay
		if (this.nextStartTime < this.audioContext.currentTime) {
			this.nextStartTime = this.audioContext.currentTime;
		}

		// Create a source to play this buffer
		const source = this.audioContext.createBufferSource();
		source.buffer = buffer;
		source.connect(this.audioContext.destination);

		// Schedule it to play
		source.start(this.nextStartTime);

		// Advance the pointer so the next chunk plays right after this one
		this.nextStartTime += buffer.duration;
	}
}
