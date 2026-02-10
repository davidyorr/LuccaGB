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
		getSerializedState: () => Uint8Array;
		loadSerializedState: (data: Uint8Array) => void;
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

export interface GameboyDebugInfo {
	apu: ApuDebugInfo;
	cartridge: CartridgeDebugInfo;
	cpu: CpuDebugInfo;
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

export async function initWasm() {
	const go = new Go();
	const wasmModule = await WebAssembly.instantiateStreaming(
		fetch("main.wasm"),
		go.importObject,
	);
	go.run(wasmModule.instance);
}
