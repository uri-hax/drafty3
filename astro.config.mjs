// @ts-check
import { defineConfig } from 'astro/config';
import react from '@astrojs/react';

const basePath = process.env.PUBLIC_BASE_PATH || '/drafty3';
const localDevAPI = process.env.PUBLIC_LOCAL_DEV_API || 'http://localhost:8081';

export default defineConfig({
  integrations: [react()],
  output: 'static',
  site: 'https://uri-hax.github.io',
  base: basePath,

  // DEV ONLY proxy
  vite: {
    server: {
      proxy: {
        '/api': localDevAPI,
      },
    },
  },
});
