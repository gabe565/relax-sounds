<template>
  <v-btn
    :color="props.sound.isStopped ? 'cardBackground' : 'accent'"
    :loading="sound.isLoading"
    size="x-large"
    class="w-100 justify-start text-none font-weight-regular"
    :aria-label="sound.isPlaying ? `Stop ${sound.name}` : `Play ${sound.name}`"
    @click="player.playStop({ sound })"
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

<style scoped>
.v-btn {
  letter-spacing: initial;
}
</style>
