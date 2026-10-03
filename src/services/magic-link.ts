import { generateBackupData, importData } from "./storage";

export async function createMagicLink(): Promise<string> {
	const backup = await generateBackupData();

	const response = await fetch("/api/transfer", {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(backup),
	});

	if (!response.ok) {
		throw new Error("Failed to generate link");
	}

	const { id } = await response.json();

	const url = new URL(window.location.href);
	url.searchParams.set("import", id);
	return url.toString();
}

export async function fetchMagicLinkData(id: string): Promise<string> {
	const response = await fetch(`/api/transfer?id=${id}`);

	if (!response.ok) {
		throw new Error("Magic link expired or invalid");
	}

	return await response.text();
}
