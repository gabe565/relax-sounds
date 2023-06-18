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
        <v-icon aria-hidden="true">fas fa-stop-circle</v-icon>
      </v-btn>
    </template>
    <span>Stop All</span>
  </v-tooltip>
</template>

<script setup>
import { ref } from "vue";
import { usePlayerStore } from "../plugins/store/player";

const showTooltip = ref(false);

const player = usePlayerStore();

const stopAll = () => {
  showTooltip.value = false;
  usePlayerStore().stopAll({});
};
</script>
