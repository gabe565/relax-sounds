<template>
  <div>
    <v-dialog
      v-model="show"
      max-width="500"
    >
      <v-card>
        <v-card-title class="text-h5">
          Restore
        </v-card-title>
        <v-card-text>
          If you previously backed up your presets, you can restore them here.
        </v-card-text>
        <v-card-text>
          <v-file-input
            outlined
            dense
            accept="application/json"
            :error="error"
            @change="file = $event"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn
            text
            @click="show = false"
          >
            <v-icon aria-hidden="true">
              fal fa-times fa-fw
            </v-icon>
            Close
          </v-btn>
          <v-btn
            text
            :disabled="!file"
            @click="restore"
          >
            <v-icon aria-hidden="true">
              fal fa-file-upload fa-fw
            </v-icon>
            Import
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-snackbar
      v-model="showSnackbar"
      timeout="5000"
      bottom
      class="pb-14 pb-md-0"
    >
      Imported {{ imported }} preset{{ imported !== 1 ? 's' : '' }}.
    </v-snackbar>
  </div>
</template>

<script>
export default {
  name: 'RestorePresets',

  props: {
    value: Boolean,
  },

  data: () => ({
    show: false,
    file: null,
    error: false,
    showSnackbar: false,
    imported: 0,
  }),

  watch: {
    value(val) {
      this.show = val;
    },
    show(val) {
      this.$emit('input', val);
    },
  },

  methods: {
    async restore() {
      try {
        const presets = JSON.parse(await this.file.text());
        presets.forEach((preset) => {
          this.$store.commit('presets/add', { preset, playing: false });
          this.show = false;
        });
        this.imported = presets.length;
        this.showSnackbar = true;
      } catch (error) {
        console.error(error);
        this.error = true;
      }
    },
  },
};
</script>

<style scoped>

</style>
