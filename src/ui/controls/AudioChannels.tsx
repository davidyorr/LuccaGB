import styles from "./AudioChannels.module.css";

import {
	createSignal,
	For,
	onCleanup,
	onMount,
	type Component,
} from "solid-js";
import { store } from "../../core/store";

export const AudioChannels: Component = () => {
	let dropdownElement!: HTMLDivElement;
	let buttonElement!: HTMLButtonElement;

	const [dropdownOpen, setDropdownOpen] = createSignal(false);

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
				setDropdownOpen(false);
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
		<div class={styles.dropdownContainer}>
			<button
				class={styles.audioChannelsButton}
				type="button"
				ref={buttonElement}
				onClick={() => setDropdownOpen(true)}
			>
				Audio Channels ▼
			</button>
			<div
				classList={{
					[styles.dropdownPanel]: true,
					[styles.open]: dropdownOpen(),
				}}
				ref={dropdownElement}
			>
				<For each={[1, 2, 3, 4]}>
					{(i) => (
						<label class={styles.channelLabel}>
							<input
								type="checkbox"
								class={styles.channelLabelInput}
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
