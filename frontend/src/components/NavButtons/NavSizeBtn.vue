<template>
  <v-tooltip :text="text" location="top">
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        :icon="preferences.shrinkLeftPanel ? LeftPanelOpenIcon : LeftPanelCloseIcon"
        color="transparent"
        variant="flat"
        aria-label="Toggle left panel"
        @click="preferences.shrinkLeftPanel = !preferences.shrinkLeftPanel"
      />
    </template>
  </v-tooltip>
</template>

<script setup>
import { computed } from "vue";
import LeftPanelCloseIcon from "~icons/material-symbols/left-panel-close-rounded";
import LeftPanelOpenIcon from "~icons/material-symbols/left-panel-open-rounded";
import { usePreferencesStore } from "@/plugins/store/preferences";

const preferences = usePreferencesStore();

const text = computed(() => (preferences.shrinkLeftPanel ? "Expand panel" : "Shrink panel"));
</script>

<style scoped lang="scss"></style>
