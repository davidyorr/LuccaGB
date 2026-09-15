import { onMount, onCleanup, type Component, createEffect } from "solid-js";
import { createStore } from "solid-js/store";
import styles from "./OnScreenJoypad.module.css";
import { inputManager } from "../services/input-manager";
import { useIsMobile } from "./hooks/useIsMobile";

export const OnScreenJoypad: Component = () => {
	const activePointers = new Map<number, string>();
	const buttonRefs = new Map<string, HTMLElement>();
	const buttonBounds = new Map<string, DOMRect>();
	const isMobile = useIsMobile();

	const [pressedButtons, setPressedButtons] = createStore<
		Record<string, boolean>
	>({});

	let containerRef!: HTMLDivElement;

	// Cache the dimensions and positions of all virtual buttons
	const updateBounds = () => {
		for (const [name, el] of buttonRefs.entries()) {
			buttonBounds.set(name, el.getBoundingClientRect());
		}
	};

	onMount(() => {
		// Calculate initial bounds
		updateBounds();

		// Use ResizeObserver to automatically update bounds if layout changes
		const resizeObserver = new ResizeObserver(() => updateBounds());
		if (containerRef) {
			resizeObserver.observe(containerRef);
		}

		// Fallbacks for layout shifts (e.g., orientation changes)
		window.addEventListener("resize", updateBounds);
		window.addEventListener("scroll", updateBounds, { passive: true });

		// Catch ANY released touches globally to prevent stuck buttons
		const handleGlobalRelease = (e: PointerEvent) => {
			if (activePointers.has(e.pointerId)) {
				setPointerButton(e.pointerId, undefined);
			}
		};
		window.addEventListener("pointerup", handleGlobalRelease);
		window.addEventListener("pointercancel", handleGlobalRelease);

		onCleanup(() => {
			resizeObserver.disconnect();
			window.removeEventListener("resize", updateBounds);
			window.removeEventListener("scroll", updateBounds);
			window.removeEventListener("pointerup", handleGlobalRelease);
			window.removeEventListener("pointercancel", handleGlobalRelease);
		});
	});

	// Recompute the bounds when switching to mobile
	createEffect(() => {
		if (isMobile()) {
			updateBounds();
		}
	});

	const setPointerButton = (
		pointerId: number,
		buttonName: string | undefined,
	) => {
		const previousButton = activePointers.get(pointerId);

		if (previousButton === buttonName) return;

		// Release the previous button
		if (previousButton) {
			inputManager.releaseTouchscreenButton(previousButton, pointerId);
			activePointers.delete(pointerId);

			// Only visually un-press if no other fingers are currently holding it down
			const isStillPressed = [...activePointers.values()].includes(
				previousButton,
			);
			setPressedButtons(previousButton, isStillPressed);
		}

		// Press the new button
		if (buttonName) {
			activePointers.set(pointerId, buttonName);
			inputManager.pressTouchscreenButton(buttonName, pointerId);
			setPressedButtons(buttonName, true);
		}
	};

	// Determine which button (if any) the finger is currently over based on cached bounds
	const getButtonFromPoint = (x: number, y: number): string | undefined => {
		for (const [name, rect] of buttonBounds.entries()) {
			if (
				x >= rect.left &&
				x <= rect.right &&
				y >= rect.top &&
				y <= rect.bottom
			) {
				return name;
			}
		}
		return undefined;
	};

	const handlePointerDown = (e: PointerEvent) => {
		// Only allow left-click/touch
		if (e.pointerType === "mouse" && e.buttons !== 1) {
			return;
		}

		e.preventDefault();

		// Lock the pointer to the container so we can track the finger even if it leaves the container area
		try {
			if (containerRef) {
				containerRef.setPointerCapture(e.pointerId);
			}
		} catch (error) {
			console.debug(error);
		}

		const buttonName = getButtonFromPoint(e.clientX, e.clientY);
		setPointerButton(e.pointerId, buttonName);
	};

	const handlePointerMove = (e: PointerEvent) => {
		// Ignore hover events where no fingers are down
		if (e.buttons === 0) {
			return;
		}

		e.preventDefault();

		const buttonName = getButtonFromPoint(e.clientX, e.clientY);
		setPointerButton(e.pointerId, buttonName);
	};

	const handlePointerUp = (e: PointerEvent) => {
		e.preventDefault();

		try {
			if (containerRef && containerRef.hasPointerCapture(e.pointerId)) {
				containerRef.releasePointerCapture(e.pointerId);
			}
		} catch (err) {}

		setPointerButton(e.pointerId, undefined);
	};

	const handlePointerCancel = handlePointerUp;

	const button = (name: string, className: string = "") => (
		<div
			ref={(el) => buttonRefs.set(name, el)}
			class={`${styles.button} ${className}`}
			classList={{ [styles.pressed]: pressedButtons[name] }}
			data-joypad-button={name}
			aria-label={name}
		/>
	);

	return (
		<div
			ref={containerRef}
			class={styles.container}
			onPointerDown={handlePointerDown}
			onPointerMove={handlePointerMove}
			onPointerUp={handlePointerUp}
			onPointerCancel={handlePointerCancel}
			// Stop mobile browsers from opening menus on rapid/long taps
			onContextMenu={(e) => e.preventDefault()}
		>
			<div class={styles.topRow}>
				<div class={styles.dpad}>
					{button("UP", styles.up)}
					{button("LEFT", styles.left)}
					<div class={styles.center} />
					{button("RIGHT", styles.right)}
					{button("DOWN", styles.down)}
				</div>

				<div class={styles.actionButtons}>
					<div class={styles.actionButton}>{button("B")}</div>
					<div class={`${styles.actionButton} ${styles.aOffset}`}>
						{button("A")}
					</div>
				</div>
			</div>

			<div class={styles.bottomRow}>
				{button("SELECT", styles.systemButton)}
				{button("START", styles.systemButton)}
			</div>
		</div>
	);
};
