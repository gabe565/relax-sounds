<template>
  <PageLayout :actions="actions">
    <filter-section />

    <v-row>
      <template v-if="loading">
        <v-overlay v-model="loading" class="align-center justify-center" persistent>
          <v-progress-circular color="primary" indeterminate size="64" />
        </v-overlay>
      </template>
      <template v-else>
        <v-fade-transition group leave-absolute hide-on-leave>
          <v-col
            v-for="sound of filters.filteredSounds"
            :key="sound.id"
            cols="12"
            sm="6"
            md="4"
            lg="3"
          >
            <SoundCard :sound="sound" />
          </v-col>
        </v-fade-transition>
      </template>
    </v-row>
  </PageLayout>
</template>

<script setup>
import { onMounted, ref } from "vue";
import PreloadAllIcon from "~icons/material-symbols/cloud-sync-rounded";
import PageLayout from "../layouts/PageLayout.vue";
import SoundCard from "../components/Sounds/SoundCard.vue";
import FilterSection from "../components/Sounds/FilterSection.vue";
import { usePlayerStore } from "../plugins/store/player";
import { useFiltersStore } from "../plugins/store/filters";
import { toast } from "vue3-toastify";

const loading = ref(true);

const actions = [
  {
    title: "Preload All",
    icon: PreloadAllIcon,
    on: {
      click: () => {
        usePlayerStore().prefetch();
      },
    },
  },
];

const filters = useFiltersStore();

onMounted(async () => {
  try {
    await usePlayerStore().initSounds();
  } catch (err) {
    console.error(err);
    toast.error("Failed to fetch sounds.");
  } finally {
    loading.value = false;
  }
});
</script>
