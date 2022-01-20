import Vue from 'vue';
import Vuex from 'vuex';
import filters from './filters';
import presets from './presets';
import player from './player';

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    filters,
    presets,
    player,
  },
});
