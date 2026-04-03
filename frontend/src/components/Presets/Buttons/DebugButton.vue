<template>
  <v-list-item
    :prepend-icon="DebugIcon"
    title="Debug"
    target="_blank"
    :disabled="disabled"
    aria-label="Debug"
    @click.stop="openUrl"
  />
</template>

<script setup>
import { computed } from "vue";
import DebugIcon from "~icons/material-symbols/bug-report-rounded";
import { usePlayer } from "@/plugins/store/player";
import { Preset } from "@/util/Preset";

const player = usePlayer();

const openUrl = async () => {
  let preset;
  if (player.isPaused) {
    preset = new Preset({ sounds: player.soundsNotStopped });
  } else {
    preset = new Preset({ sounds: player.soundsPlaying });
  }
  player.pauseAll();
  const url = await preset.mixUrlAs("mp3");
  window.open(url, "_blank");
};

const disabled = computed(() => player.isStopped);
</script>
