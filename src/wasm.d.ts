export interface WasmExports extends WebAssembly.Exports {
	processEmulatorStep: () => void;
}

declare global {
	interface Window {
		loadRom: (data: Uint8Array) => void;
	}
}

export {};
