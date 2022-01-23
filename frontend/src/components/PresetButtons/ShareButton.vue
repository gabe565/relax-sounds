<template>
  <v-col class="shrink">
    <v-btn @click.stop="show = true" elevation="0" icon aria-label="Share">
      <v-icon dense aria-hidden="true">
        fas fa-fw fa-share
      </v-icon>
    </v-btn>

    <v-dialog v-model="show" max-width="500">
      <v-card>
        <v-card-title class="headline">Share Preset</v-card-title>
        <v-card-text>
          <v-text-field :value="shareUrl"/>
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
import { encode } from '../../util/shareUrl';

export default {
  name: 'ShareButton',

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
    shareUrl() {
      const { name, sounds } = encode(this.preset);
      return `${window.location.origin}/import/${name}/${sounds}`;
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
