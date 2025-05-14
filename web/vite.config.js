import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig(() => {
  return {
    server: {
        open: true,
        port: 3000, 
        headers: {
          "Cross-Origin-Opener-Policy": "cross-origin",
          // "Cross-Origin-Embedder-Policy": "require-corp",
          'Cross-Origin-Embedder-Policy': 'credentialless',
          'Access-Control-Allow-Origin' :'*'
        },
      },
    build: {
      outDir: 'build',
    },
    plugins: [react()],
    optimizeDeps: {
      exclude: ["@ffmpeg/ffmpeg", "@ffmpeg/util"],
    },
  };
});