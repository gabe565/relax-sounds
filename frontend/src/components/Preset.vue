<template>
  <v-fade-transition>
    <v-card flat outlined :dark="preset.new"
            :color="preset.new ? 'deep-purple darken-2' : ''"
            transition="fade-transition"
    >
      <v-row align="center" justify="center" dense>
        <v-col class="grow">
          <v-card-title class="headline">
            {{ preset.name }}
          </v-card-title>
        </v-col>
        <v-col class="shrink" v-if="DEBUG_ENABLED">
          <v-btn @click.stop="debugDialog = true" elevation="0" icon aria-label="Debug">
            <v-icon dense aria-hidden="true">
              fas fa-fw fa-bug
            </v-icon>
          </v-btn>
        </v-col>
        <v-col class="shrink">
          <v-btn @click.stop="shareDialog = true" elevation="0" icon aria-label="Share">
            <v-icon dense aria-hidden="true">
              fas fa-fw fa-share
            </v-icon>
          </v-btn>
        </v-col>
        <v-col class="shrink">
          <v-btn @click.stop="deleteDialog = true" elevation="0" icon aria-label="Delete">
            <v-icon dense aria-hidden="true">
              fas fa-fw fa-trash
            </v-icon>
          </v-btn>
        </v-col>
        <v-col class="shrink pr-4">
          <v-btn @click.stop="play" elevation="0" icon aria-label="Play">
            <v-icon dense aria-hidden="true">
              fas fa-fw fa-play
            </v-icon>
          </v-btn>
        </v-col>
      </v-row>

      <v-dialog v-model="debugDialog" max-width="500">
        <v-card>
          <v-card-title class="headline">Debug</v-card-title>
          <v-card-text>
            <v-btn :href="downloadUrl" target="_blank">Mix URL</v-btn>
          </v-card-text>
          <v-card-actions>
            <v-spacer/>
            <v-btn text @click="debugDialog = false">
              <v-icon aria-hidden="true">fal fa-times fa-fw</v-icon>
              Close
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

      <v-dialog v-model="shareDialog" max-width="500">
        <v-card>
          <v-card-title class="headline">Share Preset</v-card-title>
          <v-card-text>
            <v-text-field :value="shareUrl"/>
          </v-card-text>
          <v-card-actions>
            <v-spacer/>
            <v-btn text @click="shareDialog = false">
              <v-icon aria-hidden="true">fal fa-times fa-fw</v-icon>
              Close
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

      <v-dialog v-model="deleteDialog" max-width="500">
        <v-card>
          <v-card-title class="headline">Delete Preset?</v-card-title>
          <v-card-text>
            Are you sure you want to delete "{{ preset.name }}"?
          </v-card-text>
          <v-card-actions>
            <v-spacer/>
            <v-btn text @click="deleteDialog = false">
              <v-icon aria-hidden="true">fal fa-times fa-fw</v-icon>
              Cancel
            </v-btn>
            <v-btn color="red" text @click="remove">
              <v-icon aria-hidden="true">fal fa-trash fa-fw</v-icon>
              Confirm
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-card>
  </v-fade-transition>
</template>

<script>
import { encode, encodeSounds } from '../util/shareUrl';

export default {
  name: 'Preset',

  props: {
    preset: {
      type: Object,
    },
  },

  data: () => ({
    deleteDialog: false,
    shareDialog: false,
    debugDialog: false,
  }),

  created() {
    this.DEBUG_ENABLED = process.env.NODE_ENV === 'development';
  },

  computed: {
    shareUrl() {
      const { name, sounds } = encode(this.preset);
      return `${window.location.origin}/import/${name}/${sounds}`;
    },
    downloadUrl() {
      const sounds = encodeSounds(this.preset.sounds);
      return `${window.location.origin}/api/mix/${sounds}`;
    },
  },

  methods: {
    play() {
      this.$store.dispatch('presets/play', { preset: this.preset });
    },
    remove() {
      this.$store.commit('presets/remove', { preset: this.preset });
    },
  },
};
</script>

<style scoped>
  .v-card {
    overflow: hidden;
    transition: background-color 0.5s cubic-bezier(0.215, 0.61, 0.355, 1);
  }
  .v-progress-linear {
    z-index: 0;
    left: 0;
    right: 0;
    height: 100%;
    width: 100%;
  }
  .row {
    position: relative;
    z-index: 1;
    pointer-events: none;
  }
  .v-btn {
    pointer-events: all;
  }
  .no-transition {
    transition: none;
  }
</style>
