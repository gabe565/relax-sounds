<template>
  <v-tooltip :text="player.isPlaying ? 'Pause' : 'Play'" location="top">
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        icon
        size="large"
        :disabled="player.isStopped"
        :aria-label="player.isPlaying ? 'Pause All' : 'Play All'"
        variant="text"
        @click="player.playPauseAll"
      >
        <v-icon :icon="player.isPlaying ? PauseIcon : PlayIcon" size="36" />
      </v-btn>
    </template>
  </v-tooltip>
</template>

<script setup>
import PauseIcon from "~icons/material-symbols/pause-circle-sharp";
import PlayIcon from "~icons/material-symbols/play-circle-sharp";
import { usePlayer } from "@/plugins/store/player";

const player = usePlayer();
</script>
