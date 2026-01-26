import { createEffect, onCleanup, onMount } from "solid-js";
import { gameLoop } from "../../core/game-loop";
import { store } from "../../core/store";

export const ViewportScale = () => {
	createEffect(function syncOnGlobalStateChange() {
		gameLoop.renderer()?.setScale(store.state.settings.scale);
	});

	onMount(function handleWindowResize() {
		const handleResize = () => {
			gameLoop.renderer()?.setScale(store.state.settings.scale);
		};

		window.addEventListener("resize", handleResize);
		onCleanup(() => window.removeEventListener("resize", handleResize));
	});

	const handleChange = (e: Event) => {
		const target = e.currentTarget as HTMLSelectElement;
		const value = target.value;
		const newScale = value === "fit" ? "fit" : parseInt(value, 10);
		store.actions.setScale(newScale);
	};

	return (
		<>
			<label for="scale-select">Scale:</label>
			<select
				id="scale-select"
				value={store.state.settings.scale}
				onChange={handleChange}
			>
				<option value="1">1×</option>
				<option value="2">2×</option>
				<option value="3">3×</option>
				<option value="4">4×</option>
				<option value="fit">Fit</option>
			</select>
		</>
	);
};
