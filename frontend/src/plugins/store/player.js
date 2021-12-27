import { soundConfig, SoundState } from '../../util/sounds';
import { formatError, getCastSession } from '../../util/googleCast';
import { decodeSounds, encodeSounds } from '../../util/shareUrl';

export default {
  namespaced: true,
  state: {
    sounds: soundConfig,
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
    soundById(state) {
      return (id) => state.sounds.find((sound) => sound.id === id);
    },
    encodedSounds(state, getter) {
      return encodeSounds(getter.soundsPlaying
        .map((sound) => ({
          id: sound.id,
          volume: sound.volume,
        })));
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
      if (rootState.playlists.currentName) {
        commit('playlists/disableCurrent', null, { root: true });
      }
      if (!local) dispatch('updateCast');
    },
    pauseAll({ commit, getters, state }, { local = false } = {}) {
      getters.soundsPlaying.forEach((sound) => {
        commit('pause', { sound });
      });
      if (!local && state.remotePlayerController) {
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
      if (!local && state.remotePlayerController) {
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
      if (!local && state.remotePlayerController) {
        commit('playlists/disableCurrent', null, { root: true });
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
        ({ value }) => {
          if (getters.state !== SoundState.STOPPED) {
            dispatch('pauseAll');
            setTimeout(async () => {
              await dispatch('playPauseAll');
              dispatch('updateCast');
            }, 0);
          }
          commit('castConnectedChanged', { value });
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
            await dispatch('playPauseAll');
          }
        },
      );

      remotePlayerController.addEventListener(
        cast.framework.RemotePlayerEventType.MEDIA_INFO_CHANGED,
        async ({ value }) => {
          if (value && getters.state === SoundState.STOPPED) {
            const encoded = value.contentId.match(/\/mix\/(.+?)$/)[1];
            const sounds = decodeSounds(encoded);
            await Promise.all(sounds.map((savedSound) => {
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
      if (getters.state === SoundState.PLAYING) {
        const castSession = getCastSession();
        if (castSession) {
          const { chrome } = window;
          const url = `${window.location.origin}/mix/${getters.encodedSounds}`;

          const mediaInfo = new chrome.cast.media.MediaInfo(url, 'music');
          mediaInfo.metadata = new chrome.cast.media.MusicTrackMediaMetadata();
          if (rootState.playlists.currentName) {
            mediaInfo.metadata.title = rootState.playlists.currentName;
          } else {
            mediaInfo.metadata.title = getters.soundsPlaying
              .map((sound) => sound.name)
              .sort((a, b) => a.localeCompare(b))
              .join(', ');
          }
          mediaInfo.metadata.artist = 'Relax Sounds';

          const request = new chrome.cast.media.LoadRequest(mediaInfo);

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
  },
};
