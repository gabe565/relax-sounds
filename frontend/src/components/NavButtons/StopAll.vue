<template>
  <component
    :is="component"
    :prepend-icon="StopIcon"
    :disabled="player.isStopped"
    aria-label="Stop"
    @click="player.stopAll({})"
  >
    Stop
  </component>
</template>

<script setup>
import { VBtn, VListItem } from "vuetify/components";
import StopIcon from "~icons/material-symbols/stop-circle-rounded";
import { usePlayerStore } from "@/plugins/store/player";

const props = defineProps({
  button: {
    type: Boolean,
    default: false,
  },
});

const player = usePlayerStore();
const component = props.button ? VBtn : VListItem;
</script>
