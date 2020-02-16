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
      <v-btn icon @click="playPauseAll" :disabled="audioState === 'stopped'">
        <v-icon>{{ audioState === 'playing' ? 'mdi-pause-circle' : 'mdi-play-circle' }}</v-icon>
      </v-btn>
    </v-app-bar>

    <v-content>
      <SoundsPage/>
    </v-content>

  </v-app>
</template>

<script>
import SoundsPage from './components/SoundsPage.vue';

export default {
  name: 'App',

  components: {
    SoundsPage,
  },

  created() {
    this.$vuetify.theme.dark = true;
  },

  methods: {
    playPauseAll() {
      this.$store.commit('playPauseAll');
    },
  },

  computed: {
    audioState() {
      return this.$store.getters.state;
    },
  },
};
</script>
