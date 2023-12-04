<template>
  <component
    :is="component"
    :prepend-icon="StopIcon"
    :disabled="player.isStopped"
    title="Stop All"
    text="Stop All"
    @click="player.stopAll({})"
  />
</template>

<script setup>
import StopIcon from "~icons/material-symbols/stop-circle-rounded";
import { usePlayerStore } from "../../plugins/store/player";
import { VBtn, VListItem } from "vuetify/components";

const props = defineProps({
  button: {
    type: Boolean,
    default: false,
  },
});

const player = usePlayerStore();
const component = props.button ? VBtn : VListItem;
</script>
