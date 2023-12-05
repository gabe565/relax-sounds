<template>
  <v-app>
    <v-app-bar v-if="isMobile" theme="dark" color="accent" flat>
      <v-btn to="/" class="text-body-2 text-none px-2">
        <template #prepend>
          <v-icon :icon="AppIcon" aria-hidden="true" size="28" />
        </template>
        <v-app-bar-title>Relax Sounds</v-app-bar-title>
      </v-btn>

      <v-spacer />
      <save-preset button />
      <cast-icon button />
    </v-app-bar>

    <v-navigation-drawer
      v-else
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
        <debug-button v-if="debugEnabled" list-item />
      </v-list>

      <template #append>
        <v-divider />
        <v-btn
          :icon="preferences.shrinkLeftPanel ? LeftPanelOpenIcon : LeftPanelCloseIcon"
          color="transparent"
          variant="flat"
          :aria-label="preferences.shrinkLeftPanel ? 'Expand Left Panel' : 'Shrink Left Panel'"
          @click="preferences.shrinkLeftPanel = !preferences.shrinkLeftPanel"
        />
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
import { useDisplay, useTheme } from "vuetify";
import { computed, onBeforeMount } from "vue";
import { useRouter } from "vue-router";
import SavePreset from "./components/NavButtons/SavePreset.vue";
import PlayPauseAll from "./components/NavButtons/PlayPauseAll.vue";
import StopAll from "./components/NavButtons/StopAll.vue";
import AppIcon from "~icons/relax-sounds/icon-white";
import CastIcon from "./components/NavButtons/CastIcon.vue";
import LeftPanelCloseIcon from "~icons/material-symbols/left-panel-close-rounded";
import LeftPanelOpenIcon from "~icons/material-symbols/left-panel-open-rounded";
import { usePreferencesStore } from "./plugins/store/preferences";
import DebugButton from "./components/Presets/Buttons/DebugButton.vue";
import { registerSW } from "./plugins/pwa";

const { smAndDown: isMobile } = useDisplay();
const preferences = usePreferencesStore();
const debugEnabled = import.meta.env.DEV;

const routes = computed(() => {
  return useRouter().options.routes.filter((route) => route.meta?.showInNav);
});

onBeforeMount(() => {
  // check for browser support
  if (window.matchMedia && window.matchMedia("(prefers-color-scheme)").media !== "not all") {
    // set to preferred scheme
    const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
    const theme = useTheme();
    theme.name.value = mediaQuery.matches ? "dark" : "light";
    // react to changes
    mediaQuery.addEventListener("change", (e) => {
      theme.name.value = e.matches ? "dark" : "light";
    });
  }
});

registerSW();
</script>
