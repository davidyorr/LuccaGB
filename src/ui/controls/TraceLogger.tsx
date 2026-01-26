import { createSignal, Show } from "solid-js";
import {
	downloadTraceLogs as downloadTraceLog,
	parseTraceLogs,
} from "../../utils/trace-logger";

export const TraceLogger = () => {
	// Hide in production
	if (import.meta.env.PROD) {
		return null;
	}

	const [isLogging, setIsLogging] = createSignal(false);

	const handleInputToggle = (event: Event) => {
		const checked = (event.target as HTMLInputElement).checked;
		setIsLogging(checked);

		if (checked) {
			window.enableTraceLogging();
		} else {
			window.disableTraceLogging();
		}
	};

	const handleDownloadClick = () => {
		const buffer = window.getTraceLogs();

		if (!buffer || buffer.length === 0) {
			alert("No logs available.");
			return;
		}

		const text = parseTraceLogs(buffer);
		downloadTraceLog(text);
	};

	return (
		<>
			<label id="trace-log-toggle-container">
				<input
					type="checkbox"
					id="trace-log-checkbox"
					onChange={handleInputToggle}
					checked={isLogging()}
				/>
				Trace Log
			</label>

			<Show when={isLogging()}>
				<button id="download-trace-log-button" onClick={handleDownloadClick}>
					Download trace log
				</button>
			</Show>
		</>
	);
};
