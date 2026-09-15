import { createSignal, onCleanup } from "solid-js";

/**
 * Reactive signal tracking whether the viewport matches the mobile breakpoint.
 */
export const useIsMobile = () => {
	const mediaQuery = window.matchMedia("(max-width: 768px)");
	const [isMobile, setIsMobile] = createSignal(mediaQuery.matches);

	const handleChange = (event: MediaQueryListEvent) => {
		setIsMobile(event.matches);
	};

	mediaQuery.addEventListener("change", handleChange);
	onCleanup(() => mediaQuery.removeEventListener("change", handleChange));

	return isMobile;
};
