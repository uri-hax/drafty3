// @ts-check
import { defineConfig } from 'astro/config';
import react from '@astrojs/react';

// https://astro.build/config
export default defineConfig({
    integrations: [react()],
    output: 'static',
    site: 'https://uri-hax.github.io',
    base: "/drafty3", 
    // temporary proxy to backend during development
    vite: {
      server: {
        proxy: {
          "/api": "http://localhost:8081",
        },
      },
    },
  });
