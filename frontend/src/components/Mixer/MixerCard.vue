<template>
  <v-card variant="flat" color="cardBackground">
    <v-card-title class="text-h5 text-truncate pb-1">
      <v-icon aria-hidden="true" class="mr-4" size="x-small" :color="iconColor">
        <Icon :icon="sound.icon" />
      </v-icon>
      <span>{{ sound.name }}</span>
    </v-card-title>

    <v-card-actions>
      <v-row no-gutters dense class="pr-2">
        <v-col cols="12">
          <v-slider
            v-model="volumePercentage"
            :prepend-icon="VolumeIcon"
            :min="0"
            :max="100"
            :step="1"
            thumb-size="12"
            thumb-label
            track-size="1"
            hide-details
            aria-label="Volume"
          />
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
          />
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
    </v-card-actions>
  </v-card>
</template>

<script setup>
import { computed } from "vue";
import { Icon } from "@iconify/vue";
import PlayIcon from "~icons/material-symbols/play-arrow-rounded";
import PauseIcon from "~icons/material-symbols/pause-rounded";
import StopIcon from "~icons/material-symbols/stop-rounded";
import VolumeIcon from "~icons/material-symbols/volume-up-rounded";
import SpeedIcon from "~icons/material-symbols/speed-rounded";
import { usePlayerStore } from "../../plugins/store/player";

const props = defineProps({
  sound: {
    type: Object,
    required: true,
  },
});

const player = usePlayerStore();
const iconColor = computed(() => (props.sound.isPlaying ? "primary" : ""));

// const showVolume = computed(() => !props.sound.isStopped);

const volumePercentage = computed({
  get() {
    return props.sound.volume * 100;
  },
  set(value) {
    value /= 100;
    player.volume({ sound: props.sound, value });
    usePlayerStore().updateCast();
  },
});

const rate = computed({
  get() {
    return props.sound.rate;
  },
  set(value) {
    player.rate({ sound: props.sound, value });
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
