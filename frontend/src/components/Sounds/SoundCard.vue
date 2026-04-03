<template>
  <v-dialog
    max-width="400"
    location-strategy="connected"
    location="bottom center"
    scroll-strategy="reposition"
  >
    <template #activator="{ props: dialogProps }">
      <v-btn
        :active="sound.isPlaying"
        :loading="sound.isLoading"
        size="x-large"
        class="group card-btn bg-card-background w-full justify-start border transition border-transparent v-active:bg-accent v-active:border-secondary/35 v-active:shadow-[0_0_12px] v-active:shadow-secondary/25"
        :aria-label="sound.isPlaying ? `Stop ${sound.name}` : `Play ${sound.name}`"
        variant="flat"
        @click="playStop"
        @contextmenu.prevent="dialogProps.onClick"
      >
        <template #prepend>
          <v-icon size="x-large" class="group-v-active:text-secondary">
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
import { toast } from "vue-sonner";
import MixerCard from "@/components/Mixer/MixerCard.vue";
import { usePlayer } from "@/plugins/store/player";

const props = defineProps({
  sound: {
    type: Object,
    required: true,
  },
});

const player = usePlayer();

const playStop = async () => {
  try {
    await player.playStop({ sound: props.sound });
  } catch (err) {
    toast.error(`Failed to load sound:\n${err}`);
  }
};
</script>
