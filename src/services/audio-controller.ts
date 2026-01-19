export class AudioController {
	private audioContext: AudioContext;
	private nextStartTime = 0;

	constructor() {
		this.audioContext = new AudioContext();
	}

	public scheduleAudioSamples(samples: number[]) {
		if (samples.length === 0) {
			return;
		}

		const buffer = this.audioContext.createBuffer(1, samples.length, 48000);
		const channelData = buffer.getChannelData(0);

		for (let i = 0; i < samples.length; i++) {
			// Normalize sample to range [-1.0, 1.0]
			// Divide by 32768.0 to convert int16 range to float range
			// 32768.0 = max positive int16 value (2^15)
			channelData[i] = samples[i] / 32768.0;
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
