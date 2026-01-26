import { onMount, onCleanup, type Component } from "solid-js";
import { store } from "./core/store";
import { Controls } from "./ui/controls/Controls";
import { Debugger } from "./ui/Debugger";
import { DragAndDropOverlay } from "./ui/DragAndDropOverlay";
import { gameLoop } from "./core/game-loop";
import { CanvasRenderer } from "./services/canvas-renderer";

export const App: Component = () => {
	onMount(() => {
		const handleVisibilityChange = () => {
			store.actions.setHidden(document.hidden);
		};
		document.addEventListener("visibilitychange", handleVisibilityChange);

		onCleanup(() => {
			document.removeEventListener("visibilitychange", handleVisibilityChange);
		});
	});

	onMount(function initializeApp() {
		gameLoop.attachRenderer(new CanvasRenderer());
	});

	return (
		<>
			<DragAndDropOverlay />

			<Controls />

			<div id="canvas-container">
				<canvas id="canvas" width="160" height="144" />
			</div>

			<div id="debug-container">
				<Debugger />
			</div>
		</>
	);
};
