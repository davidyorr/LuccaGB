import { createStore, unwrap } from "solid-js/store";
import { createEffect, on, batch } from "solid-js";
import {
	loadAppSettings,
	persistCartridgeRam,
	saveAppSettings,
} from "../services/storage";
import { audioController } from "../services/audio-controller";
import { debounce } from "../utils/debounce";
import type { CartridgeInfo } from "../core/wasm";
import { gameLoop } from "./game-loop";
import { updateDebugger } from "../ui/Debugger";

export type GameBoyPalette = {
	id: string;
	name: string;
	colors: [string, string, string, string];
	isReadonly?: boolean;
};

export type State = {
	isPaused: boolean;
	isRomLoaded: boolean;
	isRewinding: boolean;
	currentRomHash: string;
	cartridgeInfo: CartridgeInfo | null;

	/** The settings that get saved to IndexedDB */
	settings: {
		/** The value that gets passed to the gain node between `[0, 1.0]` */
		audioVolume: number;
		audioChannelsEnabled: boolean[];
		isDebuggerOpen: boolean;
		scale: number | "fit";
		updatedAt?: number;
		activePaletteId: string;
		palettes: GameBoyPalette[];
		/** How many frames to save */
		rewindBufferSize: number;
		/** How many frames to rewind per tick */
		rewindIncrement: number;
	};

	ui: {
		isFileInputOpen: boolean;
		isHidden: boolean;
		isControlsOpen: boolean;
	};
};

const defaultPalettes: GameBoyPalette[] = [
	{
		id: "default-green",
		name: "OG Green",
		colors: ["#d0e040", "#a0a830", "#607028", "#384828"],
		isReadonly: true,
	},
	{
		id: "lucca-gb",
		name: "LuccaGB",
		colors: ["#FFFDF1", "#FFCE99", "#FF9644", "#562F00"],
		isReadonly: true,
	},
	{
		id: "grayscale",
		name: "Grayscale",
		colors: ["#ffffff", "#aaaaaa", "#555555", "#000000"],
		isReadonly: true,
	},
];

const defaultSettings = {
	audioVolume: 0.5,
	audioChannelsEnabled: [false, true, true, true, true],
	isDebuggerOpen: false,
	scale: 3 as const,
	activePaletteId: "default-green",
	palettes: defaultPalettes,
	rewindBufferSize: 600,
	rewindIncrement: 1,
};

const [state, setState] = createStore<State>({
	isPaused: false,
	isRomLoaded: false,
	isRewinding: false,
	currentRomHash: "",
	cartridgeInfo: null,
	settings: { ...defaultSettings },
	ui: {
		isFileInputOpen: false,
		isHidden: false,
		isControlsOpen: false,
	},
});

const actions = {
	initializeAppSettings: async () => {
		const settings = await loadAppSettings();
		if (settings) {
			// Extract only the user's custom palettes from their save data
			const savedCustomPalettes = (settings.palettes || []).filter(
				(palette) => !palette.isReadonly,
			);

			// Merge the latest hardcoded defaults with the user's custom ones
			const mergedPalettes = [...defaultPalettes, ...savedCustomPalettes];

			setState("settings", {
				...defaultSettings,
				...settings,
				palettes: mergedPalettes,
			});
		}
	},

	setPaused: (paused: boolean) => {
		setState("isPaused", paused);
	},

	togglePaused: () => {
		if (state.isRomLoaded) {
			setState("isPaused", (paused) => !paused);
		}
	},

	setRomLoaded: (loaded: boolean) => {
		batch(() => {
			setState("isRomLoaded", loaded);
			if (loaded) {
				setState("isPaused", false);
				setState("ui", { isFileInputOpen: false, isControlsOpen: false });
			}
		});
	},

	setCurrentRomHash: (hash: string) => {
		setState("currentRomHash", hash);
	},

	setCartridgeInfo: (info: CartridgeInfo) => {
		setState("cartridgeInfo", info);
	},

	setRewinding: (rewinding: boolean) => {
		setState("isRewinding", rewinding);
	},

	setAudioVolume: (vol: number) => {
		setState("settings", "audioVolume", vol);
	},

	setAudioChannelsEnabled: (channel: number, enabled: boolean) => {
		if (channel >= 1 && channel <= 4) {
			setState("settings", "audioChannelsEnabled", channel, enabled);
		}
	},

	setDebuggerOpen: (isOpen: boolean) => {
		setState("settings", "isDebuggerOpen", isOpen);
	},

	setScale: (scale: number | "fit") => {
		setState("settings", "scale", scale);
	},

	setActivePaletteId: (id: string) => {
		setState("settings", "activePaletteId", id);
	},

	addPalette: (palette: GameBoyPalette) => {
		setState("settings", "palettes", (prev) => [...prev, palette]);
		setState("settings", "activePaletteId", palette.id);
	},

	updateActivePaletteColor: (index: number, hexColor: string) => {
		const activeId = state.settings.activePaletteId;
		const activeIndex = state.settings.palettes.findIndex(
			(palette) => palette.id === activeId,
		);
		if (activeIndex === -1 || state.settings.palettes[activeIndex].isReadonly)
			return;

		setState("settings", "palettes", activeIndex, "colors", index, hexColor);
	},

	renameActivePalette: (newName: string) => {
		const activeId = state.settings.activePaletteId;
		const activeIndex = state.settings.palettes.findIndex(
			(palette) => palette.id === activeId,
		);

		if (activeIndex === -1 || state.settings.palettes[activeIndex].isReadonly)
			return;

		setState("settings", "palettes", activeIndex, "name", newName);
	},

	deletePalette: (id: string) => {
		const index = state.settings.palettes.findIndex(
			(palette) => palette.id === id,
		);
		if (index > -1 && !state.settings.palettes[index].isReadonly) {
			setState("settings", "palettes", (prev) =>
				prev.filter((palette) => palette.id !== id),
			);
			// Fallback to the default green if they delete their active palette
			setState("settings", "activePaletteId", "default-green");
		}
	},

	setFileInputOpen: (isOpen: boolean) => {
		setState("ui", "isFileInputOpen", isOpen);
	},

	setHidden: (hidden: boolean) => {
		setState("ui", "isHidden", hidden);
	},

	setControlsOpen: (isOpen: boolean) => {
		setState("ui", "isControlsOpen", isOpen);
	},

	setRewindBufferSize: (size: number) => {
		setState("settings", "rewindBufferSize", size);
	},

	setRewindIncrement: (increment: number) => {
		setState("settings", "rewindIncrement", increment);
	},
};

const getters = {
	activePalette: () => {
		return (
			state.settings.palettes.find(
				(palette) => palette.id === store.state.settings.activePaletteId,
			) ?? state.settings.palettes[0]
		);
	},
};

// =========================================
// EFFECTS
// =========================================

const isRunning = () =>
	state.isRomLoaded &&
	!state.isPaused &&
	!state.isRewinding &&
	!state.ui.isHidden &&
	!state.ui.isFileInputOpen;

createEffect(function handleIsRunning() {
	if (isRunning()) {
		gameLoop.start();
	} else {
		gameLoop.stop();
	}
});

createEffect(
	on(
		() => state.isPaused,
		function saveRamOnPause(paused) {
			if (
				paused &&
				state.cartridgeInfo?.hasBattery &&
				state.cartridgeInfo.ramSize > 0
			) {
				const ram = window.getCartridgeRam();
				persistCartridgeRam(state.currentRomHash, ram, {
					name: state.cartridgeInfo.title,
				});
			}
		},
	),
);

createEffect(
	on(
		() => state.isPaused,
		function updateDebuggerOnPause(paused) {
			if (paused) {
				updateDebugger();
			}
		},
	),
);

const debouncedSave = debounce((settings) => saveAppSettings(settings), 333);
createEffect(function syncSettingsToStorage() {
	JSON.stringify(state.settings);
	const snapshot = unwrap(state.settings);
	debouncedSave({
		...snapshot,
		updatedAt: Date.now(),
	});
});

createEffect(function syncAudioPauseState() {
	if (state.isRomLoaded) {
		if (state.isPaused || state.ui.isHidden) {
			audioController.pause();
		} else {
			audioController.resume();
		}
	}
});

createEffect(function pauseWhenFileInputOpen() {
	if (state.ui.isFileInputOpen) {
		setState("isPaused", true);
	}
});

createEffect(
	on(
		() => [state.settings.rewindBufferSize, state.isRomLoaded],
		([rewindBufferSize, isRomLoaded]) => {
			if (isRomLoaded && window.setRewindBufferSize) {
				window.setRewindBufferSize(rewindBufferSize as number);
			}
		},
	),
);

createEffect(
	on(
		() => {
			const activeId = state.settings.activePaletteId;
			const activePalette = state.settings.palettes.find(
				(palette) => palette.id === activeId,
			);
			const colors = activePalette
				? activePalette.colors
				: state.settings.palettes[0].colors;

			return [[...colors], state.isRomLoaded];
		},
		([colors, isRomLoaded]) => {
			if (isRomLoaded && window.setPalette) {
				// Convert ["#hex", "#hex", "#hex", "#hex"] -> [r,g,b,a, r,g,b,a, r,g,b,a, r,g,b,a]
				const flatRgbaArray = (colors as string[]).flatMap((hex) => {
					const cleanHex = hex.replace("#", "");
					const r = parseInt(cleanHex.substring(0, 2), 16);
					const g = parseInt(cleanHex.substring(2, 4), 16);
					const b = parseInt(cleanHex.substring(4, 6), 16);
					return [r, g, b, 255];
				});
				window.setPalette(flatRgbaArray);

				// If emulation is paused, the game loop isn't running to pull this frame automatically,
				// so we need to manually do it
				if (state.isPaused) {
					const frame = window.pollFrame();
					if (frame) {
						gameLoop.forceDraw(frame);
					}
				}
			}
		},
	),
);

export const store = {
	state,
	getters,
	actions,
};
