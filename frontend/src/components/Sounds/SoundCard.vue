<template>
  <v-dialog
    max-width="400"
    location-strategy="connected"
    location="bottom center"
    scroll-strategy="reposition"
  >
    <template #activator="{ props: dialogProps }">
      <v-btn
        :active="!sound.isStopped"
        :loading="sound.isLoading"
        size="x-large"
        class="card-btn card-group-item w-full justify-start transition-colors"
        :color="stateColor"
        :aria-label="sound.isPlaying ? `Stop ${sound.name}` : `Play ${sound.name}`"
        variant="flat"
        @click="playStop"
        @contextmenu.prevent="dialogProps.onClick"
      >
        <template #prepend>
          <sound-icon :icon="sound.icon" :size="36" />
        </template>
        <span class="truncate">{{ sound.name }}</span>
      </v-btn>
    </template>

    <template #default="{ isActive }">
      <mixer-card :sound="sound" closable @close="isActive.value = false" />
    </template>
  </v-dialog>
</template>

<script setup>
import { computed } from "vue";
import { toast } from "vue-sonner";
import MixerCard from "@/components/Mixer/MixerCard.vue";
import SoundIcon from "@/components/Sounds/SoundIcon.vue";
import { usePlayer } from "@/plugins/store/player";

const props = defineProps({
  sound: {
    type: Object,
    required: true,
  },
});

const player = usePlayer();

const stateColor = computed(() => {
  if (props.sound.isPlaying) return "primary-container";
  if (props.sound.isPaused) return "secondary-container";
  return "surface-container-high";
});

const playStop = async () => {
  try {
    await player.playStop({ sound: props.sound });
  } catch (err) {
    toast.error(`Failed to load sound:\n${err}`);
  }
};
</script>
