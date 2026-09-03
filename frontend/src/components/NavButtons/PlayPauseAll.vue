<template>
  <v-tooltip :text="player.isPlaying ? 'Pause' : 'Play'" location="top">
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        icon
        :size="size"
        :disabled="player.isStopped"
        :aria-label="player.isPlaying ? 'Pause All' : 'Play All'"
        variant="flat"
        color="primary"
        @click="player.playPauseAll"
      >
        <v-icon :icon="player.isPlaying ? PauseIcon : PlayIcon" size="28" />
      </v-btn>
    </template>
  </v-tooltip>
</template>

<script setup>
import PauseIcon from "~icons/material-symbols/pause-rounded";
import PlayIcon from "~icons/material-symbols/play-arrow-rounded";
import { usePlayer } from "@/plugins/store/player";

defineProps({
  size: {
    type: String,
    default: "large",
  },
});

const player = usePlayer();
</script>
