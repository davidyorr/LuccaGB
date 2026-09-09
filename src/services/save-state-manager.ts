import { store } from "../core/store";
import { updateDebugger } from "../ui/Debugger";
import { persistSaveState, getSaveState } from "./storage";

export async function quickSaveState() {
	const romHash = store.state.currentRomHash;
	if (!romHash) {
		return;
	}

	try {
		const serializedState = window.getSerializedState();
		const status = `Saved State ${store.state.currentSaveStateSlot}`;
		await persistSaveState(
			romHash,
			store.state.currentSaveStateSlot,
			serializedState,
			{
				name: status,
			},
		);
		store.actions.setSaveStateStatus(status, 2000);
		updateDebugger();
	} catch (error) {
		console.error("Save state error:", error);
		store.actions.setSaveStateStatus("Quick save failed", 2000);
	}
}

export async function quickLoadState() {
	const romHash = store.state.currentRomHash;
	if (!romHash) {
		return;
	}

	try {
		const stateData = await getSaveState(
			romHash,
			store.state.currentSaveStateSlot,
		);
		if (!stateData) {
			console.warn(
				`No save state in slot ${store.state.currentSaveStateSlot} for`,
				romHash,
			);
			store.actions.setSaveStateStatus(
				`No Save State in Slot ${store.state.currentSaveStateSlot}`,
				2000,
			);
			return;
		}
		window.loadSerializedState(stateData);
		store.actions.setSaveStateStatus(
			`Loaded State ${store.state.currentSaveStateSlot}`,
			2000,
		);
		updateDebugger();
	} catch (error) {
		console.error("Load state error:", error);
		store.actions.setSaveStateStatus("Quick load failed", 2000);
	}
}
