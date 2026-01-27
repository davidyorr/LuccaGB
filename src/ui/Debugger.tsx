import styles from "./Debugger.module.css";

import { createSignal, Show, For, type Component } from "solid-js";
import { store } from "../core/store";
import type { GameboyDebugInfo } from "../wasm";

export const Debugger: Component = () => {
	const toHex = (val: number | undefined) =>
		val === undefined
			? "N/A"
			: "0x" + val.toString(16).toUpperCase().padStart(4, "0");

	const toBinary = (val: number | undefined) =>
		val === undefined
			? "N/A"
			: "0b" +
				val
					.toString(2)
					.padStart(8, "0")
					.replace(/(\d{4})(?=\d)/g, "$1_");

	return (
		<Show when={store.state.settings.isDebuggerOpen}>
			<Show
				when={debugInfo()}
				fallback={
					<div class={styles.debugPanel}>Waiting for Debug Info...</div>
				}
			>
				{(debugData) => (
					<div class={styles.debugPanel}>
						<table class={styles.debugTable}>
							<tbody>
								<tr>
									<td>Title</td>
									<td class={styles.cartridgeTitle}>
										{debugData().cartridge.title}
									</td>
								</tr>
								<tr>
									<td>Cartridge Type</td>
									<td>{debugData().cartridge.cartridgeType}</td>
								</tr>
								<tr>
									<td>ROM Size Code</td>
									<td>{debugData().cartridge.romSizeCode}</td>
								</tr>
								<tr>
									<td>RAM Size Code</td>
									<td>{debugData().cartridge.ramSizeCode}</td>
								</tr>
							</tbody>
						</table>

						<h2>CPU State</h2>
						<table class={styles.debugTable}>
							<thead>
								<tr>
									<th>Register</th>
									<th>Value (Hex)</th>
								</tr>
							</thead>
							<tbody>
								<tr>
									<td>AF</td>
									<td>{toHex(debugData().cpu.registers16.AF)}</td>
								</tr>
								<tr>
									<td>BC</td>
									<td>{toHex(debugData().cpu.registers16.BC)}</td>
								</tr>
								<tr>
									<td>DE</td>
									<td>{toHex(debugData().cpu.registers16.DE)}</td>
								</tr>
								<tr>
									<td>HL</td>
									<td>{toHex(debugData().cpu.registers16.HL)}</td>
								</tr>
								<tr>
									<td>SP</td>
									<td>{toHex(debugData().cpu.registers16.SP)}</td>
								</tr>
								<tr>
									<td>PC</td>
									<td>{toHex(debugData().cpu.registers16.PC)}</td>
								</tr>
							</tbody>
						</table>

						<h3>Flags (F Register)</h3>
						<table class={styles.debugTable}>
							<thead>
								<tr>
									<th>Flag</th>
									<th>Value</th>
								</tr>
							</thead>
							<tbody>
								<tr>
									<td>Z (Zero)</td>
									<td>{debugData().cpu.flags.Z}</td>
								</tr>
								<tr>
									<td>N (Subtraction)</td>
									<td>{debugData().cpu.flags.N}</td>
								</tr>
								<tr>
									<td>H (Half Carry)</td>
									<td>{debugData().cpu.flags.H}</td>
								</tr>
								<tr>
									<td>C (Carry)</td>
									<td>{debugData().cpu.flags.C}</td>
								</tr>
							</tbody>
						</table>

						<h2>APU</h2>
						<h3>Registers</h3>
						<table class={styles.debugTable}>
							<thead>
								<tr>
									<th>Register</th>
									<th>Value</th>
								</tr>
							</thead>
							<tbody>
								<tr>
									<td>NR30</td>
									<td>{toBinary(debugData().apu.registers.NR30)}</td>
								</tr>
								<tr>
									<td>NR32</td>
									<td>{toBinary(debugData().apu.registers.NR32)}</td>
								</tr>
							</tbody>
						</table>

						<h3>Wave RAM</h3>
						<table class={styles.debugTable}>
							<thead>
								<tr>
									<th>Address</th>
									<th>Value (Hex)</th>
								</tr>
							</thead>
							<tbody>
								<For each={debugData().apu.waveRam}>
									{(val, i) => (
										<tr>
											<td>0x{(0xff30 + i()).toString(16).toUpperCase()}</td>
											<td>{toHex(val)}</td>
										</tr>
									)}
								</For>
							</tbody>
						</table>
					</div>
				)}
			</Show>
		</Show>
	);
};

export const DebuggerToggle: Component = () => {
	return (
		<label>
			<input
				type="checkbox"
				checked={store.state.settings.isDebuggerOpen}
				onChange={(event) => {
					store.actions.setDebuggerOpen(event.currentTarget.checked);
				}}
			/>
			Show Debugger
		</label>
	);
};

const [debugInfo, setDebugInfo] = createSignal<GameboyDebugInfo | null>(null);

export function updateDebugger() {
	const data = window.getDebugInfo?.();
	if (data) {
		setDebugInfo(data);
	}
}
