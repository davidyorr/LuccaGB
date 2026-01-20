import { GameLoop } from "./core/game-loop";
import { emulatorState } from "./core/state";
import { AudioController } from "./services/audio-controller";
import { CanvasRenderer } from "./services/canvas-renderer";
import { InputManager } from "./services/input-manager";
import { loadCartridgeRam, persistCartridgeRam } from "./services/storage";
import { Debugger } from "./ui/debugger";
import {
	downloadTraceLogs as downloadTraceLog,
	parseTraceLogs,
} from "./utils/trace-logger";
import type { CartridgeInfo } from "./wasm";

let currentScale: number | "fit" = 1;
let cartridgeInfo: CartridgeInfo | null = null;

const go = new Go();
const canvasRenderer = new CanvasRenderer("canvas");
const audioController = new AudioController();
new InputManager({
	Space: emulatorState.togglePaused,
});
const debug = new Debugger();

const gameLoop = new GameLoop(audioController, canvasRenderer);

const romFiles = import.meta.glob("../roms/**/*.gb", {
	query: "?url",
	import: "default",
});

document.addEventListener("DOMContentLoaded", () => {
	WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then(
		(wasmModule) => {
			go.run(wasmModule.instance);
		},
	);

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
		emulatorState.setFileInputOpen(false);
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
		emulatorState.setFileInputOpen(true);
	});

	fileInput?.addEventListener("cancel", () => {
		emulatorState.setFileInputOpen(false);
	});

	// ==================================
	// ====== set up Drag and Drop ======
	// ==================================
	const dragOverlay = document.getElementById("drag-overlay");
	if (dragOverlay) {
		window.addEventListener("dragenter", (e) => {
			// Only show if dragging files
			if (e.dataTransfer?.types.includes("Files")) {
				dragOverlay.style.display = "flex";
			}
		});

		dragOverlay.addEventListener("dragover", (e) => {
			e.preventDefault();
			if (e.dataTransfer) {
				e.dataTransfer.dropEffect = "copy";
			}
		});

		// Hide overlay if user leaves the window (dragging out)
		dragOverlay.addEventListener("dragleave", (e) => {
			if (e.target === dragOverlay) {
				dragOverlay.style.display = "none";
			}
		});

		dragOverlay.addEventListener("drop", async (e) => {
			e.preventDefault();
			dragOverlay.style.display = "none";

			const files = e.dataTransfer?.files;
			if (files && files.length > 0) {
				const file = files[0];

				if (file.name.toLowerCase().endsWith(".gb")) {
					try {
						const arrayBuffer = await file.arrayBuffer();
						await handleRomLoad(arrayBuffer);

						// Clear UI inputs to match state
						const romSelect = document.getElementById(
							"rom-select",
						) as HTMLSelectElement;
						const fileInput = document.getElementById(
							"rom-input",
						) as HTMLInputElement;
						if (romSelect) {
							romSelect.value = "";
						}
						if (fileInput) {
							fileInput.value = "";
						}
					} catch (err) {
						console.error("Error reading dropped file", err);
						alert("Failed to load dropped file.");
					}
				} else {
					alert("Please drop a valid .gb file.");
				}
			}
		});
	}

	// ==================================================
	// ====== set up tab visibility event listener ======
	// ==================================================
	document.addEventListener("visibilitychange", () => {
		if (document.hidden) {
			emulatorState.setHidden(true);

			// suspend audio context
		} else {
			emulatorState.setHidden(false);

			if (!emulatorState.isPaused) {
				gameLoop.startAnimationLoop();
			}

			// resume audio context
		}
	});

	// =================================
	// ====== set up ROM dropdown ======
	// =================================
	const romSelect = document.getElementById(
		"rom-select",
	) as HTMLSelectElement | null;

	if (!romSelect) {
		console.warn("Could not find #rom-select element");
		return;
	}

	// populate options
	const sortedPaths = Object.keys(romFiles).sort();

	for (const path of sortedPaths) {
		const option = document.createElement("option");
		option.text = path.replace("../roms/", "");
		option.value = path;
		romSelect.appendChild(option);
	}

	// handle selection
	romSelect.addEventListener("change", async (event) => {
		const target = event.target as HTMLSelectElement;
		const path = target.value;

		if (!path || !romFiles[path]) {
			return;
		}

		// remove focus so keyboard controls don't toggle the dropdown
		target.blur();

		try {
			// fetch the URL from Vite
			const getUrl = romFiles[path] as () => Promise<string>;
			const url = await getUrl();

			const response = await fetch(url);
			if (!response.ok) {
				throw new Error(`Failed to fetch ${path}`);
			}

			const arrayBuffer = await response.arrayBuffer();
			await handleRomLoad(arrayBuffer);

			// reset file input so it doesn't look like two things are selected
			const fileInput = document.getElementById(
				"rom-input",
			) as HTMLInputElement;
			if (fileInput) {
				fileInput.value = "";
			}
		} catch (err) {
			console.error("Error loading ROM from dropdown:", err);
			alert("Failed to load selected ROM.");
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

	// ====================================
	// ====== set up display scaling ======
	// ====================================
	const scaleSelect = document.getElementById(
		"scale-select",
	) as HTMLSelectElement | null;
	currentScale = Number.parseInt(scaleSelect!.value) ?? 1;

	scaleSelect?.addEventListener("change", () => {
		const value = scaleSelect.value;

		if (value === "fit") {
			currentScale = "fit";
		} else {
			currentScale = parseInt(value, 10);
		}

		canvasRenderer.setScale(currentScale);
	});

	// apply default scale on load
	canvasRenderer.setScale(currentScale);
	window.addEventListener("resize", () => {
		canvasRenderer.setScale(currentScale);
	});

	emulatorState.subscribe((state) => {
		if (state.isPaused) {
			debug.update();

			if (!cartridgeInfo?.hasBattery || cartridgeInfo?.ramSize === 0) {
				return;
			}
			const ram = window.getCartridgeRam();
			persistCartridgeRam(state.currentRomHash, ram, {
				name: cartridgeInfo.title,
			});
		} else {
			gameLoop.startAnimationLoop();
		}
	});

	// ===================================
	// ====== set up debug checkbox ======
	// ===================================
	const debugCheckbox = document.getElementById(
		"debug-checkbox",
	) as HTMLInputElement | null;
	const debugPanel = document.getElementById("debug-panel");

	function syncDebugVisibility() {
		if (!debugPanel || !debugCheckbox) {
			return;
		}
		debugPanel.style.display = debugCheckbox.checked ? "block" : "none";
	}

	debugCheckbox?.addEventListener("change", () => {
		if (!debugPanel) {
			return;
		}
		syncDebugVisibility();
		debug.update();
	});

	syncDebugVisibility();
});

async function handleRomLoad(arrayBuffer: ArrayBuffer) {
	const romData = new Uint8Array(arrayBuffer);

	// Compute the ROM Hash
	const hashBuffer = await window.crypto.subtle.digest(
		"SHA-256",
		romData.buffer,
	);
	const hashArray = Array.from(new Uint8Array(hashBuffer));
	const hashHex = hashArray
		.map((b) => b.toString(16).padStart(2, "0"))
		.join("");
	emulatorState.setCurrentRomHash(hashHex);

	// Load into Go
	cartridgeInfo = window.loadRom(romData);
	console.log("Cartridge Info:", cartridgeInfo);

	// Attempt to load existing RAM
	if (cartridgeInfo.hasBattery && cartridgeInfo.ramSize > 0) {
		try {
			const ram = await loadCartridgeRam(emulatorState.currentRomHash);
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

	// Start the animation loop
	emulatorState.setRomLoaded(true);
	debug.update();
	gameLoop.startAnimationLoop();
}
