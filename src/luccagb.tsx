import { gameLoop } from "./core/game-loop";
import { store } from "./core/store";
import { InputManager } from "./services/input-manager";
import { setUpDragAndDropHandlers } from "./ui/drag-and-drop";
import { setUpControlsHandlers } from "./ui/controls";
import { render } from "solid-js/web";
import { VolumeControl } from "./ui/VolumeControl";
import { AudioChannels } from "./ui/AudioChannels";
import { DataManager } from "./ui/DataManager";
import { TestRoms } from "./ui/TestRoms";
import { handleRomLoad } from "./services/rom-loader";
import { TraceLogger } from "./ui/TraceLogger";
import { RomFileInput } from "./ui/RomFileInput";

const go = new Go();
const canvasRenderer = gameLoop.renderer();
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

	const romFileInput = document.getElementById("rom-file-input");
	if (romFileInput) {
		render(() => <RomFileInput />, romFileInput);
	}
	const solidVolumeControl = document.getElementById("solid-volume-control");
	if (solidVolumeControl) {
		render(() => <VolumeControl />, solidVolumeControl);
	}
	const solidAudioChannels = document.getElementById("solid-audio-channels");
	if (solidAudioChannels) {
		render(() => <AudioChannels />, solidAudioChannels);
	}
	const dataManager = document.getElementById("data-manager");
	if (dataManager) {
		render(() => <DataManager />, dataManager);
	}
	const testRoms = document.getElementById("test-roms");
	if (testRoms) {
		render(() => <TestRoms />, testRoms);
	}
	const traceLogger = document.getElementById("trace-logger");
	if (traceLogger) {
		render(() => <TraceLogger />, traceLogger);
	}

	setUpControlsHandlers({
		panelToggleId: "panel-toggle",
		controlsPanelId: "controls-panel",
		debugCheckboxId: "debug-checkbox",
		scaleSelectId: "scale-select",
	});

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

	// ======================================
	// ====== set up screenshot button ======
	// ======================================
	document
		.getElementById("screenshot-button")
		?.addEventListener("click", () => {
			canvasRenderer.takeScreenshot();
		});
});
