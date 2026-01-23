import { appState } from "../core/state";

export function setUpAudioChannelHandlers({
	buttonId,
	dropdownId,
}: {
	buttonId: string;
	dropdownId: string;
}) {
	const button = document.getElementById(buttonId);
	const dropdown = document.getElementById(dropdownId);

	if (button === null || dropdown === null) {
		return;
	}

	// Toggle dropdown
	button.addEventListener("click", (event) => {
		event.stopPropagation();
		dropdown.classList.toggle("show");
	});

	// Close dropdown when clicking outside
	document.addEventListener("click", (event) => {
		if (!dropdown.contains(event.target as any) && event.target !== button) {
			dropdown.classList.remove("show");
		}
	});

	// Handle channel toggles
	for (let i = 1; i <= 4; i++) {
		const checkbox = document.getElementById(
			`audio-channel-${i}`,
		) as HTMLInputElement | null;
		// set the initial state
		appState.setAudioChannelEnabled(i, checkbox?.checked ?? true);
		// handle on change
		checkbox?.addEventListener("change", (event) => {
			appState.setAudioChannelEnabled(
				i,
				(event.target as HTMLInputElement)?.checked,
			);
		});
	}
}
