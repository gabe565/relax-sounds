<template>
  <v-tooltip v-if="player.castEnabled" text="Cast" :location="tooltipLocation">
    <template #activator="{ props }">
      <v-btn v-bind="props" icon title="Cast" aria-label="Cast">
        <v-icon>
          <google-cast-launcher class="absolute inset-0 transition-colors" />
        </v-icon>
      </v-btn>
    </template>
  </v-tooltip>
</template>

<script setup>
import { usePlayer } from "@/plugins/store/player";

defineProps({
  tooltipLocation: {
    type: String,
    default: "top",
  },
});

const player = usePlayer();
</script>

<style>
google-cast-launcher {
  --disconnected-color: rgb(var(--v-theme-on-surface));
  --connected-color: rgb(var(--v-theme-primary));
}
</style>
