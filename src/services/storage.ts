import { appState } from "../core/state";

const DB_NAME = "LuccaGB-Database";
const DB_VERSION = 2;
const STORE_NAME = "cartridgeRam";

const SETTINGS_STORE = "appSettings";
const SETTINGS_KEY = "settings";

export type AppSettings = {
	audioVolume: number;
	audioChannelsEnabled: boolean[];
	isDebuggerOpen: boolean;
	scale: number | "fit";
	updatedAt: number;
};

type SaveData = {
	romHash: string;
	ram: Uint8Array;
	updatedAt: number;
	meta?: {
		name?: string;
		playTime?: number;
	};
};

type BackupFile = {
	version: number; // Schema version
	timestamp: number;
	app: "LuccaGB";
	saves: {
		romHash: string;
		ramBase64: string;
		updatedAt: number;
		meta?: SaveData["meta"];
	}[];
	settings?: Record<string, any>; // Placeholder for future settings
};

export type ImportStats = {
	added: number;
	updated: number;
	skipped: number;
	errors: number;
};

let dbPromise: Promise<IDBDatabase> | null = null;

export function openDatabase(): Promise<IDBDatabase> {
	if (dbPromise) {
		return dbPromise;
	}

	dbPromise = new Promise((resolve, reject) => {
		const request = indexedDB.open(DB_NAME, DB_VERSION);

		request.onupgradeneeded = (event) => {
			const db = (event.target as IDBOpenDBRequest).result;
			// Create stores if missing
			if (!db.objectStoreNames.contains(STORE_NAME)) {
				db.createObjectStore(STORE_NAME, {
					keyPath: "romHash",
				});
			}
			if (!db.objectStoreNames.contains(SETTINGS_STORE)) {
				db.createObjectStore(SETTINGS_STORE, {
					keyPath: "key",
				});
			}
		};

		request.onsuccess = () => {
			const db = request.result;
			db.onversionchange = () => {
				console.warn("IndexedDB version change detected, closing DB");
				db.close();
				dbPromise = null;
			};
			resolve(db);
		};

		request.onerror = () => {
			reject(request.error);
		};

		request.onblocked = () => {
			console.warn(
				"IndexedDB upgrade blocked. Close other tabs using this site.",
			);
		};
	});

	return dbPromise;
}

export async function loadCartridgeRam(
	romHash: string,
): Promise<Uint8Array | null> {
	if (romHash === "") {
		console.error("no ROM hash, unable to load RAM");
		return null;
	}

	const db = await openDatabase();

	return new Promise((resolve, reject) => {
		const transaction = db.transaction([STORE_NAME], "readonly");
		const objectStore = transaction.objectStore(STORE_NAME);
		const request = objectStore.get(romHash);

		request.onsuccess = () => {
			const result = request.result as SaveData | undefined;

			if (result) {
				console.log(
					`Save found from: ${new Date(result.updatedAt).toLocaleString()}`,
				);

				// IndexedDB might return an ArrayBuffer or Uint8Array depending on browser
				// Ensure we return a Uint8Array
				resolve(
					result.ram instanceof Uint8Array
						? result.ram
						: new Uint8Array(result.ram),
				);
			} else {
				console.log("No existing save file found.");
				resolve(null);
			}
		};

		request.onerror = () => {
			console.error("loadCartridgeRam() failed");
			reject(new Error("loadCartridgeRam() failed"));
		};
	});
}

export async function persistCartridgeRam(
	romHash: string,
	ram: Uint8Array,
	meta?: SaveData["meta"],
): Promise<void> {
	if (romHash === "") {
		return;
	}

	const db = await openDatabase();

	return new Promise((resolve, reject) => {
		const saveData: SaveData = {
			romHash: romHash,
			ram: ram,
			updatedAt: Date.now(),
			meta: meta,
		};
		const transaction = db.transaction([STORE_NAME], "readwrite");
		const store = transaction.objectStore(STORE_NAME);

		transaction.oncomplete = () => {
			resolve();
		};

		transaction.onerror = () => {
			console.error("persist RAM failed", transaction.error);
			reject(transaction.error);
		};

		store.put(saveData);
	});
}

export async function loadAppSettings(): Promise<AppSettings | null> {
	const db = await openDatabase();

	return new Promise((resolve, reject) => {
		const tx = db.transaction([SETTINGS_STORE], "readonly");
		const store = tx.objectStore(SETTINGS_STORE);
		const req = store.get(SETTINGS_KEY);

		req.onsuccess = () => {
			resolve(req.result?.value ?? null);
		};

		req.onerror = () => {
			reject(req.error);
		};
	});
}

export async function saveAppSettings(settings: AppSettings): Promise<void> {
	const db = await openDatabase();

	return new Promise((resolve, reject) => {
		const tx = db.transaction([SETTINGS_STORE], "readwrite");
		const store = tx.objectStore(SETTINGS_STORE);

		store.put({
			key: SETTINGS_KEY,
			value: settings,
		});

		tx.oncomplete = () => {
			resolve();
		};

		tx.onerror = () => {
			reject(tx.error);
		};
	});
}

export async function importData(jsonContent: string): Promise<ImportStats> {
	let backup: any;
	try {
		backup = JSON.parse(jsonContent);
	} catch (e) {
		throw new Error("Invalid JSON file");
	}

	if (backup.app !== "LuccaGB" || !Array.isArray(backup.saves)) {
		throw new Error("Unrecognized backup file format.");
	}

	const db = await openDatabase();

	return new Promise(async (resolve, reject) => {
		const transaction = db.transaction([STORE_NAME], "readwrite");
		const store = transaction.objectStore(STORE_NAME);

		const stats: ImportStats = {
			added: 0,
			updated: 0,
			skipped: 0,
			errors: 0,
		};

		transaction.oncomplete = () => {
			resolve(stats);
		};

		transaction.onerror = () => {
			reject(transaction.error);
		};

		for (const saveBackup of backup.saves as BackupFile["saves"]) {
			try {
				const ramData = base64ToBuffer(saveBackup.ramBase64);
				const importedUpdatedAt = saveBackup.updatedAt || 0;
				const hash = saveBackup.romHash;

				// Get the existing record for this ROM
				const request = store.get(hash);

				request.onsuccess = () => {
					const existingSave = request.result as SaveData | undefined;

					let shouldWrite = false;
					let isUpdate = false;

					if (!existingSave) {
						// Save doesn't exist locally -> Add it
						shouldWrite = true;
						console.log(
							`[Import] Adding ${hash} ${saveBackup.meta?.name}: (${new Date(importedUpdatedAt).toLocaleString()})`,
						);
					} else {
						// Save exists -> Compare timestamps
						const localUpdatedAt = existingSave.updatedAt || 0;

						// Keep the most recent save
						if (importedUpdatedAt > localUpdatedAt) {
							shouldWrite = true;
							isUpdate = true;
							console.log(
								`[Import] Upgrading ${hash} ${saveBackup.meta?.name}: Local(${new Date(localUpdatedAt).toLocaleString()}) < SaveFile(${new Date(importedUpdatedAt).toLocaleString()})`,
							);
						} else {
							stats.skipped++;
							console.log(
								`[Import] Skipping ${saveBackup.romHash} ${saveBackup.meta?.name}: save file is older or same`,
							);
						}
					}

					if (shouldWrite) {
						const newSave: SaveData = {
							romHash: hash,
							updatedAt: importedUpdatedAt,
							meta: saveBackup.meta,
							ram: ramData,
						};

						const putRequest = store.put(newSave);
						putRequest.onsuccess = () => {
							if (isUpdate) {
								stats.updated++;
							} else {
								stats.added++;
							}
						};
					}
				};
			} catch (error) {
				console.error("Error processing save entry", error);
				stats.errors++;
			}
		}

		if (backup.settings) {
			const settings = backup.settings as AppSettings;

			await saveAppSettings({
				...settings,
				updatedAt: Date.now(),
			});

			appState.initializeAppSettings();
		}
	});
}

export async function exportData(): Promise<void> {
	const db = await openDatabase();

	const saves = await new Promise<SaveData[]>((resolve, reject) => {
		const transaction = db.transaction([STORE_NAME], "readonly");
		const store = transaction.objectStore(STORE_NAME);
		const request = store.getAll();

		request.onsuccess = () => {
			resolve(request.result as SaveData[]);
		};
		request.onerror = () => {
			reject(request.error);
		};
	});

	const settings = await loadAppSettings();

	// Transform SaveData (Uint8Array) -> Backup format (Base64)
	const backup: BackupFile = {
		version: 1,
		timestamp: Date.now(),
		app: "LuccaGB",
		saves: saves.map((save) => ({
			romHash: save.romHash,
			updatedAt: save.updatedAt,
			meta: save.meta,
			// Convert raw bytes to Base64 string for JSON compatibility
			ramBase64: bufferToBase64(
				save.ram instanceof Uint8Array ? save.ram : new Uint8Array(save.ram),
			),
		})),
		settings: settings ?? undefined,
	};

	// Create a blob and trigger download
	const blob = new Blob([JSON.stringify(backup, null, 2)], {
		type: "application/json",
	});
	const url = URL.createObjectURL(blob);
	const a = document.createElement("a");
	const dateStr = new Date().toISOString().split("T")[0];

	a.href = url;
	a.download = `LuccaGB-Backup-${dateStr}.json`;
	document.body.appendChild(a);
	a.click();

	// Cleanup
	document.body.removeChild(a);
	URL.revokeObjectURL(url);
}

function bufferToBase64(buffer: Uint8Array): string {
	let binary = "";
	const len = buffer.byteLength;
	for (let i = 0; i < len; i++) {
		binary += String.fromCharCode(buffer[i]);
	}
	return window.btoa(binary);
}

function base64ToBuffer(base64: string): Uint8Array {
	const binaryString = window.atob(base64);
	const len = binaryString.length;
	const bytes = new Uint8Array(len);
	for (let i = 0; i < len; i++) {
		bytes[i] = binaryString.charCodeAt(i);
	}
	return bytes;
}
