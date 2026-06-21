import { animate } from 'motion';

export interface MotionParams {
	keyframes: Parameters<typeof animate>[1];
	options?: Parameters<typeof animate>[2];
}

/**
 * Svelte Action to apply motion.dev animations to a DOM element.
 * Usage: <div use:motion={{ keyframes: { opacity: [0, 1] }, options: { duration: 0.5 } }}>
 */
export function motion(node: HTMLElement, params: MotionParams) {
	if (!params || !params.keyframes) return;

	let controls = animate(node, params.keyframes, params.options);

	// Clear inline styles once entrance animation finishes
	// so that they don't override standard CSS hover rules
	controls.then(() => {
		node.style.transform = '';
		node.style.opacity = '';
	});

	return {
		update(newParams: MotionParams) {
			controls.stop();
			if (!newParams || !newParams.keyframes) return;
			controls = animate(node, newParams.keyframes, newParams.options);
			controls.then(() => {
				node.style.transform = '';
				node.style.opacity = '';
			});
		},
		destroy() {
			controls.stop();
		}
	};
}
