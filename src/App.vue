<template>
  <v-app>

    <v-navigation-drawer app bottom temporary v-model="drawerOpen">
      <v-list nav>
        <v-list-item-group>
          <v-list-item link @click="prefetch">
            <v-list-item-avatar>
              <v-icon aria-hidden="true">fal fa-download</v-icon>
            </v-list-item-avatar>
            <v-list-item-title>
              Preload All
            </v-list-item-title>
          </v-list-item>
        </v-list-item-group>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar app dark hide-on-scroll color="accent">
      <v-app-bar-nav-icon @click.stop="toggleDrawer"/>
      <v-toolbar-title>
        <v-icon aria-hidden="true" class="mr-2">fas fa-bed-alt</v-icon>
        Relax Sounds
      </v-toolbar-title>

      <v-spacer/>

      <SavePlaylist/>
      <PlayPauseAll/>
      <StopAll/>

      <template v-slot:extension>
        <v-tabs centered>
          <v-tab v-for="route in routes" :key="route.path" :to="route.path" exact>
            <v-icon class="pr-2">fas {{ route.meta.icon }} fa-fw</v-icon>
            {{ route.name }}
          </v-tab>
        </v-tabs>
      </template>
    </v-app-bar>

    <UpdateSnackbar/>

    <v-content>
      <keep-alive>
        <router-view/>
      </keep-alive>
    </v-content>

  </v-app>
</template>

<script>
import SavePlaylist from './components/SavePlaylist.vue';
import PlayPauseAll from './components/PlayPauseAll.vue';
import StopAll from './components/StopAll.vue';
import UpdateSnackbar from './components/UpdateSnackbar.vue';

export default {
  name: 'App',

  components: {
    SavePlaylist,
    StopAll,
    PlayPauseAll,
    UpdateSnackbar,
  },

  data: () => ({
    drawerOpen: false,
  }),

  computed: {
    routes() {
      return this.$router.options.routes.filter((route) => !!route.name);
    },
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
