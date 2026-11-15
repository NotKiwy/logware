import { defineConfig } from 'vite'

export default defineConfig({
  resolve: { alias: { '@web': new URL('./src', import.meta.url).pathname } },
  build: { outDir: 'dist', target: 'es2022' },
  server: { proxy: { '/api': 'http://localhost:8080' } },
})
