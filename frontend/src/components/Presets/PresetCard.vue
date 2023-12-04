<template>
  <v-fade-transition>
    <v-card
      variant="flat"
      :dark="preset.new"
      :color="preset.new ? 'newPresetCardBackground' : 'cardBackground'"
      transition="fade-transition"
      @click="presets.play({ preset })"
    >
      <div class="d-flex flex-row flex-nowrap align-center">
        <div class="text-truncate flex-grow-1">
          <v-card-title class="text-h5 pa-4 d-block text-truncate">
            {{ preset.name }}
          </v-card-title>
        </div>
        <debug-button v-if="debugEnabled" :preset="preset" />
        <share-button :preset="preset" />
        <delete-button :preset="preset" />
        <play-button :preset="preset" class="d-sr-only" />
      </div>
    </v-card>
  </v-fade-transition>
</template>

<script setup>
import ShareButton from "./Buttons/ShareButton.vue";
import DeleteButton from "./Buttons/DeleteButton.vue";
import DebugButton from "./Buttons/DebugButton.vue";
import PlayButton from "./Buttons/PlayButton.vue";
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
