declare global {
	interface Window {
		loadRom: (data: Uint8Array) => void;
		processEmulatorCycles: (cycles: number) => {
			tCyclesUsed: number;
		};
	}
}

export {};
