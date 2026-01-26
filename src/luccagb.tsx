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
import { ViewportScale } from "./ui/ViewportScale";
import { Debugger, DebuggerToggle } from "./ui/Debugger";
import { ScreenshotButton } from "./ui/ScreenshotButton";

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
	const screenshotButton = document.getElementById(
		"screenshot-button-container",
	);
	if (screenshotButton) {
		render(() => <ScreenshotButton />, screenshotButton);
	}
	const viewportScaleContainer = document.getElementById(
		"viewport-scale-container",
	);
	if (viewportScaleContainer) {
		render(() => <ViewportScale />, viewportScaleContainer);
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
	const debugToggle = document.getElementById("debug-toggle");
	if (debugToggle) {
		render(() => <DebuggerToggle />, debugToggle);
	}
	const debugContainer = document.getElementById("debug-container");
	if (debugContainer) {
		render(() => <Debugger />, debugContainer);
	}

	setUpControlsHandlers({
		panelToggleId: "panel-toggle",
		controlsPanelId: "controls-panel",
		debugCheckboxId: "debug-checkbox",
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
});
