<template>
  <v-btn
    :color="props.sound.isStopped ? 'cardBackground' : 'accent'"
    :loading="sound.isLoading"
    size="x-large"
    class="w-100 justify-start text-none font-weight-regular"
    :aria-label="sound.isPlaying ? `Stop ${sound.name}` : `Play ${sound.name}`"
    variant="flat"
    @click="playStop"
  >
    <template #prepend>
      <v-icon aria-hidden="true" class="mr-4" size="x-large" :color="iconColor">
        <icon :icon="sound.icon" />
      </v-icon>
    </template>
    {{ sound.name }}
  </v-btn>
</template>

<script setup>
import { Icon } from "@iconify/vue";
import { computed } from "vue";
import { useToast } from "vue-toastification";
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
}
</style>
