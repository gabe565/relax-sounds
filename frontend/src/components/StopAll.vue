<template>
  <v-tooltip v-model="showTooltip" location="bottom">
    <template #activator="{ props }">
      <v-btn v-bind="props" :disabled="isStopped" icon aria-label="Stop All" @click="stopAll">
        <v-icon aria-hidden="true">fas fa-stop-circle</v-icon>
      </v-btn>
    </template>
    <span>Stop All</span>
  </v-tooltip>
</template>

<script>
import { mapGetters } from "vuex";

export default {
  name: "StopAll",

  data: () => ({
    showTooltip: false,
  }),

  computed: mapGetters("player", ["isStopped"]),

  methods: {
    stopAll() {
      this.showTooltip = false;
      this.$store.dispatch("player/stopAll", {});
    },
  },
};
</script>
