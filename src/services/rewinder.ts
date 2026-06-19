import { gameLoop } from "../core/game-loop";
import { store } from "../core/store";
import { audioController } from "./audio-controller";
import { inputManager } from "./input-manager";

let rewindIntervalId: number | undefined;
let rewindDelayTimeoutId: number | undefined;
const REWIND_INTERVAL_MS = 16; // ~60fps updates
const REWIND_DELAY_MS = 100;

export function startRewind() {
	if (store.state.isRewinding) {
		return;
	}

	const currentBuffer = store.state.settings.rewindBufferSize;
	if (currentBuffer === 0 || !window.rewindFrames) {
		return;
	}

	// Clear any existing timeouts or intervals
	clearInterval(rewindIntervalId);
	clearTimeout(rewindDelayTimeoutId);

	store.actions.setRewinding(true);
	audioController.onRewind();

	// Do the first rewind immediately
	performRewindStep();

	// Wait before starting continuous rewind
	rewindDelayTimeoutId = setTimeout(() => {
		rewindIntervalId = setInterval(() => {
			performRewindStep();
		}, REWIND_INTERVAL_MS);
	}, REWIND_DELAY_MS);
}

export function stopRewind() {
	// Clear both the delay timeout and the interval
	if (rewindDelayTimeoutId !== undefined) {
		clearTimeout(rewindDelayTimeoutId);
		rewindDelayTimeoutId = undefined;
	}
	if (rewindIntervalId !== undefined) {
		clearInterval(rewindIntervalId);
		rewindIntervalId = undefined;
	}
	store.actions.setRewinding(false);

	// Sync the joypad state back to physical reality so buttons don't get stuck
	inputManager.syncJoypadState();
}

export function performRewindStep() {
	const framesRewound = window.rewindFrames(
		store.state.settings.rewindIncrement,
	);

	if (framesRewound > 0) {
		// Fetch the frame so the canvas can update
		const frameData = window.pollFrame?.();
		if (frameData) {
			gameLoop.forceDraw(frameData);
		}
	} else {
		// Stop if there are no more frames to rewind
		stopRewind();
	}
}
