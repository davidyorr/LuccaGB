import styles from "./AudioHint.module.css";

import { Show, type Component } from "solid-js";
import { store } from "../../core/store";
import { audioController } from "../../services/audio-controller";

export const AudioHint: Component = () => {
	return (
		<Show when={store.state.isRomLoaded && !audioController.unlocked}>
			<div class={styles.container}>Press any key to enable audio</div>
		</Show>
	);
};
