import { audioController } from "../services/audio-controller";
import { CanvasRenderer } from "../services/canvas-renderer";

// Decouple emulation (~59.7275 Hz) from display refresh rate:
// emulator produces frames, browser polls via requestAnimationFrame
class GameLoop {
	private animationFrameId: number | undefined;
	private lastFrameTime = 0;
	private tCycleAccumulator = 0;
	// 4,194,304 T-cycles per second
	private readonly SYSTEM_CLOCK_FREQUENCY = 4.194304 * 1_000_000;

	private _renderer: CanvasRenderer;

	constructor(renderer: CanvasRenderer) {
		this._renderer = renderer;
	}

	public start() {
		this.lastFrameTime = 0;
		this.animationFrameId = requestAnimationFrame(this.handleAnimationFrame);
	}

	public stop() {
		if (this.animationFrameId) {
			cancelAnimationFrame(this.animationFrameId);
			this.animationFrameId = undefined;
			this.lastFrameTime = 0;
		}
	}

	public renderer() {
		return this._renderer;
	}

	// timestamp is the end time of the previous frame's rendering
	private handleAnimationFrame = (timestamp: DOMHighResTimeStamp) => {
		// if we called stop()
		if (this.animationFrameId === undefined) {
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
			this._renderer.drawFrame(frame);
		}

		// play audio
		const samples = window.pollAudioBuffer();
		audioController.scheduleAudioSamples(samples);

		this.animationFrameId = requestAnimationFrame(this.handleAnimationFrame);
	};
}

export const gameLoop = new GameLoop(new CanvasRenderer());
