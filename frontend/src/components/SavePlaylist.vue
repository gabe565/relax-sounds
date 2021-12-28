<template>
  <v-btn icon @click="showDialog = !showDialog" :disabled="state === 'stopped'">
    <v-icon>fas fa-plus-circle</v-icon>
    <v-dialog v-model="showDialog" max-width="500">
      <v-card>
        <v-card-title class="headline">Playlist Name</v-card-title>
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
            Save Playlist
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-btn>
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  name: 'SavePlaylist',

  data: () => ({
    showDialog: false,
    name: '',
  }),

  computed: mapGetters('player', [
    'state',
  ]),

  methods: {
    cancel() {
      this.showDialog = false;
      this.name = '';
    },
    save() {
      this.$store.dispatch('playlists/savePlaying', { name: this.name });
      this.showDialog = false;
      this.name = '';
      this.$router.push({ name: 'Playlists' });
    },
  },
};
</script>

<style scoped>

</style>
