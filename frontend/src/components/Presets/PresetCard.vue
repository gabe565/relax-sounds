<template>
  <v-btn
    variant="flat"
    :color="preset.new ? 'newPresetCardBackground' : 'cardBackground'"
    size="x-large"
    class="w-100 d-flex justify-space-between text-none font-weight-regular"
    :aria-label="`Play ${preset.name}`"
    :loading="loading"
    @click="play"
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
import { ref } from "vue";
import ShareButton from "./Buttons/ShareButton.vue";
import DeleteButton from "./Buttons/DeleteButton.vue";
import DebugButton from "./Buttons/DebugButton.vue";
import { usePresetsStore } from "../../plugins/store/presets";
import { Preset } from "../../util/Preset";
import { useToast } from "vue-toastification";

const props = defineProps({
  preset: {
    type: Preset,
    required: true,
  },
});

const debugEnabled = import.meta.env.DEV;
const toast = useToast();
const presets = usePresetsStore();
const loading = ref(false);

const play = async () => {
  loading.value = true;
  try {
    await presets.play({ preset: props.preset });
  } catch (err) {
    toast.error("Failed to load sounds");
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped lang="scss">
.v-btn {
  letter-spacing: initial;

  &:deep(.v-btn__content) {
    max-width: 100% !important;
    white-space: nowrap !important;
    overflow: hidden !important;
    text-overflow: ellipsis !important;
  }
}
</style>
