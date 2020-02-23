<template>
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
      <v-switch v-model="filters.playing" label="Playing"/>
    </v-col>
  </v-row>
</template>

<script>
import { mapState } from 'vuex';

export default {
  name: 'Filters',

  computed: mapState('filters', [
    'filters',
  ]),
};
</script>

<style scoped>

</style>
