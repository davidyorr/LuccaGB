import type { AudioController } from "../services/audio-controller";
import type { CanvasRenderer } from "../services/canvas-renderer";
import { appState } from "./state";

// Decouple emulation (~59.7275 Hz) from display refresh rate:
// emulator produces frames, browser polls via requestAnimationFrame
export class GameLoop {
	private animationFrameId: number | undefined;
	private lastFrameTime = 0;
	private tCycleAccumulator = 0;
	// 4,194,304 T-cycles per second
	private readonly SYSTEM_CLOCK_FREQUENCY = 4.194304 * 1_000_000;

	private audio: AudioController;
	private renderer: CanvasRenderer;

	constructor(audio: AudioController, renderer: CanvasRenderer) {
		this.audio = audio;
		this.renderer = renderer;
	}

	public startAnimationLoop() {
		if (!appState.isRomLoaded) {
			return;
		}
		cancelAnimationFrame(this.animationFrameId!);
		this.lastFrameTime = 0;
		this.animationFrameId = requestAnimationFrame(this.handleAnimationFrame);
	}

	// timestamp is the end time of the previous frame's rendering
	private handleAnimationFrame = (timestamp: DOMHighResTimeStamp) => {
		if (appState.isPaused || appState.isHidden || appState.isFileInputOpen) {
			return;
		}

		// start the loop
		if (this.lastFrameTime === 0) {
			this.lastFrameTime = timestamp;
			this.animationFrameId = requestAnimationFrame(this.handleAnimationFrame);
			return;
		}

		const deltaSeconds = (timestamp - this.lastFrameTime) / 1000;
		this.lastFrameTime = timestamp;

		// run emulator steps
		const tCyclesToAdd = this.SYSTEM_CLOCK_FREQUENCY * deltaSeconds;
		this.tCycleAccumulator += tCyclesToAdd;
		const { tCyclesUsed } = window.processEmulatorCycles(
			this.tCycleAccumulator,
		);
		this.tCycleAccumulator -= tCyclesUsed;

		// render video
		const frame = window.pollFrame();
		if (frame) {
			this.renderer.drawFrame(frame);
		}

		// play audio
		const samples = window.pollAudioBuffer();
		this.audio.scheduleAudioSamples(samples);

		this.animationFrameId = requestAnimationFrame(this.handleAnimationFrame);
	};
}
