<template>
  <v-app-bar
    v-if="!isMobile"
    location="bottom"
    height="64"
    class="px-3 border-t! border-outline-variant"
    color="surface-container"
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
      class="absolute left-1/2 flex -translate-x-1/2 items-center justify-center gap-2 whitespace-nowrap"
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
      class="px-2 rounded-lg border border-outline-variant bg-surface-container/70 backdrop-blur-[10px] backdrop-saturate-200 transition-colors duration-300 pointer-events-auto"
      elevation="6"
    >
      <div
        class="absolute left-1/2 flex size-full -translate-x-1/2 items-center justify-center gap-2 whitespace-nowrap"
      >
        <save-preset button />
        <play-pause-all size="default" />
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
