(window as any).onRomLoaded = onRomLoaded;
(window as any).updateCanvas = updateCanvas;

let visibleCanvasCtx: CanvasRenderingContext2D;
let offscreenCanvasCtx: CanvasRenderingContext2D;
let imageData: ImageData;

const displayWidth = 160;
const displayHeight = 144;

const go = new Go();

document.addEventListener("DOMContentLoaded", () => {
	WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then(
		(wasmModule) => {
			go.run(wasmModule.instance);
		},
	);

	// ====== set up event listener ======
	document
		.getElementById("rom-input")
		?.addEventListener("change", async (event) => {
			const files = (event.target as HTMLInputElement | null)?.files;
			if (files?.[0]) {
				const arrayBuffer = await files?.[0].arrayBuffer();
				const romData = new Uint8Array(arrayBuffer);
				window.loadRom(romData);
			}
		});

	// ====== set up canvas ======
	const visibleCanvas = document.getElementById("canvas");
	if (visibleCanvas instanceof HTMLCanvasElement) {
		const ctx = visibleCanvas.getContext("2d");
		if (!ctx) {
			throw new Error("error getting canvas context");
		}
		visibleCanvasCtx = ctx;
		visibleCanvasCtx.imageSmoothingEnabled = false;
		imageData = visibleCanvasCtx.createImageData(displayWidth, displayHeight);

		const resizeObserver = new ResizeObserver((entries) => {
			const entry = entries[0];
			const { width, height } = entry.contentRect;

			visibleCanvas.width = width;
			visibleCanvas.height = height;
			visibleCanvasCtx.imageSmoothingEnabled = false;
		});

		const canvasContainer = document.getElementById("canvas-container");
		if (!canvasContainer) {
			throw new Error("error getting canvas container");
		}
		resizeObserver.observe(canvasContainer);
	}
	const offscreenCanvas = document.createElement("canvas");
	offscreenCanvas.width = 160;
	offscreenCanvas.height = 144;
	const ctx = offscreenCanvas.getContext("2d");
	if (!ctx) {
		throw new Error("error getting offscreen canvas context");
	}
	offscreenCanvasCtx = ctx;
});

// ============ GAME LOOP ============

let animationFrameId: number | undefined;

let lastFrameTime = 0;
let tCycleAccumulator = 0;
// 4,194,304 T-cycles per second
const systemClockFrequency = 4.194304 * 1_000_000;

// timestamp is the end time of the previous frame's rendering
function handleAnimationFrame(timestamp: DOMHighResTimeStamp) {
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
