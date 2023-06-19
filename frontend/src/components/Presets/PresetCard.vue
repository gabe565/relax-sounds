<template>
  <v-fade-transition>
    <v-card
      :variant="preset.new ? 'flat' : 'outlined'"
      :dark="preset.new"
      :color="preset.new ? 'deep-purple darken-2' : ''"
      transition="fade-transition"
    >
      <v-row align="center" dense>
        <v-col class="overflow-hidden">
          <v-card-title class="text-h5 pa-4 d-block text-truncate">
            {{ preset.name }}
          </v-card-title>
        </v-col>
        <debug-button v-if="debugEnabled" :preset="preset" />
        <share-button :preset="preset" />
        <delete-button :preset="preset" />
        <play-button :preset="preset" />
      </v-row>
    </v-card>
  </v-fade-transition>
</template>

<script setup>
import ShareButton from "./Buttons/ShareButton.vue";
import DeleteButton from "./Buttons/DeleteButton.vue";
import DebugButton from "./Buttons/DebugButton.vue";
import PlayButton from "./Buttons/PlayButton.vue";

defineProps({
  preset: {
    type: Object,
    required: true,
  },
});

const debugEnabled = import.meta.env.DEV;
</script>
