declare global {
	interface Window {
		loadRom: (data: Uint8Array) => void;
	}
}

export {};
