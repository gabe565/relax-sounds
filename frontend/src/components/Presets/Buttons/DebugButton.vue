<template>
  <v-list-item
    v-if="listItem"
    :prepend-icon="DebugIcon"
    title="Debug"
    target="_blank"
    :disabled="disabled"
    aria-label="Debug"
    @click.stop="openUrl"
  />
  <v-tooltip v-else text="Debug" location="bottom">
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        icon
        variant="flat"
        color="transparent"
        target="_blank"
        :disabled="disabled"
        aria-label="Debug"
        @click.stop="openUrl"
      >
        <v-icon :icon="DebugIcon" />
      </v-btn>
    </template>
  </v-tooltip>
</template>

<script setup>
import { computed } from "vue";
import { VBtn, VListItem } from "vuetify/components";
import DebugIcon from "~icons/material-symbols/bug-report-rounded";
import { usePlayerStore } from "@/plugins/store/player";
import { Preset } from "@/util/Preset";

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

const openUrl = async () => {
  let preset;
  if (props.preset) {
    preset = props.preset;
  } else if (player.isPaused) {
    preset = new Preset({ sounds: player.soundsNotStopped });
  } else {
    preset = new Preset({ sounds: player.soundsPlaying });
  }
  player.pauseAll();
  const url = await preset.mixUrlAs("mp3");
  window.open(url, "_blank");
};

const disabled = computed(() => !props.preset && player.isStopped);
</script>
