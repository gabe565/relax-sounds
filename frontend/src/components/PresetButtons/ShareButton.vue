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
          <v-text-field
            readonly
            :value="shareUrl"
            @focus="select($event.target)"
            @click="select($event.target)"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer/>
          <v-btn text @click="show = false">
            <v-icon aria-hidden="true">fal fa-times fa-fw</v-icon>
            Close
          </v-btn>
          <v-btn text @click="copy">
            <v-icon aria-hidden="true">fal fa-copy fa-fw</v-icon>
            Copy
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-snackbar
      v-model="showSnackbar" timeout="5000"
      bottom class="pb-14 pb-md-0"
    >
      Copied to clipboard.
    </v-snackbar>
  </v-col>
</template>

<script>
import { encode } from '../../util/shareUrl';
import { wait } from '../../util/helpers';

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
    showSnackbar: false,
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

  methods: {
    async select(e) {
      await wait(0);
      e.select();
      e.scrollLeft = 0;
    },
    async copy() {
      await navigator.clipboard.writeText(this.shareUrl);
      if (this.showSnackbar) {
        this.showSnackbar = false;
        await this.$nextTick();
      }
      this.showSnackbar = true;
    },
  },
};
</script>
