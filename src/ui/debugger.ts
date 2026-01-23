import { appState } from "../core/state";

export class Debugger {
	private debugPanel: HTMLElement | null;

	private cartridgeTitle: HTMLElement | null;
	private cartridgeType: HTMLElement | null;
	private romSize: HTMLElement | null;
	private ramSize: HTMLElement | null;
	private regAF: HTMLElement | null;
	private regBC: HTMLElement | null;
	private regDE: HTMLElement | null;
	private regHL: HTMLElement | null;
	private regSP: HTMLElement | null;
	private regPC: HTMLElement | null;
	private flagZ: HTMLElement | null;
	private flagN: HTMLElement | null;
	private flagH: HTMLElement | null;
	private flagC: HTMLElement | null;

	// APU
	private waveRamHex: Array<HTMLElement | null> = [];
	private regNR30: HTMLElement | null;
	private regNR32: HTMLElement | null;

	constructor() {
		this.debugPanel = document.getElementById("debug-panel");
		this.cartridgeTitle = document.getElementById("cartridge-title");
		this.cartridgeType = document.getElementById("cartridge-cartridge-type");
		this.romSize = document.getElementById("cartridge-rom-size-code");
		this.ramSize = document.getElementById("cartridge-ram-size-code");
		this.regAF = document.getElementById("reg-af");
		this.regBC = document.getElementById("reg-bc");
		this.regDE = document.getElementById("reg-de");
		this.regHL = document.getElementById("reg-hl");
		this.regSP = document.getElementById("reg-sp");
		this.regPC = document.getElementById("reg-pc");
		this.flagZ = document.getElementById("flag-z");
		this.flagN = document.getElementById("flag-n");
		this.flagH = document.getElementById("flag-h");
		this.flagC = document.getElementById("flag-c");
		this.regNR30 = document.getElementById("nr30");
		this.regNR32 = document.getElementById("nr32");

		for (let i = 0; i < 16; i++) {
			this.waveRamHex.push(document.getElementById(`wave-hex-${i}`));
		}

		// Handle state changes
		appState.subscribe((state) => {
			this.handleVisibility(state.isDebuggerOpen);

			if (state.isDebuggerOpen) {
				this.update();
			}
		});
	}

	private handleVisibility(isOpen: boolean) {
		if (!this.debugPanel) {
			return;
		}

		this.debugPanel.style.display = isOpen ? "block" : "none";
	}

	public update() {
		if (!appState.isDebuggerOpen) {
			return;
		}

		const debugInfo = window.getDebugInfo?.();
		if (!debugInfo) {
			return;
		}

		const { apu, cartridge, cpu } = debugInfo;

		const toHex = (val: number) =>
			"0x" + val.toString(16).toUpperCase().padStart(4, "0");
		const toBinary = (val: number) =>
			"0b" +
			val
				.toString(2)
				.padStart(8, "0")
				.replace(/(\d{4})(?=\d)/g, "$1_");

		if (this.cartridgeTitle) this.cartridgeTitle.textContent = cartridge.title;
		if (this.cartridgeType)
			this.cartridgeType.textContent = cartridge.cartridgeType.toString();
		if (this.romSize)
			this.romSize.textContent = cartridge.romSizeCode.toString();
		if (this.ramSize)
			this.ramSize.textContent = cartridge.ramSizeCode.toString();

		if (this.regAF) this.regAF.textContent = toHex(cpu.registers16.AF);
		if (this.regBC) this.regBC.textContent = toHex(cpu.registers16.BC);
		if (this.regDE) this.regDE.textContent = toHex(cpu.registers16.DE);
		if (this.regHL) this.regHL.textContent = toHex(cpu.registers16.HL);
		if (this.regSP) this.regSP.textContent = toHex(cpu.registers16.SP);
		if (this.regPC) this.regPC.textContent = toHex(cpu.registers16.PC);

		if (this.flagZ) this.flagZ.textContent = cpu.flags.Z.toString();
		if (this.flagN) this.flagN.textContent = cpu.flags.N.toString();
		if (this.flagH) this.flagH.textContent = cpu.flags.H.toString();
		if (this.flagC) this.flagC.textContent = cpu.flags.C.toString();

		if (this.regNR30) this.regNR30.textContent = toBinary(apu.registers.NR30);
		if (this.regNR32) this.regNR32.textContent = toBinary(apu.registers.NR32);

		// wave RAM
		if (apu.waveRam && apu.waveRam.length === 16) {
			for (let i = 0; i < 16; i++) {
				const value = apu.waveRam[i];

				if (this.waveRamHex[i]) {
					this.waveRamHex[i]!.textContent = `0x${toHex(value)}`;
				}
			}
		}
	}
}
