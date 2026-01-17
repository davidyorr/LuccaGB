const DB_NAME = "LuccaGB-Database";
const DB_VERSION = 1;
const STORE_NAME = "cartridgeRam";

type SaveData = {
	romHash: string;
	ram: Uint8Array;
	updatedAt: number;
	meta?: {
		name?: string;
		playTime?: number;
	};
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
			// Create store if missing
			if (!db.objectStoreNames.contains(STORE_NAME)) {
				db.createObjectStore(STORE_NAME, {
					keyPath: "romHash",
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
