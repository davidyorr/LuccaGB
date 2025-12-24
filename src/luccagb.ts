(window as any).onRomLoaded = onRomLoaded;
(window as any).updateCanvas = updateCanvas;

let visibleCanvasCtx: CanvasRenderingContext2D;
let offscreenCanvasCtx: CanvasRenderingContext2D;
let imageData: ImageData;
let currentScale: number | "fit" = 1;
let isPaused = false;
let isHidden = false;

const displayWidth = 160;
const displayHeight = 144;

const go = new Go();

// ==========================
// ====== for debugger ======
// ==========================
const debugElements = {
	cartridgeTitle: null as HTMLElement | null,
	mbcType: null as HTMLElement | null,
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
	document
		.getElementById("rom-input")
		?.addEventListener("change", async (event) => {
			const files = (event.target as HTMLInputElement | null)?.files;
			if (files?.[0]) {
				const arrayBuffer = await files?.[0].arrayBuffer();
				const romData = new Uint8Array(arrayBuffer);
				window.loadRom(romData);
			}

			// reset the dropdown to the default so it doesn't look like two things are selected
			const romSelect = document.getElementById(
				"rom-select",
			) as HTMLSelectElement | null;
			if (romSelect) {
				romSelect.value = "";
			}
		});

	// ==================================================
	// ====== set up tab visibility event listener ======
	// ==================================================
	document.addEventListener("visibilitychange", () => {
		if (document.hidden) {
			isHidden = true;

			// suspend audio context
		} else {
			isHidden = false;
			lastFrameTime = 0;

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

		if (!path || !romFiles[path]) return;

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
			const romData = new Uint8Array(arrayBuffer);

			window.loadRom(romData);
			console.log(`Loaded ROM: ${path}`);

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

	// ====================================
	// ====== set up display scaling ======
	// ====================================
	const scaleSelect = document.getElementById(
		"scale-select",
	) as HTMLSelectElement | null;
	scaleSelect!.value = "1"; // default scale

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

		if (!container || !canvas) return;

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

	// ==================================
	// ====== set up download logs ======
	// ==================================
	document
		.getElementById("download-trace-log-button")
		?.addEventListener("click", () => {
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

	// =================================
	// ====== set up joypad input ======
	// =================================
	const keyToJoypadButton: { [key: string]: string } = {
		Enter: "START",
		Backspace: "SELECT",
		z: "B",
		x: "A",
		ArrowDown: "DOWN",
		ArrowUp: "UP",
		ArrowLeft: "LEFT",
		ArrowRight: "RIGHT",
	};
	window.addEventListener("keydown", (event) => {
		if (!keyToJoypadButton[event.key] || event.repeat) {
			return;
		}

		window.handleJoypadButtonPressed(keyToJoypadButton[event.key]);
	});
	window.addEventListener("keyup", (event) => {
		if (!keyToJoypadButton[event.key]) {
			return;
		}

		window.handleJoypadButtonReleased(keyToJoypadButton[event.key]);
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

		if (debugCheckbox.checked && isPaused) {
			updateDebugView();
		}
	});

	syncDebugVisibility();

	// =============================
	// ====== set up debugger ======
	// =============================
	debugElements.cartridgeTitle = document.getElementById("cartridge-title");
	debugElements.mbcType = document.getElementById("cartridge-mbc-type");
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

	const pauseButton = document.getElementById(
		"pause-button",
	) as HTMLButtonElement;
	const stepButton = document.getElementById(
		"step-button",
	) as HTMLButtonElement;

	pauseButton?.addEventListener("click", () => {
		isPaused = !isPaused;
		if (isPaused) {
			pauseButton.textContent = "Run";
			stepButton.disabled = false;
			updateDebugView();
		} else {
			pauseButton.textContent = "Pause";
			stepButton.disabled = true;
			lastFrameTime = 0;
			startAnimationLoop();
		}
	});

	stepButton?.addEventListener("click", () => {
		if (isPaused) {
			window.processEmulatorCycles(4);
			updateDebugView();
		}
	});
});

// ===================================
// ============ GAME LOOP ============
// ===================================

let animationFrameId: number | undefined;

let lastFrameTime = 0;
let tCycleAccumulator = 0;
// 4,194,304 T-cycles per second
const systemClockFrequency = 4.194304 * 1_000_000;

// timestamp is the end time of the previous frame's rendering
function handleAnimationFrame(timestamp: DOMHighResTimeStamp) {
	if (isPaused || isHidden) {
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

	const debuggerEnabled =
		(document.getElementById("debug-checkbox") as HTMLInputElement | null)
			?.checked ?? false;
	if (debuggerEnabled) {
		updateDebugView();
	}

	animationFrameId = requestAnimationFrame(handleAnimationFrame);
}

function startAnimationLoop() {
	cancelAnimationFrame(animationFrameId!);
	animationFrameId = requestAnimationFrame(handleAnimationFrame);
}

function onRomLoaded() {
	startAnimationLoop();
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
	if (debugElements.mbcType)
		debugElements.mbcType.textContent = cartridge.mbcType.toString();
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
