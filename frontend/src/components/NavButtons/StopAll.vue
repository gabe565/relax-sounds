<template>
  <component
    :is="component"
    :prepend-icon="StopIcon"
    :disabled="player.isStopped"
    title="Stop All"
    text="Stop All"
    aria-label="Stop All"
    @click="player.stopAll({})"
  />
</template>

<script setup>
import { VBtn, VListItem } from "vuetify/components";
import { usePlayerStore } from "../../plugins/store/player";
import StopIcon from "~icons/material-symbols/stop-circle-rounded";

const props = defineProps({
  button: {
    type: Boolean,
    default: false,
  },
});

const player = usePlayerStore();
const component = props.button ? VBtn : VListItem;
</script>
