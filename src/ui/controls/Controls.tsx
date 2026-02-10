import styles from "./Controls.module.css";

import { type Component, onCleanup, onMount } from "solid-js";
import { TestRoms } from "./TestRoms";
import { RomFileInput } from "./RomFileInput";
import { DataManager } from "./DataManager";
import { ScreenshotButton } from "./ScreenshotButton";
import { ViewportScale } from "./ViewportScale";
import { VolumeControl } from "./VolumeControl";
import { AudioChannels } from "./AudioChannels";
import { TraceLogger } from "./TraceLogger";
import { DebuggerToggle } from "../Debugger";
import { store } from "../../core/store";
import { SaveStateControls } from "./SaveStateControls";

export const Controls: Component = () => {
	let panelRef: HTMLDivElement | undefined;
	let toggleRef: HTMLButtonElement | undefined;

	onMount(() => {
		const handleOutsideClick = (event: MouseEvent) => {
			const target = event.target as Node;

			// If the panel is open, and the click is NOT inside the panel
			// and NOT on the toggle button, close the panel
			if (
				store.state.ui.isControlsOpen &&
				panelRef &&
				!panelRef.contains(target) &&
				toggleRef &&
				!toggleRef.contains(target)
			) {
				store.actions.setControlsOpen(false);
			}
		};

		document.addEventListener("click", handleOutsideClick);

		onCleanup(() => {});
	});

	return (
		<>
			<button
				ref={toggleRef}
				class={styles.panelToggle}
				onClick={() =>
					store.actions.setControlsOpen(!store.state.ui.isControlsOpen)
				}
			>
				☰
			</button>

			<div
				ref={panelRef}
				class={styles.controlsPanel}
				classList={{ [styles.open]: store.state.ui.isControlsOpen }}
			>
				<div class={styles.controls}>
					<TestRoms />
					<RomFileInput />
					<DataManager />
					<ScreenshotButton />
					<ViewportScale />
					<VolumeControl />
					<AudioChannels />
					<TraceLogger />
					<DebuggerToggle />
					<SaveStateControls />
				</div>
			</div>
		</>
	);
};
