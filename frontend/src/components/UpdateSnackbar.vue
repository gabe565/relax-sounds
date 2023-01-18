<template>
  <v-snackbar
    :model-value="offlineReady || needRefresh"
    :timeout="needRefresh ? -1 : 5000"
    location="bottom"
    class="pb-14 pb-md-0"
  >
    <span v-if="needRefresh">New content available, click on reload button to update.</span>
    <span v-else>App ready to work offline</span>
    <template #action="{ attrs }">
      <v-btn
        v-if="needRefresh"
        v-bind="attrs"
        variant="text"
        color="primary"
        :loading="loading"
        @click="updateServiceWorker"
      >
        Refresh
      </v-btn>
      <v-btn icon @click="closePromptUpdateSW">
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
