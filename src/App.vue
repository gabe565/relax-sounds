<template>
  <v-app>

    <v-app-bar
      app
      color="primary"
      dark
      hide-on-scroll
    >
      <v-app-bar-nav-icon>
        <v-icon>fal fa-bars</v-icon>
      </v-app-bar-nav-icon>
      <v-toolbar-title>Relax Sounds</v-toolbar-title>
      <v-spacer/>

      <GlobalPlayPause/>
    </v-app-bar>

    <UpdateSnackbar/>

    <v-content>
      <v-container>

        <Filters/>

        <v-row>
          <v-col
            cols="12"
            sm="6"
            md="4"
            v-for="sound of $store.getters['filters/sounds']"
            :key="sound.id"
          >
            <Sound :sound="sound"/>
          </v-col>
        </v-row>
      </v-container>
    </v-content>

  </v-app>
</template>

<script>
import GlobalPlayPause from './components/GlobalPlayPause.vue';
import UpdateSnackbar from './components/UpdateSnackbar.vue';
import Filters from './components/Filters.vue';
import Sound from './components/Sound.vue';

export default {
  name: 'App',

  components: {
    Sound,
    Filters,
    GlobalPlayPause,
  },

  created() {
    this.$vuetify.theme.dark = true;
  },
};
</script>

<style>
  .fa-spin-2x {
    animation: fa-spin 750ms infinite linear;
  }
</style>
