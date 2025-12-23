declare global {
	interface Window {
		loadRom: (data: Uint8Array) => void;
		processEmulatorCycles: (cycles: number) => {
			tCyclesUsed: number;
		};
		getTraceLogs: () => Uint8Array;
		getDebugInfo: () => GameboyDebugInfo | null;
	}
}

interface GameboyDebugInfo {
	cartridge: CartridgeDebugInfo;
	cpu: CpuDebugInfo;
	// ppu: PpuDebugInfo;
}

interface CartridgeDebugInfo {
	title: string;
	mbcType: number;
	romSizeCode: number;
	ramSizeCode: number;
}

interface CpuDebugInfo {
	registers16: {
		AF: number;
		BC: number;
		DE: number;
		HL: number;
		SP: number;
		PC: number;
	};
	flags: {
		Z: number;
		N: number;
		H: number;
		C: number;
	};
}

export {};
