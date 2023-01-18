<template>
  <PageLayout
    :alert="alert"
    :actions="actions"
  >
    <v-row>
      <v-col
        v-for="(preset, key) of presets"
        :key="key"
        cols="12"
        md="6"
        xl="4"
      >
        <PresetCard :preset="preset" />
      </v-col>
      <v-col v-if="presets.length === 0">
        <v-alert
          prominent
          text
          color="info"
          icon="fal fa-info-circle"
        >
          No Presets Saved Yet!
        </v-alert>
      </v-col>
    </v-row>

    <restore-presets v-model="showRestore" />
    <remove-all v-model="showRemoveAll" />
  </PageLayout>
</template>

<script>
import { mapActions, mapState } from 'vuex';
import { saveAs } from 'file-saver/src/FileSaver';
import PresetCard from '../components/PresetCard.vue';
import PageLayout from '../layouts/PageLayout.vue';
import RestorePresets from '../components/RestorePresets.vue';
import RemoveAll from '../components/RemoveAll.vue';

export default {
  name: 'PresetsPage',

  components: {
    RemoveAll,
    RestorePresets,
    PageLayout,
    PresetCard,
  },

  props: {
    alert: {
      type: Object,
      default: null,
    },
  },

  data: () => ({
    showRestore: false,
    showRemoveAll: false,
  }),

  computed: {
    actions() {
      return [
        {
          title: 'Backup',
          icon: 'fas fa-file-download',
          on: {
            click: this.export,
          },
        },
        {
          title: 'Restore',
          icon: 'fas fa-file-upload',
          on: {
            click: () => { this.showRestore = true; },
          },
        },
        {
          title: 'Remove All',
          icon: 'fas fa-trash',
          on: {
            click: () => { this.showRemoveAll = true; },
          },
        },
      ];
    },
    ...mapState('presets', ['presets']),
  },

  async created() {
    await this.initSounds();
  },

  methods: {
    export() {
      const blob = new Blob(
        [JSON.stringify(this.presets)],
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
