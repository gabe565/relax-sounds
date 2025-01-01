<template>
  <v-card variant="flat" color="cardBackground">
    <template #title>
      <v-icon aria-hidden="true" class="mr-4" size="x-small" :color="iconColor">
        <icon :icon="sound.icon" />
      </v-icon>
      <span>{{ sound.name }}</span>
    </template>

    <template #actions>
      <v-row no-gutters dense class="pr-2">
        <v-col cols="12">
          <v-slider
            v-model="volume"
            :prepend-icon="VolumeIcon"
            :min="0"
            :max="1"
            :step="0.01"
            thumb-size="12"
            thumb-label
            track-size="1"
            hide-details
            aria-label="Volume"
          >
            <template #thumb-label>{{ Math.round(volume * 100) }}%</template>
          </v-slider>
        </v-col>
        <v-col cols="12">
          <v-slider
            v-model="rate"
            :prepend-icon="SpeedIcon"
            :min="0.5"
            :max="1.5"
            :step="0.05"
            thumb-size="12"
            thumb-label
            track-size="1"
            hide-details
            aria-label="Speed"
          >
            <template #thumb-label>{{ Math.round(rate * 100) }}%</template>
          </v-slider>
        </v-col>
        <v-col cols="12">
          <v-slider
            v-model="pan"
            :prepend-icon="WidthIcon"
            :min="-1"
            :max="1"
            :step="0.05"
            thumb-size="12"
            thumb-label
            track-size="1"
            hide-details
            aria-label="Pan"
          >
            <template #thumb-label>{{ Math.round(pan * 100) }}%</template>
          </v-slider>
        </v-col>
      </v-row>

      <v-btn
        elevation="0"
        icon
        :aria-label="sound.isPlaying ? 'Pause' : 'Play'"
        :loading="sound.isLoading"
        @click.stop="player.playPause({ sound })"
      >
        <v-icon v-bind="iconProps" aria-hidden="true" />
      </v-btn>

      <v-btn icon aria-label="Stop" @click.stop="player.stop({ sound })">
        <v-icon :icon="StopIcon" aria-hidden="true" />
      </v-btn>
    </template>
  </v-card>
</template>

<script setup>
import { Icon } from "@iconify/vue";
import { computed } from "vue";
import PauseIcon from "~icons/material-symbols/pause-rounded";
import PlayIcon from "~icons/material-symbols/play-arrow-rounded";
import SpeedIcon from "~icons/material-symbols/speed-rounded";
import StopIcon from "~icons/material-symbols/stop-rounded";
import VolumeIcon from "~icons/material-symbols/volume-up-rounded";
import WidthIcon from "~icons/material-symbols/width-rounded";
import { usePlayerStore } from "@/plugins/store/player";

const props = defineProps({
  sound: {
    type: Object,
    required: true,
  },
});

const player = usePlayerStore();
const iconColor = computed(() => (props.sound.isPlaying ? "primary" : ""));

const volume = computed({
  get: () => props.sound.volume,
  set(value) {
    player.volume({ sound: props.sound, value });
    usePlayerStore().updateCast();
  },
});

const rate = computed({
  get: () => props.sound.rate,
  set(value) {
    player.rate({ sound: props.sound, value });
    usePlayerStore().updateCast();
  },
});

const pan = computed({
  get: () => props.sound.pan,
  set(value) {
    player.pan({ sound: props.sound, value });
    usePlayerStore().updateCast();
  },
});

const iconProps = computed(() => {
  if (!props.sound.isPlaying) {
    return { icon: PlayIcon, size: "x-large" };
  } else {
    return { icon: PauseIcon, size: "large" };
  }
});
</script>
