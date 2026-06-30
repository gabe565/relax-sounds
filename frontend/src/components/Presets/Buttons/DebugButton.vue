<template>
  <v-list-item
    :prepend-icon="DebugIcon"
    title="Copy HLS URL"
    :disabled="disabled"
    aria-label="Copy HLS URL"
    @click.stop="copyUrl"
  />
</template>

<script setup>
import { computed } from "vue";
import { toast } from "vue-sonner";
import DebugIcon from "~icons/material-symbols/bug-report-rounded";
import { usePlayer } from "@/plugins/store/player";
import { Preset } from "@/util/Preset";

const player = usePlayer();

const copyUrl = async () => {
  let preset;
  if (player.isPaused) {
    preset = new Preset({ sounds: player.soundsNotStopped });
  } else {
    preset = new Preset({ sounds: player.soundsPlaying });
  }
  player.pauseAll();
  const url = await preset.hlsUrl();
  try {
    await navigator.clipboard.writeText(url);
    toast.success("HLS URL copied");
  } catch (err) {
    console.error(err);
    toast.error(`Failed to copy URL:\n${err}`);
  }
};

const disabled = computed(() => player.isStopped);
</script>
