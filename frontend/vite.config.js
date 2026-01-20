import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  build: {
    outDir: 'dist'
  },
  server: {
    proxy: {
      '/api': 'http://localhost:8443',
      '/login': 'http://localhost:8443',
      '/signup': 'http://localhost:8443',
      '/hasadm': 'http://localhost:8443'
    }
  }
})
