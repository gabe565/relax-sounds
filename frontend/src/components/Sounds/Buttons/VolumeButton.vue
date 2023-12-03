<template>
  <v-dialog v-if="show" location="bottom" location-strategy="connected" max-width="400">
    <template #activator="{ props }">
      <v-btn v-bind="props" elevation="0" icon color="transparent" aria-label="Volume">
        <v-icon :icon="VolumeIcon" aria-hidden="true" />
      </v-btn>
    </template>

    <v-card class="pa-8">
      <v-slider
        v-model="volumePercentage"
        :min="0"
        :max="100"
        :step="1"
        thumb-label
        color="deep-orange-lighten-1"
        class="pb-1"
        hide-details
      />
    </v-card>
  </v-dialog>
</template>

<script setup>
import VolumeIcon from "~icons/material-symbols/volume-up-rounded";
import { computed } from "vue";
import { usePlayerStore } from "../../../plugins/store/player";

const props = defineProps({
  sound: {
    type: Object,
    required: true,
  },
  show: {
    type: Boolean,
    default: false,
  },
});

const volumePercentage = computed({
  get() {
    return props.sound.volume * 100;
  },
  set(newValue) {
    // eslint-disable-next-line vue/no-mutating-props
    props.sound.volume = newValue / 100;
    usePlayerStore().updateCast();
  },
});
</script>
