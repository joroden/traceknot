import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  build: {
    rollupOptions: {
      input: {
        index: new URL("./index.html", import.meta.url).pathname,
        select: new URL("./select.html", import.meta.url).pathname,
      },
    },
  },
  server: {
    proxy: {
      "/api": "http://127.0.0.1:4318",
    },
  },
});
