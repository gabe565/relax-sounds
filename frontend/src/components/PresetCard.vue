<template>
  <v-fade-transition>
    <v-card
      flat
      outlined
      :dark="preset.new"
      :color="preset.new ? 'deep-purple darken-2' : ''"
      transition="fade-transition"
    >
      <v-row
        align="center"
        dense
      >
        <v-col class="overflow-hidden">
          <v-card-title class="text-h5 d-block text-truncate">
            {{ preset.name }}
          </v-card-title>
        </v-col>
        <debug-button
          v-if="debugEnabled"
          :preset="preset"
        />
        <share-button :preset="preset" />
        <delete-button :preset="preset" />
        <play-button :preset="preset" />
      </v-row>
    </v-card>
  </v-fade-transition>
</template>

<script>
import ShareButton from './PresetButtons/ShareButton.vue';
import DeleteButton from './PresetButtons/DeleteButton.vue';
import DebugButton from './PresetButtons/DebugButton.vue';
import PlayButton from './PresetButtons/PlayButton.vue';

export default {
  name: 'PresetCard',

  components: {
    PlayButton,
    DebugButton,
    DeleteButton,
    ShareButton,
  },

  props: {
    preset: {
      type: Object,
      required: true,
    },
  },

  data: () => ({
    deleteDialog: false,
    debugEnabled: import.meta.env.DEV,
    debugDialog: false,
  }),
};
</script>
