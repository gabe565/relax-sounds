<template>
  <v-container>

    <v-text-field
      v-model="filter"
      label="Filter"
      prepend-inner-icon="mdi-magnify"
      clearable
      autofocus
    />

    <v-layout row wrap>
      <v-col
        cols="12"
        sm="6"
        md="4"
        v-for="sound of filteredSounds"
        :key="sound.id"
      >
        <Sound :sound="sound"/>
      </v-col>
    </v-layout>
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
    filter: '',
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
      ],
    });
  },

  computed: {
    filteredSounds() {
      if (this.sounds) {
        if (this.filter) {
          return this.fuse.search(this.filter);
        }
        return this.sounds.sort( // eslint-disable-line
          (left, right) => left.name.localeCompare(right.name),
        );
      }
      return [];
    },
  },

  beforeDestroy() {
    this.filter = '';
    this.sounds.forEach((e) => {
      e.player.unload();
      e.playing = false;
    });
  },
};
</script>

<style scoped>

</style>
