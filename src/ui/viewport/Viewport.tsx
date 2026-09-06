import styles from "./Viewport.module.css";

import type { Component } from "solid-js";
import { store } from "../../core/store";
import { ViewportStatus } from "./ViewportStatus";

export const Viewport: Component = () => {
	const activePalette = store.getters.activePalette;

	return (
		<>
			<div
				id="canvas-container"
				class={styles.container}
				style={{
					// lightest color
					"background-color": activePalette().colors[0],
				}}
			>
				<canvas id="canvas" class={styles.canvas} width="160" height="144" />
				<ViewportStatus />
			</div>
		</>
	);
};
