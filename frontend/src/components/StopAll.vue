<template>
  <v-tooltip v-model="showTooltip" location="bottom" aria-label="Stop All">
    <template #activator="{ props }">
      <v-btn v-bind="props" :disabled="isStopped" icon aria-label="Stop All" @click="stopAll">
        <v-icon aria-hidden="true">fas fa-stop-circle</v-icon>
      </v-btn>
    </template>
    <span>Stop All</span>
  </v-tooltip>
</template>

<script>
import { mapState } from "pinia";
import { usePlayerStore } from "../plugins/store/player";

export default {
  name: "StopAll",

  data: () => ({
    showTooltip: false,
  }),

  computed: mapState(usePlayerStore, ["isStopped"]),

  methods: {
    stopAll() {
      this.showTooltip = false;
      usePlayerStore().stopAll({});
    },
  },
};
</script>
