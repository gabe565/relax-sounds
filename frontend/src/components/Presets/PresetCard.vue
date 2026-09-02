<template>
  <div
    class="card-group-item flex min-h-16 items-center gap-1 overflow-hidden pr-2 transition-colors"
    :class="isActive ? 'bg-primary-container' : 'bg-surface-container-high'"
    @contextmenu.prevent="menuOpen = true"
  >
    <v-btn
      variant="text"
      rounded="0"
      size="x-large"
      class="card-btn h-auto min-h-16 min-w-0 flex-1 justify-start rounded-none"
      :loading="loading"
      :aria-label="`Play ${preset.name}`"
      @click="play"
    >
      <template #prepend>
        <sound-icon :icon="PresetIcon" :size="36" />
      </template>
      <span class="min-w-0 grow text-left">
        <span class="block truncate">{{ preset.name }}</span>
        <span class="block truncate text-xs opacity-70">{{ soundCount }}</span>
      </span>
    </v-btn>

    <share-button :preset="preset" />
    <delete-button :preset="preset" />

    <v-menu v-model="menuOpen" activator="parent" :open-on-click="false" location="bottom start">
      <v-list>
        <v-list-item title="Copy stream URL" :prepend-icon="CopyUrlIcon" @click="copyStreamUrl" />
      </v-list>
    </v-menu>
  </div>
</template>

<script setup>
import { computed, ref } from "vue";
import { toast } from "vue-sonner";
import CopyUrlIcon from "~icons/material-symbols/content-copy-rounded";
import PresetIcon from "~icons/material-symbols/playlist-play-rounded";
import DeleteButton from "@/components/Presets/Buttons/DeleteButton.vue";
import ShareButton from "@/components/Presets/Buttons/ShareButton.vue";
import SoundIcon from "@/components/Sounds/SoundIcon.vue";
import { usePlayer } from "@/plugins/store/player";
import { usePresets } from "@/plugins/store/presets";
import { Preset } from "@/util/Preset";

const props = defineProps({
  preset: {
    type: Preset,
    required: true,
  },
});

const presets = usePresets();
const player = usePlayer();
const loading = ref(false);
const menuOpen = ref(false);

const isActive = computed(() => !player.isStopped && player.currentName === props.preset.name);

const soundCount = computed(() => {
  const count = props.preset.sounds?.length ?? 0;
  return `${count} ${count === 1 ? "sound" : "sounds"}`;
});

const copyStreamUrl = async () => {
  try {
    const url = await props.preset.hlsUrl({ fresh: true });
    await navigator.clipboard.writeText(url);
    toast.success("Stream URL copied");
  } catch (err) {
    console.error(err);
    toast.error(`Failed to copy stream URL:\n${err}`);
  }
};

const play = async () => {
  loading.value = true;
  try {
    await presets.play({ preset: props.preset });
  } catch (err) {
    toast.error(`Failed to load sounds:\n${err}`);
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.v-btn {
  color: inherit;
}

.v-btn :deep(.v-btn__content) {
  min-width: 0;
  flex: 1 1 auto;
}
</style>
