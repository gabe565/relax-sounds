<template>
  <v-card
    variant="flat"
    color="surface-container-high"
    rounded="xl"
    class="border transition-colors"
    :class="sound.isPlaying ? 'border-primary/50' : 'border-outline-variant'"
  >
    <div class="flex items-center gap-3 px-4 pt-4 pb-1">
      <sound-icon :icon="sound.icon" />
      <span class="grow truncate">{{ sound.name }}</span>

      <v-btn
        v-if="props.closable"
        :icon="CloseIcon"
        variant="text"
        density="comfortable"
        aria-label="Close Mixer"
        @click="emit('close')"
      />
    </div>

    <v-card-text class="flex flex-col gap-2 px-4 pt-2 pb-0">
      <div v-for="control in controls" :key="control.label">
        <div class="text-on-surface-variant flex items-center gap-2 text-sm">
          <v-icon :icon="control.icon" size="small" />
          <span class="grow">{{ control.label }}</span>
          <span class="tabular-nums">{{ control.display }}</span>
        </div>
        <v-slider
          v-model="control.model.value"
          :min="control.min"
          :max="control.max"
          :step="control.step"
          :ticks="control.detent === undefined ? undefined : { [control.detent]: '' }"
          :show-ticks="control.detent === undefined ? false : 'always'"
          thumb-size="12"
          hide-details
          :aria-label="control.label"
          color="primary"
        />
      </div>
    </v-card-text>

    <v-card-actions class="justify-center gap-2 pt-0 pb-4">
      <v-btn
        :icon="sound.isPlaying ? PauseIcon : PlayIcon"
        variant="flat"
        color="primary-container"
        :loading="sound.isLoading"
        :aria-label="sound.isPlaying ? 'Pause' : 'Play'"
        @click.stop="player.playPause({ sound })"
      />

      <v-btn
        :icon="StopIcon"
        variant="tonal"
        :disabled="sound.isStopped"
        aria-label="Stop"
        @click.stop="
          player.stop({ sound });
          emit('close');
        "
      />
    </v-card-actions>
  </v-card>
</template>

<script setup>
import { computed } from "vue";
import CloseIcon from "~icons/material-symbols/close-rounded";
import PauseIcon from "~icons/material-symbols/pause-rounded";
import PlayIcon from "~icons/material-symbols/play-arrow-rounded";
import SpeedIcon from "~icons/material-symbols/speed-rounded";
import StopIcon from "~icons/material-symbols/stop-rounded";
import VolumeIcon from "~icons/material-symbols/volume-up-rounded";
import WidthIcon from "~icons/material-symbols/width-rounded";
import SoundIcon from "@/components/Sounds/SoundIcon.vue";
import { usePlayer } from "@/plugins/store/player";

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

const player = usePlayer();

const volume = computed({
  get: () => props.sound.volume,
  set(value) {
    player.volume({ sound: props.sound, value });
    player.updateCast();
  },
});

const rate = computed({
  get: () => props.sound.rate,
  set(value) {
    player.rate({ sound: props.sound, value });
    player.updateCast();
  },
});

const pan = computed({
  get: () => props.sound.pan,
  set(value) {
    player.pan({ sound: props.sound, value });
    player.updateCast();
  },
});

const panDisplay = computed(() => {
  const value = Math.round(pan.value * 100);
  if (value === 0) return "Center";
  return `${Math.abs(value)}% ${value < 0 ? "L" : "R"}`;
});

const controls = computed(() => [
  {
    label: "Volume",
    icon: VolumeIcon,
    model: volume,
    min: 0,
    max: 1,
    step: 0.01,
    display: `${Math.round(volume.value * 100)}%`,
  },
  {
    label: "Speed",
    icon: SpeedIcon,
    model: rate,
    min: 0.5,
    max: 1.5,
    step: 0.05,
    detent: 1,
    display: `${rate.value.toFixed(2)}×`,
  },
  {
    label: "Pan",
    icon: WidthIcon,
    model: pan,
    min: -1,
    max: 1,
    step: 0.05,
    detent: 0,
    display: panDisplay.value,
  },
]);
</script>
