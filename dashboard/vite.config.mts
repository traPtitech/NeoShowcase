import path from "node:path";
import { visualizer } from "rollup-plugin-visualizer";
import UnoCSS from "unocss/vite";
import Unfonts from "unplugin-fonts/vite";
import { defineConfig, type PluginOption } from "vite";
import viteCompression from "vite-plugin-compression";
import solid from "@solidjs/vite-plugin";

export default defineConfig(({ mode }) => ({
  plugins: [
    UnoCSS(),
    solid(),
    Unfonts({
      google: {
        families: ["Lato"],
      },
    }),
    viteCompression({ algorithm: "gzip" }),
    viteCompression({ algorithm: "brotliCompress" }),
  ],
  server: {
    port: 5173,
    proxy: {
      "/neoshowcase.protobuf.APIService": {
        target: "http://ns.local.trapti.tech",
        changeOrigin: true,
      },
    },
    allowedHosts: ["ns.local.trapti.tech"],
  },
  resolve: {
    alias: {
      "/@": path.resolve(__dirname, "/src"),
    },
  },
  build: {
    target: "esnext",
    rolldownOptions: {
      plugins: [
        mode === "analyze" &&
        (visualizer({
          open: true,
          filename: "dist/stats.html",
          gzipSize: true,
          brotliSize: true,
        }) as PluginOption),
      ],
    },
  },
  test: {
    environment: "jsdom",
  },
}));
