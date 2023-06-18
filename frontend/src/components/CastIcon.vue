<template>
  <v-tooltip v-model="showTooltip" location="bottom" aria-label="Stop All">
    <template #activator="{ props }">
      <google-cast-launcher
        ref="element"
        v-bind="props"
        class="v-btn v-btn--icon v-theme--dark v-btn--density-compact v-btn--size-default mr-4"
      />
    </template>
    <span>Cast</span>
  </v-tooltip>
</template>

<script setup>
import { ref } from "vue";

const showTooltip = ref(false);
const element = ref();
</script>

<style scoped lang="scss">
google-cast-launcher {
  --disconnected-color: #fff;
  cursor: pointer;
}
</style>
