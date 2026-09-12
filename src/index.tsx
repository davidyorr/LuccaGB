import "./global.css";

import { store } from "./core/store";
import { inputManager } from "./services/input-manager";
import { render } from "solid-js/web";
import { App } from "./App";
import { initWasm } from "./core/wasm";
import { gameLoop } from "./core/game-loop";
import { startRewind, stopRewind } from "./services/rewinder";
import { toggleRecordingGameplay } from "./services/gameplay-recorder";
import { quickLoadState, quickSaveState } from "./services/save-state-manager";

inputManager.registerShortcuts({
	togglePaused: { keydown: store.actions.togglePaused },
	rewind: { keydown: startRewind, keyup: stopRewind },
	toggleRecording: { keydown: toggleRecordingGameplay },
	quickSave: { keydown: quickSaveState },
	quickLoad: { keydown: quickLoadState },
	prevSaveStateSlot: { keydown: store.actions.setSaveStateSlotToPrev },
	nextSaveStateSlot: { keydown: store.actions.setSaveStateSlotToNext },
});

gameLoop.attachInputManager(inputManager);

store.actions.initializeAppSettings();

const root = document.getElementById("app");
if (root) {
	render(() => <App />, root);
} else {
	console.error("root element not found");
}

document.addEventListener("DOMContentLoaded", initWasm);
