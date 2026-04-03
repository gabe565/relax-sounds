<template>
  <v-app-bar
    v-if="!isMobile"
    location="bottom"
    height="64"
    class="px-3 border-t! border-on-surface/15"
    elevation="10"
    theme="dark"
    color="surface"
    order="-1"
    flat
  >
    <template #prepend>
      <div class="flex items-center">
        <nav-size-btn />
        <theme-btn />
      </div>
    </template>

    <div
      class="absolute left-1/2 -translate-x-1/2 whitespace-nowrap flex items-center justify-center"
    >
      <save-preset button />
      <play-pause-all />
      <stop-all />
    </div>

    <template #append>
      <div class="flex items-center mr-1">
        <cast-icon />
      </div>
    </template>
  </v-app-bar>

  <div v-else class="fixed bottom-14 left-0 z-1000 pointer-events-none px-2 pb-2 w-full">
    <v-toolbar
      height="56"
      class="px-2 rounded-lg border border-on-surface/20 bg-surface/40 backdrop-blur-[10px] backdrop-saturate-200 transition-colors duration-300 pointer-events-auto"
      elevation="6"
    >
      <div
        class="absolute left-1/2 -translate-x-1/2 whitespace-nowrap flex items-center justify-center size-full"
      >
        <save-preset button />
        <play-pause-all />
        <stop-all />
      </div>

      <div class="absolute right-1 flex items-center">
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
