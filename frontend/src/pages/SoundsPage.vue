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
          <v-col v-for="(sound, key) of filters.sounds" :key="key" cols="12" md="6" lg="4">
            <SoundCard :sound="sound" />
          </v-col>
        </template>
      </v-row>
    </FilterSection>
  </PageLayout>
</template>

<script setup>
import { onMounted, ref } from "vue";
import PageLayout from "../layouts/PageLayout.vue";
import SoundCard from "../components/SoundCard.vue";
import FilterSection from "../components/FilterSection.vue";
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
    icon: "fas fa-sync",
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
