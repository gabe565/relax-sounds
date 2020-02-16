import { Howl } from 'howler';
import defaultSounds from './assets/sounds.json';

const sounds = defaultSounds.map((sound) => ({
  ...sound,
  state: 'stopped',
  volume: 1,
  loading: false,
  player: new Howl({
    src: [`/audio/${sound.id}.ogg`],
    loop: true,
    preload: false,
    volume: 0,
  }),
})).sort(
  (left, right) => left.name.localeCompare(right.name),
);

export default {
  state: {
    sounds,
    filter: {
      word: '',
      playing: false,
    },
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
    icon(state) {
      return (id) => {
        const sound = state.sounds.find((e) => e.id === id);
        return sound.state === 'playing' ? 'mdi-stop' : 'mdi-play';
      };
    },
  },

  mutations: {
    load(state, { id, fade = 500 }) {
      const sound = state.sounds.find((e) => e.id === id);
      sound.loading = true;
      sound.player.once('load', () => {
        sound.state = 'playing';
        sound.loading = false;
        this.commit('play', { id, fade });
      });
      sound.player.load();
    },
    play(state, { id, fade = 500 }) {
      const sound = state.sounds.find((e) => e.id === id);
      sound.state = 'playing';
      sound.player.play();
      if (fade) {
        sound.player.fade(0, sound.volume, fade);
      }
    },
    pause(state, { id }) {
      const sound = state.sounds.find((e) => e.id === id);
      sound.state = 'paused';
      sound.player.pause();
    },
    stop(state, { id }) {
      const sound = state.sounds.find((e) => e.id === id);
      sound.state = 'stopped';
      sound.player.once('fade', async () => {
        sound.player.stop();
      });
      sound.player.fade(sound.player.volume(), 0, 500);
    },
    volume(state, { id, value }) {
      const sound = state.sounds.find((e) => e.id === id);
      sound.player.volume(value);
      sound.volume = value;
    },
    playPause(state, { id }) {
      const sound = state.sounds.find((e) => e.id === id);
      if (sound.player.playing()) {
        this.commit('stop', { id });
      } else if (sound.player.state() === 'loaded') {
        this.commit('play', { id });
      } else {
        this.commit('load', { id });
      }
    },
    playPauseAll(state) {
      const newState = this.getters.state === 'playing' ? 'paused' : 'playing';
      state.sounds.filter(
        (sound) => sound.state !== 'stopped',
      ).forEach((sound) => {
        sound.state = newState;
        if (newState === 'paused') {
          sound.player.pause();
        } else if (newState === 'playing') {
          sound.player.play();
        }
      });
    },
  },
};
