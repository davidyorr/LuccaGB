import { gameLoop } from "../core/game-loop";
import { store } from "../core/store";
import { audioController } from "./audio-controller";

let mediaRecorder: MediaRecorder | null = null;
let canvasStream: MediaStream | null = null;
let chunks: Blob[] = [];

let recordingStartTime = 0;
let recordingTimer: number | undefined;
let recordingSize = 0;

export function toggleRecordingGameplay() {
	if (!store.state.isRomLoaded) {
		return;
	}

	if (!store.state.isRecordingGameplay) {
		startRecordingGameplay();
	} else {
		stopRecordingGameplay();
	}
}

function startRecordingGameplay() {
	const renderer = gameLoop.renderer();

	if (!renderer) {
		console.error("Cannot record: no renderer attached");
		return;
	}

	// get the first supported MIME type, preferring VP9 > VP8 > generic WebM
	const mimeType = [
		"video/webm;codecs=vp9,opus",
		"video/webm;codecs=vp8,opus",
		"video/webm",
	].find((type) => MediaRecorder.isTypeSupported(type));

	if (!mimeType) {
		console.error("No supported recording format");
		return;
	}

	// video stream
	const canvas = renderer.canvas();
	canvasStream = canvas.captureStream();

	// audio stream
	const audioStream = audioController.recordingStream();

	// combine streams
	const stream = new MediaStream([
		...canvasStream.getVideoTracks(),
		...audioStream.getAudioTracks(),
	]);

	chunks = [];

	mediaRecorder = new MediaRecorder(stream, { mimeType });

	mediaRecorder.addEventListener("dataavailable", (event) => {
		if (event.data.size === 0) {
			return;
		}

		chunks.push(event.data);

		recordingSize += event.data.size;
		store.actions.setRecordingSize(recordingSize);
	});

	mediaRecorder.addEventListener("stop", () => {
		const blob = new Blob(chunks, { type: mimeType });
		const url = URL.createObjectURL(blob);

		// create a temporary link element to trigger the download
		const link = document.createElement("a");
		link.href = url;
		link.download = `luccagb-video-${new Date().toISOString()}.webm`;

		// trigger the click and cleanup
		document.body.appendChild(link);
		link.click();
		document.body.removeChild(link);
		URL.revokeObjectURL(url);

		canvasStream?.getTracks().forEach((track) => track.stop());

		canvasStream = null;
		mediaRecorder = null;
		chunks = [];

		clearInterval(recordingTimer);
		recordingTimer = undefined;

		store.actions.setRecordingGameplay(false);
	});

	mediaRecorder.start(1000);

	recordingStartTime = performance.now();

	store.actions.setRecordingSize(0);
	store.actions.setRecordingDuration(0);

	recordingTimer = window.setInterval(() => {
		store.actions.setRecordingDuration(performance.now() - recordingStartTime);
	}, 250);

	store.actions.setRecordingGameplay(true);
}

function stopRecordingGameplay() {
	if (!mediaRecorder || mediaRecorder.state === "inactive") {
		return;
	}

	mediaRecorder.stop();
}
