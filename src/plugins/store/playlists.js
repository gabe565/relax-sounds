import { castPlaylist, getCastSession } from '../../util/googleCast';
import { SoundState } from '../../util/sounds';

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
      const sounds = rootState.player.sounds
        .filter((sound) => sound.state === SoundState.PLAYING)
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
      commit('play', { playlist });
      dispatch('player/stopAll', { fade: 0 }, { root: true });
      const castSession = getCastSession();
      if (castSession) {
        try {
          const result = await castPlaylist(castSession, playlist);
          console.log(result);
        } catch (error) {
          console.log(`Error code: ${error}`);
        }
      } else {
        await Promise.all(playlist.sounds.map(async (savedSound) => {
          const sound = rootGetters['player/soundById'](savedSound.id);
          sound.volume = savedSound.volume;
          const fade = rootGetters['player/state'] === 'stopped' ? 250 : false;
          return dispatch('player/playStop', { sound, fade }, { root: true });
        }));
      }
    },
  },
};
