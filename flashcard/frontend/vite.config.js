import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  server: {
    port: 8080,
    strictPort: false,
    cors: true,
    host: '0.0.0.0',
    hmr: {
      protocol: 'ws',
      host: 'localhost',
    }
  },
  optimizeDeps: {
    include: ['svelte-routing']
  },
  build: {
    outDir: 'dist',
    assetsDir: 'assets',
    emptyOutDir: true
  }
})
