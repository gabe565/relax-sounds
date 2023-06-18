<template>
  <v-tooltip
    :key="player.isPlaying"
    location="bottom"
    :aria-label="player.isPlaying ? 'Pause' : 'Play'"
  >
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        :disabled="player.isStopped"
        icon
        :aria-label="player.isPlaying ? 'Pause' : 'Play'"
        @click="player.playPauseAll"
      >
        <v-icon v-if="player.isPlaying" aria-hidden="true">fas fa-pause-circle</v-icon>
        <v-icon v-else aria-hidden="true">fas fa-play-circle</v-icon>
      </v-btn>
    </template>
    <span v-if="player.isPlaying">Pause</span>
    <span v-else>Play</span>
  </v-tooltip>
</template>

<script setup>
import { usePlayerStore } from "../plugins/store/player";

const player = usePlayerStore();
</script>
