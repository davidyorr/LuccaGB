import { defineConfig } from "vite";
import solidPlugin from "vite-plugin-solid";
import basicSsl from "@vitejs/plugin-basic-ssl";

const useSsl = process.env.USE_SSL === "true";

export default defineConfig({
	plugins: [solidPlugin(), ...(useSsl ? [basicSsl()] : [])],
	server: {
		host: "0.0.0.0",
	},
});
