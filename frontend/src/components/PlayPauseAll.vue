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
        <v-icon v-if="player.isPlaying" :icon="PauseIcon" aria-hidden="true" />
        <v-icon v-else :icon="PlayIcon" aria-hidden="true" />
      </v-btn>
    </template>
    <span v-if="player.isPlaying">Pause</span>
    <span v-else>Play</span>
  </v-tooltip>
</template>

<script setup>
import PlayIcon from "~icons/solar/play-circle-bold";
import PauseIcon from "~icons/solar/pause-circle-bold";
import { usePlayerStore } from "../plugins/store/player";

const player = usePlayerStore();
</script>
