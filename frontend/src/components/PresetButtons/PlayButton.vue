<template>
  <v-col class="shrink pr-4">
    <v-btn @click.stop="play" elevation="0" icon aria-label="Play">
      <v-icon dense aria-hidden="true">
        fas fa-fw fa-play
      </v-icon>
    </v-btn>
  </v-col>
</template>

<script>
export default {
  name: 'PlayButton',

  props: {
    preset: {
      type: Object,
      required: true,
    },
  },

  methods: {
    play() {
      this.$store.dispatch('presets/play', { preset: this.preset });
    },
  },
};
</script>
