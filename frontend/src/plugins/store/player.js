import * as _ from "lodash-es";
import pLimit from "p-limit";
import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { toast } from "vue-sonner";
import { ApiPath } from "@/config/api";
import { usePocketBase } from "@/plugins/store/pocketbase.js";
import { Preset } from "@/util/Preset";
import { SoundState } from "@/util/Sound";
import { Filetype } from "@/util/filetype";
import { formatError, getCastSession } from "@/util/googleCast";

export const usePlayer = defineStore("player", () => {
  const pb = usePocketBase();
  const sounds = ref([]);
  const currentName = ref();
  let remotePlayer;
  let remotePlayerController;
  const castConnected = ref(false);

  const soundsReady = pb.loadSounds().then((data) => {
    sounds.value = data;
  });

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

  const play = async ({ sound, fade = 250 }) => {
    if (!castConnected.value && sound.isUnloaded) {
      try {
        await sound.load();
      } catch (err) {
        console.error(err);
        toast.error(`Failed to load ${sound.name}:\n${err}`);
        return;
      }
    }
    if (sound.isPaused) {
      fade = false;
    }
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

  const rate = ({ sound, value }) => {
    sound.rate = value;
  };

  const pan = ({ sound, value }) => {
    sound.pan = value;
  };

  const castConnectedChanged = ({ value }) => {
    castConnected.value = value;
  };

  const syncFromCast = async (mediaInfo) => {
    await soundsReady;
    const preset = new Preset();
    await preset.setMixUrl(mediaInfo.contentId);
    await Promise.all(
      preset.sounds.map((savedSound) => {
        const sound = soundById(savedSound.id);
        sound.volume = savedSound.volume || 1;
        sound.rate = savedSound.rate || 1;
        sound.pan = savedSound.pan || 0;
        const fade = state.value === SoundState.STOPPED ? 250 : false;
        return playStop({
          sound,
          fade,
          local: true,
        });
      }),
    );
  };

  const updateCast = _.debounce(
    async () => {
      if (isPlaying.value) {
        const castSession = getCastSession();
        if (castSession) {
          const { cast } = window.chrome;
          const preset = new Preset({ sounds: soundsPlaying.value });

          const mixUrl = await preset.mixUrlAs(Filetype.Mp3);
          const mediaInfo = new cast.media.MediaInfo(mixUrl, "audio/mpeg");
          mediaInfo.streamType = cast.media.StreamType.LIVE;
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
            new cast.Image(
              `${window.location.origin}/img/icons/android-chrome-maskable-512x512.png`,
            ),
          ];

          const queue = new cast.media.QueueLoadRequest([new cast.media.QueueItem(mediaInfo)]);
          queue.repeatMode = cast.media.RepeatMode.SINGLE;

          const request = new cast.media.LoadRequest(mediaInfo);
          request.queueData = queue;

          try {
            await castSession.loadMedia(request);
          } catch (error) {
            console.error(`Remote media load error: ${formatError(error)}`);
            toast.error(`Failed to cast:\n${error}`);
            return;
          }
        }
      } else {
        await stopCast();
      }
    },
    1000,
    { leading: true },
  );

  const stopCast = async () => {
    if (remotePlayerController) {
      remotePlayerController.stop();

      if (castConnected.value && isStopped.value) {
        await deleteStream();
      }
    }
  };

  const deleteStream = async () => {
    const uuid = sessionStorage.getItem("uuid");
    if (uuid) {
      try {
        const resp = await fetch(ApiPath(`/api/mix/${uuid}`), { method: "DELETE" });
        if (!resp.ok && resp.status !== 404) {
          console.error(`Remote media stop error: ${resp.statusText}`);
          toast.error(`Failed to stop cast:\n${resp.statusText}`);
        }
      } catch (error) {
        console.error(`Remote media stop error: ${formatError(error)}`);
        toast.error(`Failed to stop cast:\n${error}`);
      }
    }
  };

  const playStop = async ({ sound, fade = 250, local = false }) => {
    if (sound.state === SoundState.PLAYING) {
      stop({ sound, fade });
    } else {
      await play({ sound, fade });
    }
    currentName.value = null;
    if (!local && castConnected.value) {
      updateCast();
    }
  };

  const playPause = async ({ sound, fade = 250, local = false }) => {
    if (sound.state === SoundState.PLAYING) {
      pause({ sound, fade });
    } else {
      await play({ sound, fade });
    }
    currentName.value = null;
    if (!local && castConnected.value) {
      updateCast();
    }
  };

  const pauseAll = ({ local = false } = {}) => {
    soundsPlaying.value.forEach((sound) => {
      pause({ sound });
    });
    if (!local && castConnected.value && remotePlayer && !remotePlayer.isPaused) {
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
          await play({ sound, fade: 0 });
        }
      }),
    );
    if (!local && castConnected.value && remotePlayerController) {
      remotePlayerController.playOrPause();
    }
  };

  const stopAll = ({ fade = 250, local = false }) => {
    soundsNotStopped.value.forEach((sound) => {
      stop({ sound, fade });
    });
    if (!local && castConnected.value) {
      updateCast();
      currentName.value = null;
    }
  };

  const initializeCastApi = () => {
    const { framework: castFramework } = window.cast;
    const { cast } = window.chrome;

    remotePlayer = new castFramework.RemotePlayer();
    remotePlayerController = new castFramework.RemotePlayerController(remotePlayer);

    remotePlayerController.addEventListener(
      castFramework.RemotePlayerEventType.IS_CONNECTED_CHANGED,
      async ({ value }) => {
        const wasPlaying = state.value !== SoundState.STOPPED;
        if (wasPlaying) {
          pauseAll();
        }
        castConnectedChanged({ value });
        if (wasPlaying) {
          await playPauseAll({ local: true });
          updateCast();
          if (!value) {
            await deleteStream();
          }
        }
      },
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
      },
    );

    remotePlayerController.addEventListener(
      castFramework.RemotePlayerEventType.MEDIA_INFO_CHANGED,
      async ({ value }) => {
        if (value && isStopped.value) {
          await syncFromCast(value);
        }
      },
    );

    // Handle session resume (e.g. page refresh while casting).
    castFramework.CastContext.getInstance().addEventListener(
      castFramework.CastContextEventType.SESSION_STATE_CHANGED,
      async ({ sessionState }) => {
        if (sessionState === castFramework.SessionState.SESSION_RESUMED && isStopped.value) {
          // The SDK doesn't request media status on session resume,
          // so MEDIA_INFO_CHANGED never fires. Send an explicit
          // GET_STATUS to the receiver to trigger it.
          getCastSession()?.getSessionObj().sendMessage("urn:x-cast:com.google.cast.media", {
            type: "GET_STATUS",
            requestId: 1,
          });
        }
      },
    );

    // setOptions triggers auto-join, so it must come after all listeners
    // are registered. Otherwise RemotePlayer misses the media session.
    castFramework.CastContext.getInstance().setOptions({
      receiverApplicationId: cast.media.DEFAULT_MEDIA_RECEIVER_APP_ID,
      autoJoinPolicy: cast.AutoJoinPolicy.ORIGIN_SCOPED,
    });
  };

  const prefetch = async () => {
    const id = toast.loading("Preloading sounds...");
    const limit = pLimit(8);
    try {
      await Promise.all(
        sounds.value.map((sound) =>
          limit(async () => {
            sound.isLoading = true;
            await sound.load();
            sound.howl.unload();
            sound.isLoading = false;
          }),
        ),
      );
      toast.success("Preloaded all sounds.", { id });
    } catch (err) {
      console.error(err);
      toast.error(`Failed to preload sounds.\n${err}`, { id });
    }
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
    rate,
    pan,
    castConnectedChanged,
    playStop,
    playPause,
    pauseAll,
    playPauseAll,
    stopAll,
    initializeCastApi,
    updateCast,
    prefetch,
  };
});
