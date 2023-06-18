<template>
  <v-col class="flex-grow-0">
    <v-btn elevation="0" icon variant="plain" aria-label="Debug" @click.stop="show = true">
      <v-icon :icon="DebugIcon" aria-hidden="true" />
    </v-btn>

    <v-dialog v-model="show" max-width="400">
      <v-card>
        <v-card-title class="text-h5">Debug</v-card-title>
        <v-card-text>
          <v-btn :href="preset.mixUrlAs('mp3')" target="_blank" class="mr-2">Mix MP3</v-btn>
          <v-btn :href="preset.mixUrlAs('wav')" target="_blank" class="mr-2">Mix WAV</v-btn>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="show = false">
            <v-icon :icon="CloseIcon" aria-hidden="true" />
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-col>
</template>

<script setup>
import { ref } from "vue";
import DebugIcon from "~icons/solar/bug-minimalistic-bold";
import CloseIcon from "~icons/solar/close-circle-bold";
import { Preset } from "../../util/Preset";

defineProps({
  preset: {
    type: Preset,
    required: true,
  },
});

const show = ref(false);
</script>
