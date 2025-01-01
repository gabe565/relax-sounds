<template>
  <component
    :is="component"
    :prepend-icon="player.isPlaying ? PauseIcon : PlayIcon"
    :disabled="player.isStopped"
    :aria-label="player.isPlaying ? 'Pause All' : 'Play All'"
    @click="player.playPauseAll"
  >
    {{ player.isPlaying ? "Pause" : "Play" }}
  </component>
</template>

<script setup>
import { VBtn, VListItem } from "vuetify/components";
import PauseIcon from "~icons/material-symbols/pause-circle-rounded";
import PlayIcon from "~icons/material-symbols/play-circle-rounded";
import { usePlayerStore } from "@/plugins/store/player";

const props = defineProps({
  button: {
    type: Boolean,
    default: false,
  },
});

const player = usePlayerStore();
const component = props.button ? VBtn : VListItem;
</script>
