<template>
  <Page :alert="alert" :actions="actions">
    <Filters>
      <v-row>
        <template v-if="loading">
          <v-col
            cols="12" md="6" lg="4"
            v-for="i in 24"
            :key="i"
          >
            <v-card flat outlined>
              <v-skeleton-loader
                type="card-heading"
                class="my-1 transparent"
              />
            </v-card>
          </v-col>
        </template>
        <template v-else>
          <v-col
            cols="12" md="6" lg="4"
            v-for="(sound, key) of sounds"
            :key="key"
          >
            <Sound :sound="sound"/>
          </v-col>
        </template>
      </v-row>
    </Filters>
  </Page>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import Page from '../layouts/Page.vue';
import Sound from '../components/Sound.vue';
import Filters from '../components/Filters.vue';
import { prefetch } from '../data/sounds';

export default {
  name: 'Sounds',
  components: { Filters, Page, Sound },
  props: {
    alert: Object,
  },
  data: () => ({
    loading: true,
    error: null,
    page: 1,
  }),
  async created() {
    await this.initSounds();
    this.loading = false;
  },
  computed: {
    actions() {
      return [
        { title: 'Preload All', icon: 'fas fa-sync', on: { click: prefetch } },
      ];
    },
    ...mapGetters('filters', ['sounds']),
  },
  methods: mapActions('player', ['initSounds']),
};
</script>
