import sounds from '../../sounds';

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
    soundById: (state) => (id) => state.sounds.find((sound) => sound.id === id),
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
    stop(state, { sound, fade = 500 }) {
      sound.state = 'stopped';
      if (fade) {
        sound.player.once('fade', async () => {
          sound.player.stop();
          sound.player.unload();
        });
        sound.player.fade(sound.player.volume(), 0, 500);
      } else {
        sound.player.stop();
        sound.player.unload();
      }
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
    async playStop({ commit, dispatch }, { sound, fade = 500 }) {
      if (sound.player.playing()) {
        commit('stop', { sound, fade });
      } else {
        if (sound.player.state() === 'unloaded') {
          await dispatch('load', { sound });
        }
        commit('play', { sound, fade });
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
          commit('play', { sound, fade: false });
        }
      });
    },
    stopAll({ commit, state }, { fade = 500 }) {
      state.sounds.filter(
        (sound) => sound.state !== 'stopped',
      ).forEach((sound) => {
        commit('stop', { sound, fade });
      });
    },
    async prefetch({ state }) {
      const cache = await window.caches.open('audio-cache');
      await Promise.all(state.sounds.map(async (sound) => {
        sound.loading = true;
        await cache.add(sound.src);
        sound.loading = false;
      }));
    },
  },
};
