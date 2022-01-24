<template>
  <v-col class="shrink">
    <v-btn @click.stop="show = true" elevation="0" icon aria-label="Share">
      <v-icon dense aria-hidden="true">
        fas fa-fw fa-trash
      </v-icon>
    </v-btn>

    <v-dialog v-model="show" max-width="500">
      <v-card>
        <v-card-title class="headline">Delete Preset?</v-card-title>
        <v-card-text>
          Are you sure you want to delete "{{ preset.name }}"?
        </v-card-text>
        <v-card-actions>
          <v-spacer/>
          <v-btn text @click="show = false">
            <v-icon aria-hidden="true">fal fa-times fa-fw</v-icon>
            Close
          </v-btn>
          <v-btn color="red" text @click="remove">
            <v-icon aria-hidden="true">fal fa-trash fa-fw</v-icon>
            Delete
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-col>
</template>

<script>
export default {
  name: 'DeleteButton',

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
    remove() {
      this.$store.commit('presets/remove', { preset: this.preset });
      this.show = false;
    },
  },
};
</script>
