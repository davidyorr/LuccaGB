import { store } from "./core/store";
import { InputManager } from "./services/input-manager";
import { setUpDragAndDropHandlers } from "./ui/drag-and-drop";
import { render } from "solid-js/web";
import { handleRomLoad } from "./services/rom-loader";
import { Debugger } from "./ui/Debugger";
import { Controls } from "./ui/Controls";

const go = new Go();
new InputManager({
	Space: store.legacyAppState.togglePaused,
});

document.addEventListener("DOMContentLoaded", async () => {
	WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then(
		(wasmModule) => {
			go.run(wasmModule.instance);
		},
	);

	await store.legacyAppState.initializeAppSettings();

	const controlsRoot = document.getElementById("controls-root");
	if (controlsRoot) {
		render(() => <Controls />, controlsRoot);
	}

	const debugContainer = document.getElementById("debug-container");
	if (debugContainer) {
		render(() => <Debugger />, debugContainer);
	}

	setUpDragAndDropHandlers({
		overlayId: "drag-overlay",
		onRomLoaded: handleRomLoad,
	});

	// ==================================================
	// ====== set up tab visibility event listener ======
	// ==================================================
	document.addEventListener("visibilitychange", () => {
		if (document.hidden) {
			store.actions.setHidden(true);
		} else {
			store.actions.setHidden(false);
		}
	});
});
