import vue from "@vitejs/plugin-vue";
import vuetify from "vite-plugin-vuetify";
import { defineConfig } from "vite";
import { VitePWA } from "vite-plugin-pwa";

export default defineConfig({
  plugins: [
    vue(),
    vuetify({
      styles: {
        configFile: "src/scss/variables.scss",
      },
    }),
    VitePWA({
      devOptions: {
        enabled: true,
      },
      includeAssets: ["favicon.ico", "img/icons/*"],
      manifest: {
        name: "Relax Sounds",
        short_name: "Relax Sounds",
        id: "/",
        description: "Stream relaxing sound mixes",
        theme_color: "#673AB7",
        background_color: "#673AB7",
        icons: [
          {
            src: "img/icons/android-chrome-192x192.png",
            sizes: "192x192",
            type: "image/png",
          },
          {
            src: "img/icons/android-chrome-512x512.png",
            sizes: "512x512",
            type: "image/png",
          },
          {
            src: "img/icons/android-chrome-maskable-192x192.png",
            sizes: "192x192",
            type: "image/png",
            purpose: "maskable",
          },
          {
            src: "img/icons/android-chrome-maskable-512x512.png",
            sizes: "512x512",
            type: "image/png",
            purpose: "maskable",
          },
        ],
      },
      workbox: {
        clientsClaim: true,
        globPatterns: ["**/*{js,css,html,woff2,svg}"],
        navigateFallbackDenylist: [/^\/api\//, /^\/data\//],
        runtimeCaching: [
          {
            urlPattern: /\/api\/(?!mix\/)/,
            handler: "NetworkFirst",
            options: {
              cacheName: "api-cache",
              cacheableResponse: {
                statuses: [0, 200],
              },
            },
          },
          {
            urlPattern: /\/data\/audio/,
            handler: "CacheFirst",
            options: {
              cacheName: "data-cache",
              expiration: {
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
