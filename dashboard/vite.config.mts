import path from "node:path";
import { visualizer } from "rollup-plugin-visualizer";
import UnoCSS from "unocss/vite";
import Unfonts from "unplugin-fonts/vite";
import { defineConfig, type PluginOption } from "vite";
import viteCompression from "vite-plugin-compression";
import solidPlugin from "vite-plugin-solid";
import solidSvg from "vite-plugin-solid-svg";

export default defineConfig(({ mode }) => ({
  plugins: [
    UnoCSS(),
    solidPlugin(),
    solidSvg({ defaultAsComponent: true }),
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
    conditions: ["module", "browser", "default"],
  },
  build: {
    target: "esnext",
    rollupOptions: {
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
    server: {
      deps: {
        inline: ["@solidjs/router"],
      },
    },
  },
}));
