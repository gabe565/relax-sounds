<template>
  <v-btn icon @click="playPauseAll" :disabled="state === 'stopped'">
    <v-icon>fas {{ state === 'playing' ? 'fa-pause-circle' : 'fa-play-circle' }}</v-icon>
  </v-btn>
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
