import { createEffect, onCleanup, onMount } from "solid-js";
import { appState } from "../core/store";
import { gameLoop } from "../core/game-loop";

export const ViewportScale = () => {
	const canvasRenderer = gameLoop.renderer();

	createEffect(function syncOnGlobalStateChange() {
		canvasRenderer.setScale(appState.scale);
	});

	onMount(function handleWindowResize() {
		const handleResize = () => {
			canvasRenderer.setScale(appState.scale);
		};

		window.addEventListener("resize", handleResize);
		onCleanup(() => window.removeEventListener("resize", handleResize));
	});

	const handleChange = (e: Event) => {
		const target = e.currentTarget as HTMLSelectElement;
		const value = target.value;
		const newScale = value === "fit" ? "fit" : parseInt(value, 10);
		appState.setScale(newScale);
	};

	return (
		<>
			<label for="scale-select">Scale:</label>
			<select id="scale-select" value={appState.scale} onChange={handleChange}>
				<option value="1">1×</option>
				<option value="2">2×</option>
				<option value="3">3×</option>
				<option value="4">4×</option>
				<option value="fit">Fit</option>
			</select>
		</>
	);
};
