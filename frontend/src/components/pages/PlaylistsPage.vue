<template>
  <Page>
    <v-row>
      <v-col
        cols="12" lg="6"
        v-for="(playlist, key) of playlists"
        :key="key"
      >
        <Playlist :playlist="playlist"/>
      </v-col>
      <v-col v-if="playlists.length === 0">
        <v-alert outlined color="warning" icon="fal fa-info-circle">
          No Playlists Saved Yet!
        </v-alert>
      </v-col>
    </v-row>
  </Page>
</template>

<script>
import { mapState } from 'vuex';
import Playlist from '../Playlist.vue';
import Page from './Page.vue';

export default {
  name: 'Playlists',
  components: { Page, Playlist },
  computed: mapState('playlists', [
    'playlists',
  ]),
};
</script>

<style scoped>

</style>
