<template>
  <component
    :is="component"
    :prepend-icon="player.isPlaying ? PauseIcon : PlayIcon"
    :disabled="player.isStopped"
    :aria-label="player.isPlaying ? 'Pause All' : 'Play All'"
    @click="player.playPauseAll"
  >
    <span v-if="player.isPlaying">Pause</span>
    <span v-else>Play</span>
  </component>
</template>

<script setup>
import PlayIcon from "~icons/material-symbols/play-circle-rounded";
import PauseIcon from "~icons/material-symbols/pause-circle-rounded";
import { usePlayerStore } from "../../plugins/store/player";
import { VBtn, VListItem } from "vuetify/components";

const props = defineProps({
  button: {
    type: Boolean,
    default: false,
  },
});

const player = usePlayerStore();
const component = props.button ? VBtn : VListItem;
</script>
