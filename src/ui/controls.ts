import { appState } from "../core/state";
import type { CanvasRenderer } from "../services/canvas-renderer";

export function setUpControlsHandlers({
	panelToggleId,
	controlsPanelId,
	debugCheckboxId,
	scaleSelectId,
	canvasRenderer,
}: {
	panelToggleId: string;
	controlsPanelId: string;
	debugCheckboxId: string;
	scaleSelectId: string;
	canvasRenderer: CanvasRenderer;
}) {
	const panelToggle = document.getElementById(panelToggleId);
	const controlsPanel = document.getElementById(controlsPanelId);
	const debugCheckbox = document.getElementById(
		debugCheckboxId,
	) as HTMLInputElement | null;
	const scaleSelect = document.getElementById(
		scaleSelectId,
	) as HTMLSelectElement | null;

	if (
		panelToggle === null ||
		controlsPanel === null ||
		debugCheckbox === null ||
		scaleSelect === null
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
	// DISPLAY SCALING
	// =========================================================

	// initial sync
	canvasRenderer.setScale(appState.scale);
	scaleSelect.value = appState.scale.toString();

	// User selects option
	scaleSelect.addEventListener("change", () => {
		const value = scaleSelect.value;
		const newScale = value === "fit" ? "fit" : parseInt(value, 10);
		appState.setScale(newScale);
	});

	let previousScale = appState.scale;

	appState.subscribe((state) => {
		if (state.scale !== previousScale) {
			canvasRenderer.setScale(state.scale);

			if (scaleSelect.value !== state.scale.toString()) {
				scaleSelect.value = state.scale.toString();
			}

			previousScale = state.scale;
		}
	});

	window.addEventListener("resize", () => {
		canvasRenderer.setScale(appState.scale);
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
