<template>
    <v-snackbar v-model="show" :timeout="timeout" bottom left class="pr-4">
      New version available!
      <v-spacer/>
      <v-btn icon @click.native="refreshApp">
        <v-icon>fal fa-sync</v-icon>
      </v-btn>
      <v-btn icon @click="show = false">
        <v-icon>$close</v-icon>
      </v-btn>
    </v-snackbar>
</template>

<script>
export default {
  name: 'UpdateSnackbar',

  data: () => ({
    refreshing: false,
    registration: null,
    show: false,
    timeout: 0,
  }),

  created() {
    // Listen for swUpdated event and display refresh snackbar as required.
    document.addEventListener('swUpdated', this.showRefreshUI, { once: true });
    // Refresh all open app tabs when a new service worker is installed.
    navigator.serviceWorker.addEventListener('controllerchange', () => {
      if (this.refreshing) return;
      this.refreshing = true;
      window.location.reload();
    });
  },

  methods: {
    showRefreshUI(e) {
      // Display a snackbar inviting the user to refresh/reload the app due
      // to an app update being available.
      // The new service worker is installed, but not yet active.
      // Store the ServiceWorkerRegistration instance for later use.
      this.registration = e.detail;
      this.show = true;
    },

    refreshApp() {
      this.show = false;
      // Protect against missing registration.waiting.
      if (!this.registration || !this.registration.waiting) { return; }
      this.registration.waiting.postMessage('skipWaiting');
    },
  },
};
</script>

<style scoped>

</style>
