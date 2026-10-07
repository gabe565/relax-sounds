<template>
  <v-tooltip v-if="player.castEnabled" text="Cast" :location="tooltipLocation">
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        icon
        title="Cast"
        aria-label="Cast"
        :color="player.castConnected ? 'primary' : undefined"
        @click="requestSession"
      >
        <v-icon :icon="player.castConnected ? CastConnectedIcon : CastIcon" />
      </v-btn>
    </template>
  </v-tooltip>
</template>

<script setup>
import { toast } from "vue-sonner";
import CastConnectedIcon from "~icons/material-symbols/cast-connected";
import CastIcon from "~icons/material-symbols/cast-outline";
import { usePlayer } from "@/plugins/store/player";

defineProps({
  tooltipLocation: {
    type: String,
    default: "top",
  },
});

const player = usePlayer();

const requestSession = async () => {
  try {
    await globalThis.cast.framework.CastContext.getInstance().requestSession();
  } catch (error) {
    if (error !== globalThis.chrome.cast.ErrorCode.CANCEL) {
      toast.error(`Failed to cast:\n${error}`);
    }
  }
};
</script>
