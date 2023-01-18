<template>
  <v-col class="flex-grow-0">
    <v-btn elevation="0" icon variant="plain" aria-label="Share" @click.stop="show = true">
      <v-icon aria-hidden="true">fas fa-fw fa-trash</v-icon>
    </v-btn>

    <v-dialog v-model="show" width="400">
      <v-card>
        <v-card-title class="text-h5">Confirm</v-card-title>
        <v-card-text>Delete the preset "{{ preset.name }}"?</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="show = false">
            <v-icon aria-hidden="true">fal fa-times fa-fw</v-icon>
            Close
          </v-btn>
          <v-btn color="red" variant="text" @click="remove">
            <v-icon aria-hidden="true">fal fa-trash fa-fw</v-icon>
            Delete
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-col>
</template>

<script>
import { wait } from "../../util/helpers";

export default {
  name: "DeleteButton",

  props: {
    preset: {
      type: Object,
      required: true,
    },
  },

  data: () => ({
    show: false,
  }),

  methods: {
    async remove() {
      this.show = false;
      await wait(300);
      this.$store.commit("presets/remove", { preset: this.preset });
    },
  },
};
</script>
