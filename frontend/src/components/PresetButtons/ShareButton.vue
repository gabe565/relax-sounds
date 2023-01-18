<template>
  <v-col class="flex-grow-0">
    <v-btn elevation="0" icon variant="plain" aria-label="Share" @click.stop="show = true">
      <v-icon aria-hidden="true">fas fa-fw fa-share-alt</v-icon>
    </v-btn>

    <v-dialog v-model="show" width="400">
      <v-card>
        <v-card-title class="text-h5">Share</v-card-title>
        <v-card-text>
          <v-text-field
            readonly
            :model-value="preset.shareUrl"
            label="Share URL"
            @focus="select($event.target)"
            @click="select($event.target)"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="show = false">
            <v-icon aria-hidden="true">fal fa-times fa-fw</v-icon>
            Close
          </v-btn>
          <v-btn variant="text" @click="copy">
            <v-icon aria-hidden="true">fal fa-copy fa-fw</v-icon>
            Copy
          </v-btn>
          <v-btn v-if="canShare" variant="text" @click="share">
            <v-icon aria-hidden="true">fal fa-share-alt fa-fw</v-icon>
            Share
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-snackbar
      v-model="showSnackbar"
      timeout="5000"
      location="bottom"
      content-class="pb-14 pb-md-0"
    >
      Copied to clipboard.
    </v-snackbar>
  </v-col>
</template>

<script>
import { wait } from "../../util/helpers";
import { Preset } from "../../util/Preset";

export default {
  name: "ShareButton",

  props: {
    preset: {
      type: Preset,
      required: true,
    },
  },

  emits: ["update:modelValue"],

  data: () => ({
    show: false,
    showSnackbar: false,
  }),

  computed: {
    shareData() {
      return {
        title: "Relax Sounds",
        text: `Import my Relax Sounds preset called "${this.preset.name}"`,
        url: this.preset.shareUrl,
      };
    },
    canShare() {
      return navigator.canShare && navigator.canShare(this.shareData);
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
      this.$emit("update:modelValue", val);
    },
  },

  methods: {
    async select(e) {
      await wait(0);
      e.select();
      e.scrollLeft = 0;
    },
    async copy() {
      await navigator.clipboard.writeText(this.preset.shareUrl);
      if (this.showSnackbar) {
        this.showSnackbar = false;
        await this.$nextTick();
      }
      this.showSnackbar = true;
      this.show = false;
    },
    async share() {
      await navigator.share(this.shareData);
      this.show = false;
    },
  },
};
</script>
