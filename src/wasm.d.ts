declare global {
	interface Window {
		loadRom: (data: Uint8Array) => CartridgeInfo;
		setCartridgeRam: (data: Uint8Array | null) => void;
		getCartridgeRam: () => Uint8Array;
		processEmulatorCycles: (cycles: number) => {
			tCyclesUsed: number;
		};
		pollFrame: () => Uint8Array;
		pollAudioBuffer: () => Array<number>;
		handleJoypadButtonPressed: (button: string) => void;
		handleJoypadButtonReleased: (button: string) => void;
		enableTraceLogging: () => void;
		disableTraceLogging: () => void;
		getTraceLogs: () => Uint8Array;
		setAudioChannelEnabled: (channel: number, enabled: boolean) => void;
		getAudioChannelEnabled: (channel: number) => boolean;
		getDebugInfo: () => GameboyDebugInfo | null;
	}
}

export interface CartridgeInfo {
	title: string;
	ramSize: number;
	hasBattery: boolean;
}

interface GameboyDebugInfo {
	apu: ApuDebugInfo;
	cartridge: CartridgeDebugInfo;
	cpu: CpuDebugInfo;
	// ppu: PpuDebugInfo;
}

interface ApuDebugInfo {
	waveRam: number[];
	registers: {
		NR30: number;
		NR32: number;
	};
}

interface CartridgeDebugInfo {
	title: string;
	cartridgeType: number;
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
