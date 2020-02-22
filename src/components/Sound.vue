<template>
  <v-card raised
          outlined
  >
    <v-progress-linear
      v-model="volumePercentage"
      absolute
      height="100%"
      color="secondary"
      v-if="sound.state !== 'stopped'"
    />
    <v-row align="center" justify="center" dense>
      <v-col class="grow">
        <v-card-title class="headline">
          <v-icon class="pr-2">{{ sound.icon }}</v-icon>
          {{ sound.name }}
        </v-card-title>
      </v-col>
      <v-col class="shrink pr-4">
        <v-btn @click.stop="playPause" elevation="0" outlined :loading="sound.loading">
          <v-icon>{{ icon }}</v-icon>
        </v-btn>
      </v-col>
    </v-row>
  </v-card>
</template>

<script>
export default {
  name: 'Sound',

  props: {
    sound: {
      type: Object,
    },
  },

  computed: {
    volumePercentage: {
      get() {
        return this.sound.volume * 100;
      },
      set(newValue) {
        this.$store.commit('sounds/volume', { id: this.sound.id, value: newValue / 100 });
      },
    },
    icon() {
      return this.sound.state === 'playing' ? 'mdi-stop' : 'mdi-play';
    },
  },

  methods: {
    async playPause() {
      this.$store.commit('sounds/playPause', { id: this.sound.id });
    },
  },
};
</script>

<style scoped>
  .v-card {
    overflow: hidden;
  }
  .v-progress-linear {
    z-index: 0;
    left: 0;
    right: 0;
    height: 100%;
    width: 100%;
  }
  .row {
    position: relative;
    z-index: 1;
    pointer-events: none;
  }
  .v-btn {
    pointer-events: all;
  }
</style>
