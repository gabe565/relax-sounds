<template>
  <page-layout>
    <template #menu>
      <v-list-item title="Preload All" :prepend-icon="PreloadAllIcon" @click="player.prefetch" />
    </template>

    <filter-section />

    <v-overlay v-model="isLoading" class="self-center justify-center" persistent>
      <v-progress-circular color="primary" indeterminate size="64" />
    </v-overlay>

    <div
      class="card-group card-cols-4 mt-6 grid grid-cols-1 gap-0.5 sm:gap-2 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4"
    >
      <v-fade-transition group leave-absolute hide-on-leave>
        <div v-for="sound of filters.filteredSounds" :key="sound.id">
          <sound-card :sound="sound" />
        </div>
      </v-fade-transition>
    </div>
  </page-layout>
</template>

<script setup>
import { ref } from "vue";
import { toast } from "vue-sonner";
import PreloadAllIcon from "~icons/material-symbols/cloud-sync-rounded";
import FilterSection from "@/components/Sounds/FilterSection.vue";
import SoundCard from "@/components/Sounds/SoundCard.vue";
import PageLayout from "@/layouts/PageLayout.vue";
import { useFilters } from "@/plugins/store/filters";
import { usePlayer } from "@/plugins/store/player";
import { getErrorMessage, usePocketBase } from "@/plugins/store/pocketbase.js";

const pb = usePocketBase();
const player = usePlayer();
const filters = useFilters();
const isLoading = ref(true);

(async () => {
  try {
    await pb.loadSounds();
  } catch (err) {
    console.error(err);
    toast.error(`Failed to fetch sounds:\n${getErrorMessage(err)}`);
  } finally {
    isLoading.value = false;
  }
})();
</script>
