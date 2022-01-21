<template>
  <Page :alert="alert">
    <v-row>
      <v-col
        cols="12" lg="6"
        v-for="(preset, key) of presets"
        :key="key"
      >
        <Preset :preset="preset"/>
      </v-col>
      <v-col v-if="presets.length === 0">
        <v-alert prominent text color="info" icon="fal fa-info-circle">
          No Presets Saved Yet!
        </v-alert>
      </v-col>
    </v-row>
  </Page>
</template>

<script>
import { mapActions, mapState } from 'vuex';
import Preset from '../components/Preset.vue';
import Page from '../layouts/Page.vue';

export default {
  name: 'Presets',
  components: { Page, Preset },
  props: {
    alert: Object,
  },
  async created() {
    await this.initSounds();
  },
  computed: mapState('presets', [
    'presets',
  ]),
  methods: mapActions('player', ['initSounds']),
};
</script>
