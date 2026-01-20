export function parseTraceLogs(buffer: Uint8Array): string {
	let lines = [];
	for (let i = 0; i < buffer.length; ) {
		const type = buffer[i];

		if (type === 0) {
			// Instruction
			const pc = (buffer[i + 1] << 8) | buffer[i + 2];
			const opcode = buffer[i + 3];
			lines.push(
				`EXEC PC:0x${pc.toString(16).padStart(4, "0")} OP:0x${opcode.toString(16).padStart(2, "0")}`,
			);
			// if we log a 3rd thing, use this instead
			// const view = new DataView(buffer.buffer);
			// const data = view.getUint32(i + 4, false);
			// lines.push(
			// 	`EXEC PC:0x${pc.toString(16).padStart(4, "0")} OP:0x${opcode.toString(16).padStart(2, "0")} OTHER:${data}`,
			// );
			i += 8;
		} else if (type === 1) {
			// Memory read
			const addr = (buffer[i + 1] << 8) | buffer[i + 2];
			const value = buffer[i + 3];
			lines.push(
				`READ [0x${addr.toString(16).padStart(4, "0")}] = 0x${value.toString(16).padStart(2, "0")}`,
			);
			i += 4;
		} else if (type === 2) {
			// Memory write
			const addr = (buffer[i + 1] << 8) | buffer[i + 2];
			const value = buffer[i + 3];
			lines.push(
				`WRITE [0x${addr.toString(16).padStart(4, "0")}] = 0x${value.toString(16).padStart(2, "0")}`,
			);
			i += 4;
		} else {
			i++; // Skip unknown types
		}
	}

	return lines.join("\n");
}

export function downloadTraceLogs(content: string) {
	const blob = new Blob([content], { type: "text/plain" });
	const url = URL.createObjectURL(blob);
	const a = document.createElement("a");
	a.href = url;
	a.download = `luccagb-logs-${new Date().toISOString()}.txt`;
	a.click();
}
