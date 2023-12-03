<template>
  <v-card
    variant="flat"
    :color="props.sound.isStopped ? 'cardBackground' : 'accent'"
    class="px-1 py-2"
    :loading="sound.isLoading"
    @click="player.playStop({ sound })"
  >
    <v-card-title class="text-h5">
      <v-icon aria-hidden="true" class="mr-4" size="x-small" :color="iconColor">
        <Icon :icon="sound.icon" />
      </v-icon>
      <span>{{ sound.name }}</span>
    </v-card-title>
  </v-card>
</template>

<script setup>
import { Icon } from "@iconify/vue";
import { usePlayerStore } from "../../plugins/store/player";
import { computed } from "vue";

const props = defineProps({
  sound: {
    type: Object,
    required: true,
  },
});

const iconColor = computed(() => (props.sound.isPlaying ? "primary" : ""));

const player = usePlayerStore();
</script>
