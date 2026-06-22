import styles from "./PalettePicker.module.css";
import {
	createSignal,
	For,
	Show,
	onCleanup,
	onMount,
	type Component,
} from "solid-js";
import { store } from "../../core/store";

export const PalettePicker: Component = () => {
	let dropdownElement!: HTMLDivElement;
	let buttonElement!: HTMLButtonElement;

	const [dropdownOpen, setDropdownOpen] = createSignal(false);
	const activePalette = store.getters.activePalette;

	onMount(() => {
		const handleClickOutside = (event: MouseEvent) => {
			const target = event.target as Node;
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

	const handleCloneClick = () => {
		const current = activePalette();
		store.actions.addPalette({
			id: crypto.randomUUID(),
			name: `${current.name} (Custom)`,
			colors: [...current.colors],
			isReadonly: false,
		});
	};

	const colorNames = [
		"Color 0 (Lightest)",
		"Color 1 (Light)",
		"Color 2 (Dark)",
		"Color 3 (Darkest)",
	];

	return (
		<div class={styles.dropdownContainer}>
			<button
				class={styles.paletteButton}
				type="button"
				ref={buttonElement}
				onClick={() => setDropdownOpen((prev) => !prev)}
			>
				Color Palette [{activePalette().name}] ▼
			</button>
			<div
				classList={{
					[styles.dropdownPanel]: true,
					[styles.open]: dropdownOpen(),
				}}
				ref={dropdownElement}
			>
				{/* Palette Selector Dropdown */}
				<select
					class={styles.paletteSelect}
					onChange={(e) => store.actions.setActivePaletteId(e.target.value)}
				>
					<For each={store.state.settings.palettes}>
						{(palette) => (
							<option
								value={palette.id}
								selected={palette.id === store.state.settings.activePaletteId}
							>
								{palette.name}
							</option>
						)}
					</For>
				</select>

				{/* Rename Input (Only shows for custom palettes) */}
				<Show when={!activePalette().isReadonly}>
					<input
						type="text"
						class={styles.nameInput}
						value={activePalette().name}
						onInput={(event) =>
							store.actions.renameActivePalette(event.currentTarget.value)
						}
						placeholder="Palette Name"
						maxlength={30}
					/>
				</Show>

				<hr class={styles.divider} />

				{/* Color Inputs */}
				<For each={[0, 1, 2, 3]}>
					{(i) => (
						<label class={styles.colorLabel}>
							<input
								type="color"
								class={styles.colorInput}
								value={activePalette().colors[i]}
								disabled={activePalette().isReadonly}
								onInput={(event) =>
									store.actions.updateActivePaletteColor(i, event.target.value)
								}
							/>
							{colorNames[i]}
						</label>
					)}
				</For>

				{/* Action Buttons */}
				<Show when={activePalette().isReadonly}>
					<button
						type="button"
						class={styles.actionButton}
						onClick={handleCloneClick}
					>
						Clone to Edit
					</button>
				</Show>

				<Show when={!activePalette().isReadonly}>
					<button
						type="button"
						class={`${styles.actionButton} ${styles.deleteButton}`}
						onClick={() => store.actions.deletePalette(activePalette().id)}
					>
						Delete Palette
					</button>
				</Show>
			</div>
		</div>
	);
};
