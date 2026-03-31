<template>
  <v-app-bar
    v-if="!isMobile"
    location="bottom"
    height="64"
    class="px-3 player-bar-desktop"
    elevation="10"
    color="surface"
    order="-1"
    flat
  >
    <template #prepend>
      <div class="d-flex align-center">
        <nav-size-btn />
        <theme-btn />
      </div>
    </template>

    <div class="player-bar-desktop-content d-flex align-center justify-center">
      <save-preset button />
      <play-pause-all />
      <stop-all />
    </div>

    <template #append>
      <div class="d-flex align-center mr-1">
        <cast-icon />
      </div>
    </template>
  </v-app-bar>

  <div v-else class="player-bar-mobile-container px-2 pb-2 w-100">
    <v-toolbar height="56" class="px-2 rounded-lg border player-bar-glass" elevation="6">
      <div class="player-bar-mobile-content d-flex align-center justify-center w-100 h-100">
        <save-preset button />
        <play-pause-all />
        <stop-all />
      </div>

      <div class="player-bar-mobile-right d-flex align-center">
        <cast-icon />
      </div>
    </v-toolbar>
  </div>
</template>

<script setup>
import { useDisplay } from "vuetify";
import CastIcon from "@/components/NavButtons/CastIcon.vue";
import NavSizeBtn from "@/components/NavButtons/NavSizeBtn.vue";
import PlayPauseAll from "@/components/NavButtons/PlayPauseAll.vue";
import SavePreset from "@/components/NavButtons/SavePreset.vue";
import StopAll from "@/components/NavButtons/StopAll.vue";
import ThemeBtn from "@/components/NavButtons/ThemeBtn.vue";

const { smAndDown: isMobile } = useDisplay();
</script>

<style scoped>
.v-toolbar {
  transition: none;
}

.player-bar-desktop-content,
.player-bar-mobile-content {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  white-space: nowrap;
}

.player-bar-mobile-container {
  position: fixed;
  bottom: 56px;
  left: 0;
  z-index: 1000;
  pointer-events: none;
}

.player-bar-mobile-container > .v-toolbar {
  pointer-events: auto;
}

.player-bar-mobile-right {
  position: absolute;
  right: 4px;
}

.player-bar-desktop {
  border-top: 1px solid rgba(var(--v-border-color), 0.15) !important;
}

.player-bar-glass {
  background: rgba(var(--v-theme-surface), 0.6) !important;
  backdrop-filter: blur(12px) saturate(180%);
  -webkit-backdrop-filter: blur(12px) saturate(180%);
  border: 1px solid rgba(var(--v-border-color), 0.2) !important;
  transition: background-color 0.3s ease;
}
</style>
