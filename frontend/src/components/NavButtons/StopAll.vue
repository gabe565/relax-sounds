<template>
  <v-tooltip text="Stop" location="top">
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        icon
        :disabled="player.isStopped"
        aria-label="Stop"
        variant="text"
        @click="player.stopAll({})"
      >
        <v-icon :icon="StopIcon" />
      </v-btn>
    </template>
  </v-tooltip>
</template>

<script setup>
import StopIcon from "~icons/material-symbols/stop-sharp";
import { usePlayer } from "@/plugins/store/player";

const player = usePlayer();
</script>
