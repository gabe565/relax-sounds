<template>
  <v-tooltip bottom>
    <template #activator="{ on, attrs }">
      <v-btn
        v-bind="attrs"
        v-on="on"
        @click="playPauseAll"
        :disabled="state === 'stopped'"
        icon
      >
        <v-icon>fas {{ state === 'playing' ? 'fa-pause-circle' : 'fa-play-circle' }}</v-icon>
      </v-btn>
    </template>
    <span>
      {{ state === 'playing' ? 'Pause' : 'Play' }}
    </span>
  </v-tooltip>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';

export default {
  name: 'PlayPauseAll',

  computed: mapGetters('player', [
    'state',
  ]),

  methods: mapActions('player', [
    'playPauseAll',
  ]),
};
</script>

<style scoped>

</style>
