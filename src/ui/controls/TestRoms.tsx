import {
	createEffect,
	createSignal,
	For,
	onMount,
	type Component,
} from "solid-js";
import { handleRomLoad } from "../../services/rom-loader";

type RomFile = {
	path: string;
	displayName: string;
	getUrl: () => Promise<string>;
};

export const TestRoms: Component = () => {
	const [roms, setRomFiles] = createSignal<Record<string, RomFile>>({});
	const [selectedRom, setSelectedRom] = createSignal<RomFile>();
	let romSelectElement!: HTMLSelectElement | undefined;

	onMount(async () => {
		const modules = import.meta.glob<string>("../../../roms/**/*.gb", {
			query: "?url",
			import: "default",
		});

		const sortedPaths = Object.keys(modules).sort();
		const options: Record<string, RomFile> = {};

		for (const path of sortedPaths) {
			options[path] = {
				path: path,
				displayName: path.replace("../../../roms/", ""),
				getUrl: modules[path],
			};
		}

		setRomFiles(options);
	});

	const handleRomChange = async (event: Event) => {
		const romPath = (event.currentTarget as HTMLSelectElement).value;
		if (romPath && roms()[romPath]) {
			setSelectedRom(roms()[romPath]);
		}
	};

	createEffect(async function loadRomOnSelect() {
		const rom = selectedRom();
		if (!rom) {
			return;
		}

		// remove focus so keyboard controls don't toggle the dropdown
		romSelectElement?.blur();

		// reset file input so it doesn't look like two things are selected

		const loadRom = async () => {
			try {
				const url = await rom.getUrl();
				const response = await fetch(url);

				if (!response.ok) {
					throw new Error(`Failed to fetch ${url}`);
				}

				const arrayBuffer = await response.arrayBuffer();
				if (arrayBuffer) {
					await handleRomLoad(arrayBuffer);
				}
			} catch (err) {
				console.error("Error loading ROM:", err);
				alert("Failed to load selected ROM.");
			}
		};

		loadRom();
	});

	return (
		<>
			<select ref={romSelectElement} onChange={handleRomChange}>
				<option value="">-- Load a Test ROM --</option>
				<For each={Object.values(roms())}>
					{(rom) => <option value={rom.path}>{rom.displayName}</option>}
				</For>
			</select>
		</>
	);
};
