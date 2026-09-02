<template>
  <page-layout>
    <v-overlay v-model="isLoading" class="self-center justify-center" persistent>
      <v-progress-circular color="primary" indeterminate size="64" />
    </v-overlay>

    <v-row>
      <v-fade-transition group leave-absolute>
        <template v-if="player.soundsNotStopped.length !== 0">
          <v-col
            v-for="sound of player.soundsNotStopped"
            :key="sound.id"
            cols="12"
            sm="6"
            lg="4"
            xl="3"
          >
            <mixer-card :sound="sound" />
          </v-col>
        </template>
        <v-col v-else>
          <div class="text-on-surface-variant flex flex-col items-center gap-3 py-16 text-center">
            <v-icon :icon="EmptyIcon" size="48" class="opacity-60" />
            <span>No sounds are playing</span>
            <v-btn to="/sounds" variant="tonal" rounded="pill">Browse sounds</v-btn>
          </div>
        </v-col>
      </v-fade-transition>
    </v-row>
  </page-layout>
</template>

<script setup>
import { ref } from "vue";
import { toast } from "vue-sonner";
import EmptyIcon from "~icons/material-symbols/instant-mix-rounded";
import MixerCard from "@/components/Mixer/MixerCard.vue";
import PageLayout from "@/layouts/PageLayout.vue";
import { usePlayer } from "@/plugins/store/player";
import { getErrorMessage, usePocketBase } from "@/plugins/store/pocketbase.js";

const pb = usePocketBase();
const player = usePlayer();
const isLoading = ref(true);

(async () => {
  try {
    await pb.loadSounds();
  } catch (err) {
    console.error(err);
    toast.error(`Failed to load:\n${getErrorMessage(err)}`);
  } finally {
    isLoading.value = false;
  }
})();
</script>
