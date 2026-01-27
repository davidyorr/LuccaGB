import styles from "./VolumeControl.module.css";

import type { Component } from "solid-js";
import { store } from "../../core/store";
import { audioController } from "../../services/audio-controller";

export const VolumeControl: Component = () => {
	const displayVolume = () =>
		Math.round(store.state.settings.audioVolume * 100);

	return (
		<>
			<label for="volume-slider">Volume:</label>
			<input
				type="range"
				id="volume-slider"
				class={styles.volumeSlider}
				min="0"
				max="100"
				value={displayVolume()}
				onInput={(event) => {
					const value = Number.parseInt(event.target.value);
					if (Number.isInteger(value)) {
						const engineValue = value / 100;
						store.actions.setAudioVolume(engineValue);
						audioController.setVolume(engineValue);
					}
				}}
			/>
			<span class={styles.volumeValue}>{displayVolume()}%</span>
		</>
	);
};
