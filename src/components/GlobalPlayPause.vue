<template>
  <v-btn icon @click="playPauseAll" :disabled="audioState === 'stopped'">
    <v-icon>fal {{ audioState === 'playing' ? 'fa-pause-circle' : 'fa-play-circle' }}</v-icon>
  </v-btn>
</template>

<script>
export default {
  name: 'GlobalPlayPause',

  methods: {
    playPauseAll() {
      this.$store.commit('sounds/playPauseAll');
    },
  },

  computed: {
    audioState() {
      return this.$store.getters['sounds/state'];
    },
  },
};
</script>

<style scoped>

</style>
