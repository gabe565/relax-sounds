<template>
  <v-btn
    variant="flat"
    :color="preset.new ? 'newPresetCardBackground' : 'cardBackground'"
    size="x-large"
    class="w-100 d-flex justify-space-between text-none font-weight-regular"
    :aria-label="`Play ${preset.name}`"
    @click="presets.play({ preset })"
  >
    <span class="text-truncate">
      {{ preset.name }}
    </span>
    <template #append>
      <debug-button v-if="debugEnabled" :preset="preset" />
      <share-button :preset="preset" />
      <delete-button :preset="preset" />
    </template>
  </v-btn>
</template>

<script setup>
import ShareButton from "./Buttons/ShareButton.vue";
import DeleteButton from "./Buttons/DeleteButton.vue";
import DebugButton from "./Buttons/DebugButton.vue";
import { usePresetsStore } from "../../plugins/store/presets";

defineProps({
  preset: {
    type: Object,
    required: true,
  },
});

const debugEnabled = import.meta.env.DEV;
const presets = usePresetsStore();
</script>

<style scoped lang="scss">
.v-btn {
  letter-spacing: initial;
}
.v-btn::v-deep .v-btn__content {
  max-width: 100% !important;
  white-space: nowrap !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
}
</style>
