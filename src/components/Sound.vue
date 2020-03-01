<template>
  <v-card flat
          outlined
  >
    <v-progress-linear
      v-model="volumePercentage"
      absolute
      height="100%"
      color="deep-purple darken-2"
      v-if="sound.state !== 'stopped'"
    />
    <v-row align="center" justify="center" dense>
      <v-col class="grow">
        <v-card-title class="headline">
          <v-icon aria-hidden="true" class="mr-4" :color="iconColor">
            {{ iconStyle }} {{ sound.icon }} fa-fw no-transition
          </v-icon>
          {{ sound.name }}
        </v-card-title>
      </v-col>
      <v-col class="shrink pr-4">
        <v-btn @click.stop="playStop" elevation="0" outlined>
          <v-icon dense>
            fas fa-fw {{ icon }}
          </v-icon>
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
        this.$store.commit('sounds/volume', { sound: this.sound, value: newValue / 100 });
      },
    },

    iconStyle() {
      return this.sound.state === 'stopped' ? 'fal' : 'fas';
    },

    iconColor() {
      return this.sound.state !== 'stopped' ? 'primary' : '';
    },

    icon() {
      if (this.sound.loading) {
        return 'fa-spinner-third fa-spin-2x';
      }
      if (this.sound.state === 'playing') {
        return 'fa-stop';
      }
      return 'fa-play';
    },
  },

  methods: {
    async playStop() {
      return this.$store.dispatch('sounds/playStop', { sound: this.sound });
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
  .no-transition {
    transition: none;
  }
</style>
