import { appState } from "../core/state";

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

	// Handle volume changes
	volumeSlider.addEventListener("input", (event) => {
		const volume = Number.parseInt((event.target as HTMLInputElement).value);
		if (Number.isInteger(volume)) {
			volumeValue.textContent = `${volume}%`;
			appState.setAudioVolume(volume / 100);
		}
	});

	appState.subscribe((state) => {
		const volume = state.audioVolume * 100;
		volumeSlider.value = volume.toString();
		volumeValue.textContent = `${volume}%`;
	});
}
