<template>
  <Page :alert="alert" :actions="actions">
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
import { saveAs } from 'file-saver/src/FileSaver';
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
  computed: {
    actions() {
      return [
        { title: 'Backup', icon: 'fas fa-file-download', on: { click: this.export } },
      ];
    },
    ...mapState('presets', ['presets']),
  },
  methods: {
    export() {
      const { presets } = this.$store.state.presets;
      const blob = new Blob(
        [JSON.stringify(presets)],
        { type: 'application/json;charset=utf-8' },
      );
      const offset = (new Date()).getTimezoneOffset() * 60000; // Offset in milliseconds
      const localISOTime = (new Date(Date.now() - offset))
        .toISOString()
        .slice(0, -5) // Remove ".000Z"
        .replaceAll(':', '');
      saveAs(blob, `relax-sounds-presets-${localISOTime}.json`);
    },
    ...mapActions('player', ['initSounds']),
  },
};
</script>
