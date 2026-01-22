export function setUpControlsHandlers({
	panelToggleId,
	controlsPanelId,
}: {
	panelToggleId: string;
	controlsPanelId: string;
}) {
	const panelToggle = document.getElementById(panelToggleId);
	const controlsPanel = document.getElementById(controlsPanelId);

	if (panelToggle === null || controlsPanel === null) {
		return;
	}

	// Toggle the panel on click
	panelToggle.addEventListener("click", () => {
		controlsPanel.classList.toggle("open");
	});

	// Close the panel when clicking outside
	document.addEventListener("click", (event) => {
		const target = event.target as HTMLInputElement;
		if (!controlsPanel.contains(target) && target !== panelToggle) {
			controlsPanel.classList.remove("open");
		}
	});
}
