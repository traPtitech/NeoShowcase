import { vanillaExtractPlugin } from '@vanilla-extract/vite-plugin'
import path from 'path'
import { defineConfig } from 'vite'
import solidPlugin from 'vite-plugin-solid'

export default defineConfig({
	plugins: [solidPlugin(), vanillaExtractPlugin()],
	server: {
		port: 3000,
	},
	resolve: {
		alias: {
			'/@': path.resolve(__dirname, '/src'),
		},
	},
	build: {
		target: 'esnext',
	},
})
