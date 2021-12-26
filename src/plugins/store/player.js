import { sounds, SoundState } from '../../util/sounds';

export default {
  namespaced: true,
  state: {
    sounds,
  },

  getters: {
    playing(state) {
      return state.sounds.filter((sound) => sound.isPlaying);
    },
    state(state) {
      const states = new Set(state.sounds.map((sound) => sound.state));
      if (states.has(SoundState.PLAYING)) {
        return SoundState.PLAYING;
      }
      if (states.has(SoundState.PAUSED)) {
        return SoundState.PAUSED;
      }
      return SoundState.STOPPED;
    },
    soundById: (state) => (id) => state.sounds.find((sound) => sound.id === id),
  },

  mutations: {
    play(state, { sound, fade = 250 }) {
      sound.play(true, fade);
    },
    pause(state, { sound }) {
      sound.pause(true);
    },
    stop(state, { sound, fade = 250 }) {
      sound.stop(true, fade);
    },
    volume(state, { sound, value }) {
      sound.volume(value);
    },
  },

  actions: {
    async playStop({ commit }, { sound, fade = 250 }) {
      if (sound.state === SoundState.PLAYING) {
        commit('stop', { sound, fade });
      } else {
        if (sound.isUnloaded) {
          await sound.load();
        }
        if (sound.isPaused) {
          fade = false;
        }
        commit('play', { sound, fade });
      }
    },
    playPauseAll({ commit, state }) {
      const newState = this.getters['player/state'] === SoundState.PLAYING ? SoundState.PAUSED : SoundState.PLAYING;
      state.sounds.filter(
        (sound) => !sound.isStopped,
      ).forEach((sound) => {
        sound.state = newState;
        if (newState === SoundState.PAUSED) {
          commit('pause', { sound });
        } else {
          commit('play', { sound, fade: false });
        }
      });
    },
    stopAll({ commit, state }, { fade = 250 }) {
      state.sounds.filter(
        (sound) => !sound.isStopped,
      ).forEach((sound) => {
        commit('stop', { sound, fade });
      });
    },
  },
};
