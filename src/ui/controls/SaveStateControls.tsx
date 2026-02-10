import { store } from "../../core/store";
import { getSaveState, persistSaveState } from "../../services/storage";
import { updateDebugger } from "../Debugger";

/**
 * For testing out save states. Ideally this would be hotkey combos rather than
 * buttons in the UI.
 */
export const SaveStateControls = () => {
	const handleSaveState = async () => {
		try {
			const romHash = store.state.currentRomHash;

			if (!romHash) {
				return;
			}

			const serializedState = window.getSerializedState();

			await persistSaveState(romHash, 1, serializedState, {
				name: "Quick Save Slot 1",
			});

			updateDebugger();
		} catch (error) {
			console.error("Save state error:", error);
		}
	};

	const handleLoadState = async () => {
		try {
			const romHash = store.state.currentRomHash;

			if (!romHash) {
				return;
			}

			const stateData = await getSaveState(romHash, 1);

			if (!stateData) {
				console.warn("No save state found in slot 1 for ROM:", romHash);
				return;
			}

			window.loadSerializedState(stateData);
			updateDebugger();
		} catch (error) {
			console.error("Load state error:", error);
		}
	};

	return (
		<>
			<button onClick={handleSaveState} disabled={!store.state.currentRomHash}>
				Save State
			</button>
			<button onClick={handleLoadState} disabled={!store.state.currentRomHash}>
				Load State
			</button>
		</>
	);
};
