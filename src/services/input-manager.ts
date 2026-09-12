type Shortcut = {
	keydown?: () => void;
	keyup?: () => void;
};

// Helper function to generate a consistent key combo string (e.g., "Ctrl+Shift+KeyX")
export function getEventCombo(e: KeyboardEvent): string {
	const modifiers = [];

	// Only add modifier to list if the pressed key itself is NOT this modifier.
	// (e.g., if you only press "Shift", we want "ShiftLeft", not "Shift+ShiftLeft")
	if (e.ctrlKey && !e.code.startsWith("Control")) modifiers.push("Ctrl");
	if (e.metaKey && !e.code.startsWith("Meta")) modifiers.push("Meta");
	if (e.altKey && !e.code.startsWith("Alt")) modifiers.push("Alt");
	if (e.shiftKey && !e.code.startsWith("Shift")) modifiers.push("Shift");

	if (modifiers.length > 0) {
		return `${modifiers.join("+")}+${e.code}`;
	}
	return e.code;
}

export class InputManager {
	// Mapping of Combo String (e.g. "Ctrl+KeyX") -> Joypad Action ("A")
	private keyToJoypadButton: Record<string, string> = {};
	// Joypad actions used for input state synchronization
	private allInputs: Array<string> = [
		"START",
		"SELECT",
		"B",
		"A",
		"DOWN",
		"UP",
		"LEFT",
		"RIGHT",
	];

	// https://w3c.github.io/gamepad/#remapping
	private buttonIndexToJoypadButton: Record<number, string> = {
		0: "B",
		1: "A",
		8: "SELECT",
		9: "START",
		12: "UP",
		13: "DOWN",
		14: "LEFT",
		15: "RIGHT",
	};

	private joypadButtonToButtonIndex: Record<string, number>;

	// Action name -> shortcut handlers
	private registeredShortcuts: Record<string, Shortcut> = {};
	// Key combo -> action name (e.g. "Ctrl+Space" -> "togglePaused")
	private keyToShortcutAction: Record<string, string> = {};

	// Map<PhysicalKeyCode, JoypadAction/ShortcutAction> to handle out-of-order modifier keyups:
	// Track the action activated by each physical key so keyup can release it
	// even if the modifier state has changed since keydown.
	private activeCodeToJoypadAction: Map<string, string> = new Map();
	private activeCodeToShortcutAction: Map<string, string> = new Map();

	// Map<Button, isPressed>
	private keyboardState: Map<string, boolean> = new Map();

	// the state sent to the emulator
	// Map<Button, isPressed>
	private emulatedButtonState: Map<string, boolean> = new Map();

	constructor() {
		this.joypadButtonToButtonIndex = Object.fromEntries(
			Object.entries(this.buttonIndexToJoypadButton).map(([k, v]) => [
				v,
				parseInt(k),
			]),
		);

		window.addEventListener("keydown", this.handleKeyDown);
		window.addEventListener("keyup", this.handleKeyUp);
	}

	public setJoypadBindings(bindings: Record<string, string>) {
		this.keyToJoypadButton = {};
		for (const [action, combo] of Object.entries(bindings)) {
			this.keyToJoypadButton[combo] = action;
		}
	}

	public setShortcutBindings(bindings: Record<string, string>) {
		this.keyToShortcutAction = {};
		for (const [action, combo] of Object.entries(bindings)) {
			this.keyToShortcutAction[combo] = action;
		}
	}

	public registerShortcuts(shortcuts: Record<string, Shortcut>) {
		this.registeredShortcuts = {
			...this.registeredShortcuts,
			...shortcuts,
		};
	}

	public poll(forceSync: boolean = false) {
		// Even if no gamepad is connected, we must run this loop
		// to process the keyboard state captured in handleKeyDown/Up

		const gamepads = navigator.getGamepads ? navigator.getGamepads() : [];
		const gamepad = Array.from(gamepads).find((gamepad) => gamepad !== null);

		this.allInputs.forEach((action) => {
			let gamepadPressed = false;

			// Check buttons
			if (gamepad) {
				const btnIndex = this.joypadButtonToButtonIndex[action];
				if (btnIndex !== undefined && gamepad.buttons[btnIndex]?.pressed) {
					gamepadPressed = true;
				}
			}

			// Check keyboard
			const keyboardPressed = this.keyboardState.get(action) || false;

			// Combine them
			const isPressedNow = gamepadPressed || keyboardPressed;
			const wasPressedBefore = this.emulatedButtonState.get(action) || false;

			// If forcing sync, inform the emulator of the current state no matter what
			if (forceSync) {
				if (isPressedNow) {
					window.handleJoypadButtonPressed(action);
				} else {
					window.handleJoypadButtonReleased(action);
				}
			} else {
				// Otherwise, only trigger if the result changed
				if (isPressedNow && !wasPressedBefore) {
					window.handleJoypadButtonPressed(action);
				} else if (!isPressedNow && wasPressedBefore) {
					window.handleJoypadButtonReleased(action);
				}
			}

			this.emulatedButtonState.set(action, isPressedNow);
		});
	}

	public syncJoypadState() {
		this.poll(true);
	}

	private handleKeyDown = (event: KeyboardEvent) => {
		if (event.repeat) {
			return;
		}

		const target = event.target as HTMLElement;
		if (
			target.tagName === "INPUT" ||
			target.tagName === "TEXTAREA" ||
			target.isContentEditable
		) {
			return;
		}

		const combo = getEventCombo(event);

		// Handle Shortcuts
		const shortcutAction = this.keyToShortcutAction[combo];
		const shortcut = shortcutAction
			? this.registeredShortcuts[shortcutAction]
			: null;

		if (shortcut) {
			event.preventDefault();
			this.activeCodeToShortcutAction.set(event.code, shortcutAction);
			shortcut.keydown?.();
			return;
		}

		// Handle Joypad
		const joypadAction = this.keyToJoypadButton[combo];
		if (joypadAction) {
			this.keyboardState.set(joypadAction, true);
			this.activeCodeToJoypadAction.set(event.code, joypadAction);
		}
	};

	private handleKeyUp = (event: KeyboardEvent) => {
		const target = event.target as HTMLElement;
		if (
			target.tagName === "INPUT" ||
			target.tagName === "TEXTAREA" ||
			target.isContentEditable
		) {
			return;
		}

		// Use the physical key code to find what was activated, so it successfully
		// releases even if modifiers (Ctrl, Shift, etc) were released out of order.
		const shortcutAction = this.activeCodeToShortcutAction.get(event.code);
		if (shortcutAction) {
			const shortcut = this.registeredShortcuts[shortcutAction];
			if (shortcut && shortcut.keyup) {
				event.preventDefault();
				shortcut.keyup();
			}
			this.activeCodeToShortcutAction.delete(event.code);
			return;
		}

		const joypadAction = this.activeCodeToJoypadAction.get(event.code);
		if (joypadAction) {
			this.keyboardState.set(joypadAction, false);
			this.activeCodeToJoypadAction.delete(event.code);
		}
	};
}

export const inputManager = new InputManager();
