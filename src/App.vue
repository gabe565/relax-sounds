<template>
  <v-app>

    <v-app-bar
      app
      color="primary"
      dark
      hide-on-scroll
    >
      <v-app-bar-nav-icon></v-app-bar-nav-icon>
      <v-toolbar-title>Relax Sounds</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-btn icon @click="playPauseAll" :disabled="state === 'stopped'">
        <v-icon>{{ state === 'playing' ? 'mdi-pause-circle' : 'mdi-play-circle' }}</v-icon>
      </v-btn>
    </v-app-bar>

    <v-content>
      <SoundsPage :sounds="sounds" ref="SoundsPage"/>
    </v-content>

  </v-app>
</template>

<script>
import sounds from './assets/sounds.json';
import SoundsPage from './components/SoundsPage.vue';

export default {
  name: 'App',

  components: {
    SoundsPage,
  },

  data: () => ({
    sounds,
  }),

  created() {
    this.$vuetify.theme.dark = true;
    this.sounds = this.sounds.map((sound) => ({
      ...sound,
      state: 'stopped',
      volume: 1,
    })).sort(
      (left, right) => left.name.localeCompare(right.name),
    );
  },

  computed: {
    state() {
      const states = new Set(this.sounds.map((sound) => sound.state));
      if (states.has('playing')) {
        return 'playing';
      }
      if (states.has('paused')) {
        return 'paused';
      }
      return 'stopped';
    },
  },

  methods: {
    playPauseAll() {
      const newState = this.state === 'playing' ? 'paused' : 'playing';
      this.sounds.filter(
        (sound) => sound.state !== 'stopped',
      ).forEach((sound) => {
        sound.state = newState;
        if (newState === 'paused') {
          sound.player.pause();
        } else if (newState === 'playing') {
          sound.player.play();
        }
      });
    },
  },
};
</script>
