import type { Component } from "solid-js";
import { Portal } from "solid-js/web";
import { exportData, importData, type ImportStats } from "../services/storage";

export const DataManager: Component = () => {
	let dataModalElement!: HTMLDialogElement | undefined;

	const handleManageDataClick = () => {
		dataModalElement?.showModal();
	};

	const handleCloseClick = () => {
		dataModalElement?.close();
	};

	const handleExportClick = async () => {
		try {
			await exportData();
		} catch (error) {
			alert("Failed to export data: " + error);
		}
	};

	const handleInputChange = async (event: Event) => {
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

			dataModalElement?.close();
		} catch (err) {
			alert("Error importing file: " + err);
		} finally {
			// reset the input
			target.value = "";
		}
	};

	return (
		<>
			<button id="data-manager-button" onClick={handleManageDataClick}>
				Manage Data
			</button>

			<Portal>
				<dialog id="data-modal" ref={dataModalElement}>
					<form method="dialog">
						<h2>Data Management</h2>

						<div class="data-section">
							<h3>Export</h3>
							<p>Download all saves and settings to a JSON file.</p>
							<button
								type="button"
								id="export-data-button"
								onClick={handleExportClick}
							>
								Download Backup
							</button>
						</div>

						<hr />

						<div class="data-section">
							<h3>Import</h3>
							<p>
								Restore from a JSON file. <br />
								<small>
									<em>
										Note: Newer saves in the file will overwrite local data.
									</em>
								</small>
							</p>
							<input
								type="file"
								id="import-data-input"
								accept=".json"
								onChange={handleInputChange}
							/>
						</div>

						<div class="dialog-footer">
							<button value="close" onClick={handleCloseClick}>
								Close
							</button>
						</div>
					</form>
				</dialog>
			</Portal>
		</>
	);
};
