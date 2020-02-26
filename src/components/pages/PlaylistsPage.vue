<template>
  <Page>
    <v-row>
      <v-col
        cols="12" md="6"
        v-for="playlist of playlists"
        :key="playlist.name"
      >
        <Playlist :playlist="playlist"/>
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
