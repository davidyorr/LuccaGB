import type { WasmExports } from "./wasm";

let wasmExports: WasmExports;

(window as any).onRomLoaded = onRomLoaded;

document.addEventListener("DOMContentLoaded", () => {
	const go = new Go();

	WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then(
		(wasmModule) => {
			wasmExports = wasmModule.instance.exports as WasmExports;
			go.run(wasmModule.instance);
		},
	);

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
});

// ============ GAME LOOP ============

let animationFrameId: number | undefined;

function handleAnimationFrame() {
	wasmExports.processEmulatorStep();
}

function startAnimationLoop() {
	cancelAnimationFrame(animationFrameId!);
	handleAnimationFrame();
}

function onRomLoaded() {
	startAnimationLoop();
}
