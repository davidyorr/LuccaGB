import { emulatorState } from "../core/state";

export function setUpAudioVolumeHandlers({
	volumeSliderId,
	volumeValueId,
}: {
	volumeSliderId: string;
	volumeValueId: string;
}) {
	const volumeSlider = document.getElementById(
		volumeSliderId,
	) as HTMLInputElement | null;
	const volumeValue = document.getElementById(
		volumeValueId,
	) as HTMLSpanElement | null;

	if (volumeSlider === null || volumeValue === null) {
		return;
	}

	const setVolume = (volume: number) => {
		if (Number.isInteger(volume)) {
			volumeValue.textContent = `${volume}%`;
			emulatorState.setAudioVolume(volume / 100);
		}
	};

	// Set initial state
	const volume = Number.parseInt(volumeSlider.value);
	setVolume(volume);

	// Handle volume changes
	volumeSlider.addEventListener("input", (event) => {
		const volume = Number.parseInt((event.target as HTMLInputElement).value);
		setVolume(volume);
	});
}
