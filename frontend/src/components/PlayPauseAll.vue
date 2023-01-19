<template>
  <v-tooltip :key="isPlaying" location="bottom" :aria-label="isPlaying ? 'Pause' : 'Play'">
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        :disabled="isStopped"
        icon
        :aria-label="isPlaying ? 'Pause' : 'Play'"
        @click="playPauseAll"
      >
        <v-icon v-if="isPlaying" aria-hidden="true">fas fa-pause-circle</v-icon>
        <v-icon v-else aria-hidden="true">fas fa-play-circle</v-icon>
      </v-btn>
    </template>
    <span v-if="isPlaying">Pause</span>
    <span v-else>Play</span>
  </v-tooltip>
</template>

<script>
import { mapActions, mapGetters } from "vuex";

export default {
  name: "PlayPauseAll",

  computed: mapGetters("player", ["isPlaying", "isStopped"]),

  methods: mapActions("player", ["playPauseAll"]),
};
</script>
