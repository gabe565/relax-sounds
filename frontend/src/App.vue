<template>
  <v-app>
    <v-app-bar
      app
      dark
      hide-on-scroll
      color="accent"
      flat
    >
      <v-btn
        to="/"
        text
        class="text-body-2 text-none px-2"
      >
        <v-toolbar-title>
          <v-icon
            aria-hidden="true"
            class="mr-2"
          >
            fas fa-bed-alt
          </v-icon>
          Relax Sounds
        </v-toolbar-title>
      </v-btn>

      <v-spacer />

      <google-cast-launcher class="v-btn v-btn--icon theme--dark v-size--default" />
      <SavePreset />
      <PlayPauseAll />
      <StopAll />

      <template
        v-if="$vuetify.breakpoint.mdAndUp"
        #extension
      >
        <v-tabs centered>
          <v-tab
            v-for="route in routes"
            :key="route.path"
            :to="route.path"
            exact
          >
            <v-icon class="pr-2">
              fas {{ route.meta.icon }} fa-fw
            </v-icon>
            {{ route.name }}
          </v-tab>
        </v-tabs>
      </template>
    </v-app-bar>

    <UpdateSnackbar />

    <v-main>
      <keep-alive>
        <router-view />
      </keep-alive>
      <v-spacer :style="{ height: $vuetify.breakpoint.smAndDown ? '56px' : '28px' }" />
    </v-main>

    <v-bottom-navigation
      v-if="$vuetify.breakpoint.smAndDown"
      fixed
      background-color="accent"
      color="primary"
      dark
      shift
    >
      <v-btn
        v-for="route in routes"
        :key="route.path"
        :to="route.path"
        :value="route.name"
      >
        <span>{{ route.name }}</span>
        <v-icon>fas {{ route.meta.icon }} fa-fw</v-icon>
      </v-btn>
    </v-bottom-navigation>
  </v-app>
</template>

<script>
import SavePreset from './components/SavePreset.vue';
import PlayPauseAll from './components/PlayPauseAll.vue';
import StopAll from './components/StopAll.vue';
import UpdateSnackbar from './components/UpdateSnackbar.vue';

export default {
  name: 'App',

  components: {
    SavePreset,
    StopAll,
    PlayPauseAll,
    UpdateSnackbar,
  },

  computed: {
    routes() {
      return this.$router.options.routes.filter((route) => route.meta?.showInNav);
    },
  },

  beforeMount() {
    // check for browser support
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme)').media !== 'not all') {
      // set to preferred scheme
      const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
      this.$vuetify.theme.dark = mediaQuery.matches;
      // react to changes
      mediaQuery.addEventListener('change', (e) => {
        this.$vuetify.theme.dark = e.matches;
      });
    }
  },
};
</script>

<style lang="scss">
  html {
    --disconnected-color: #fff;
  }

  .v-toolbar__content {
    padding-left: 8px;
  }

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

  .theme--dark {
    &.v-skeleton-loader.transparent > div {
      background: transparent !important;
    }
  }
</style>
