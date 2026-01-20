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

	private shortcuts: Record<string, () => void>;

	constructor(shortcuts: Record<string, () => void>) {
		this.shortcuts = shortcuts;
		window.addEventListener("keydown", this.handleKeyDown);
		window.addEventListener("keyup", this.handleKeyUp);
	}

	private handleKeyDown = (event: KeyboardEvent) => {
		if (event.repeat) {
			return;
		}

		if (this.shortcuts[event.code]) {
			this.shortcuts[event.code]();
			return;
		}

		if (!this.keyToJoypadButton[event.code]) {
			return;
		}

		window.handleJoypadButtonPressed(this.keyToJoypadButton[event.code]);
	};

	private handleKeyUp = (event: KeyboardEvent) => {
		if (!this.keyToJoypadButton[event.code]) {
			return;
		}

		window.handleJoypadButtonReleased(this.keyToJoypadButton[event.code]);
	};
}
