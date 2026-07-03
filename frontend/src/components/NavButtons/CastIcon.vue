<template>
  <v-tooltip v-if="player.castEnabled" text="Cast" :location="tooltipLocation">
    <template #activator="{ props }">
      <v-btn v-show="show" v-bind="props" icon title="Cast" aria-label="Cast">
        <v-icon>
          <google-cast-launcher ref="launcherRef" class="absolute inset-0 transition-colors" />
        </v-icon>
      </v-btn>
    </template>
  </v-tooltip>
</template>

<script setup>
import { useMutationObserver } from "@vueuse/core";
import { ref, useTemplateRef, watch } from "vue";
import { usePlayer } from "@/plugins/store/player";

defineProps({
  tooltipLocation: {
    type: String,
    default: "top",
  },
});

const player = usePlayer();
const launcherRef = useTemplateRef("launcherRef");
const show = ref(false);

const update = () => {
  show.value = launcherRef.value ? getComputedStyle(launcherRef.value).display !== "none" : false;
};

watch(launcherRef, update, { once: true });
useMutationObserver(launcherRef, update, { attributes: true });
</script>

<style>
google-cast-launcher {
  --disconnected-color: rgb(var(--v-theme-on-surface));
  --connected-color: rgb(var(--v-theme-primary));
}
</style>
