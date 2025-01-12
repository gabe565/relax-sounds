<template>
  <page-layout>
    <template #actions>
      <template v-if="isMobile">
        <save-preset button />
        <cast-icon button />
      </template>
    </template>

    <v-overlay v-if="loading" v-model="loading" class="align-center justify-center" persistent>
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
          <v-alert prominent type="info">
            <template #prepend>
              <v-icon :icon="InfoIcon" />
            </template>
            No sounds are playing
          </v-alert>
        </v-col>
      </v-fade-transition>
    </v-row>
  </page-layout>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useToast } from "vue-toastification";
import { useDisplay } from "vuetify";
import InfoIcon from "~icons/material-symbols/info-rounded";
import MixerCard from "@/components/Mixer/MixerCard.vue";
import CastIcon from "@/components/NavButtons/CastIcon.vue";
import SavePreset from "@/components/NavButtons/SavePreset.vue";
import PageLayout from "@/layouts/PageLayout.vue";
import { usePlayerStore } from "@/plugins/store/player";

const { smAndDown: isMobile } = useDisplay();
const player = usePlayerStore();
const loading = ref(true);
const toast = useToast();

onMounted(async () => {
  try {
    await usePlayerStore().initSounds();
  } catch (err) {
    console.error(err);
    toast.error(`Failed to fetch sounds:\n${err}`);
  } finally {
    loading.value = false;
  }
});
</script>
