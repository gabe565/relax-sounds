<template>
  <v-card variant="outlined" :dark="showProgress" class="px-1 py-2">
    <div class="d-flex flex-row flex-nowrap align-center">
      <div class="text-truncate flex-grow-1">
        <v-card-title class="text-h5">
          <v-icon aria-hidden="true" class="mr-4" size="x-small" :color="iconColor">
            <Icon :icon="sound.icon" />
          </v-icon>
          {{ sound.name }}
        </v-card-title>
      </div>

      <v-dialog v-if="showProgress" location="bottom" location-strategy="connected" max-width="400">
        <template #activator="{ props }">
          <v-btn v-bind="props" elevation="0" icon variant="plain" aria-label="Volume">
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

      <v-btn
        elevation="0"
        icon
        variant="plain"
        :aria-label="sound.isPlaying ? 'Stop' : 'Play'"
        :loading="sound.isLoading"
        @click.stop="playStop"
      >
        <v-icon v-if="sound.isPlaying" :icon="PauseIcon" aria-hidden="true" size="large" />
        <v-icon v-else :icon="PlayIcon" aria-hidden="true" size="x-large" />
      </v-btn>
    </div>
  </v-card>
</template>

<script setup>
import { computed } from "vue";
import { Icon } from "@iconify/vue";
import PlayIcon from "~icons/material-symbols/play-arrow-rounded";
import PauseIcon from "~icons/material-symbols/pause-rounded";
import VolumeIcon from "~icons/material-symbols/volume-up-rounded";
import { usePlayerStore } from "../../plugins/store/player";

const props = defineProps({
  sound: {
    type: Object,
    required: true,
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

const iconColor = computed(() => (props.sound.isStopped ? "" : "primary"));

const showProgress = computed(() => !props.sound.isStopped);

const playStop = async () => {
  return usePlayerStore().playStop({ sound: props.sound });
};
</script>
