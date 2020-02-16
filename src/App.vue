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
    </v-app-bar>

    <v-content>
      <SoundsPage :sounds="sounds"/>
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

  created() {
    this.$vuetify.theme.dark = true;
    this.sounds = this.sounds.sort(
      (left, right) => left.name.localeCompare(right.name),
    );
  },

  data: () => ({
    sounds,
  }),
};
</script>
