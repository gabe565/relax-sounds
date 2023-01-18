<template>
  <div>
    <v-dialog v-model="show" max-width="400">
      <v-card>
        <v-card-title class="text-h5">Confirm</v-card-title>
        <v-card-text>Delete {{ count }} presets?</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="show = false">
            <v-icon aria-hidden="true">fal fa-times fa-fw</v-icon>
            Close
          </v-btn>
          <v-btn variant="text" :disabled="countdown > 0" color="red" @click="remove">
            <v-icon aria-hidden="true">fal fa-trash fa-fw</v-icon>
            Delete
            <template v-if="countdown > 0">({{ countdown }})</template>
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-snackbar
      v-model="showSnackbar"
      timeout="5000"
      location="bottom"
      content-class="mb-14 mb-md-0"
    >
      All presets have been removed.
    </v-snackbar>
  </div>
</template>

<script>
import { wait } from "../util/helpers";

const Countdown = 5;

let timeout;

export default {
  name: "RemoveAll",

  props: {
    modelValue: Boolean,
  },

  emits: ["update:modelValue"],

  data: () => ({
    show: false,
    showSnackbar: false,
    countdown: 0,
  }),

  computed: {
    count() {
      return this.$store.state.presets.presets.length || 0;
    },
  },

  watch: {
    modelValue(val) {
      this.show = val;
    },
    show(val) {
      this.$emit("update:modelValue", val);

      if (val) {
        clearTimeout(timeout);
        this.countdown = Countdown;
        timeout = setTimeout(this.doCountdown, 1000);
      }
    },
  },

  methods: {
    async remove() {
      this.show = false;
      await wait(300);
      this.$store.commit("presets/removeAll");
      this.showSnackbar = true;
    },
    doCountdown() {
      this.countdown -= 1;
      if (this.countdown > 0) {
        timeout = setTimeout(this.doCountdown, 1000);
      }
    },
  },
};
</script>
