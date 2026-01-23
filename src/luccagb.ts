import { GameLoop } from "./core/game-loop";
import { appState } from "./core/state";
import { AudioController } from "./services/audio-controller";
import { CanvasRenderer } from "./services/canvas-renderer";
import { InputManager } from "./services/input-manager";
import { TestRomLibrary } from "./services/test-rom-library";
import {
	loadCartridgeRam,
	persistCartridgeRam,
	saveAppSettings,
} from "./services/storage";
import { Debugger } from "./ui/debugger";
import {
	downloadTraceLogs as downloadTraceLog,
	parseTraceLogs,
} from "./utils/trace-logger";
import type { CartridgeInfo } from "./wasm";
import { setUpDataManagerHandlers } from "./ui/data-manager";
import { setUpDragAndDropHandlers } from "./ui/drag-and-drop";
import { setUpAudioChannelHandlers } from "./ui/audio-channels";
import { setUpAudioVolumeHandlers } from "./ui/audio-volume";
import { setUpControlsHandlers } from "./ui/controls";
import { debounce } from "./utils/debounce";

let cartridgeInfo: CartridgeInfo | null = null;

const go = new Go();
const canvasRenderer = new CanvasRenderer("canvas");
const audioController = new AudioController();
const testRomLibrary = new TestRomLibrary();
new InputManager({
	Space: appState.togglePaused,
});
const debug = new Debugger();

const gameLoop = new GameLoop(audioController, canvasRenderer);

document.addEventListener("DOMContentLoaded", async () => {
	WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then(
		(wasmModule) => {
			go.run(wasmModule.instance);
		},
	);

	await appState.initializeAppSettings();

	const debouncedSave = debounce((state: typeof appState) => {
		saveAppSettings(state.getSettingsSnapshot());
	}, 333);

	appState.subscribe((state) => {
		debouncedSave(state);
	});

	setUpControlsHandlers({
		panelToggleId: "panel-toggle",
		controlsPanelId: "controls-panel",
		debugCheckboxId: "debug-checkbox",
		scaleSelectId: "scale-select",
		canvasRenderer: canvasRenderer,
	});

	setUpDragAndDropHandlers({
		overlayId: "drag-overlay",
		onRomLoaded: handleRomLoad,
	});

	setUpDataManagerHandlers({
		dataModalId: "data-modal",
		dataButtonId: "data-manager-button",
		exportButtonId: "export-data-button",
		importId: "import-data-input",
	});

	setUpAudioVolumeHandlers({
		volumeSliderId: "volume-slider",
		volumeValueId: "volume-value",
	});

	setUpAudioChannelHandlers({
		buttonId: "audio-channels-button",
		dropdownId: "audio-channels-dropdown",
		audioController: audioController,
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
		appState.setFileInputOpen(false);
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
		appState.setFileInputOpen(true);
	});

	fileInput?.addEventListener("cancel", () => {
		appState.setFileInputOpen(false);
	});

	// ==================================================
	// ====== set up tab visibility event listener ======
	// ==================================================
	document.addEventListener("visibilitychange", () => {
		if (document.hidden) {
			appState.setHidden(true);

			// suspend audio context
		} else {
			appState.setHidden(false);

			if (!appState.isPaused) {
				gameLoop.startAnimationLoop();
			}

			// resume audio context
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

	// ==================================
	// ====== handle state changes ======
	// ==================================
	appState.subscribe(async (state) => {
		if (state.isPaused) {
			if (!cartridgeInfo?.hasBattery || cartridgeInfo?.ramSize === 0) {
				return;
			}
			const ram = window.getCartridgeRam();
			persistCartridgeRam(state.currentRomHash, ram, {
				name: cartridgeInfo.title,
			});

			await audioController.pause();
		} else {
			await audioController.resume();
			gameLoop.startAnimationLoop();
		}
	});
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
	appState.setCurrentRomHash(hashHex);

	// Load into Go
	cartridgeInfo = window.loadRom(romData);
	console.log("Cartridge Info:", cartridgeInfo);

	// Attempt to load existing RAM
	if (cartridgeInfo.hasBattery && cartridgeInfo.ramSize > 0) {
		try {
			const ram = await loadCartridgeRam(appState.currentRomHash);
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
	window.setAudioChannelEnabled(1, appState.audioChannelsEnabled[1]);
	window.setAudioChannelEnabled(2, appState.audioChannelsEnabled[2]);
	window.setAudioChannelEnabled(3, appState.audioChannelsEnabled[3]);
	window.setAudioChannelEnabled(4, appState.audioChannelsEnabled[4]);

	// resume the audio context
	await audioController.resume();

	// Start the animation loop
	appState.setRomLoaded(true);
	debug.update();
	gameLoop.startAnimationLoop();
}
