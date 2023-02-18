import { vanillaExtractPlugin } from '@vanilla-extract/vite-plugin'
import path from 'path'
import { defineConfig } from 'vite'
import solidPlugin from 'vite-plugin-solid'
import solidSvg from 'vite-plugin-solid-svg'

export default defineConfig({
  plugins: [
    solidPlugin(),
    vanillaExtractPlugin(),
    solidSvg({ defaultAsComponent: true }),
  ],
  server: {
    port: 5173,
    open: true,
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