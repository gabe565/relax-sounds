<template>
  <template v-if="loading">
    <slot />
  </template>
  <template v-else>
    <v-row class="filters">
      <v-col class="pb-0">
        <v-text-field
          v-model="filters.word"
          label="Search"
          prepend-icon="fal fa-search"
          clearable
          variant="underlined"
          hide-details
        />
      </v-col>
      <v-col class="flex-grow-0 pb-0">
        <v-switch
          v-model="filters.playing"
          label="Playing"
          inset
          hide-details
        />
      </v-col>
    </v-row>
    <v-row class="pb-5">
      <v-chip-group
        v-model="filters.word"
        column
      >
        <v-chip
          v-for="(tag, key) in tags"
          :key="key"
          :value="tag.id"
          variant="outlined"
          active-class="deep-orange"
          class="ma-2"
          filter
        >
          <v-icon
            v-if="tag.icon"
            size="x-small"
            class="mr-2"
          >
            far {{ tag.icon }}
          </v-icon>
          {{ tag.name }}
        </v-chip>
      </v-chip-group>
    </v-row>
    <v-row>
      <v-col>
        <v-divider />
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-pagination
          v-model="filters.page"
          :length="pages"
          variant="outlined"
          size="small"
          active-color="primary"
        />
      </v-col>
    </v-row>
    <slot />
    <v-row>
      <v-col>
        <v-pagination
          v-model="filters.page"
          :length="pages"
          variant="outlined"
          size="small"
          active-color="primary"
        />
      </v-col>
    </v-row>
  </template>
</template>

<script>
import { mapState, mapGetters, mapActions } from 'vuex';
import { getTags } from '../data/tags';

export default {
  name: 'FilterSection',

  data: () => ({
    tags: null,
    loading: true,
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

  async created() {
    this.tags = await getTags();
    this.loading = false;
  },

  methods: mapActions('filters', ['initSounds']),
};
</script>

<style scoped>
.filters :deep(.v-input__control) {
  grid-area: auto;
}
</style>
