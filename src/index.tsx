import "./global.css";

import { store } from "./core/store";
import { InputManager } from "./services/input-manager";
import { render } from "solid-js/web";
import { App } from "./App";

const go = new Go();

new InputManager({
	Space: store.actions.togglePaused,
});

store.actions.initializeAppSettings();

const root = document.getElementById("app");
if (root) {
	render(() => <App />, root);
} else {
	console.error("root element not found");
}

document.addEventListener("DOMContentLoaded", async () => {
	WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then(
		(wasmModule) => {
			go.run(wasmModule.instance);
		},
	);
});
