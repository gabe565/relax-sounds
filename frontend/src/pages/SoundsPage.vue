<template>
  <PageLayout :alert="alert" :actions="actions">
    <FilterSection>
      <v-row>
        <template v-if="loading">
          <v-overlay v-model="loading" class="align-center justify-center" persistent>
            <v-progress-circular color="primary" indeterminate size="64" />
          </v-overlay>
        </template>
        <template v-else>
          <v-col v-for="(sound, key) of filters.sounds" :key="key" cols="12" sm="6" md="4" lg="3">
            <SoundCard :sound="sound" />
          </v-col>
        </template>
      </v-row>
    </FilterSection>
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

defineProps({
  alert: {
    type: Object,
    default: null,
  },
});

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
  await usePlayerStore().initSounds();
  loading.value = false;
});
</script>
