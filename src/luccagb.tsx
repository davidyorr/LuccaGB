import { gameLoop } from "./core/game-loop";
import { store } from "./core/store";
import { audioController } from "./services/audio-controller";
import { InputManager } from "./services/input-manager";
import { TestRomLibrary } from "./services/test-rom-library";
import { loadCartridgeRam } from "./services/storage";
import { Debugger } from "./ui/debugger";
import {
	downloadTraceLogs as downloadTraceLog,
	parseTraceLogs,
} from "./utils/trace-logger";
import type { CartridgeInfo } from "./wasm";
import { setUpDragAndDropHandlers } from "./ui/drag-and-drop";
import { setUpControlsHandlers } from "./ui/controls";
import { render } from "solid-js/web";
import { VolumeControl } from "./ui/VolumeControl";
import { AudioChannels } from "./ui/AudioChannels";
import { DataManager } from "./ui/DataManager";

let cartridgeInfo: CartridgeInfo | null = null;

const go = new Go();
const canvasRenderer = gameLoop.renderer();
const testRomLibrary = new TestRomLibrary();
new InputManager({
	Space: store.legacyAppState.togglePaused,
});
const debug = new Debugger();

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

	// =================================
	// ====== set up ROM dropdown ======
	// =================================
	testRomLibrary.populateSelect("rom-select");

	// handle selection
	document
		.getElementById("rom-select")
		?.addEventListener("change", async (event) => {
			const target = event.target as HTMLSelectElement;

			// remove focus so keyboard controls don't toggle the dropdown
			target.blur();
			// reset file input so it doesn't look like two things are selected
			const fileInput = document.getElementById(
				"rom-input",
			) as HTMLInputElement;
			if (fileInput) {
				fileInput.value = "";
			}

			const buffer = await testRomLibrary.loadRomByPath(target.value);
			if (buffer) {
				await handleRomLoad(buffer);
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

async function handleRomLoad(arrayBuffer: ArrayBuffer) {
	const romData = new Uint8Array(arrayBuffer);

	// pause the previous audio context
	await audioController.pause();

	// Compute the ROM Hash
	const hashBuffer = await window.crypto.subtle.digest(
		"SHA-256",
		romData.buffer,
	);
	const hashArray = Array.from(new Uint8Array(hashBuffer));
	const hashHex = hashArray
		.map((b) => b.toString(16).padStart(2, "0"))
		.join("");
	store.actions.setCurrentRomHash(hashHex);

	// Load into Go
	cartridgeInfo = window.loadRom(romData);
	store.actions.setCartridgeInfo(cartridgeInfo);
	console.log("Cartridge Info:", cartridgeInfo);

	// Attempt to load existing RAM
	if (cartridgeInfo.hasBattery && cartridgeInfo.ramSize > 0) {
		try {
			const ram = await loadCartridgeRam(store.legacyAppState.currentRomHash);
			if (ram) {
				// Ensure the loaded RAM size matches what the cartridge expects
				if (ram.length !== cartridgeInfo.ramSize) {
					console.warn(
						`Save file size mismatch. Expected ${cartridgeInfo.ramSize}, got ${ram.length}`,
					);
				}
				window.setCartridgeRam(ram);
			}
		} catch (e) {
			console.error("Failed to load save data:", e);
		}
	}

	// Focus the canvas so keyboard controls work immediately
	const canvas = document.getElementById("canvas");
	if (canvas) {
		canvas.tabIndex = 0;
		canvas.focus();
	}

	// set the initial audio channels state
	window.setAudioChannelEnabled(
		1,
		store.state.settings.audioChannelsEnabled[1],
	);
	window.setAudioChannelEnabled(
		2,
		store.state.settings.audioChannelsEnabled[2],
	);
	window.setAudioChannelEnabled(
		3,
		store.state.settings.audioChannelsEnabled[3],
	);
	window.setAudioChannelEnabled(
		4,
		store.state.settings.audioChannelsEnabled[4],
	);

	// resume the audio context
	await audioController.resume();

	// Start the animation loop
	store.actions.setRomLoaded(true);
	debug.update();
	gameLoop.start();
}
