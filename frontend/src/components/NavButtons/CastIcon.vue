<template>
  <v-tooltip v-if="castEnabled" text="Cast" :location="tooltipLocation">
    <template #activator="{ props }">
      <v-btn v-bind="props" icon title="Cast" aria-label="Cast" @click.stop="castLauncher.click()">
        <v-icon>
          <google-cast-launcher ref="castLauncher" class="relative transition-colors" />
        </v-icon>
      </v-btn>
    </template>
  </v-tooltip>
</template>

<script setup>
import { useTemplateRef } from "vue";
import { castEnabled } from "@/cast";

const castLauncher = useTemplateRef("castLauncher");

defineProps({
  tooltipLocation: {
    type: String,
    default: "top",
  },
});
</script>

<style>
google-cast-launcher {
  --disconnected-color: rgb(var(--v-theme-on-surface));
  --connected-color: rgb(var(--v-theme-primary));
}
</style>
