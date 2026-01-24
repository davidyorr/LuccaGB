import { createStore, unwrap } from "solid-js/store";
import { createRoot, createEffect, on, batch } from "solid-js";
import {
	loadAppSettings,
	persistCartridgeRam,
	saveAppSettings,
} from "../services/storage";
import { audioController } from "../services/audio-controller";
import { debounce } from "../utils/debounce";
import type { CartridgeInfo } from "../wasm";
import { gameLoop } from "./game-loop";

type State = {
	isPaused: boolean;
	isRomLoaded: boolean;
	currentRomHash: string;
	cartridgeInfo: CartridgeInfo | null;
	settings: {
		/** The value that gets passed to the gain node between `[0, 1.0]` */
		audioVolume: number;
		audioChannelsEnabled: boolean[];
		isDebuggerOpen: boolean;
		scale: number | "fit";
	};
	ui: {
		isFileInputOpen: boolean;
		isHidden: boolean;
		isControlsOpen: boolean;
	};
};

const defaultSettings = {
	audioVolume: 0.5,
	audioChannelsEnabled: [false, true, true, true, true],
	isDebuggerOpen: false,
	scale: 3 as const,
};

export const store = createRoot(() => {
	const [state, setState] = createStore<State>({
		isPaused: false,
		isRomLoaded: false,
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

		setFileInputOpen: (isOpen: boolean) => {
			setState("ui", "isFileInputOpen", isOpen);
		},

		setHidden: (hidden: boolean) => {
			setState("ui", "isHidden", hidden);
		},

		setControlsOpen: (isOpen: boolean) => {
			setState("ui", "isControlsOpen", isOpen);
		},
	};

	// =========================================
	// EFFECTS
	// =========================================

	const isRunning = () =>
		state.isRomLoaded &&
		!state.isPaused &&
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
		if (state.isPaused || state.ui.isHidden) {
			audioController.pause();
		} else {
			audioController.resume();
		}
	});

	createEffect(function pauseWhenFileInputOpen() {
		if (state.ui.isFileInputOpen) {
			setState("isPaused", true);
		}
	});

	// =========================================
	// LEGACY BRIDGE
	// =========================================

	const legacyAppState = {
		// Getters proxy to the store
		get isPaused() {
			return state.isPaused;
		},
		get isRomLoaded() {
			return state.isRomLoaded;
		},
		get currentRomHash() {
			return state.currentRomHash;
		},
		get cartridgeInfo() {
			return state.cartridgeInfo;
		},
		get audioVolume() {
			return state.settings.audioVolume;
		},
		get audioChannelsEnabled() {
			return state.settings.audioChannelsEnabled;
		},
		get isDebuggerOpen() {
			return state.settings.isDebuggerOpen;
		},
		get scale() {
			return state.settings.scale;
		},
		get isFileInputOpen() {
			return state.ui.isFileInputOpen;
		},
		get isHidden() {
			return state.ui.isHidden;
		},
		get isControlsOpen() {
			return state.ui.isControlsOpen;
		},

		// Methods proxy to actions
		setPaused: actions.setPaused,
		togglePaused: actions.togglePaused,
		setRomLoaded: actions.setRomLoaded,
		setCurrentRomHash: actions.setCurrentRomHash,
		setAudioVolume: actions.setAudioVolume,
		setAudioChannelsEnabled: actions.setAudioChannelsEnabled,
		setDebuggerOpen: actions.setDebuggerOpen,
		setScale: actions.setScale,
		setFileInputOpen: actions.setFileInputOpen,
		setHidden: actions.setHidden,
		setControlsOpen: actions.setControlsOpen,

		initializeAppSettings: async () => {
			const appSettings = await loadAppSettings();
			if (appSettings) {
				setState("settings", appSettings);
			}
		},

		getSettingsSnapshot: () => ({
			...unwrap(state.settings),
			updatedAt: Date.now(),
		}),

		// Backward Compatibility Subscribe
		subscribe: (listener: (s: any) => void) => {
			// We create a nested effect that tracks the whole state object
			// and calls the legacy listener when something changes.
			createEffect(() => {
				// access state to track it
				const rawState = unwrap(state);

				listener(legacyAppState);
			});
		},
	};

	return {
		state,
		actions,
		legacyAppState,
	};
});

export const appState = store.legacyAppState;
