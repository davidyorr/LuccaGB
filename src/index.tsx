import "./global.css";

import { store } from "./core/store";
import { inputManager } from "./services/input-manager";
import { render } from "solid-js/web";
import { App } from "./App";
import { initWasm } from "./core/wasm";
import { gameLoop } from "./core/game-loop";
import { startRewind, stopRewind } from "./services/rewinder";

inputManager.registerShortcuts({
	Space: {
		keydown: store.actions.togglePaused,
	},
	Comma: {
		keydown: startRewind,
		keyup: stopRewind,
	},
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
