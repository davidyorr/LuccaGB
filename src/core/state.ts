type StateListener = (state: EmulatorState) => void;

class EmulatorState {
	/** Read-only */
	public isPaused = false;
	/** Read-only */
	public isHidden = false;
	/** Read-only */
	public isFileInputOpen = false;
	/** Read-only */
	public isRomLoaded = false;
	/** Read-only */
	public currentRomHash = "";
	/** Read-only */
	public audioChannelsEnabled = [false, true, true, true, true];

	// Simple event system if you need UI updates
	private listeners: Set<StateListener> = new Set();

	public setPaused(paused: boolean) {
		this.isPaused = paused;
		this.notify();
	}

	public togglePaused = () => {
		if (!emulatorState.isRomLoaded) {
			return;
		}
		this.setPaused(!this.isPaused);
	};

	public setHidden(hidden: boolean) {
		this.isHidden = hidden;
		this.notify();
	}

	public setFileInputOpen(isOpen: boolean) {
		this.isFileInputOpen = isOpen;

		// Force pause when input is open
		if (isOpen) {
			this.setPaused(true);
		}
		this.notify();
	}

	public setRomLoaded(loaded: boolean) {
		this.isRomLoaded = loaded;

		// Force resume after ROM is loaded
		this.setPaused(false);
		this.notify();
	}

	public setCurrentRomHash(hash: string) {
		this.currentRomHash = hash;
		this.notify();
	}

	public setAudioChannelEnabled(channel: number, enabled: boolean) {
		if (channel >= 1 && channel <= 4) {
			this.audioChannelsEnabled[channel] = enabled;
		}
		this.notify();
	}

	// Allow UI components to subscribe to changes
	public subscribe(listener: StateListener) {
		this.listeners.add(listener);
		return () => this.listeners.delete(listener); // unsubscribe function
	}

	private notify() {
		this.listeners.forEach((l) => l(this));
	}
}

// Export a single instance
export const emulatorState = new EmulatorState();
