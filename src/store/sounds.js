import axios from 'axios';

import sounds from '../sounds';

export default {
  namespaced: true,
  state: {
    sounds,
  },

  getters: {
    playing(state) {
      return state.sounds.filter((sound) => sound.state === 'playing');
    },
    state(state) {
      const states = new Set(state.sounds.map((sound) => sound.state));
      if (states.has('playing')) {
        return 'playing';
      }
      if (states.has('paused')) {
        return 'paused';
      }
      return 'stopped';
    },
  },

  mutations: {
    play(state, { sound, fade = 500 }) {
      sound.state = 'playing';
      sound.player.play();
      if (fade) {
        sound.player.fade(0, sound.volume, fade);
      }
    },
    pause(state, { sound }) {
      sound.state = 'paused';
      sound.player.pause();
    },
    stop(state, { sound }) {
      sound.state = 'stopped';
      sound.player.once('fade', async () => {
        sound.player.stop();
      });
      sound.player.fade(sound.player.volume(), 0, 500);
    },
    volume(state, { sound, value }) {
      sound.player.volume(value);
      sound.volume = value;
    },
  },

  actions: {
    async load(_, { sound }) {
      if (sound.player.state() === 'unloaded') {
        sound.loading = true;
        return new Promise((resolve, reject) => {
          sound.player.once('load', () => {
            sound.loading = false;
            resolve();
          });
          sound.player.once('loaderror', () => {
            sound.loading = false;
            reject();
          });
          sound.player.load();
        });
      }
      return true;
    },
    playPause({ commit }, { sound }) {
      if (sound.player.playing()) {
        commit('stop', { sound });
      } else if (sound.player.state() === 'loaded') {
        commit('play', { sound });
      } else {
        commit('load', { sound });
      }
    },
    playPauseAll({ commit, state }) {
      const newState = this.getters['sounds/state'] === 'playing' ? 'paused' : 'playing';
      state.sounds.filter(
        (sound) => sound.state !== 'stopped',
      ).forEach((sound) => {
        sound.state = newState;
        if (newState === 'paused') {
          commit('pause', { sound });
        } else {
          commit('play', { sound });
        }
      });
    },
    async prefetch({ state }) {
      return Promise.all(state.sounds.map(async (sound) => {
        sound.loading = true;
        await axios.get(sound.src);
        sound.loading = false;
      }));
    },
  },
};
