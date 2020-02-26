<template>
  <Page>
    <Filters>
      <v-row>
        <v-col
          cols="12" sm="6" md="4"
          v-for="sound of sounds"
          :key="sound.id"
        >
          <Sound :sound="sound"/>
        </v-col>
      </v-row>
    </Filters>
  </Page>
</template>

<script>
import { mapGetters } from 'vuex';
import Page from './Page.vue';
import Sound from '../Sound.vue';
import Filters from '../Filters.vue';

export default {
  name: 'Sounds',
  components: { Filters, Page, Sound },
  data: () => ({
    page: 1,
  }),
  computed: mapGetters('filters', ['sounds']),
};
</script>

<style scoped>

</style>
