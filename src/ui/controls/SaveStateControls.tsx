import { store } from "../../core/store";
import {
	quickLoadState,
	quickSaveState,
} from "../../services/save-state-manager";

/**
 * For testing out save states. Ideally this would be hotkey combos rather than
 * buttons in the UI.
 */
export const SaveStateControls = () => {
	const handleSaveState = async () => {
		await quickSaveState();
	};

	const handleLoadState = async () => {
		await quickLoadState();
	};

	return (
		<>
			<div>Save state slot: {store.state.currentSaveStateSlot}</div>
			<button onClick={handleSaveState} disabled={!store.state.currentRomHash}>
				Save State
			</button>
			<button onClick={handleLoadState} disabled={!store.state.currentRomHash}>
				Load State
			</button>
		</>
	);
};
