<template>
  <v-tooltip bottom>
    <template #activator="{ on, attrs }">
      <v-btn
        v-bind="attrs"
        v-on="on"
        @click="showDialog = !showDialog"
        :disabled="isStopped"
        icon
        aria-label="Save Preset"
      >
        <v-icon aria-hidden="true">fas fa-plus-circle</v-icon>
        <v-dialog v-model="showDialog" max-width="500">
          <v-card>
            <v-card-title class="headline">Preset Name</v-card-title>
            <v-card-text>
              <v-form @submit.prevent="save">
                <v-text-field required autofocus label="Name" v-model="name"/>
              </v-form>
            </v-card-text>
            <v-card-actions>
              <v-spacer/>
              <v-btn text @click="cancel">
                <v-icon aria-hidden="true">fal fa-times fa-fw</v-icon>
                Cancel
              </v-btn>
              <v-btn color="green" text @click="save">
                <v-icon aria-hidden="true">fal fa-save fa-fw</v-icon>
                Save Preset
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-btn>
    </template>
    <span>Save Preset</span>
  </v-tooltip>
</template>

<script>
import { mapGetters } from 'vuex';
import { wait } from '../util/helpers';

export default {
  name: 'SavePreset',

  data: () => ({
    showDialog: false,
    name: '',
  }),

  computed: mapGetters('player', [
    'isStopped',
  ]),

  methods: {
    cancel() {
      this.showDialog = false;
      this.name = '';
    },
    async save() {
      this.showDialog = false;
      await wait(300);

      let params;
      try {
        this.$store.dispatch('presets/savePlaying', { name: this.name });
        params = { alert: { type: 'info', text: `Preset "${this.name}" saved successfully.` } };
        this.name = '';
      } catch (error) {
        console.error(error);
        params = { alert: { type: 'error', text: 'Failed to save preset. Please try again later.' } };
      }
      return this.$router.push({ name: 'Presets', params });
    },
  },
};
</script>
