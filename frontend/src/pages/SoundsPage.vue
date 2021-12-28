<template>
  <Page>
    <Filters>
      <v-row>
        <v-col
          cols="12" md="6" lg="4"
          v-for="(sound, key) of sounds"
          :key="key"
        >
          <Sound :sound="sound"/>
        </v-col>
      </v-row>
    </Filters>
  </Page>
</template>

<script>
import { mapGetters } from 'vuex';
import Page from '../layouts/Page.vue';
import Sound from '../components/Sound.vue';
import Filters from '../components/Filters.vue';

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
