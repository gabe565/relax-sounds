<template>
  <v-row>
    <v-col class="grow">
      <v-text-field
        v-model="word"
        label="Filter"
        prepend-inner-icon="fa-search v-icon--dense"
        clearable
        autofocus
      >
      </v-text-field>
    </v-col>
    <v-col class="shrink">
      <v-switch v-model="playing" label="Playing"/>
    </v-col>
  </v-row>
</template>

<script>
export default {
  name: 'Filters',

  computed: {
    filteredSounds() {
      if (this.$store.state.sounds.sounds) {
        let result;
        if (this.filter.word) {
          result = this.fuse.search(this.filter.word);
        } else {
          result = this.$store.state.sounds.sounds;
        }
        if (this.filter.playing) {
          result = result.filter((e) => e.state !== 'stopped');
        }
        return result;
      }
      return [];
    },
    word: {
      get() {
        return this.$store.getters['filters/word'];
      },
      set(value) {
        this.$store.commit('filters/byWord', value);
      },
    },
    playing: {
      get() {
        return this.$store.getters['filters/playing'];
      },
      set() {
        this.$store.commit('filters/byPlaying');
      },
    },
  },
};
</script>

<style scoped>

</style>
