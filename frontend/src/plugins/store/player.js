import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { SoundState } from "../../util/Sound";
import { getSounds } from "../../data/sounds";
import { formatError, getCastSession } from "../../util/googleCast";
import { Preset } from "../../util/Preset";
import pb from "../pocketbase";
import { wait } from "../../util/helpers";

export const usePlayerStore = defineStore("player", () => {
  const sounds = ref([]);
  const currentName = ref();
  let remotePlayer;
  let remotePlayerController;
  const castConnected = ref(false);

  const soundsPlaying = computed(() => {
    return sounds.value.filter((sound) => sound.isPlaying);
  });

  const soundsNotStopped = computed(() => {
    return sounds.value.filter((sound) => !sound.isStopped);
  });

  const state = computed(() => {
    const states = new Set(sounds.value.map((sound) => sound.state));
    if (states.has(SoundState.PLAYING)) {
      return SoundState.PLAYING;
    }
    if (states.has(SoundState.PAUSED)) {
      return SoundState.PAUSED;
    }
    return SoundState.STOPPED;
  });

  const isPlaying = computed(() => {
    return state.value === SoundState.PLAYING;
  });

  const isPaused = computed(() => {
    return state.value === SoundState.PAUSED;
  });

  const isStopped = computed(() => {
    return state.value === SoundState.STOPPED;
  });

  const soundById = (id) => sounds.value.find((sound) => sound.id === id);

  const play = ({ sound, fade = 250 }) => {
    sound.play(!castConnected.value, fade);
  };

  const pause = ({ sound }) => {
    sound.pause(!castConnected.value, true);
  };

  const stop = ({ sound, fade = 250 }) => {
    sound.stop(!castConnected.value, fade);
  };

  const volume = ({ sound, value }) => {
    sound.volume = value;
  };

  const castConnectedChanged = ({ value }) => {
    castConnected.value = value;
  };

  const updateCast = async () => {
    if (isPlaying.value) {
      const castSession = getCastSession();
      if (castSession) {
        const { cast } = window.chrome;
        const preset = new Preset({ sounds: soundsPlaying.value });

        const mediaInfo = new cast.media.MediaInfo(preset.mixUrl, "music");
        mediaInfo.metadata = new cast.media.MusicTrackMediaMetadata();
        mediaInfo.metadata.title = currentName.value;
        if (!mediaInfo.metadata.title) {
          mediaInfo.metadata.title = soundsPlaying.value
            .map((sound) => sound.name)
            .sort((a, b) => a.localeCompare(b))
            .join(", ");
        }
        mediaInfo.metadata.artist = "Relax Sounds";
        mediaInfo.metadata.images = [
          new cast.Image(`${window.location.origin}/img/icons/android-chrome-maskable-512x512.png`),
        ];

        const queue = new cast.media.QueueLoadRequest([new cast.media.QueueItem(mediaInfo)]);
        queue.repeatMode = cast.media.RepeatMode.SINGLE;

        const request = new cast.media.LoadRequest(mediaInfo);
        request.queueData = queue;

        try {
          await castSession.loadMedia(request);
        } catch (error) {
          console.error(`Remote media load error: ${formatError(error)}`);
        }
      }
    } else {
      remotePlayerController.stop();
    }
  };

  const playStop = async ({ sound, fade = 250, local = false }) => {
    if (sound.state === SoundState.PLAYING) {
      stop({ sound, fade });
    } else {
      if (!castConnected.value && sound.isUnloaded) {
        await sound.load();
      }
      if (sound.isPaused) {
        fade = false;
      }
      play({ sound, fade });
    }
    currentName.value = null;
    if (!local && castConnected) {
      await updateCast();
    }
  };

  const pauseAll = ({ local = false } = {}) => {
    soundsPlaying.value.forEach((sound) => {
      pause({ sound });
    });
    if (!local && castConnected) {
      remotePlayerController.playOrPause();
    }
  };

  const playPauseAll = async ({ local = false } = {}) => {
    const newState = state.value === SoundState.PLAYING ? SoundState.PAUSED : SoundState.PLAYING;
    await Promise.all(
      soundsNotStopped.value.map(async (sound) => {
        sound.state = newState;
        if (newState === SoundState.PAUSED) {
          pause({ sound });
        } else {
          if (!castConnected.value && sound.isUnloaded) {
            await sound.load();
          }
          play({ sound, fade: 0 });
        }
      })
    );
    if (!local && castConnected) {
      remotePlayerController.playOrPause();
    }
  };

  const stopAll = ({ fade = 250, local = false }) => {
    soundsNotStopped.value.forEach((sound) => {
      stop({ sound, fade });
    });
    if (remotePlayerController) {
      remotePlayerController.stop();
    }
    if (!local && castConnected) {
      currentName.value = null;
    }
  };

  const initSounds = async () => {
    if (!sounds.value.length) {
      sounds.value = await getSounds();
    }
  };

  const initializeCastApi = () => {
    const { framework: castFramework } = window.cast;
    const { cast } = window.chrome;

    castFramework.CastContext.getInstance().setOptions({
      receiverApplicationId: cast.media.DEFAULT_MEDIA_RECEIVER_APP_ID,
      autoJoinPolicy: cast.AutoJoinPolicy.ORIGIN_SCOPED,
    });

    remotePlayer = new castFramework.RemotePlayer();
    remotePlayerController = new castFramework.RemotePlayerController(remotePlayer);

    remotePlayerController.addEventListener(
      castFramework.RemotePlayerEventType.IS_CONNECTED_CHANGED,
      async ({ value }) => {
        let shouldPlay;
        if (state.value !== SoundState.STOPPED) {
          pauseAll();
          shouldPlay = true;
        }
        castConnectedChanged({ value });
        if (shouldPlay) {
          await playPauseAll();
          await updateCast();
        }
      }
    );

    remotePlayerController.addEventListener(
      castFramework.RemotePlayerEventType.IS_PAUSED_CHANGED,
      async () => {
        if (remotePlayer.isPaused) {
          pauseAll({ local: true });
        } else if (state.value !== SoundState.PLAYING) {
          // If currently not playing, start to play.
          // This occurs if starting to play from local, but this check is
          // required if the state is changed remotely.
          await playPauseAll({ local: true });
        }
      }
    );

    remotePlayerController.addEventListener(
      castFramework.RemotePlayerEventType.MEDIA_INFO_CHANGED,
      async ({ value }) => {
        if (value && isStopped.value) {
          let waitMs = 100;
          while (sounds.value.length === 0) {
            console.warn(`Sounds not loaded. Waiting ${waitMs}ms.`);
            // eslint-disable-next-line no-await-in-loop
            await wait(waitMs);
            waitMs *= 2;
          }
          const preset = new Preset();
          preset.mixUrl = value.contentId;
          await Promise.all(
            preset.sounds.map((savedSound) => {
              const sound = soundById(savedSound.id);
              sound.volume = savedSound.volume;
              const fade = state.value === SoundState.STOPPED ? 250 : false;
              return playStop({
                sound,
                fade,
                local: true,
              });
            })
          );
        }
      }
    );
  };

  const prefetch = async () => {
    return Promise.all(
      sounds.value.map(async (sound) => {
        sound.isLoading = true;
        try {
          const url = pb.getFileUrl(sound, sound.file);
          await fetch(url);
        } catch (error) {
          console.error(error);
        }
        sound.isLoading = false;
      })
    );
  };

  return {
    sounds,
    currentName,
    castConnected,
    soundsPlaying,
    soundsNotStopped,
    state,
    isPlaying,
    isPaused,
    isStopped,
    soundById,
    play,
    pause,
    stop,
    volume,
    castConnectedChanged,
    playStop,
    pauseAll,
    playPauseAll,
    stopAll,
    initializeCastApi,
    updateCast,
    initSounds,
    prefetch,
  };
});
