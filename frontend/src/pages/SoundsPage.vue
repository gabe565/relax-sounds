<template>
  <page-layout>
    <template #actions>
      <template v-if="isMobile">
        <save-preset button />
        <cast-icon button />
      </template>
    </template>
    <template #menu>
      <v-list-item title="Preload All" :prepend-icon="PreloadAllIcon" @click="player.prefetch" />
    </template>

    <filter-section />

    <v-overlay v-model="isLoading" class="align-center justify-center" persistent>
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
import { useAsyncState } from "@vueuse/core";
import { useToast } from "vue-toastification";
import { useDisplay } from "vuetify";
import PreloadAllIcon from "~icons/material-symbols/cloud-sync-rounded";
import CastIcon from "@/components/NavButtons/CastIcon.vue";
import SavePreset from "@/components/NavButtons/SavePreset.vue";
import FilterSection from "@/components/Sounds/FilterSection.vue";
import SoundCard from "@/components/Sounds/SoundCard.vue";
import PageLayout from "@/layouts/PageLayout.vue";
import { useFiltersStore } from "@/plugins/store/filters";
import { usePlayerStore } from "@/plugins/store/player";

const { smAndDown: isMobile } = useDisplay();
const player = usePlayerStore();
const toast = useToast();
const filters = useFiltersStore();

const { isLoading } = useAsyncState(player.initSounds, undefined, {
  onError(e) {
    toast.error(`Failed to fetch sounds:\n${e}`);
  },
});
</script>
