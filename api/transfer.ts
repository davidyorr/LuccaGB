import { randomBytes } from "node:crypto";
import { Redis } from "@upstash/redis";
import type { VercelRequest, VercelResponse } from "@vercel/node";

const redis = Redis.fromEnv();

const EXPIRY_SECONDS = 600; // 10 minutes

export default async function handler(req: VercelRequest, res: VercelResponse) {
	if (req.method === "POST") {
		const contentLength = req.headers["content-length"];
		if (contentLength && parseInt(contentLength) > 1024 * 1024) {
			return res.status(413).json({ error: "Payload too large. Max 1MB." });
		}

		try {
			const id = randomBytes(16).toString("base64url");
			await redis.set(`magic_link:${id}`, req.body, { ex: EXPIRY_SECONDS });
			return res.status(200).json({ id });
		} catch (error) {
			console.error(error);
			return res.status(500).json({ error: "Failed to create magic link" });
		}
	}

	if (req.method === "GET") {
		try {
			const { id } = req.query;
			if (typeof id !== "string" || !id) {
				return res.status(400).json({ error: "Missing ID" });
			}

			const data = await redis.get(`magic_link:${id}`);
			if (!data) {
				return res.status(404).json({ error: "Link expired or not found" });
			}

			return res.status(200).send(data);
		} catch (error) {
			console.error(error);
			return res.status(500).json({ error: "Failed to retrieve data" });
		}
	}

	return res.status(405).json({ error: "Method not allowed" });
}
