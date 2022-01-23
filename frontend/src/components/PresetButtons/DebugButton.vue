<template>
  <v-col class="shrink">
    <v-btn @click.stop="show = true" elevation="0" icon aria-label="Debug">
      <v-icon dense aria-hidden="true">
        fas fa-fw fa-bug
      </v-icon>
    </v-btn>

    <v-dialog v-model="show" max-width="500">
      <v-card>
        <v-card-title class="headline">Debug</v-card-title>
        <v-card-text>
          <v-btn :href="downloadUrl" target="_blank">Mix URL</v-btn>
        </v-card-text>
        <v-card-actions>
          <v-spacer/>
          <v-btn text @click="show = false">
            <v-icon aria-hidden="true">fal fa-times fa-fw</v-icon>
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-col>
</template>

<script>
import { encodeSounds } from '../../util/shareUrl';

export default {
  name: 'DebugButton',

  props: {
    preset: {
      type: Object,
      required: true,
    },
  },

  data: () => ({
    show: false,
  }),

  computed: {
    downloadUrl() {
      const sounds = encodeSounds(this.preset.sounds);
      return `${window.location.origin}/api/mix/${sounds}`;
    },
  },

  watch: {
    value: {
      handler(val) {
        this.show = val;
      },
      immediate: true,
    },
    show(val) {
      this.$emit('input', val);
    },
  },
};
</script>
