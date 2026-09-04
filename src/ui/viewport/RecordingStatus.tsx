import styles from "./ViewportStatus.module.css";

import { Show, type Component } from "solid-js";
import { store } from "../../core/store";

export const RecordingStatus: Component = () => {
	return (
		<Show when={store.state.isRomLoaded && store.state.isRecordingGameplay}>
			<div class={styles.container}>
				🔴 Recording {formatDuration(store.state.recordingDuration)}
				{store.state.recordingSize > 0
					? " · " + formatBytes(store.state.recordingSize)
					: ""}
			</div>
		</Show>
	);
};

function formatDuration(ms: number) {
	const totalSeconds = Math.floor(ms / 1000);
	const minutes = Math.floor(totalSeconds / 60);
	const seconds = totalSeconds % 60;

	return `${minutes}:${seconds.toString().padStart(2, "0")}`;
}

function formatBytes(bytes: number) {
	if (bytes < 1024) {
		return `${bytes} B`;
	}

	const kb = bytes / 1024;

	if (kb < 1024) {
		return `${kb.toFixed(1)} KB`;
	}

	const mb = kb / 1024;

	return `${mb.toFixed(1)} MB`;
}
