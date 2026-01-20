export class TestRomLibrary {
	private romFiles = import.meta.glob("../../roms/**/*.gb", {
		query: "?url",
		import: "default",
	});

	public populateSelect(selectElementId: string) {
		const romSelect = document.getElementById(
			selectElementId,
		) as HTMLSelectElement | null;
		if (romSelect === null) {
			return;
		}

		// populate options
		const sortedPaths = Object.keys(this.romFiles).sort();

		for (const path of sortedPaths) {
			const option = document.createElement("option");
			option.text = path.replace("../../roms/", "");
			option.value = path;
			romSelect.appendChild(option);
		}
	}

	public async loadRomByPath(path: string): Promise<ArrayBuffer | undefined> {
		if (!path || !this.romFiles[path]) {
			return;
		}

		try {
			// fetch the URL from Vite
			const getUrl = this.romFiles[path] as () => Promise<string>;
			const url = await getUrl();

			const response = await fetch(url);
			if (!response.ok) {
				throw new Error(`Failed to fetch ${path}`);
			}

			const arrayBuffer = await response.arrayBuffer();
			return arrayBuffer;
		} catch (err) {
			console.error("Error loading ROM from dropdown:", err);
			alert("Failed to load selected ROM.");
		}
	}
}
