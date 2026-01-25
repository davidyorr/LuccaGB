import { store } from "../core/store";
import { updateDebugger } from "../ui/Debugger";
import type { CartridgeInfo } from "../wasm";
import { audioController } from "./audio-controller";
import { loadCartridgeRam } from "./storage";

let cartridgeInfo: CartridgeInfo | null = null;

export async function handleRomLoad(arrayBuffer: ArrayBuffer) {
	const romData = new Uint8Array(arrayBuffer);

	// pause the previous audio context
	await audioController.pause();

	// Compute the ROM Hash
	const hashBuffer = await window.crypto.subtle.digest(
		"SHA-256",
		romData.buffer,
	);
	const hashArray = Array.from(new Uint8Array(hashBuffer));
	const hashHex = hashArray
		.map((b) => b.toString(16).padStart(2, "0"))
		.join("");
	store.actions.setCurrentRomHash(hashHex);

	// Load into Go
	cartridgeInfo = window.loadRom(romData);
	store.actions.setCartridgeInfo(cartridgeInfo);
	console.log("Cartridge Info:", cartridgeInfo);

	// Attempt to load existing RAM
	if (cartridgeInfo.hasBattery && cartridgeInfo.ramSize > 0) {
		try {
			const ram = await loadCartridgeRam(store.legacyAppState.currentRomHash);
			if (ram) {
				// Ensure the loaded RAM size matches what the cartridge expects
				if (ram.length !== cartridgeInfo.ramSize) {
					console.warn(
						`Save file size mismatch. Expected ${cartridgeInfo.ramSize}, got ${ram.length}`,
					);
				}
				window.setCartridgeRam(ram);
			}
		} catch (e) {
			console.error("Failed to load save data:", e);
		}
	}

	// Focus the canvas so keyboard controls work immediately
	const canvas = document.getElementById("canvas");
	if (canvas) {
		canvas.tabIndex = 0;
		canvas.focus();
	}

	// set the initial audio channels state
	window.setAudioChannelEnabled(
		1,
		store.state.settings.audioChannelsEnabled[1],
	);
	window.setAudioChannelEnabled(
		2,
		store.state.settings.audioChannelsEnabled[2],
	);
	window.setAudioChannelEnabled(
		3,
		store.state.settings.audioChannelsEnabled[3],
	);
	window.setAudioChannelEnabled(
		4,
		store.state.settings.audioChannelsEnabled[4],
	);

	// resume the audio context
	await audioController.resume();

	// Start the animation loop
	store.actions.setRomLoaded(true);
	updateDebugger();
}
