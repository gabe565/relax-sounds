<template>
  <v-btn icon @click="playPauseAll" :disabled="state === 'stopped'">
    <v-icon>fal {{ state === 'playing' ? 'fa-pause-circle' : 'fa-play-circle' }}</v-icon>
  </v-btn>
</template>

<script>
import { mapGetters, mapMutations } from 'vuex';

export default {
  name: 'GlobalPlayPause',

  methods: mapMutations('sounds', [
    'playPauseAll',
  ]),

  computed: mapGetters('sounds', [
    'state',
  ]),
};
</script>

<style scoped>

</style>
