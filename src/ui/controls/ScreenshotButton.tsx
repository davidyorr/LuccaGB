import type { Component } from "solid-js";
import { gameLoop } from "../../core/game-loop";
import { store } from "../../core/store";

export const ScreenshotButton: Component = () => {
	const handleClick = () => {
		if (store.state.isRomLoaded) {
			gameLoop.renderer()?.takeScreenshot();
		}
	};

	return (
		<button onClick={handleClick} style={{ width: "100%" }}>
			Screenshot
		</button>
	);
};
