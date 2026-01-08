import type { CartridgeInfo } from "./wasm";

let visibleCanvasCtx: CanvasRenderingContext2D;
let offscreenCanvasCtx: CanvasRenderingContext2D;
let imageData: ImageData;
let currentScale: number | "fit" = 1;
let isPaused = false;
let isHidden = false;
let isFileInputOpen = false;
let isRomLoaded = false;
let romHash = "";
let cartridgeInfo: CartridgeInfo | null = null;

const displayWidth = 160;
const displayHeight = 144;

const go = new Go();

// ==========================
// ====== for debugger ======
// ==========================
const debugElements = {
	cartridgeTitle: null as HTMLElement | null,
	cartridgeType: null as HTMLElement | null,
	romSize: null as HTMLElement | null,
	ramSize: null as HTMLElement | null,
	regAF: null as HTMLElement | null,
	regBC: null as HTMLElement | null,
	regDE: null as HTMLElement | null,
	regHL: null as HTMLElement | null,
	regSP: null as HTMLElement | null,
	regPC: null as HTMLElement | null,
	flagZ: null as HTMLElement | null,
	flagN: null as HTMLElement | null,
	flagH: null as HTMLElement | null,
	flagC: null as HTMLElement | null,
};

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
		isFileInputOpen = false;
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
		isFileInputOpen = true;
	});

	fileInput?.addEventListener("cancel", () => {
		isFileInputOpen = false;
		if (!isPaused) {
			startAnimationLoop();
		}
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
			isHidden = true;

			// suspend audio context
		} else {
			isHidden = false;

			if (!isPaused) {
				startAnimationLoop();
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
			if (!visibleCanvasCtx) {
				return;
			}

			// grab the current pixels as a PNG image string
			const imageURL = visibleCanvasCtx.canvas.toDataURL("image/png");

			// create a temporary link element to trigger the download
			const link = document.createElement("a");
			link.href = imageURL;
			link.download = `luccagb-screenshot-${new Date().toISOString()}.png`;

			// trigger the click and cleanup
			document.body.appendChild(link);
			link.click();
			document.body.removeChild(link);
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

		let lines = [];
		for (let i = 0; i < buffer.length; ) {
			const type = buffer[i];

			if (type === 0) {
				// Instruction
				const pc = (buffer[i + 1] << 8) | buffer[i + 2];
				const opcode = buffer[i + 3];
				lines.push(
					`EXEC PC:0x${pc.toString(16).padStart(4, "0")} OP:0x${opcode.toString(16).padStart(2, "0")}`,
				);
				// if we log a 3rd thing, use this instead
				// const view = new DataView(buffer.buffer);
				// const data = view.getUint32(i + 4, false);
				// lines.push(
				// 	`EXEC PC:0x${pc.toString(16).padStart(4, "0")} OP:0x${opcode.toString(16).padStart(2, "0")} OTHER:${data}`,
				// );
				i += 8;
			} else if (type === 1) {
				// Memory read
				const addr = (buffer[i + 1] << 8) | buffer[i + 2];
				const value = buffer[i + 3];
				lines.push(
					`READ [0x${addr.toString(16).padStart(4, "0")}] = 0x${value.toString(16).padStart(2, "0")}`,
				);
				i += 4;
			} else if (type === 2) {
				// Memory write
				const addr = (buffer[i + 1] << 8) | buffer[i + 2];
				const value = buffer[i + 3];
				lines.push(
					`WRITE [0x${addr.toString(16).padStart(4, "0")}] = 0x${value.toString(16).padStart(2, "0")}`,
				);
				i += 4;
			} else {
				i++; // Skip unknown types
			}
		}

		const blob = new Blob([lines.join("\n")], { type: "text/plain" });
		const url = URL.createObjectURL(blob);
		const a = document.createElement("a");
		a.href = url;
		a.download = `luccagb-logs-${new Date().toISOString()}.txt`;
		a.click();
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

	// ===========================
	// ====== set up canvas ======
	// ===========================
	const visibleCanvas = document.getElementById("canvas");
	if (visibleCanvas instanceof HTMLCanvasElement) {
		const ctx = visibleCanvas.getContext("2d");
		if (!ctx) {
			throw new Error("error getting canvas context");
		}
		visibleCanvasCtx = ctx;
		visibleCanvasCtx.imageSmoothingEnabled = false;
		imageData = visibleCanvasCtx.createImageData(displayWidth, displayHeight);

		visibleCanvas.width = displayWidth;
		visibleCanvas.height = displayHeight;
	}
	const offscreenCanvas = document.createElement("canvas");
	offscreenCanvas.width = 160;
	offscreenCanvas.height = 144;
	const ctx = offscreenCanvas.getContext("2d");
	if (!ctx) {
		throw new Error("error getting offscreen canvas context");
	}
	offscreenCanvasCtx = ctx;

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

		applyCanvasScale();
	});

	function applyCanvasScale() {
		const container = document.getElementById("canvas-container");
		const canvas = visibleCanvasCtx?.canvas;

		if (!container || !canvas) {
			return;
		}

		if (currentScale === "fit") {
			// Fit to viewport (minus controls height)
			const maxWidth = window.innerWidth;
			const maxHeight = window.innerHeight - 56;

			const scale = Math.floor(
				Math.min(maxWidth / displayWidth, maxHeight / displayHeight),
			);

			const finalScale = Math.max(1, scale);

			container.style.width = `${displayWidth * finalScale}px`;
			container.style.height = `${displayHeight * finalScale}px`;
		} else {
			container.style.width = `${displayWidth * currentScale}px`;
			container.style.height = `${displayHeight * currentScale}px`;
		}
	}

	// apply default scale on load
	applyCanvasScale();
	window.addEventListener("resize", applyCanvasScale);

	// =================================================
	// ====== set up joypad (and shortcuts) input ======
	// =================================================
	const keyToJoypadButton: { [key: string]: string } = {
		Enter: "START",
		Backspace: "SELECT",
		KeyZ: "B",
		KeyX: "A",
		ArrowDown: "DOWN",
		ArrowUp: "UP",
		ArrowLeft: "LEFT",
		ArrowRight: "RIGHT",
	};
	const nonJoypadShortcuts: { [key: string]: (() => void) | undefined } = {
		Space: pauseOrResume,
	};
	window.addEventListener("keydown", (event) => {
		if (event.repeat) {
			return;
		}

		const funcHandler = nonJoypadShortcuts[event.code];
		if (funcHandler) {
			funcHandler();
			return;
		}

		if (!keyToJoypadButton[event.code]) {
			return;
		}

		window.handleJoypadButtonPressed(keyToJoypadButton[event.code]);
	});
	window.addEventListener("keyup", (event) => {
		if (!keyToJoypadButton[event.code]) {
			return;
		}

		window.handleJoypadButtonReleased(keyToJoypadButton[event.code]);
	});

	function pauseOrResume() {
		if (!isRomLoaded) {
			return;
		}

		isPaused = !isPaused;
		if (isPaused) {
			updateDebugView();
			persistCartridgeRam();
		} else {
			startAnimationLoop();
		}
	}

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
		updateDebugView();
	});

	syncDebugVisibility();

	// =============================
	// ====== set up debugger ======
	// =============================
	debugElements.cartridgeTitle = document.getElementById("cartridge-title");
	debugElements.cartridgeType = document.getElementById(
		"cartridge-cartridge-type",
	);
	debugElements.romSize = document.getElementById("cartridge-rom-size-code");
	debugElements.ramSize = document.getElementById("cartridge-ram-size-code");
	debugElements.regAF = document.getElementById("reg-af");
	debugElements.regBC = document.getElementById("reg-bc");
	debugElements.regDE = document.getElementById("reg-de");
	debugElements.regHL = document.getElementById("reg-hl");
	debugElements.regSP = document.getElementById("reg-sp");
	debugElements.regPC = document.getElementById("reg-pc");
	debugElements.flagZ = document.getElementById("flag-z");
	debugElements.flagN = document.getElementById("flag-n");
	debugElements.flagH = document.getElementById("flag-h");
	debugElements.flagC = document.getElementById("flag-c");
});

// ===================================
// ============ GAME LOOP ============
// ===================================
//
// Emulation runs at the Game Boy's hardware-accurate rate (~59.7275 Hz),
// while the browser renders at the display refresh rate (typically 60 Hz).
//
// Because these rates do not match, we decouple generation from presentation.
// This implements a Pull Architecture:
//   - The Emulator acts as a Producer, latching the latest completed frame.
//   - The Browser acts as a Consumer, polling for frames during requestAnimationFrame.
//
// This prevents lag accumulation and keeps timing accurate. The trade-off is
// minor visual judder (frame doubling) when the same frame is displayed twice to
// maintain synchronization.

let animationFrameId: number | undefined;

let lastFrameTime = 0;
let tCycleAccumulator = 0;
// 4,194,304 T-cycles per second
const systemClockFrequency = 4.194304 * 1_000_000;

const audioContext = new AudioContext();
let nextStartTime = 0;

// timestamp is the end time of the previous frame's rendering
function handleAnimationFrame(timestamp: DOMHighResTimeStamp) {
	if (isPaused || isHidden || isFileInputOpen) {
		return;
	}

	// start the loop
	if (lastFrameTime === 0) {
		lastFrameTime = timestamp;
		animationFrameId = requestAnimationFrame(handleAnimationFrame);
		return;
	}

	const deltaSeconds = (timestamp - lastFrameTime) / 1000;
	lastFrameTime = timestamp;

	const tCyclesToAdd = systemClockFrequency * deltaSeconds;
	tCycleAccumulator += tCyclesToAdd;

	const { tCyclesUsed } = window.processEmulatorCycles(tCycleAccumulator);
	tCycleAccumulator -= tCyclesUsed;

	const frame = window.pollFrame();
	if (frame) {
		updateCanvas(frame);
	}

	// ===================
	// ====== AUDIO ======
	// ===================

	const samples = window.pollAudioBuffer();
	if (samples?.length > 0) {
		const buffer = audioContext.createBuffer(1, samples.length, 48000);
		const channelData = buffer.getChannelData(0);

		for (let i = 0; i < samples.length; i++) {
			// Normalize sample to range [-1.0, 1.0]
			// Divide by 32768.0 to convert int16 range to float range
			// 32768.0 = max positive int16 value (2^15)
			channelData[i] = samples[i] / 32768.0;
		}

		// Handle underrun
		// If nextStartTime is in the past, reset it to "now" so we don't delay
		if (nextStartTime < audioContext.currentTime) {
			nextStartTime = audioContext.currentTime;
		}

		// Create a source to play this buffer
		const source = audioContext.createBufferSource();
		source.buffer = buffer;
		source.connect(audioContext.destination);

		// Schedule it to play
		source.start(nextStartTime);

		// Advance the pointer so the next chunk plays right after this one
		nextStartTime += buffer.duration;
	}

	animationFrameId = requestAnimationFrame(handleAnimationFrame);
}

function startAnimationLoop() {
	if (!isRomLoaded) {
		return;
	}
	lastFrameTime = 0;
	cancelAnimationFrame(animationFrameId!);
	animationFrameId = requestAnimationFrame(handleAnimationFrame);
}

function updateCanvas(uint8Array: Uint8Array) {
	// put the 160x144 data onto the same size offscreen canvas
	imageData.data.set(uint8Array);
	offscreenCanvasCtx.putImageData(imageData, 0, 0);

	// scale the image from the offscreen canvas onto the visible canvas
	visibleCanvasCtx.drawImage(
		offscreenCanvasCtx.canvas,
		0,
		0,
		visibleCanvasCtx.canvas.width,
		visibleCanvasCtx.canvas.height,
	);
}

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
	romHash = hashHex;

	// Load into Go
	cartridgeInfo = window.loadRom(romData);
	console.log("Cartridge Info:", cartridgeInfo);

	// Attempt to load existing RAM
	if (cartridgeInfo.hasBattery && cartridgeInfo.ramSize > 0) {
		try {
			const ram = await loadCartridgeRam();
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
	isRomLoaded = true;
	isPaused = false;
	updateDebugView();
	startAnimationLoop();
}

// ========================================
// ============ PERSISTING RAM ============
// ========================================

const DB_NAME = "LuccaGB-Database";
const DB_VERSION = 1;
const STORE_NAME = "cartridgeRam";

type SaveData = {
	romHash: string;
	ram: Uint8Array;
	updatedAt: number;
	meta?: {
		name?: string;
		playTime?: number;
	};
};

let dbPromise: Promise<IDBDatabase> | null = null;

export function openDatabase(): Promise<IDBDatabase> {
	if (dbPromise) {
		return dbPromise;
	}

	dbPromise = new Promise((resolve, reject) => {
		const request = indexedDB.open(DB_NAME, DB_VERSION);

		request.onupgradeneeded = (event) => {
			const db = (event.target as IDBOpenDBRequest).result;
			// Create store if missing
			if (!db.objectStoreNames.contains(STORE_NAME)) {
				db.createObjectStore(STORE_NAME, {
					keyPath: "romHash",
				});
			}
		};

		request.onsuccess = () => {
			const db = request.result;
			db.onversionchange = () => {
				console.warn("IndexedDB version change detected, closing DB");
				db.close();
				dbPromise = null;
			};
			resolve(db);
		};

		request.onerror = () => {
			reject(request.error);
		};

		request.onblocked = () => {
			console.warn(
				"IndexedDB upgrade blocked. Close other tabs using this site.",
			);
		};
	});

	return dbPromise;
}

async function loadCartridgeRam(): Promise<Uint8Array | null> {
	if (romHash === "") {
		console.error("no ROM hash, unable to load RAM");
		return null;
	}

	const db = await openDatabase();

	return new Promise((resolve, reject) => {
		const transaction = db.transaction([STORE_NAME], "readonly");
		const objectStore = transaction.objectStore(STORE_NAME);
		const request = objectStore.get(romHash);

		request.onsuccess = () => {
			const result = request.result as SaveData | undefined;

			if (result) {
				console.log(
					`Save found from: ${new Date(result.updatedAt).toLocaleString()}`,
				);

				// IndexedDB might return an ArrayBuffer or Uint8Array depending on browser
				// Ensure we return a Uint8Array
				resolve(
					result.ram instanceof Uint8Array
						? result.ram
						: new Uint8Array(result.ram),
				);
			} else {
				console.log("No existing save file found.");
				resolve(null);
			}
		};

		request.onerror = () => {
			console.error("loadCartridgeRam() failed");
			reject(new Error("loadCartridgeRam() failed"));
		};
	});
}

async function persistCartridgeRam(): Promise<void> {
	if (romHash === "") {
		return;
	}

	if (!cartridgeInfo?.hasBattery || cartridgeInfo?.ramSize == 0) {
		return;
	}

	const ram = window.getCartridgeRam();
	const db = await openDatabase();

	return new Promise((resolve, reject) => {
		const saveData: SaveData = {
			romHash: romHash,
			ram: ram,
			updatedAt: Date.now(),
			meta: {
				name: cartridgeInfo?.title,
			},
		};
		const transaction = db.transaction([STORE_NAME], "readwrite");
		const store = transaction.objectStore(STORE_NAME);

		transaction.oncomplete = () => {
			resolve();
		};

		transaction.onerror = () => {
			console.error("persist RAM failed", transaction.error);
			reject(transaction.error);
		};

		store.put(saveData);
	});
}

// ===============================
// ============ OTHER ============
// ===============================

function updateDebugView() {
	const debugInfo = window.getDebugInfo();
	if (!debugInfo) {
		return;
	}

	const { cartridge, cpu } = debugInfo;

	const toHex = (val: number) =>
		"0x" + val.toString(16).toUpperCase().padStart(4, "0");

	if (debugElements.cartridgeTitle)
		debugElements.cartridgeTitle.textContent = cartridge.title;
	if (debugElements.cartridgeType)
		debugElements.cartridgeType.textContent =
			cartridge.cartridgeType.toString();
	if (debugElements.romSize)
		debugElements.romSize.textContent = cartridge.romSizeCode.toString();
	if (debugElements.ramSize)
		debugElements.ramSize.textContent = cartridge.ramSizeCode.toString();

	if (debugElements.regAF)
		debugElements.regAF.textContent = toHex(cpu.registers16.AF);
	if (debugElements.regBC)
		debugElements.regBC.textContent = toHex(cpu.registers16.BC);
	if (debugElements.regDE)
		debugElements.regDE.textContent = toHex(cpu.registers16.DE);
	if (debugElements.regHL)
		debugElements.regHL.textContent = toHex(cpu.registers16.HL);
	if (debugElements.regSP)
		debugElements.regSP.textContent = toHex(cpu.registers16.SP);
	if (debugElements.regPC)
		debugElements.regPC.textContent = toHex(cpu.registers16.PC);

	if (debugElements.flagZ)
		debugElements.flagZ.textContent = cpu.flags.Z.toString();
	if (debugElements.flagN)
		debugElements.flagN.textContent = cpu.flags.N.toString();
	if (debugElements.flagH)
		debugElements.flagH.textContent = cpu.flags.H.toString();
	if (debugElements.flagC)
		debugElements.flagC.textContent = cpu.flags.C.toString();
}
