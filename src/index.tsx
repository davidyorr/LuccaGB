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
	Space: { keydown: store.actions.togglePaused },
	Comma: { keydown: startRewind, keyup: stopRewind },
	F9: { keydown: toggleRecordingGameplay },
	KeyO: { keydown: quickSaveState },
	KeyP: { keydown: quickLoadState },
	BracketLeft: { keydown: store.actions.setSaveStateSlotToPrev },
	BracketRight: { keydown: store.actions.setSaveStateSlotToNext },
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
