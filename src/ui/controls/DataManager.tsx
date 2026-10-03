import styles from "./DataManager.module.css";

import { createSignal, onMount, Show, type Component } from "solid-js";
import { Portal } from "solid-js/web";
import {
	exportData,
	importData,
	type ImportStats,
} from "../../services/storage";
import { createMagicLink, fetchMagicLinkData } from "../../services/magic-link";

function showImportResult(stats: ImportStats) {
	alert(
		`Import Complete!\n` +
			`• Added: ${stats.added}\n` +
			`• Updated: ${stats.updated}\n` +
			`• Skipped (Older): ${stats.skipped}\n` +
			`• Errors: ${stats.errors}`,
	);
}

export const DataManager: Component = () => {
	let dataModalElement!: HTMLDialogElement | undefined;

	// Magic link state
	const [syncId, setSyncId] = createSignal<string | null>(null);
	const [isSyncing, setIsSyncing] = createSignal(false);
	const [magicLink, setMagicLink] = createSignal("");
	const [isGenerating, setIsGenerating] = createSignal(false);

	onMount(function checkForMagicLink() {
		const params = new URLSearchParams(window.location.search);
		const id = params.get("import");

		if (id) {
			setSyncId(id);
			dataModalElement?.showModal();
		}
	});

	const handleManageDataClick = () => {
		dataModalElement?.showModal();
	};

	const handleCloseClick = () => {
		dataModalElement?.close();
	};

	// Clears the sync URL parameter so the import prompt doesn't come back.
	const handleDialogClose = () => {
		if (syncId()) {
			window.history.replaceState({}, document.title, window.location.pathname);
			setSyncId(null);
		}
		setMagicLink("");
	};

	const handleExportClick = async () => {
		try {
			await exportData();
		} catch (error) {
			alert("Failed to export data: " + error);
		}
	};

	const handleGenerateLinkClick = async () => {
		setIsGenerating(true);
		try {
			const link = await createMagicLink();
			setMagicLink(link);
		} catch (error) {
			alert("Failed to generate link: " + error);
		} finally {
			setIsGenerating(false);
		}
	};

	const handleCopyLinkClick = () => {
		navigator.clipboard.writeText(magicLink());
	};

	const handleConfirmSyncClick = async () => {
		const id = syncId();
		if (!id) {
			return;
		}

		setIsSyncing(true);

		try {
			const text = await fetchMagicLinkData(id);
			const stats = await importData(text);

			showImportResult(stats);
			dataModalElement?.close();
		} catch (err) {
			alert(`Import failed: ${err instanceof Error ? err.message : err}`);
		} finally {
			setIsSyncing(false);
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
			const stats = await importData(text);

			showImportResult(stats);
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
			<button onClick={handleManageDataClick}>Manage Data</button>

			<Portal>
				<dialog
					class={styles.modal}
					ref={dataModalElement}
					onClose={handleDialogClose}
				>
					<form method="dialog">
						<Show
							when={syncId()}
							fallback={
								<>
									<h2>Data Management</h2>

									<div class={styles.dataSection}>
										<h3>Transfer to Another Device</h3>
										<p>
											Create a temporary link that imports your saves and
											settings on another device (valid for 10 minutes).
										</p>
										<button
											type="button"
											onClick={handleGenerateLinkClick}
											disabled={isGenerating()}
										>
											{isGenerating() ? "Creating..." : "Create Link"}
										</button>

										<Show when={magicLink()}>
											<div class={styles.linkResult}>
												<input type="text" readOnly value={magicLink()} />
												<button type="button" onClick={handleCopyLinkClick}>
													Copy
												</button>
											</div>
										</Show>
									</div>

									<hr class={styles.divider} />

									<div class={styles.dataSection}>
										<h3>Export File</h3>
										<p>Download all saves and settings to a JSON file.</p>
										<button type="button" onClick={handleExportClick}>
											Download Backup
										</button>
									</div>

									<hr class={styles.divider} />

									<div class={styles.dataSection}>
										<h3>Import File</h3>
										<p>
											Restore from a JSON file. <br />
											<small>
												<em>
													Note: Newer saves in the file will overwrite local
													data.
												</em>
											</small>
										</p>
										<input
											type="file"
											accept=".json"
											onChange={handleInputChange}
										/>
									</div>

									<div class={styles.dialogFooter}>
										<button value="close" onClick={handleCloseClick}>
											Close
										</button>
									</div>
								</>
							}
						>
							{/* Opened from a magic link: show only the import prompt */}
							<h2>Import Saves</h2>
							<p>
								Import the saves and settings from your other device? Saves on
								this device are only replaced if the incoming ones are newer.
								Your settings will also be replaced.
							</p>

							<div class={styles.dialogFooter}>
								<button
									type="button"
									onClick={handleConfirmSyncClick}
									disabled={isSyncing()}
								>
									{isSyncing() ? "Importing..." : "Import"}
								</button>
								<button type="button" onClick={handleCloseClick}>
									Cancel
								</button>
							</div>
						</Show>
					</form>
				</dialog>
			</Portal>
		</>
	);
};
