<template>
  <v-tooltip bottom :key="isPlaying">
    <template #activator="{ on, attrs }">
      <v-btn
        v-bind="attrs"
        v-on="on"
        @click="playPauseAll"
        :disabled="isStopped"
        icon
        :aria-label="isPlaying ? 'Pause' : 'Play'"
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
import { mapActions, mapGetters } from 'vuex';

export default {
  name: 'PlayPauseAll',

  computed: mapGetters('player', [
    'isPlaying', 'isStopped',
  ]),

  methods: mapActions('player', [
    'playPauseAll',
  ]),
};
</script>
