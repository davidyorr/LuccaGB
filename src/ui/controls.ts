import { appState } from "../core/state";

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
		appState.setControlsOpen(!appState.isControlsOpen);
	});

	// Close the panel when clicking outside
	document.addEventListener("click", (event) => {
		const target = event.target as HTMLInputElement;
		if (!controlsPanel.contains(target) && target !== panelToggle) {
			appState.setControlsOpen(false);
		}
	});

	// Handle app state changes
	appState.subscribe((state) => {
		if (state.isControlsOpen) {
			controlsPanel.classList.add("open");
		} else {
			controlsPanel.classList.remove("open");
		}
	});
}
