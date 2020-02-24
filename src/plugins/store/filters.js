import Fuse from 'fuse.js';

import sounds from '../../sounds';

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
    fuse,
    filters: {
      word: '',
      playing: false,
    },
  },
  getters: {
    sounds(state, _, rootState) {
      if (rootState.sounds.sounds) {
        let result;
        if (state.filters.word) {
          result = state.fuse.search(state.filters.word);
        } else {
          result = rootState.sounds.sounds;
        }
        if (state.filters.playing) {
          result = result.filter((e) => e.state !== 'stopped');
        }
        return result;
      }
      return [];
    },
  },
};
