import { VuetifyResolver } from 'unplugin-vue-components/resolvers';
import { createVuePlugin as vue } from 'vite-plugin-vue2';
import Components from 'unplugin-vue-components/vite';
import { defineConfig } from 'vite';
import { VitePWA } from 'vite-plugin-pwa';

export default defineConfig({
  plugins: [
    vue(),
    Components({
      directives: false,
      resolvers: [
        VuetifyResolver(),
      ],
      types: [
        {
          from: 'vue-router',
          names: ['RouterLink', 'RouterView'],
        },
      ],
    }),
    VitePWA({
      devOptions: {
        enabled: true,
      },
      includeAssets: [
        'favicon.ico',
        'img/icons/*',
      ],
      manifest: {
        name: 'Relax Sounds',
        short_name: 'Relax Sounds',
        id: '/',
        description: 'Stream relaxing sound mixes',
        theme_color: '#673AB7',
        background_color: '#673AB7',
        icons: [
          {
            src: 'img/icons/android-chrome-192x192.png',
            sizes: '192x192',
            type: 'image/png',
          },
          {
            src: 'img/icons/android-chrome-512x512.png',
            sizes: '512x512',
            type: 'image/png',
          },
          {
            src: 'img/icons/android-chrome-maskable-192x192.png',
            sizes: '192x192',
            type: 'image/png',
            purpose: 'maskable',
          },
          {
            src: 'img/icons/android-chrome-maskable-512x512.png',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'maskable',
          },
        ],
      },
      workbox: {
        globPatterns: ['**/*{js,css,html,woff2}'],
        runtimeCaching: [
          {
            urlPattern: /\/api\//,
            handler: 'NetworkFirst',
            options: {
              cacheName: 'api-cache',
              cacheableResponse: {
                statuses: [0, 200],
              },
            },
          },
          {
            urlPattern: /\/data\/audio/,
            handler: 'CacheFirst',
            options: {
              cacheName: 'data-cache',
              expiration: {
                maxEntries: 10,
                maxAgeSeconds: 60 * 60 * 24 * 31,
              },
              cacheableResponse: {
                statuses: [0, 200],
              },
            },
          },
        ],
      },
    }),
  ],
});
