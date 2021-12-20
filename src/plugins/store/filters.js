import Fuse from 'fuse.js';
import { sounds } from '../../util/sounds';

const PER_PAGE = 48;

const fuse = new Fuse(sounds, {
  shouldSort: true,
  threshold: 0.3,
  location: 0,
  distance: 100,
  maxPatternLength: 32,
  minMatchCharLength: 1,
  keys: [
    'name',
    'tags',
  ],
});

export default {
  namespaced: true,
  state: {
    filters: {
      word: '',
      playing: false,
      page: 1,
    },
  },
  getters: {
    filteredSounds(state, _, rootState) {
      let result;
      if (state.filters.word) {
        result = fuse.search(state.filters.word).map((e) => e.item);
      } else {
        result = rootState.player.sounds;
      }
      if (state.filters.playing) {
        result = result.filter((e) => !e.isStopped);
      }
      return result;
    },
    sounds(state, getters) {
      const result = getters.filteredSounds;
      const offset = PER_PAGE * (state.filters.page - 1);
      return result.slice(offset, offset + PER_PAGE);
    },
    pages(_, getters) {
      return Math.ceil(getters.filteredSounds.length / PER_PAGE);
    },
  },
};
