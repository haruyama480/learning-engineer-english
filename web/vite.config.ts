import adapter from '@sveltejs/adapter-static';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

const base = (process.env.BASE_PATH ?? '') as '' | `/${string}`;

export default defineConfig({
	build: {
		// Import した画像・フォントなどを data URL として HTML に埋め込む
		assetsInlineLimit: Infinity
	},
	plugins: [
		sveltekit({
			compilerOptions: {
				// Force runes mode for the project, except for libraries. Can be removed in svelte 6.
				runes: ({ filename }) =>
					filename.split(/[/\\]/).includes('node_modules') ? undefined : true
			},
			adapter: adapter(),
			output: {
				// JS / CSS を index.html にインライン化する
				bundleStrategy: 'inline'
			},
			paths: { base }
		})
	]
});
