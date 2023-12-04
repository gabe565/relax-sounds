<template>
  <PageLayout>
    <v-row>
      <v-fade-transition group leave-absolute>
        <v-col
          v-for="sound of player.soundsNotStopped"
          :key="sound.id"
          cols="12"
          sm="6"
          md="4"
          xl="3"
        >
          <MixerCard :sound="sound" />
        </v-col>
      </v-fade-transition>
      <v-col v-if="player.soundsNotStopped.length === 0">
        <v-alert prominent text type="info">No sounds are playing</v-alert>
      </v-col>
    </v-row>
  </PageLayout>
</template>

<script setup>
import PageLayout from "../layouts/PageLayout.vue";
import MixerCard from "../components/Mixer/MixerCard.vue";
import { usePlayerStore } from "../plugins/store/player";

const player = usePlayerStore();
</script>
