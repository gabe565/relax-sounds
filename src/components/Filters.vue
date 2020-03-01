<template>
  <div>
    <v-row>
      <v-col class="grow pb-0">
        <v-text-field
          v-model="filters.word"
          label="Filter"
          prepend-inner-icon="fal fa-search v-icon--dense"
          clearable
          autofocus
          single-line
        />
      </v-col>
      <v-col class="shrink pb-0">
        <v-switch v-model="filters.playing" label="Playing" flat inset/>
      </v-col>
    </v-row>
    <v-row class="pb-5">
      <v-chip-group v-model="filters.word">
        <v-chip v-for="(tag, key) in tags" :key="key"
                :value="key"
                outlined
                active-class="deep-orange"
                class="ma-2"
        >
          {{ tag.name }}
        </v-chip>
      </v-chip-group>
    </v-row>
    <v-row>
      <v-pagination v-model="filters.page" :length="pages"/>
    </v-row>
    <slot/>
    <v-row>
      <v-pagination v-model="filters.page" :length="pages"/>
    </v-row>
  </div>
</template>

<script>
import { mapState, mapGetters } from 'vuex';
import tags from '../data/tags.json';

export default {
  name: 'Filters',

  data: () => ({
    tags,
  }),

  computed: {
    ...mapState('filters', ['filters']),
    ...mapGetters('filters', ['pages']),
  },
};
</script>

<style scoped>

</style>
