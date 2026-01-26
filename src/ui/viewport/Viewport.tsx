import type { Component } from "solid-js";
import { AudioHint } from "./AudioHint";

export const Viewport: Component = () => {
	return (
		<>
			<div id="canvas-container">
				<canvas id="canvas" width="160" height="144" />
				<AudioHint />
			</div>
		</>
	);
};
