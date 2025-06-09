const go = new Go();

WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then(
	(wasmModule) => {
		go.run(wasmModule.instance);
	},
);

document.addEventListener("DOMContentLoaded", () => {
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
