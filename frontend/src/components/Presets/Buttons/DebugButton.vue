<template>
  <v-list-item
    v-if="listItem"
    :prepend-icon="DebugIcon"
    title="Debug"
    target="_blank"
    :disabled="disabled"
    @click="openUrl"
  />
  <v-btn
    v-else
    icon
    variant="flat"
    color="transparent"
    target="_blank"
    :disabled="disabled"
    @click="openUrl"
  >
    <v-icon :icon="DebugIcon" aria-hidden="true" />
  </v-btn>
</template>

<script setup>
import DebugIcon from "~icons/material-symbols/bug-report-rounded";
import { Preset } from "../../../util/Preset";
import { VBtn, VListItem } from "vuetify/components";
import { usePlayerStore } from "../../../plugins/store/player";
import { computed } from "vue";

const props = defineProps({
  listItem: {
    type: Boolean,
    default: false,
  },
  preset: {
    type: Preset,
    default: null,
  },
});

const player = usePlayerStore();

const openUrl = () => {
  let preset;
  if (props.preset) {
    preset = props.preset;
  } else {
    preset = new Preset({ sounds: player.soundsPlaying });
  }
  player.pauseAll();
  window.open(preset.mixUrlAs("mp3"), "_blank");
};

const disabled = computed(() => !props.preset && player.isStopped);
</script>
