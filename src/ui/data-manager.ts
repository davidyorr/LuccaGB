import { exportData, importData, type ImportStats } from "../services/storage";

export function setUpDataManagerHandlers({
	dataModalId,
	dataButtonId,
	exportButtonId,
	importId,
}: {
	dataModalId: string;
	dataButtonId: string;
	exportButtonId: string;
	importId: string;
}) {
	const dataModal = document.getElementById(
		dataModalId,
	) as HTMLDialogElement | null;
	const dataBtn = document.getElementById(
		dataButtonId,
	) as HTMLButtonElement | null;
	const exportBtn = document.getElementById(
		exportButtonId,
	) as HTMLButtonElement | null;
	const importInput = document.getElementById(
		importId,
	) as HTMLInputElement | null;

	if (
		dataModal === null ||
		dataBtn === null ||
		exportBtn === null ||
		importInput === null
	) {
		return;
	}

	// force clear the input
	importInput.value = "";

	dataBtn.addEventListener("click", () => {
		dataModal.showModal();
	});

	exportBtn.addEventListener("click", async () => {
		try {
			await exportData();
		} catch (error) {
			alert("Failed to export data: " + error);
		}
	});

	importInput.addEventListener("change", async (event) => {
		const target = event.target as HTMLInputElement;
		const file = target.files?.[0];
		if (!file) {
			return;
		}

		try {
			const text = await file.text();
			const stats: ImportStats = await importData(text);

			alert(
				`Import Complete!\n` +
					`• Added: ${stats.added}\n` +
					`• Updated: ${stats.updated}\n` +
					`• Skipped (Older): ${stats.skipped}\n` +
					`• Errors: ${stats.errors}`,
			);

			dataModal.close();
		} catch (err) {
			alert("Error importing file: " + err);
		} finally {
			// reset the input
			target.value = "";
		}
	});
}
