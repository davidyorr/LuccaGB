export class DragAndDropHandler {
	private overlayId: string;
	private onRomLoaded: (buffer: ArrayBuffer) => Promise<void>;

	constructor(
		overlayId: string,
		onRomLoaded: (buffer: ArrayBuffer) => Promise<void>,
	) {
		this.overlayId = overlayId;
		this.onRomLoaded = onRomLoaded;

		const dragOverlay = document.getElementById(this.overlayId);
		if (!dragOverlay) {
			return;
		}

		window.addEventListener("dragenter", (event) => {
			// Only show if dragging files
			if (event.dataTransfer?.types.includes("Files")) {
				dragOverlay.style.display = "flex";
			}
		});

		dragOverlay.addEventListener("dragover", (event) => {
			event.preventDefault();
			if (event.dataTransfer) {
				event.dataTransfer.dropEffect = "copy";
			}
		});

		// Hide overlay if user leaves the window (dragging out)
		dragOverlay.addEventListener("dragleave", (event) => {
			if (event.target === dragOverlay) {
				dragOverlay.style.display = "none";
			}
		});

		dragOverlay.addEventListener("drop", async (event) => {
			event.preventDefault();
			dragOverlay.style.display = "none";

			const files = event.dataTransfer?.files;
			if (files && files.length > 0) {
				const file = files[0];

				if (file.name.toLowerCase().endsWith(".gb")) {
					try {
						const arrayBuffer = await file.arrayBuffer();
						await this.onRomLoaded(arrayBuffer);

						// Clear UI inputs to match state
						const romSelect = document.getElementById(
							"rom-select",
						) as HTMLSelectElement;
						const fileInput = document.getElementById(
							"rom-input",
						) as HTMLInputElement;
						if (romSelect) {
							romSelect.value = "";
						}
						if (fileInput) {
							fileInput.value = "";
						}
					} catch (error) {
						console.error("Error reading dropped file", error);
						alert("Failed to load dropped file.");
					}
				} else {
					alert("Please drop a valid .gb file.");
				}
			}
		});
	}
}
