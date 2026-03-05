import "./global.css";

import { store } from "./core/store";
import { InputManager } from "./services/input-manager";
import { render } from "solid-js/web";
import { App } from "./App";
import { initWasm } from "./core/wasm";
import { gameLoop } from "./core/game-loop";

const inputManager = new InputManager({
	Space: store.actions.togglePaused,
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
