import { gameLoop } from "./core/game-loop";
import { store } from "./core/store";
import { InputManager } from "./services/input-manager";
import {
	downloadTraceLogs as downloadTraceLog,
	parseTraceLogs,
} from "./utils/trace-logger";
import { setUpDragAndDropHandlers } from "./ui/drag-and-drop";
import { setUpControlsHandlers } from "./ui/controls";
import { render } from "solid-js/web";
import { VolumeControl } from "./ui/VolumeControl";
import { AudioChannels } from "./ui/AudioChannels";
import { DataManager } from "./ui/DataManager";
import { TestRoms } from "./ui/TestRoms";
import { handleRomLoad } from "./services/rom-loader";

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

	// ==================================
	// ====== set up download logs ======
	// ==================================
	const downloadTraceLogButton = document.getElementById(
		"download-trace-log-button",
	);
	const traceLogCheckbox = document.getElementById(
		"trace-log-checkbox",
	) as HTMLInputElement | null;
	const traceLogLabel = document.getElementById("trace-log-toggle-container");

	traceLogCheckbox?.addEventListener("change", (event) => {
		const isEnabled = (event.target as HTMLInputElement).checked;

		if (isEnabled) {
			window.enableTraceLogging();
			if (downloadTraceLogButton) {
				downloadTraceLogButton.style.display = "inline-block";
			}
		} else {
			window.disableTraceLogging();
			if (downloadTraceLogButton) {
				downloadTraceLogButton.style.display = "none";
			}
		}
	});

	downloadTraceLogButton?.addEventListener("click", () => {
		const buffer = window.getTraceLogs();

		if (!buffer || buffer.length === 0) {
			alert("No logs available.");
			return;
		}

		const text = parseTraceLogs(buffer);
		downloadTraceLog(text);
	});

	// Hide the button and checkbox in production
	if (import.meta.env.PROD) {
		if (downloadTraceLogButton) {
			downloadTraceLogButton.style.display = "none";
		}
		if (traceLogLabel) {
			traceLogLabel.style.display = "none";
		}
	}
});
