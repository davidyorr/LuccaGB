import { appState } from "../core/store";

export function setUpControlsHandlers({
	panelToggleId,
	controlsPanelId,
	debugCheckboxId,
}: {
	panelToggleId: string;
	controlsPanelId: string;
	debugCheckboxId: string;
}) {
	const panelToggle = document.getElementById(panelToggleId);
	const controlsPanel = document.getElementById(controlsPanelId);
	const debugCheckbox = document.getElementById(
		debugCheckboxId,
	) as HTMLInputElement | null;

	if (
		panelToggle === null ||
		controlsPanel === null ||
		debugCheckbox === null
	) {
		return;
	}

	// =========================================================
	// SIDE PANEL
	// =========================================================

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

	// =========================================================
	// DEBUG CHECKBOX
	// =========================================================

	// Debug checkbox
	debugCheckbox.addEventListener("change", () => {
		appState.setDebuggerOpen(debugCheckbox.checked);
	});

	// Handle app state changes
	appState.subscribe((state) => {
		if (debugCheckbox.checked !== state.isDebuggerOpen) {
			debugCheckbox.checked = state.isDebuggerOpen;
		}
	});
}
