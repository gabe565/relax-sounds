<template>
  <v-dialog
    max-width="400"
    location-strategy="connected"
    location="bottom center"
    scroll-strategy="reposition"
  >
    <template #activator="{ props: dialogProps }">
      <v-btn
        :color="props.sound.isStopped ? 'cardBackground' : 'accent'"
        :loading="sound.isLoading"
        size="x-large"
        class="w-100 justify-start text-none font-weight-regular"
        :aria-label="sound.isPlaying ? `Stop ${sound.name}` : `Play ${sound.name}`"
        variant="flat"
        @click="playStop"
        @contextmenu.prevent="dialogProps.onClick"
      >
        <template #prepend>
          <v-icon class="mr-4" size="x-large" :color="iconColor">
            <icon :icon="sound.icon" />
          </v-icon>
        </template>
        <span class="text-truncate">{{ sound.name }}</span>
      </v-btn>
    </template>

    <template #default="{ isActive }">
      <mixer-card :sound="sound" closable @close="isActive.value = false" />
    </template>
  </v-dialog>
</template>

<script setup>
import { Icon } from "@iconify/vue";
import { computed } from "vue";
import { useToast } from "vue-toastification";
import MixerCard from "@/components/Mixer/MixerCard.vue";
import { usePlayerStore } from "@/plugins/store/player";

const props = defineProps({
  sound: {
    type: Object,
    required: true,
  },
});

const iconColor = computed(() => (props.sound.isPlaying ? "primary" : ""));
const toast = useToast();
const player = usePlayerStore();

const playStop = async () => {
  try {
    await player.playStop({ sound: props.sound });
  } catch (err) {
    toast.error(`Failed to load sound:\n${err}`);
  }
};
</script>

<style scoped>
.v-btn {
  letter-spacing: initial;
  padding: 0 6px 0 24px;
}
:deep(.v-btn__content) {
  min-width: 0;
}
</style>
