const defaultState = {
  playlists: [],
};

const saveState = (state) => localStorage.setItem('playlists', JSON.stringify(state));
const loadState = () => JSON.parse(localStorage.getItem('playlists'));

export default {
  namespaced: true,
  state: loadState() || defaultState,

  mutations: {
    add(state, { playlist }) {
      state.playlists.push(playlist);
      saveState(state);
    },
    remove(state, { playlist }) {
      const index = state.playlists.indexOf(playlist);
      state.playlists.splice(index, 1);
      saveState(state);
    },
    play(state, { playlist }) {
      playlist.new = false;
      saveState(state);
    },
  },

  actions: {
    savePlaying({ commit, rootState }, { name }) {
      const sounds = rootState.sounds.sounds
        .filter((sound) => sound.state === 'playing')
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
        const sound = rootGetters['sounds/soundById'](savedSound.id);
        return dispatch('sounds/load', { sound }, { root: true });
      }));
    },

    async play({ commit, dispatch, rootGetters }, { playlist }) {
      commit('play', { playlist });
      dispatch('sounds/stopAll', { fade: 0 }, { root: true });
      await Promise.all(playlist.sounds.map(async (savedSound) => {
        const sound = rootGetters['sounds/soundById'](savedSound.id);
        commit('sounds/volume', { sound, value: savedSound.volume }, { root: true });
        const fade = rootGetters['sounds/state'] === 'stopped' ? 500 : false;
        return dispatch('sounds/playStop', { sound, fade }, { root: true });
      }));
    },
  },
};
