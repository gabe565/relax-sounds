<template>
  <v-btn icon @click="playPauseAll" :disabled="state === 'stopped'">
    <v-icon>fal {{ state === 'playing' ? 'fa-pause-circle' : 'fa-play-circle' }}</v-icon>
  </v-btn>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';

export default {
  name: 'GlobalPlayPause',

  methods: mapActions('sounds', [
    'playPauseAll',
  ]),

  computed: mapGetters('sounds', [
    'state',
  ]),
};
</script>

<style scoped>

</style>
