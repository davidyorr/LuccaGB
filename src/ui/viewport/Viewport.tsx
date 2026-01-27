import styles from "./Viewport.module.css";

import type { Component } from "solid-js";
import { AudioHint } from "./AudioHint";

export const Viewport: Component = () => {
	return (
		<>
			<div id="canvas-container" class={styles.container}>
				<canvas id="canvas" class={styles.canvas} width="160" height="144" />
				<AudioHint />
			</div>
		</>
	);
};
