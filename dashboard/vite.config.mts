import path from 'node:path'
import { defineConfig, type PluginOption } from 'vite'
import solidPlugin from 'vite-plugin-solid'
import solidSvg from 'vite-plugin-solid-svg'
import Unfonts from 'unplugin-fonts/vite'
import { macaronVitePlugin } from '@macaron-css/vite'
import viteCompression from 'vite-plugin-compression'
import { visualizer } from 'rollup-plugin-visualizer'
import UnoCSS from 'unocss/vite'

export default defineConfig(({ mode }) => ({
  plugins: [
    UnoCSS(),
    macaronVitePlugin(), // comes first
    solidPlugin(),
    solidSvg({ defaultAsComponent: true }),
    Unfonts({
      google: {
        families: [
          'Lato',
          {
            name: 'Material+Symbols+Rounded',
            styles: 'opsz,wght,FILL,GRAD@20..24,300,0..1,0',
          },
        ],
      },
    }),
    viteCompression({ algorithm: 'gzip' }),
    viteCompression({ algorithm: 'brotliCompress' }),
  ],
  server: {
    port: 5173,
    proxy: {
      '/neoshowcase.protobuf.APIService': {
        target: 'http://ns.local.trapti.tech',
      },
    },
  },
  resolve: {
    alias: {
      '/@': path.resolve(__dirname, '/src'),
    },
  },
  build: {
    target: 'esnext',
    rollupOptions: {
      plugins: [
        mode === 'analyze' &&
          (visualizer({
            open: true,
            filename: 'dist/stats.html',
            gzipSize: true,
            brotliSize: true,
          }) as PluginOption),
      ],
    },
  },
}))
