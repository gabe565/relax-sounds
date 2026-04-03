<template>
  <v-dialog
    max-width="400"
    location-strategy="connected"
    location="bottom center"
    scroll-strategy="reposition"
  >
    <template #activator="{ props: dialogProps }">
      <v-btn
        :color="props.sound.isStopped ? 'card-background' : 'accent'"
        :loading="sound.isLoading"
        size="x-large"
        class="card-btn w-full justify-start border transition"
        :class="[
          props.sound.isStopped
            ? 'border-transparent'
            : 'border-secondary/35 shadow-[0_0_12px] shadow-secondary/25',
        ]"
        :aria-label="sound.isPlaying ? `Stop ${sound.name}` : `Play ${sound.name}`"
        variant="flat"
        @click="playStop"
        @contextmenu.prevent="dialogProps.onClick"
      >
        <template #prepend>
          <v-icon size="x-large" :color="iconColor">
            <icon :icon="sound.icon" />
          </v-icon>
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
import { Icon } from "@iconify/vue";
import { computed } from "vue";
import { toast } from "vue-sonner";
import MixerCard from "@/components/Mixer/MixerCard.vue";
import { usePlayer } from "@/plugins/store/player";

const props = defineProps({
  sound: {
    type: Object,
    required: true,
  },
});

/* Use secondary (Amber) for playing icons to break up the purple */
const iconColor = computed(() => (props.sound.isPlaying ? "secondary" : ""));
const player = usePlayer();

const playStop = async () => {
  try {
    await player.playStop({ sound: props.sound });
  } catch (err) {
    toast.error(`Failed to load sound:\n${err}`);
  }
};
</script>
