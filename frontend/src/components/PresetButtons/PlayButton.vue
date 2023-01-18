<template>
  <v-col class="flex-grow-0 pr-4">
    <v-btn
      elevation="0"
      icon
      variant="plain"
      aria-label="Play"
      @click.stop="play"
    >
      <v-icon aria-hidden="true">
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
