<template>
  <v-card flat
          outlined
          :color="playlist.new ? 'deep-purple darken-2' : ''"
          transition="fade-transition"
  >
    <v-row align="center" justify="center" dense>
      <v-col class="grow">
        <v-card-title class="headline">
          {{ playlist.name }}
        </v-card-title>
      </v-col>
      <v-col class="shrink">
        <v-btn @click.stop="shareDialog = true" elevation="0" icon>
          <v-icon dense>
            fas fa-fw fa-share
          </v-icon>
        </v-btn>
      </v-col>
      <v-col class="shrink">
        <v-btn @click.stop="deleteDialog = true" elevation="0" icon>
          <v-icon dense>
            fas fa-fw fa-trash
          </v-icon>
        </v-btn>
      </v-col>
      <v-col class="shrink pr-4">
        <v-btn @click.stop="play" elevation="0" icon>
          <v-icon dense>
            fas fa-fw fa-play
          </v-icon>
        </v-btn>
      </v-col>
    </v-row>

    <v-dialog v-model="shareDialog" max-width="500">
      <v-card>
        <v-card-title class="headline">Share Playlist</v-card-title>
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
        <v-card-title class="headline">Delete Playlist?</v-card-title>
        <v-card-text>
          Are you sure you want to delete "{{ playlist.name }}"?
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
</template>

<script>
import { encode } from '../snippets/shareUrl';

export default {
  name: 'Playlist',

  props: {
    playlist: {
      type: Object,
    },
  },

  data: () => ({
    deleteDialog: false,
    shareDialog: false,
  }),

  computed: {
    shareUrl() {
      const { name, sounds } = encode(this.playlist);
      return `${window.location.origin}/import/${name}/${sounds}`;
    },
    downloadUrl() {
      const { sounds } = encode(this.playlist);
      return `${window.location.origin}/mix/${sounds}`;
    },
  },

  methods: {
    play() {
      this.$store.dispatch('playlists/play', { playlist: this.playlist });
    },
    remove() {
      this.$store.commit('playlists/remove', { playlist: this.playlist });
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
