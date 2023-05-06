import path from 'path'
import { defineConfig } from 'vite'
import solidPlugin from 'vite-plugin-solid'
import solidSvg from 'vite-plugin-solid-svg'
import VitePluginFonts from 'vite-plugin-fonts'
import { macaronVitePlugin } from '@macaron-css/vite'

export default defineConfig({
  plugins: [
    macaronVitePlugin(), // comes first
    solidPlugin(),
    solidSvg({ defaultAsComponent: true }),
    VitePluginFonts({
      google: {
        families: [
          'Mulish',
          'Noto Sans JP',
        ],
      }
    }),
  ],
  server: {
    port: 5173,
    proxy: {
      '/neoshowcase.protobuf.APIService': {
        target: 'http://ns.local.trapti.tech'
      }
    }
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
