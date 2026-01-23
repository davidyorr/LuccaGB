export function debounce<T extends (...args: any[]) => void>(
	func: T,
	wait: number,
): T {
	let timeout: number | undefined;

	return function (this: any, ...args: any[]) {
		window.clearTimeout(timeout);
		timeout = window.setTimeout(() => {
			func.apply(this, args);
		}, wait);
	} as T;
}
