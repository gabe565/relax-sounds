<template>
  <v-tooltip v-model="showTooltip" location="bottom" aria-label="Stop All">
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        :disabled="player.isStopped"
        icon
        aria-label="Stop All"
        @click="stopAll"
      >
        <v-icon :icon="StopIcon" aria-hidden="true" />
      </v-btn>
    </template>
    <span>Stop All</span>
  </v-tooltip>
</template>

<script setup>
import { ref } from "vue";
import StopIcon from "~icons/material-symbols/stop-circle-rounded";
import { usePlayerStore } from "../../plugins/store/player";

const showTooltip = ref(false);

const player = usePlayerStore();

const stopAll = () => {
  showTooltip.value = false;
  usePlayerStore().stopAll({});
};
</script>
