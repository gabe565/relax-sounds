<template>
  <v-app>

    <v-navigation-drawer app bottom temporary v-model="drawerOpen">
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

    <v-app-bar app dark collapse-on-scroll color="accent">
      <v-app-bar-nav-icon @click.stop="toggleDrawer"/>
      <v-toolbar-title>
        <v-icon class="mr-2">fas fa-bed-alt</v-icon>
        Relax Sounds
      </v-toolbar-title>
      <v-spacer/>
      <GlobalPlayPause/>
    </v-app-bar>

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
