<template>
  <div>
    <v-row>
      <v-col class="grow">
        <v-text-field
          v-model="filters.word"
          label="Filter"
          prepend-inner-icon="fal fa-search v-icon--dense"
          clearable
          autofocus
        />
      </v-col>
      <v-col class="shrink">
        <v-switch v-model="filters.playing" label="Playing" flat inset/>
      </v-col>
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

export default {
  name: 'Filters',

  computed: {
    ...mapState('filters', ['filters']),
    ...mapGetters('filters', ['pages']),
  },
};
</script>

<style scoped>

</style>
