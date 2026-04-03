<template>
  <v-btn
    :active="player.currentName === props.preset.name"
    variant="flat"
    size="x-large"
    class="card-btn bg-card-background w-full flex justify-between pr-2 v-active:bg-accent v-active:light:text-white"
    :aria-label="`Play ${preset.name}`"
    :loading="loading"
    @click="play"
  >
    <span class="truncate">
      {{ preset.name }}
    </span>
    <template #append>
      <v-btn-group>
        <share-button :preset="preset" />
        <delete-button :preset="preset" />
      </v-btn-group>
    </template>
  </v-btn>
</template>

<script setup>
import { ref } from "vue";
import { toast } from "vue-sonner";
import DeleteButton from "@/components/Presets/Buttons/DeleteButton.vue";
import ShareButton from "@/components/Presets/Buttons/ShareButton.vue";
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
