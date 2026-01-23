import { appState } from "../core/state";
import type { AudioController } from "../services/audio-controller";

export function setUpAudioChannelHandlers({
	buttonId,
	dropdownId,
	audioController,
}: {
	buttonId: string;
	dropdownId: string;
	audioController: AudioController;
}) {
	const button = document.getElementById(buttonId);
	const dropdown = document.getElementById(dropdownId);

	if (button === null || dropdown === null) {
		return;
	}

	// Toggle dropdown
	button.addEventListener("click", (event) => {
		event.stopPropagation();
		dropdown.classList.toggle("show");
	});

	// Close dropdown when clicking outside
	document.addEventListener("click", (event) => {
		if (!dropdown.contains(event.target as any) && event.target !== button) {
			dropdown.classList.remove("show");
		}
	});

	// Handle channel toggles
	for (let i = 1; i <= 4; i++) {
		const checkbox = document.getElementById(
			`audio-channel-${i}`,
		) as HTMLInputElement | null;
		// set the initial state
		appState.setAudioChannelEnabled(i, checkbox?.checked ?? true);
		// handle on change
		checkbox?.addEventListener("change", (event) => {
			appState.setAudioChannelEnabled(
				i,
				(event.target as HTMLInputElement)?.checked,
			);
		});
	}

	let prevAudioChannelsEnabled = [...appState.audioChannelsEnabled];

	// Handle app state changes
	appState.subscribe(async (state) => {
		// check if any audio channels were enabled/disabled
		for (let i = 1; i <= 4; i++) {
			if (state.audioChannelsEnabled[i] !== prevAudioChannelsEnabled[i]) {
				window.setAudioChannelEnabled(i, state.audioChannelsEnabled[i]);
				prevAudioChannelsEnabled[i] = state.audioChannelsEnabled[i];
				audioController.resetTime();
			}
		}
	});
}
