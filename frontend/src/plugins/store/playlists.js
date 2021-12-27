import { SoundState } from '../../util/sounds';

const defaultState = {
  playlists: [],
  currentName: null,
};

const saveState = (state) => localStorage.setItem('playlists', JSON.stringify(state));
const loadState = () => JSON.parse(localStorage.getItem('playlists'));

export default {
  namespaced: true,
  state: loadState() || defaultState,

  mutations: {
    add(state, { playlist }) {
      state.playlists.push(playlist);
      state.currentName = playlist.name;
      saveState(state);
    },
    remove(state, { playlist }) {
      const index = state.playlists.indexOf(playlist);
      state.playlists.splice(index, 1);
      saveState(state);
    },
    play(state, { playlist }) {
      state.currentName = playlist.name;
      if (playlist.new) {
        playlist.new = false;
        saveState(state);
      }
    },
    disableCurrent(state) {
      state.currentName = null;
    },
  },

  actions: {
    savePlaying({ commit, rootGetters }, { name }) {
      const sounds = rootGetters['player/soundsPlaying']
        .map((sound) => ({
          id: sound.id,
          volume: sound.volume,
        }));

      commit('add', {
        playlist: {
          name,
          sounds,
          new: true,
        },
      });
    },

    load({ dispatch, rootGetters }, { playlist }) {
      return Promise.all(playlist.sounds.map((savedSound) => {
        const sound = rootGetters['player/soundById'](savedSound.id);
        return dispatch('player/load', { sound }, { root: true });
      }));
    },

    async play({ commit, dispatch, rootGetters }, { playlist }) {
      if (rootGetters['player/state'] !== SoundState.STOPPED) {
        dispatch('player/stopAll', { fade: 0, local: true }, { root: true });
      }
      await Promise.all(playlist.sounds.map((savedSound) => {
        const sound = rootGetters['player/soundById'](savedSound.id);
        sound.volume = savedSound.volume;
        const fade = rootGetters['player/state'] === SoundState.STOPPED ? 500 : false;
        return dispatch('player/playStop', { sound, fade, local: true }, { root: true });
      }));
      commit('play', { playlist });
      await dispatch('player/updateCast', null, { root: true });
    },
  },
};
