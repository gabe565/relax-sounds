<template>
  <v-card variant="flat" color="cardBackground">
    <v-card-title class="text-h5 text-truncate pb-1">
      <v-icon aria-hidden="true" class="mr-4" size="x-small" :color="iconColor">
        <Icon :icon="sound.icon" />
      </v-icon>
      <span>{{ sound.name }}</span>
    </v-card-title>

    <v-card-actions class="pb-0">
      <v-slider
        v-model="volumePercentage"
        :min="0"
        :max="100"
        :step="1"
        thumb-size="12"
        track-size="1"
        hide-details
        class="pl-10 pr-4"
        aria-label="Volume"
      />

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
import StopIcon from "~icons/material-symbols/stop-rounded.vue";
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
  set(newValue) {
    // eslint-disable-next-line vue/no-mutating-props
    props.sound.volume = newValue / 100;
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
