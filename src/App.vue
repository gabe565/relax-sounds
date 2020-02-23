<template>
  <v-app>

    <v-app-bar app color="primary" dark hide-on-scroll>
      <v-app-bar-nav-icon @click.stop="toggleDrawer"/>
      <v-toolbar-title>Relax Sounds</v-toolbar-title>
      <v-spacer/>
      <GlobalPlayPause/>
    </v-app-bar>


    <v-navigation-drawer v-model="drawerOpen" absolute temporary>
      <v-list nav>
        <v-list-item-group>
          <v-list-item link @click="prefetch">
            <v-list-item-avatar>
              <v-icon>fal fa-download</v-icon>
            </v-list-item-avatar>
            <v-list-item-title>
              Preload All
            </v-list-item-title>
          </v-list-item>
        </v-list-item-group>
      </v-list>
    </v-navigation-drawer>

    <UpdateSnackbar/>

    <v-content>
      <v-container>
        <Filters/>
        <v-row>
          <v-col
            cols="12" sm="6" md="4"
            v-for="sound of filteredSounds"
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
import { mapGetters } from 'vuex';
import GlobalPlayPause from './components/GlobalPlayPause.vue';
import UpdateSnackbar from './components/UpdateSnackbar.vue';
import Filters from './components/Filters.vue';
import Sound from './components/Sound.vue';

export default {
  name: 'App',

  components: {
    GlobalPlayPause,
    UpdateSnackbar,
    Sound,
    Filters,
  },

  data: () => ({
    drawerOpen: false,
  }),

  computed: mapGetters('filters', {
    filteredSounds: 'sounds',
  }),

  created() {
    this.$vuetify.theme.dark = true;
  },

  methods: {
    toggleDrawer() {
      this.drawerOpen = !this.drawerOpen;
    },

    prefetch() {
      this.toggleDrawer();
      this.$store.dispatch('sounds/prefetch');
    },
  },
};
</script>

<style>
  .fa-spin-2x {
    animation: fa-spin 750ms infinite linear;
  }
</style>
