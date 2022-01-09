<template>
  <div>
    <v-row>
      <v-col class="grow pb-0">
        <v-text-field
          v-model="filters.word"
          label="Search"
          prepend-inner-icon="fal fa-search v-icon--dense"
          clearable
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
      <v-col>
        <v-divider/>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-pagination v-model="filters.page" :length="pages"/>
      </v-col>
    </v-row>
    <slot/>
    <v-row>
      <v-col>
        <v-pagination v-model="filters.page" :length="pages"/>
      </v-col>
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

  watch: {
    pages() {
      this.filters.page = 1;
    },
  },
};
</script>
