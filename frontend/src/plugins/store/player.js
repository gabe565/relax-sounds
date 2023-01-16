import axios from 'axios';
import { SoundState } from '../../util/Sound';
import { getSounds } from '../../data/sounds';
import { formatError, getCastSession } from '../../util/googleCast';
import { Preset } from '../../util/Preset';

export default {
  namespaced: true,
  state: {
    sounds: [],
    remotePlayer: null,
    remotePlayerController: null,
    castConnected: false,
  },

  getters: {
    soundsPlaying(state) {
      return state.sounds.filter((sound) => sound.isPlaying);
    },
    soundsNotStopped(state) {
      return state.sounds.filter((sound) => !sound.isStopped);
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
    isPlaying(_, getters) {
      return getters.state === SoundState.PLAYING;
    },
    isPaused(_, getters) {
      return getters.state === SoundState.PAUSED;
    },
    isStopped(_, getters) {
      return getters.state === SoundState.STOPPED;
    },
    soundById(state) {
      return (id) => state.sounds.find((sound) => sound.id === id);
    },
  },

  mutations: {
    play(state, { sound, fade = 250 }) {
      sound.play(!state.castConnected, fade);
    },
    pause(state, { sound }) {
      sound.pause(!state.castConnected, true);
    },
    stop(state, { sound, fade = 250 }) {
      sound.stop(!state.castConnected, fade);
    },
    volume(state, { sound, value }) {
      sound.volume = value;
    },
    initCastApi(state, { remotePlayer, remotePlayerController }) {
      state.remotePlayer = remotePlayer;
      state.remotePlayerController = remotePlayerController;
    },
    castConnectedChanged(state, { value }) {
      state.castConnected = value;
    },
    initSounds(state, value) {
      state.sounds = value;
    },
  },

  actions: {
    async playStop({
      state, commit, dispatch, rootState,
    }, { sound, fade = 250, local = false }) {
      if (sound.state === SoundState.PLAYING) {
        commit('stop', { sound, fade });
      } else {
        if (!state.castConnected && sound.isUnloaded) {
          await sound.load();
        }
        if (sound.isPaused) {
          fade = false;
        }
        commit('play', { sound, fade });
      }
      if (rootState.presets.currentName) {
        commit('presets/disableCurrent', null, { root: true });
      }
      if (!local && state.castConnected) dispatch('updateCast');
    },
    pauseAll({ commit, getters, state }, { local = false } = {}) {
      getters.soundsPlaying.forEach((sound) => {
        commit('pause', { sound });
      });
      if (!local && state.castConnected) {
        state.remotePlayerController.playOrPause();
      }
    },
    async playPauseAll({ commit, getters, state }, { local = false } = {}) {
      const newState = this.getters['player/state'] === SoundState.PLAYING ? SoundState.PAUSED : SoundState.PLAYING;
      await Promise.all(getters.soundsNotStopped
        .map(async (sound) => {
          sound.state = newState;
          if (newState === SoundState.PAUSED) {
            commit('pause', { sound });
          } else {
            if (!state.castConnected && sound.isUnloaded) {
              await sound.load();
            }
            commit('play', { sound, fade: false });
          }
        }));
      if (!local && state.castConnected) {
        state.remotePlayerController.playOrPause();
      }
    },
    stopAll({ commit, getters, state }, { fade = 250, local = false }) {
      getters.soundsNotStopped.forEach((sound) => {
        commit('stop', { sound, fade });
      });
      if (state.remotePlayerController) {
        state.remotePlayerController.stop();
      }
      if (!local && state.castConnected) {
        commit('presets/disableCurrent', null, { root: true });
      }
    },
    initializeCastApi({ commit, dispatch, getters }) {
      const { cast, chrome } = window;

      cast.framework.CastContext.getInstance().setOptions({
        receiverApplicationId: chrome.cast.media.DEFAULT_MEDIA_RECEIVER_APP_ID,
        autoJoinPolicy: chrome.cast.AutoJoinPolicy.ORIGIN_SCOPED,
      });

      const remotePlayer = new cast.framework.RemotePlayer();
      const remotePlayerController = new cast.framework.RemotePlayerController(remotePlayer);
      commit('initCastApi', { remotePlayer, remotePlayerController });

      remotePlayerController.addEventListener(
        cast.framework.RemotePlayerEventType.IS_CONNECTED_CHANGED,
        async ({ value }) => {
          let shouldPlay;
          if (getters.state !== SoundState.STOPPED) {
            dispatch('pauseAll');
            shouldPlay = true;
          }
          commit('castConnectedChanged', { value });
          if (shouldPlay) {
            await dispatch('playPauseAll');
            dispatch('updateCast');
          }
        },
      );

      remotePlayerController.addEventListener(
        cast.framework.RemotePlayerEventType.IS_PAUSED_CHANGED,
        async () => {
          if (remotePlayer.isPaused) {
            dispatch('pauseAll', { local: true });
          } else if (getters.state !== SoundState.PLAYING) {
            // If currently not playing, start to play.
            // This occurs if starting to play from local, but this check is
            // required if the state is changed remotely.
            await dispatch('playPauseAll', { local: true });
          }
        },
      );

      remotePlayerController.addEventListener(
        cast.framework.RemotePlayerEventType.MEDIA_INFO_CHANGED,
        async ({ value }) => {
          if (value && getters.isStopped) {
            const preset = new Preset();
            preset.mixUrl = value.contentId;
            await Promise.all(preset.sounds.map((savedSound) => {
              const sound = getters.soundById(savedSound.id);
              sound.volume = savedSound.volume;
              const fade = getters.state === SoundState.STOPPED ? 250 : false;
              return dispatch('playStop', {
                sound,
                fade,
                local: true,
              });
            }));
          }
        },
      );
    },
    async updateCast({ getters, rootState, state }) {
      if (getters.isPlaying) {
        const castSession = getCastSession();
        if (castSession) {
          const { chrome } = window;
          const preset = new Preset({ sounds: getters.soundsPlaying });

          const mediaInfo = new chrome.cast.media.MediaInfo(preset.mixUrl, 'music');
          mediaInfo.metadata = new chrome.cast.media.MusicTrackMediaMetadata();
          if (rootState.presets.currentName) {
            mediaInfo.metadata.title = rootState.presets.currentName;
          } else {
            mediaInfo.metadata.title = getters.soundsPlaying
              .map((sound) => sound.name)
              .sort((a, b) => a.localeCompare(b))
              .join(', ');
          }
          mediaInfo.metadata.artist = 'Relax Sounds';
          mediaInfo.metadata.images = [
            new chrome.cast.Image(`${window.location.origin}/img/icons/android-chrome-maskable-512x512.png`),
          ];

          const queue = new chrome.cast.media.QueueLoadRequest([
            new chrome.cast.media.QueueItem(mediaInfo),
          ]);
          queue.repeatMode = chrome.cast.media.RepeatMode.SINGLE;

          const request = new chrome.cast.media.LoadRequest(mediaInfo);
          request.queueData = queue;

          try {
            await castSession.loadMedia(request);
          } catch (error) {
            console.error(`Remote media load error: ${formatError(error)}`);
          }
        }
      } else {
        state.remotePlayerController.stop();
      }
    },
    async initSounds({ commit, dispatch, state }) {
      if (!state.sounds.length) {
        const conf = await getSounds();
        commit('initSounds', conf);
        dispatch('filters/initSounds', conf, { root: true });
      }
    },
    async prefetch({ state }) {
      return Promise.all(state.sounds.map(async (sound) => {
        sound.isLoading = true;
        try {
          await axios.get(sound.src);
        } catch (error) {
          console.error(error);
        }
        sound.isLoading = false;
      }));
    },
  },
};
