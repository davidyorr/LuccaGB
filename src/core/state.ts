type StateListener = (state: AppState) => void;

class AppState {
	// =========================================
	// GAME STATE
	// =========================================

	/** Read-only */
	public isPaused = false;
	/** Read-only */
	public isRomLoaded = false;
	/** Read-only */
	public currentRomHash = "";
	/** Read-only */
	public audioVolume = 0.5;
	/** Read-only */
	public audioChannelsEnabled = [false, true, true, true, true];

	// =========================================
	// UI STATE
	// =========================================

	/** Read-only */
	public isFileInputOpen = false;
	/** Read-only */
	public isHidden = false;
	/** Read-only */
	public isControlsOpen = false;

	// Event system
	private listeners: Set<StateListener> = new Set();

	// =========================================
	// GAME STATE
	// =========================================

	public setPaused(paused: boolean) {
		this.isPaused = paused;
		this.notify();
	}

	public togglePaused = () => {
		if (!appState.isRomLoaded) {
			return;
		}
		this.setPaused(!this.isPaused);
	};

	public setRomLoaded(loaded: boolean) {
		this.isRomLoaded = loaded;

		// When game starts, clear the UI
		if (loaded) {
			this.setPaused(false);
			this.setFileInputOpen(false);
			this.setControlsOpen(false);
		}
		this.notify();
	}

	public setCurrentRomHash(hash: string) {
		this.currentRomHash = hash;
		this.notify();
	}

	public setAudioVolume(volume: number) {
		this.audioVolume = volume;
		this.notify();
	}

	public setAudioChannelEnabled(channel: number, enabled: boolean) {
		if (channel >= 1 && channel <= 4) {
			this.audioChannelsEnabled[channel] = enabled;
		}
		this.notify();
	}

	// =========================================
	// UI STATE
	// =========================================

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

	public setControlsOpen(isOpen: boolean) {
		this.isControlsOpen = isOpen;
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
export const appState = new AppState();
