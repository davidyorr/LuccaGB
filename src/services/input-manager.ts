type Shortcut = {
	keydown?: () => void;
	keyup?: () => void;
};

export class InputManager {
	private keyToJoypadButton: Record<string, string> = {
		Enter: "START",
		Backspace: "SELECT",
		KeyZ: "B",
		KeyX: "A",
		ArrowDown: "DOWN",
		ArrowUp: "UP",
		ArrowLeft: "LEFT",
		ArrowRight: "RIGHT",
	};
	private allInputs: Array<string> = Object.values(this.keyToJoypadButton);

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

	private shortcuts: Record<string, Shortcut> = {};

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

	public registerShortcuts(shortcuts: Record<string, Shortcut>) {
		this.shortcuts = {
			...this.shortcuts,
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

		const shortcut = this.shortcuts[event.code];
		if (shortcut) {
			shortcut.keydown?.();
			return;
		}

		const action = this.keyToJoypadButton[event.code];
		if (action) {
			this.keyboardState.set(action, true);
		}
	};

	private handleKeyUp = (event: KeyboardEvent) => {
		const shortcut = this.shortcuts[event.code];
		if (shortcut && shortcut.keyup) {
			shortcut.keyup();
			return;
		}

		const action = this.keyToJoypadButton[event.code];
		if (action) {
			this.keyboardState.set(action, false);
		}
	};
}

export const inputManager = new InputManager();
