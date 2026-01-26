import { For, onCleanup, onMount, type Component } from "solid-js";
import { store } from "../../core/store";

export const AudioChannels: Component = () => {
	let dropdownElement!: HTMLDivElement;
	let buttonElement!: HTMLButtonElement;

	const handleChange = (channel: number, checked: boolean) => {
		store.actions.setAudioChannelsEnabled(channel, checked);
		window.setAudioChannelEnabled(channel, checked);
	};

	onMount(() => {
		const handleClickOutside = (event: MouseEvent) => {
			const target = event.target as Node;

			// If click is NOT inside dropdown AND NOT inside button, close the dropdown
			if (
				dropdownElement &&
				!dropdownElement.contains(target) &&
				buttonElement &&
				!buttonElement.contains(target)
			) {
				dropdownElement.classList.remove("show");
			}
		};

		document.addEventListener("click", handleClickOutside);

		onCleanup(() => {
			document.removeEventListener("click", handleClickOutside);
		});
	});

	const channelNames = [
		"",
		"Ch1: Square",
		"Ch2: Square",
		"Ch3: Wave",
		"Ch4: Noise",
	];

	return (
		<div class="dropdown-container">
			<button
				id="audio-channels-button"
				type="button"
				ref={buttonElement}
				onClick={() => {
					dropdownElement.classList.toggle("show");
				}}
			>
				Audio Channels ▼
			</button>
			<div
				id="audio-channels-dropdown"
				class="dropdown-panel"
				ref={dropdownElement}
			>
				<For each={[1, 2, 3, 4]}>
					{(i) => (
						<label class="channel-label">
							<input
								type="checkbox"
								id={`audio-channel-${i}`}
								checked={store.state.settings.audioChannelsEnabled[i]}
								onChange={(event) => handleChange(i, event.target.checked)}
							/>
							{channelNames[i]}
						</label>
					)}
				</For>
			</div>
		</div>
	);
};
