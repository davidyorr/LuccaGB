import { appState } from "../core/state";

export class AudioController {
	private audioContext: AudioContext;
	private nextStartTime = 0;
	// Max allowed latency in seconds.
	private readonly MAX_LATENCY_TIME = 0.1;
	// The amount of time to schedule ahead when resetting.
	private readonly LOOKAHEAD_TIME = 0.06;

	constructor() {
		const AudioContext =
			window.AudioContext || (window as any).webkitAudioContext;
		this.audioContext = new AudioContext();
	}

	public scheduleAudioSamples(samples: number[] | null) {
		if (samples === null || samples.length === 0) {
			return;
		}

		// Latency check:
		// If the next scheduled time is too far in the future, it means we have
		// too much audio queued up. Drop this chunk to let the visual/audio sync up.
		const currentTime = this.audioContext.currentTime;
		if (this.nextStartTime > currentTime + this.MAX_LATENCY_TIME) {
			return;
		}

		// Handle underrun:
		// If nextStartTime is in the past, reset it to "now" + LOOKAHEAD_TIME so we don't delay
		if (this.nextStartTime < currentTime) {
			this.nextStartTime = currentTime + this.LOOKAHEAD_TIME;
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

		// adjust volume from UI setting
		const gainNode = this.audioContext.createGain();
		gainNode.gain.value = appState.audioVolume;
		gainNode.connect(this.audioContext.destination);

		// Create a source to play this buffer
		const source = this.audioContext.createBufferSource();
		source.buffer = buffer;
		source.connect(gainNode);

		// Schedule it to play
		source.start(this.nextStartTime);

		// Advance the pointer so the next chunk plays right after this one
		this.nextStartTime += buffer.duration;
	}

	public async pause() {
		if (this.audioContext.state === "running") {
			return this.audioContext.suspend();
		}
	}

	public async resume() {
		this.nextStartTime = 0;
		if (this.audioContext.state === "suspended") {
			return this.audioContext.resume();
		}
	}

	public resetTime() {
		this.nextStartTime = this.audioContext.currentTime + this.LOOKAHEAD_TIME;
	}
}
