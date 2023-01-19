<template>
  <v-snackbar
    :model-value="offlineReady || needRefresh"
    :timeout="needRefresh ? -1 : 5000"
    location="bottom"
    content-class="mb-14 mb-md-0"
  >
    <span v-if="needRefresh">New content available, click on reload button to update.</span>
    <span v-else>App ready to work offline</span>
    <template #actions>
      <v-btn
        v-if="needRefresh"
        variant="text"
        color="primary"
        :loading="loading"
        @click="updateServiceWorker"
      >
        Refresh
      </v-btn>
      <v-btn icon aria-label="Close" @click="closePromptUpdateSW">
        <v-icon>$close</v-icon>
      </v-btn>
    </template>
  </v-snackbar>
</template>

<script>
import useRegisterSW from "../mixins/useRegisterSW";

const intervalMS = 60 * 60 * 1000;

export default {
  name: "UpdateSnackbar",
  mixins: [useRegisterSW],
  methods: {
    handleSWManualUpdates(r) {
      if (r) {
        setInterval(() => r.update(), intervalMS);
      }
    },
  },
};
</script>
