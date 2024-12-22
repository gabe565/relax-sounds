<template>
  <v-app>
    <v-navigation-drawer
      :rail="preferences.shrinkLeftPanel"
      color="accent"
      width="200"
      mobile-breakpoint="md"
    >
      <template #prepend>
        <v-list v-if="!isMobile">
          <v-list-item to="/" title="Relax Sounds" :prepend-icon="AppIcon" />
        </v-list>
        <v-divider />
      </template>

      <v-list nav>
        <v-list-item
          v-for="route in routes"
          :key="route.path"
          :to="route.path"
          exact
          link
          :title="route.name"
          :prepend-icon="route.meta.icon"
        />
      </v-list>

      <v-divider />

      <v-list>
        <play-pause-all />
        <stop-all />
        <save-preset />
        <cast-icon />
        <debug-button v-if="DebugEnabled" list-item />
      </v-list>

      <template #append>
        <v-divider />
        <div class="d-flex overflow-hidden">
          <nav-size-btn />
          <theme-btn />
        </div>
      </template>
    </v-navigation-drawer>

    <v-main>
      <router-view v-slot="{ Component }">
        <keep-alive>
          <component :is="Component" />
        </keep-alive>
      </router-view>
      <v-spacer style="height: 28px" />
    </v-main>

    <v-bottom-navigation
      v-if="isMobile"
      v-model="botnav"
      bg-color="accent"
      color="primary"
      theme="dark"
      :model-value="null"
    >
      <v-btn v-for="route in routes" :key="route.path" :to="route.path" :value="route.name">
        <v-icon v-if="route.meta.icon" :icon="route.meta.icon" />
        <span>{{ route.name }}</span>
      </v-btn>

      <v-divider vertical />
      <play-pause-all button />
      <stop-all button />
    </v-bottom-navigation>
  </v-app>
</template>

<script setup>
import { computed, watch } from "vue";
import { useRouter } from "vue-router";
import { useDisplay, useTheme } from "vuetify";
import CastIcon from "./components/NavButtons/CastIcon.vue";
import NavSizeBtn from "./components/NavButtons/NavSizeBtn.vue";
import PlayPauseAll from "./components/NavButtons/PlayPauseAll.vue";
import SavePreset from "./components/NavButtons/SavePreset.vue";
import StopAll from "./components/NavButtons/StopAll.vue";
import ThemeBtn from "./components/NavButtons/ThemeBtn.vue";
import DebugButton from "./components/Presets/Buttons/DebugButton.vue";
import { DebugEnabled } from "./config/debug";
import { registerSW } from "./plugins/pwa";
import { Theme, usePreferencesStore } from "./plugins/store/preferences";
import AppIcon from "~icons/relax-sounds/icon-white";

const { smAndDown: isMobile } = useDisplay();
const preferences = usePreferencesStore();
const theme = useTheme();

const routes = computed(() => {
  return useRouter().options.routes.filter((route) => route.meta?.showInNav);
});

const autoTheme = (e) => (theme.name.value = e.matches ? "dark" : "light");
const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
watch(
  () => preferences.theme,
  (val) => {
    mediaQuery.removeEventListener("change", autoTheme);
    switch (val) {
      case Theme.dark:
        theme.name.value = "dark";
        break;
      case Theme.light:
        theme.name.value = "light";
        break;
      default:
        // check for browser support
        if (window.matchMedia && window.matchMedia("(prefers-color-scheme)").media !== "not all") {
          // set to preferred scheme
          autoTheme(mediaQuery);
          // react to changes
          mediaQuery.addEventListener("change", autoTheme);
        }
    }
  },
  { immediate: true },
);

registerSW();
</script>
