<template>
  <v-card
    variant="flat"
    :color="props.sound.isStopped ? 'cardBackground' : 'accent'"
    class="px-1 py-2"
  >
    <div class="d-flex flex-row flex-nowrap align-center">
      <div class="text-truncate flex-grow-1">
        <v-card-title class="text-h5">
          <v-icon aria-hidden="true" class="mr-4" size="x-small" :color="iconColor">
            <Icon :icon="sound.icon" />
          </v-icon>
          <span>{{ sound.name }}</span>
        </v-card-title>
      </div>

      <volume-button :sound="sound" :show="showVolume" />
      <play-pause-button :sound="sound" />
    </div>
  </v-card>
</template>

<script setup>
import { computed } from "vue";
import { Icon } from "@iconify/vue";
import PlayPauseButton from "./Buttons/PlayPauseButton.vue";
import VolumeButton from "./Buttons/VolumeButton.vue";

const props = defineProps({
  sound: {
    type: Object,
    required: true,
  },
});

const iconColor = computed(() => (props.sound.isPlaying ? "primary" : ""));

const showVolume = computed(() => !props.sound.isStopped);
</script>
