<template>
  <v-card variant="flat" color="cardBackground">
    <v-card-title class="d-flex align-center py-4">
      <v-icon class="mr-4" size="x-small" :color="iconColor">
        <icon :icon="sound.icon" />
      </v-icon>
      <span class="flex-grow-1">
        {{ sound.name }}
      </span>

      <v-btn
        v-if="props.closable"
        :icon="CloseIcon"
        variant="text"
        density="comfortable"
        aria-label="Close Mixer"
        @click="emit('close')"
      />
    </v-card-title>

    <v-card-text class="pt-1 pb-0 pr-8">
      <v-row class="flex-column">
        <v-col>
          <v-slider
            v-model="volume"
            :prepend-icon="VolumeIcon"
            :min="0"
            :max="1"
            :step="0.01"
            thumb-size="12"
            thumb-label
            hide-details
            aria-label="Volume"
          >
            <template #thumb-label>{{ Math.round(volume * 100) }}%</template>
          </v-slider>
        </v-col>
        <v-col>
          <v-slider
            v-model="rate"
            :prepend-icon="SpeedIcon"
            :min="0.5"
            :max="1.5"
            :step="0.05"
            thumb-size="12"
            thumb-label
            hide-details
            aria-label="Speed"
          >
            <template #thumb-label>{{ Math.round(rate * 100) }}%</template>
          </v-slider>
        </v-col>
        <v-col>
          <v-slider
            v-model="pan"
            :prepend-icon="WidthIcon"
            :min="-1"
            :max="1"
            :step="0.05"
            thumb-size="12"
            thumb-label
            hide-details
            aria-label="Pan"
          >
            <template #thumb-label>{{ Math.round(pan * 100) }}%</template>
          </v-slider>
        </v-col>
      </v-row>
    </v-card-text>

    <v-card-actions class="justify-center">
      <v-btn
        elevation="0"
        icon
        :aria-label="sound.isPlaying ? 'Pause' : 'Play'"
        :loading="sound.isLoading"
        @click.stop="player.playPause({ sound })"
      >
        <v-icon v-bind="iconProps" />
      </v-btn>

      <v-btn
        icon
        aria-label="Stop"
        :disabled="sound.isStopped"
        @click.stop="
          player.stop({ sound });
          emit('close');
        "
      >
        <v-icon :icon="StopIcon" />
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup>
import { Icon } from "@iconify/vue";
import { computed } from "vue";
import CloseIcon from "~icons/material-symbols/close-rounded";
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
  closable: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["close"]);

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
