<template>
  <page-layout>
    <template #menu>
      <v-list-item title="Preload All" :prepend-icon="PreloadAllIcon" @click="player.prefetch" />
    </template>

    <filter-section />

    <v-overlay v-model="isLoading" class="self-center justify-center" persistent>
      <v-progress-circular color="primary" indeterminate size="64" />
    </v-overlay>

    <v-row>
      <v-fade-transition group leave-absolute hide-on-leave>
        <v-col
          v-for="sound of filters.filteredSounds"
          :key="sound.id"
          cols="12"
          sm="6"
          md="4"
          lg="3"
        >
          <sound-card :sound="sound" />
        </v-col>
      </v-fade-transition>
    </v-row>
  </page-layout>
</template>

<script setup>
import { ref } from "vue";
import { toast } from "vue-sonner";
import PreloadAllIcon from "~icons/material-symbols/cloud-sync-rounded";
import FilterSection from "@/components/Sounds/FilterSection.vue";
import SoundCard from "@/components/Sounds/SoundCard.vue";
import PageLayout from "@/layouts/PageLayout.vue";
import { useFiltersStore } from "@/plugins/store/filters";
import { usePlayerStore } from "@/plugins/store/player";
import { getErrorMessage } from "@/plugins/store/pocketbase.js";

const player = usePlayerStore();
const filters = useFiltersStore();
const isLoading = ref(true);

(async () => {
  try {
    await player.loadSounds();
  } catch (err) {
    console.error(err);
    toast.error(`Failed to fetch sounds:\n${getErrorMessage(err)}`);
  } finally {
    isLoading.value = false;
  }
})();
</script>
