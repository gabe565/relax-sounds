<template>
  <v-container>

    <v-row>
      <v-col class="grow">
        <v-text-field
          v-model="filter.word"
          label="Filter"
          prepend-inner-icon="mdi-magnify"
          clearable
          autofocus
        >
        </v-text-field>
      </v-col>
      <v-col class="shrink">
        <v-switch v-model="filter.playing" label="Playing"/>
      </v-col>
    </v-row>

    <v-row>
      <v-col
        cols="12"
        sm="6"
        md="4"
        v-for="sound of filteredSounds"
        :key="sound.id"
      >
        <Sound :sound="sound" @play="playPaused"/>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import Fuse from 'fuse.js';

import Sound from './Sound.vue';

export default {
  name: 'SoundsPage',
  components: { Sound },
  props: {
    sounds: {
      type: Array,
    },
  },

  data: () => ({
    filter: {
      word: '',
      playing: false,
    },
    fuse: null,
  }),

  created() {
    this.fuse = new Fuse(this.sounds, {
      shouldSort: true,
      threshold: 0.3,
      location: 0,
      distance: 100,
      maxPatternLength: 32,
      minMatchCharLength: 1,
      keys: [
        'name',
        'tags',
      ],
    });
  },

  computed: {
    filteredSounds() {
      if (this.sounds) {
        let result;
        if (this.filter.word) {
          result = this.fuse.search(this.filter.word);
        } else {
          result = this.sounds;
        }
        if (this.filter.playing) {
          result = result.filter((e) => e.state !== 'stopped');
        }
        return result;
      }
      return [];
    },
  },

  methods: {
    playPaused() {
      this.sounds.forEach((sound) => {
        if (sound.state === 'paused') {
          sound.state = 'playing';
          sound.player.play();
        }
      });
    },
  },

  beforeDestroy() {
    this.filter.word = '';
    this.sounds.forEach((e) => {
      e.player.unload();
      e.state = 'stopped';
    });
  },
};
</script>

<style scoped>

</style>
