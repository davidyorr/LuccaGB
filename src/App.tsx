import { onMount, onCleanup, type Component } from "solid-js";
import { store } from "./core/store";
import { Controls } from "./ui/controls/Controls";
import { Debugger } from "./ui/Debugger";
import { DragAndDropOverlay } from "./ui/DragAndDropOverlay";
import { gameLoop } from "./core/game-loop";
import { CanvasRenderer } from "./services/canvas-renderer";
import { Viewport } from "./ui/viewport/Viewport";

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

			<Viewport />

			<Debugger />
		</>
	);
};
