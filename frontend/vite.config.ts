import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";
import path from "path";

// https://vitejs.dev/config/
export default defineConfig({
  optimizeDeps: {
    exclude: ['@3d-dice/dice-box'],
  },
  server: {
    host: "::",
    port: 8081,
    fs: {
      allow: [path.resolve(__dirname, "..")],
    },
    hmr: {
      overlay: false,
    },
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
});
