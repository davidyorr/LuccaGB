import { createSignal, onCleanup, onMount, type Component } from "solid-js";
import { handleRomLoad } from "../services/rom-loader";

export const DragAndDropOverlay: Component = () => {
	const [isVisible, setIsVisible] = createSignal(false);

	// Window Listener: Detects when a file enters the window
	onMount(() => {
		const handleWindowDragEnter = (event: DragEvent) => {
			// Only show if dragging files
			if (event.dataTransfer?.types.includes("Files")) {
				setIsVisible(true);
			}
		};

		window.addEventListener("dragenter", handleWindowDragEnter);

		onCleanup(() => {
			window.removeEventListener("dragenter", handleWindowDragEnter);
		});
	});

	// Overlay Handlers: Handle the drop once the overlay is visible
	const handleDragOver = (event: DragEvent) => {
		event.preventDefault();
		if (event.dataTransfer) {
			event.dataTransfer.dropEffect = "copy";
		}
	};

	// Hide overlay if user leaves the window (dragging out)
	const handleDragLeave = (event: DragEvent) => {
		if (event.target === event.currentTarget) {
			setIsVisible(false);
		}
	};

	const handleDrop = async (event: DragEvent) => {
		event.preventDefault();
		setIsVisible(false);

		const files = event.dataTransfer?.files;
		if (files && files.length > 0) {
			const file = files[0];

			if (file.name.toLowerCase().endsWith(".gb")) {
				try {
					const arrayBuffer = await file.arrayBuffer();
					await handleRomLoad(arrayBuffer);
				} catch (error) {
					console.error("Error reading dropped file", error);
					alert("Failed to load dropped file.");
				}
			} else {
				alert("Please drop a valid .gb file.");
			}
		}
	};

	return (
		<div
			id="drag-overlay"
			style={{ display: isVisible() ? "flex" : "none" }}
			onDragOver={handleDragOver}
			onDragLeave={handleDragLeave}
			onDrop={handleDrop}
		>
			Drop ROM file (.gb) to load
		</div>
	);
};
