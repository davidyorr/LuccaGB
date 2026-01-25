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

	// ========================================
	// ====== force clear the file input ======
	// ========================================
	const fileInput = document.getElementById(
		"rom-input",
	) as HTMLInputElement | null;
	if (fileInput) {
		fileInput.value = "";
	}

	// =============================================
	// ====== set up ROM input event listener ======
	// =============================================
	fileInput?.addEventListener("change", async (event) => {
		store.actions.setFileInputOpen(false);
		const files = (event.target as HTMLInputElement | null)?.files;
		if (files?.[0]) {
			const arrayBuffer = await files?.[0].arrayBuffer();
			await handleRomLoad(arrayBuffer);
		}

		// reset the dropdown to the default so it doesn't look like two things are selected
		const romSelect = document.getElementById(
			"rom-select",
		) as HTMLSelectElement | null;
		if (romSelect) {
			romSelect.value = "";
		}
	});

	fileInput?.addEventListener("click", () => {
		store.actions.setFileInputOpen(true);
	});

	fileInput?.addEventListener("cancel", () => {
		store.actions.setFileInputOpen(false);
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
