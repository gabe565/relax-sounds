<template>
  <v-btn elevation="0" icon variant="plain" aria-label="Play" @click.stop="play">
    <v-icon :icon="PlayButton" aria-hidden="true" size="x-large" />
  </v-btn>
</template>

<script setup>
import PlayButton from "~icons/material-symbols/play-arrow-rounded";
import { usePresetsStore } from "../../../plugins/store/presets";

const props = defineProps({
  preset: {
    type: Object,
    required: true,
  },
});

const play = () => usePresetsStore().play({ preset: props.preset });
</script>
