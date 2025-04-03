<template>
  <page-layout>
    <template #actions>
      <template v-if="isMobile">
        <save-preset button />
        <cast-icon button />
      </template>
    </template>

    <v-overlay v-model="isLoading" class="align-center justify-center" persistent>
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
          <v-alert prominent type="info" :icon="InfoIcon">No sounds are playing</v-alert>
        </v-col>
      </v-fade-transition>
    </v-row>
  </page-layout>
</template>

<script setup>
import { useAsyncState } from "@vueuse/core";
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
const toast = useToast();

const { isLoading } = useAsyncState(player.initSounds, undefined, {
  onError(e) {
    toast.error(`Failed to fetch sounds:\n${e}`);
  },
});
</script>
