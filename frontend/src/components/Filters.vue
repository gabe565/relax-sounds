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
        <template v-if="loading">
          <v-skeleton-loader
            type="chip"
            v-for="i in 5"
            :key="i"
            class="ma-2"
          />
        </template>
        <template v-else>
          <v-chip v-for="(tag, key) in tags" :key="key"
                  :value="key"
                  outlined
                  active-class="deep-orange"
                  class="ma-2"
          >
            <v-icon v-if="tag.icon" x-small class="mr-2">
              far {{ tag.icon }}
            </v-icon>
            {{ tag.name }}
          </v-chip>
        </template>
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
import { mapState, mapGetters, mapActions } from 'vuex';
import { getTags } from '../data/tags';

export default {
  name: 'Filters',

  data: () => ({
    tags: null,
    loading: true,
  }),

  async created() {
    this.tags = await getTags();
    this.loading = false;
  },

  computed: {
    ...mapState('filters', ['filters']),
    ...mapGetters('filters', ['pages']),
  },

  methods: mapActions('filters', ['initSounds']),

  watch: {
    pages() {
      this.filters.page = 1;
    },
  },
};
</script>
