import { fileURLToPath } from "node:url";
import { defineConfig } from "vite";

export default defineConfig({
  root: fileURLToPath(new URL("./receiver", import.meta.url)),
  build: {
    outDir: fileURLToPath(new URL("./dist-receiver", import.meta.url)),
    emptyOutDir: true,
    modulePreload: { polyfill: false },
  },
});
