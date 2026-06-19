import type { Component } from "solid-js";
import styles from "./RewindControls.module.css";
import { store } from "../../core/store";

export const RewindControls: Component = () => {
	const handleBufferSizeChange = (event: Event) => {
		const target = event.target as HTMLInputElement;
		const value = Math.max(0, Number.parseInt(target.value) || 0);
		store.actions.setRewindBufferSize(value);
	};

	const handleIncrementChange = (event: Event) => {
		const target = event.target as HTMLInputElement;
		const value = Math.max(1, Number.parseInt(target.value) || 1);
		store.actions.setRewindIncrement(value);
	};

	// Convert frames to seconds
	const displaySeconds = (frames: number) => {
		return (frames / 59.727500569606).toFixed(2);
	};

	return (
		<div class={styles.rewindControl}>
			<div class={styles.inputGroup}>
				<label for="rewind-buffer-size">Rewind History Buffer (frames):</label>
				<input
					type="number"
					id="rewind-buffer-size"
					class={styles.numberInput}
					min="0"
					value={store.state.settings.rewindBufferSize}
					onInput={handleBufferSizeChange}
				/>
				<span class={styles.timeValue}>
					~{displaySeconds(store.state.settings.rewindBufferSize)}s
				</span>
			</div>

			<div class={styles.inputGroup}>
				<label for="rewind-increment">Rewind Amount (frames):</label>
				<input
					type="number"
					id="rewind-increment"
					class={styles.numberInput}
					min="1"
					value={store.state.settings.rewindIncrement}
					onInput={handleIncrementChange}
				/>
				<span class={styles.timeValue}>
					~{displaySeconds(store.state.settings.rewindIncrement)}s
				</span>
			</div>
		</div>
	);
};
