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
  }),

  computed: {
    filteredSounds() {
      if (this.sounds) {
        let result = this.sounds;
        if (this.filter) {
          result = result.filter(
            (sound) => sound.name.toLowerCase().includes(this.filter.toLowerCase()),
          );
        }
        return result.sort(
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
