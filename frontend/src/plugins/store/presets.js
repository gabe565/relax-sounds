import { SoundState } from '../../util/Sound';
import { Preset } from '../../util/Preset';

const saveState = ({ presets }) => localStorage.setItem('presets', JSON.stringify(presets));

const loadState = () => {
  let presets = JSON.parse(localStorage.getItem('presets'));

  if (!presets) {
    // Playlist to preset migration
    const playlists = JSON.parse(localStorage.getItem('playlists'));
    if (playlists) {
      presets = playlists.playlists;
      localStorage.setItem('presets', JSON.stringify(presets));
      localStorage.removeItem('playlists');
    }
  }

  if (presets) {
    return presets.map((preset) => new Preset(preset));
  }
  return [];
};

export default {
  namespaced: true,
  state: {
    presets: loadState(),
    currentName: null,
  },

  mutations: {
    add(state, { preset, playing = true }) {
      state.presets.push(new Preset(preset));
      if (playing) {
        state.currentName = preset.name;
      }
      saveState(state);
    },
    remove(state, { preset }) {
      const index = state.presets.indexOf(preset);
      state.presets.splice(index, 1);
      saveState(state);
    },
    play(state, { preset }) {
      state.currentName = preset.name;
      if (preset.new) {
        preset.new = false;
        saveState(state);
      }
    },
    disableCurrent(state) {
      state.currentName = null;
    },
    removeAll(state) {
      state.presets = [];
      state.currentName = null;
      saveState(state);
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
        preset: {
          name,
          sounds,
          new: true,
        },
      });
    },

    load({ dispatch, rootGetters }, { preset }) {
      return Promise.all(preset.sounds.map((savedSound) => {
        const sound = rootGetters['player/soundById'](savedSound.id);
        return dispatch('player/load', { sound }, { root: true });
      }));
    },

    async play({ commit, dispatch, rootGetters }, { preset }) {
      if (rootGetters['player/state'] !== SoundState.STOPPED) {
        dispatch('player/stopAll', { fade: 0, local: true }, { root: true });
      }
      await Promise.all(preset.sounds.map((savedSound) => {
        const sound = rootGetters['player/soundById'](savedSound.id);
        sound.volume = savedSound.volume;
        const fade = rootGetters['player/state'] === SoundState.STOPPED ? 500 : false;
        return dispatch('player/playStop', { sound, fade, local: true }, { root: true });
      }));
      commit('play', { preset });
      await dispatch('player/updateCast', null, { root: true });
    },
  },
};
