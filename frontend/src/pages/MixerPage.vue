<template>
  <page-layout>
    <template #actions>
      <template v-if="isMobile">
        <save-preset button />
        <cast-icon button />
      </template>
    </template>

    <v-row>
      <template v-if="loading">
        <v-overlay v-model="loading" class="align-center justify-center" persistent>
          <v-progress-circular color="primary" indeterminate size="64" />
        </v-overlay>
      </template>
      <template v-else>
        <v-fade-transition group leave-absolute>
          <v-col
            v-for="sound of player.soundsNotStopped"
            :key="sound.id"
            cols="12"
            sm="6"
            md="4"
            xl="3"
          >
            <mixer-card :sound="sound" />
          </v-col>
          <v-col v-if="player.soundsNotStopped.length === 0">
            <v-alert prominent text type="info">No sounds are playing</v-alert>
          </v-col>
        </v-fade-transition>
      </template>
    </v-row>
  </page-layout>
</template>

<script setup>
import PageLayout from "../layouts/PageLayout.vue";
import MixerCard from "../components/Mixer/MixerCard.vue";
import { usePlayerStore } from "../plugins/store/player";
import { onMounted, ref } from "vue";
import { useToast } from "vue-toastification";
import SavePreset from "../components/NavButtons/SavePreset.vue";
import CastIcon from "../components/NavButtons/CastIcon.vue";
import { useDisplay } from "vuetify";

const { smAndDown: isMobile } = useDisplay();
const player = usePlayerStore();
const loading = ref(true);
const toast = useToast();

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
