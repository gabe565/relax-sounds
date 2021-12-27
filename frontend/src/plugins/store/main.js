import Vue from 'vue';
import Vuex from 'vuex';
import filters from './filters';
import playlists from './playlists';
import player from './player';

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    filters,
    playlists,
    player,
  },
});
