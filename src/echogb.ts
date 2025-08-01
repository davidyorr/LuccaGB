(window as any).onRomLoaded = onRomLoaded;
(window as any).updateCanvas = updateCanvas;

let canvasCtx: CanvasRenderingContext2D;
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
	const canvas = document.getElementById("canvas");
	if (canvas instanceof HTMLCanvasElement) {
		const ctx = canvas.getContext("2d");
		if (!ctx) {
			throw new Error("error getting canvas context");
		}
		canvas.width = displayWidth;
		canvas.height = displayHeight;
		ctx.imageSmoothingEnabled = false;
		imageData = ctx.createImageData(displayWidth, displayHeight);
		canvasCtx = ctx;
	}
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
	imageData.data.set(uint8Array);
	canvasCtx.putImageData(imageData, 0, 0);
}
