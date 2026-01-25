import type { Component } from "solid-js";
import { store } from "../core/store";
import { handleRomLoad } from "../services/rom-loader";

export const RomFileInput: Component = () => {
	let fileInputRef: HTMLInputElement | undefined;

	const handleInputChange = async (event: Event) => {
		store.actions.setFileInputOpen(false);

		const target = event.target as HTMLInputElement;
		const file = target.files?.[0];

		if (file) {
			const arrayBuffer = await file.arrayBuffer();
			await handleRomLoad(arrayBuffer);
		}
	};

	const handleInputClick = () => {
		store.actions.setFileInputOpen(true);
	};

	const handleFileDialogCancel = () => {
		store.actions.setFileInputOpen(false);
	};

	return (
		<>
			<label for="rom-input">Load ROM:</label>
			<input
				ref={fileInputRef}
				id="rom-input"
				type="file"
				name="rom-input"
				accept=".gb"
				onChange={handleInputChange}
				onClick={handleInputClick}
				on:cancel={handleFileDialogCancel}
			/>
		</>
	);
};
