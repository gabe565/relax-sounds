<template>
  <v-tooltip bottom :key="isPlaying">
    <template #activator="{ on, attrs }">
      <v-btn
        v-bind="attrs"
        v-on="on"
        @click="playPauseAll"
        :disabled="isStopped"
        icon
      >
        <v-icon>fas {{ isPlaying ? 'fa-pause-circle' : 'fa-play-circle' }}</v-icon>
      </v-btn>
    </template>
    <span>
      {{ isPlaying ? 'Pause' : 'Play' }}
    </span>
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
