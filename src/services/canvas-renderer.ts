export class CanvasRenderer {
	private visibleCanvasCtx: CanvasRenderingContext2D;
	private offscreenCanvasCtx: CanvasRenderingContext2D;
	private imageData: ImageData;
	private displayWidth = 160;
	private displayHeight = 144;

	constructor(canvasId: string) {
		const visibleCanvas = document.getElementById(canvasId);
		if (!(visibleCanvas instanceof HTMLCanvasElement)) {
			throw new Error("invalid canvas id");
		}
		visibleCanvas.width = this.displayWidth;
		visibleCanvas.height = this.displayHeight;
		this.visibleCanvasCtx = visibleCanvas.getContext("2d")!;
		if (this.visibleCanvasCtx === null) {
			throw new Error("error getting canvas context");
		}
		this.visibleCanvasCtx.imageSmoothingEnabled = false;
		this.imageData = this.visibleCanvasCtx.createImageData(
			this.displayWidth,
			this.displayHeight,
		);

		// offscreen setup
		const offscreenCanvas = document.createElement("canvas");
		offscreenCanvas.width = this.displayWidth;
		offscreenCanvas.height = this.displayHeight;
		this.offscreenCanvasCtx = offscreenCanvas.getContext("2d")!;
		if (this.offscreenCanvasCtx === null) {
			throw new Error("error getting offscreen canvas context");
		}
	}

	public drawFrame(frameData: Uint8Array) {
		// put the 160x144 data onto the same size offscreen canvas
		this.imageData.data.set(frameData);
		this.offscreenCanvasCtx.putImageData(this.imageData, 0, 0);

		// scale the image from the offscreen canvas onto the visible canvas
		this.visibleCanvasCtx.drawImage(
			this.offscreenCanvasCtx.canvas,
			0,
			0,
			this.visibleCanvasCtx.canvas.width,
			this.visibleCanvasCtx.canvas.height,
		);
	}

	public setScale(scale: number | "fit") {
		const container = document.getElementById("canvas-container");
		const canvas = this.visibleCanvasCtx.canvas;

		if (!container || !canvas) {
			return;
		}

		if (scale === "fit") {
			const maxWidth = window.innerWidth;
			const maxHeight = window.innerHeight;

			const fitScale = Math.floor(
				Math.min(maxWidth / this.displayWidth, maxHeight / this.displayHeight),
			);

			const finalScale = Math.max(1, fitScale);

			container.style.width = `${this.displayWidth * finalScale}px`;
			container.style.height = `${this.displayHeight * finalScale}px`;
		} else {
			container.style.width = `${this.displayWidth * scale}px`;
			container.style.height = `${this.displayHeight * scale}px`;
		}
	}

	public takeScreenshot() {
		// grab the current pixels as a PNG image string
		const imageURL = this.visibleCanvasCtx.canvas.toDataURL("image/png");

		// create a temporary link element to trigger the download
		const link = document.createElement("a");
		link.href = imageURL;
		link.download = `luccagb-screenshot-${new Date().toISOString()}.png`;

		// trigger the click and cleanup
		document.body.appendChild(link);
		link.click();
		document.body.removeChild(link);
	}
}
