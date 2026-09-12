import styles from "./KeybindingsManager.module.css";
import {
	type Component,
	createSignal,
	onCleanup,
	onMount,
	For,
} from "solid-js";
import { Portal } from "solid-js/web";
import { store } from "../../core/store";
import { getEventCombo } from "../../services/input-manager";

type ListenState = { type: "joypad" | "shortcut"; action: string } | null;

const formatKeyCode = (combo: string | undefined) => {
	if (!combo) return "Unbound";
	return combo
		.split("+")
		.map((part) => {
			if (part.startsWith("Key")) return part.replace(/^Key/, "");
			if (part.startsWith("Digit")) return part.replace(/^Digit/, "");
			if (part.startsWith("Arrow"))
				return part.replace(/^Arrow/, "") + " Arrow";
			if (part === "ControlLeft" || part === "ControlRight") return "Ctrl";
			if (part === "ShiftLeft" || part === "ShiftRight") return "Shift";
			if (part === "AltLeft" || part === "AltRight") return "Alt";
			if (part === "MetaLeft" || part === "MetaRight") return "Meta";
			return part;
		})
		.join(" + ");
};

const joypadLabels: Record<string, string> = {
	UP: "D-Pad Up",
	DOWN: "D-Pad Down",
	LEFT: "D-Pad Left",
	RIGHT: "D-Pad Right",
	A: "A Button",
	B: "B Button",
	START: "Start",
	SELECT: "Select",
};

const shortcutLabels: Record<string, string> = {
	togglePaused: "Play / Pause",
	rewind: "Rewind",
	toggleRecording: "Toggle Recording",
	quickSave: "Quick Save",
	quickLoad: "Quick Load",
	prevSaveStateSlot: "Prev Save Slot",
	nextSaveStateSlot: "Next Save Slot",
};

const modifierCodes = new Set([
	"ControlLeft",
	"ControlRight",
	"MetaLeft",
	"MetaRight",
	"AltLeft",
	"AltRight",
	"ShiftLeft",
	"ShiftRight",
]);

const isModifierCode = (code: string) => modifierCodes.has(code);

export const KeybindingsManager: Component = () => {
	let modalRef!: HTMLDialogElement | undefined;
	const [listening, setListening] = createSignal<ListenState>(null);

	const open = () => modalRef?.showModal();

	const close = () => {
		setListening(null);
		modalRef?.close();
	};

	const saveBind = (state: ListenState, combo: string) => {
		if (state?.type === "joypad") {
			store.actions.updateJoypadBinding(state.action, combo);
		} else if (state?.type === "shortcut") {
			store.actions.updateShortcutBinding(state.action, combo);
		}
		setListening(null);
	};

	const handleKeyDown = (e: KeyboardEvent) => {
		const state = listening();
		if (!state) {
			return;
		}

		e.preventDefault();
		e.stopPropagation();

		if (e.code === "Escape") {
			setListening(null);
			return;
		}

		if (isModifierCode(e.code)) {
			// They pressed a modifier. Don't bind yet, wait for them
			// to press a normal key OR release this modifier.
			return;
		}

		saveBind(state, getEventCombo(e));
	};

	const handleKeyUp = (e: KeyboardEvent) => {
		const state = listening();
		if (!state) {
			return;
		}

		e.preventDefault();
		e.stopPropagation();

		if (isModifierCode(e.code)) {
			// If they released a modifier and we are STILL listening,
			// it means they didn't press any other key (e.g. they just want to bind "Shift")
			saveBind(state, getEventCombo(e));
		}
	};

	onMount(() => {
		window.addEventListener("keydown", handleKeyDown, { capture: true });
		window.addEventListener("keyup", handleKeyUp, { capture: true });
	});

	onCleanup(() => {
		window.removeEventListener("keydown", handleKeyDown, { capture: true });
		window.removeEventListener("keyup", handleKeyUp, { capture: true });
	});

	return (
		<>
			<button onClick={open}>Keybindings</button>

			<Portal>
				<dialog class={styles.modal} ref={modalRef} onClose={close}>
					<h2>Keybindings</h2>

					<div class={styles.scrollableContent}>
						<div class={styles.section}>
							<h3>Joypad</h3>
							<div class={styles.grid}>
								<For each={Object.keys(joypadLabels)}>
									{(action) => (
										<div class={styles.row}>
											<span>{joypadLabels[action]}</span>
											<button
												type="button"
												class={styles.keyButton}
												classList={{
													[styles.listening]:
														listening()?.type === "joypad" &&
														listening()?.action === action,
												}}
												onClick={() => setListening({ type: "joypad", action })}
											>
												{listening()?.type === "joypad" &&
												listening()?.action === action
													? "Press any key..."
													: formatKeyCode(
															store.state.settings.keybindings.joypad[action],
														)}
											</button>
										</div>
									)}
								</For>
							</div>
						</div>

						<hr class={styles.divider} />

						<div class={styles.section}>
							<h3>Shortcuts</h3>
							<div class={styles.grid}>
								<For each={Object.keys(shortcutLabels)}>
									{(action) => (
										<div class={styles.row}>
											<span>{shortcutLabels[action]}</span>
											<button
												type="button"
												class={styles.keyButton}
												classList={{
													[styles.listening]:
														listening()?.type === "shortcut" &&
														listening()?.action === action,
												}}
												onClick={() =>
													setListening({ type: "shortcut", action })
												}
											>
												{listening()?.type === "shortcut" &&
												listening()?.action === action
													? "Press any key..."
													: formatKeyCode(
															store.state.settings.keybindings.shortcuts[
																action
															],
														)}
											</button>
										</div>
									)}
								</For>
							</div>
						</div>
					</div>

					<div class={styles.dialogFooter}>
						<button type="button" onClick={close}>
							Close
						</button>
					</div>
				</dialog>
			</Portal>
		</>
	);
};
