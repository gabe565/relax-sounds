<template>
  <v-btn
    elevation="0"
    icon
    color="transparent"
    :aria-label="sound.isPlaying ? 'Stop' : 'Play'"
    :loading="sound.isLoading"
    @click.stop="playStop"
  >
    <v-icon v-if="sound.isPlaying" :icon="PauseIcon" aria-hidden="true" size="large" />
    <v-icon v-else :icon="PlayIcon" aria-hidden="true" size="x-large" />
  </v-btn>
</template>

<script setup>
import PlayIcon from "~icons/material-symbols/play-arrow-rounded";
import PauseIcon from "~icons/material-symbols/pause-rounded";
import { usePlayerStore } from "../../../plugins/store/player";

const props = defineProps({
  sound: {
    type: Object,
    required: true,
  },
});

const playStop = async () => {
  return usePlayerStore().playStop({ sound: props.sound });
};
</script>
