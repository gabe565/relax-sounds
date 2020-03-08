<template>
  <v-app>

    <v-navigation-drawer app bottom temporary v-model="drawerOpen">
      <v-list nav>
        <v-list-item-group>
          <v-list-item link @click="prefetch">
            <v-list-item-avatar>
              <v-icon aria-hidden="true">fas fa-download</v-icon>
            </v-list-item-avatar>
            <v-list-item-title>
              Preload All
            </v-list-item-title>
          </v-list-item>
        </v-list-item-group>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar app dark hide-on-scroll color="accent" flat>
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

<style lang="scss">
  .theme--light {
    &.v-pagination {
      .v-pagination__navigation, .v-pagination__item {
        border: thin solid rgba(0, 0, 0, 0.12);
      }
      .v-pagination__item--active {
        color: rgba(0, 0, 0, 0.87);
      }
    }
  }

  .theme--dark {
    &.v-application {
      background: #150b29;

      & > .v-application--wrap {
        .v-card {
          background: transparent;
        }
      }
    }
    &.v-pagination {
      .v-pagination__navigation, .v-pagination__item {
        border: thin solid rgba(255, 255, 255, 0.12);
        background: transparent !important;
      }
    }
  }

  .v-application {
    .v-pagination {
      .v-pagination__item--active {
        background-color: rgba(255, 255, 255, 0.08) !important;
      }
    }
  }

  .fa-spin-2x {
    animation: fa-spin 750ms infinite linear;
  }

  .v-pagination__navigation, .v-pagination__item {
    box-shadow: none;
    outline: none;
    transition: 0.2s cubic-bezier(0, 0.5, 0.2, 1);
  }
</style>
