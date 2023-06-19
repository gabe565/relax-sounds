<template>
  <v-col class="flex-grow-0 pr-4">
    <v-btn elevation="0" icon variant="plain" aria-label="Play" @click.stop="play">
      <v-icon :icon="PlayButton" aria-hidden="true" />
    </v-btn>
  </v-col>
</template>

<script setup>
import PlayButton from "~icons/material-symbols/play-arrow-rounded";
import { usePresetsStore } from "../../plugins/store/presets";

const props = defineProps({
  preset: {
    type: Object,
    required: true,
  },
});

const play = () => usePresetsStore().play({ preset: props.preset });
</script>
