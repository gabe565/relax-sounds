import { SoundState } from '../../util/Sound';
import { Preset } from '../../util/Preset';

const version = 2;

const saveState = ({ presets }) => {
  const state = { version, presets };
  localStorage.setItem('presets', JSON.stringify(state));
};

const loadState = () => {
  let state = JSON.parse(localStorage.getItem('presets'));

  if (!state) {
    // Playlist to preset migration
    const playlists = JSON.parse(localStorage.getItem('playlists'));
    if (playlists) {
      state = { version, presets: playlists.playlists };
      localStorage.setItem('presets', JSON.stringify(state));
      localStorage.removeItem('playlists');
    }
  }

  if (state) {
    let dirty = false;
    if (Array.isArray(state)) {
      // Migrate state to object
      dirty = true;
      state = { presets: state };
    }

    if (!state.version || state.version === 1) {
      // v2 migration
      dirty = true;
      for (const preset of state.presets) {
        for (const sound of preset.sounds) {
          sound.id = sound.id.toString();
        }
      }
      state.version = version;
    }

    if (dirty) {
      saveState(state);
    }

    return state.presets.map((preset) => new Preset(preset));
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
      for (const sound of preset.sounds) {
        if (typeof sound.id === 'number') {
          sound.id = sound.id.toString();
        }
      }
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
