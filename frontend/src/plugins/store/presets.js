import { defineStore } from "pinia";
import { computed, ref, watch } from "vue";
import { toast } from "vue-sonner";
import { usePlayer } from "@/plugins/store/player";
import { getErrorMessage, usePocketBase } from "@/plugins/store/pocketbase.js";
import { Preset } from "@/util/Preset";
import { SoundState } from "@/util/Sound";

const Version = 3;

export const usePresets = defineStore(
  "presets",
  () => {
    const pb = usePocketBase();
    const player = usePlayer();
    const version = ref(Version);
    const presets = ref([]);
    const isSyncing = ref(false);
    let syncPromise = null;

    const add = async ({ preset, playing = true, sync = true }) => {
      for (const sound of preset.sounds) {
        if (typeof sound.id === "number") {
          sound.id = sound.id.toString();
        }
      }
      const newPreset = new Preset(preset);

      presets.value.push(newPreset);
      if (playing) {
        player.currentName = preset.name;
      }

      if (sync) {
        await performSync();
      }
    };

    const active = computed(() => {
      return presets.value.filter((preset) => !preset.hidden);
    });

    const hide = ({ preset }) => {
      preset.hidden = true;
    };

    const hideAll = () => {
      for (const preset of presets.value) {
        preset.hidden = true;
      }
    };

    const unhide = ({ preset }) => {
      delete preset.hidden;
    };

    const unhideAll = () => {
      for (const preset of presets.value) {
        delete preset.hidden;
      }
    };

    const remove = async ({ preset, sync = true }) => {
      const index = presets.value.indexOf(preset);
      if (index !== -1) {
        if (sync && pb.isAuthenticated && preset.synced) {
          try {
            await pb.client.collection("presets").delete(preset.id);
          } catch (e) {
            if (e.status !== 404) {
              throw e;
            }
          }
        }

        presets.value.splice(index, 1);
      }
    };

    const removeHidden = async ({ sync = true } = {}) => {
      const toRemove = presets.value.filter((preset) => preset.hidden);

      if (sync && pb.isAuthenticated) {
        const syncedToRemove = toRemove.filter((p) => p.synced);
        if (syncedToRemove.length > 0) {
          const batch = pb.createBatch();
          for (const preset of syncedToRemove) {
            batch.collection("presets").delete(preset.id);
          }
          await batch.send();
        }
      }

      presets.value = presets.value.filter((preset) => !preset.hidden);
    };

    const savePlaying = async ({ name }) => {
      const sounds = player.soundsPlaying.map((sound) => ({
        id: sound.id,
        volume: sound.volume,
        rate: sound.rate,
        pan: sound.pan,
      }));

      await add({
        preset: { name, sounds },
      });
    };

    const play = async ({ preset }) => {
      if (player.state !== SoundState.STOPPED) {
        player.stopAll({ fade: 0, local: true });
      }
      await Promise.all(
        preset.sounds.map((savedSound) => {
          const sound = player.soundById(savedSound.id);
          sound.volume = savedSound.volume || 1;
          sound.rate = savedSound.rate || 1;
          sound.pan = savedSound.pan || 0;
          const fade = player.state === SoundState.STOPPED ? 500 : false;
          return player.playStop({ sound, fade, local: true });
        }),
      );
      player.currentName = preset.name;
      if (preset.new) {
        preset.new = false;
      }
      player.updateCast();
    };

    const migrate = async () => {
      if (version.value === 2) {
        await Promise.all(presets.value.map((preset) => preset.migrate()));
        version.value++;
      }
    };

    const performSync = async () => {
      if (!pb.isAuthenticated) {
        return;
      }
      if (syncPromise) return syncPromise;

      syncPromise = (async () => {
        isSyncing.value = true;
        try {
          const batch = pb.client.createBatch();
          const remotePresets = await pb.client.collection("presets").getFullList();
          const localPresets = presets.value;

          // Resolve remote to local
          for (const remote of remotePresets) {
            const remoteMetadata = remote.metadata || {};
            const remoteSounds = (remote.sounds || []).map((id) => ({
              id,
              volume: remoteMetadata[id]?.volume ?? 1,
              pan: remoteMetadata[id]?.pan ?? 0,
              rate: remoteMetadata[id]?.rate ?? 1,
            }));

            let local = localPresets.find((p) => p.id === remote.id);
            if (!local) {
              localPresets.push(
                new Preset({
                  id: remote.id,
                  name: remote.name,
                  sounds: remoteSounds,
                  synced: true,
                }),
              );
            } else {
              local.synced = true;
              if (
                local.name !== remote.name ||
                JSON.stringify(local.sounds) !== JSON.stringify(remoteSounds)
              ) {
                local.name = remote.name;
                local.sounds = remoteSounds;
              }
            }
          }

          // Resolve local to remote
          const idsToUpdate = [];
          const idsToRemove = [];
          for (const local of localPresets) {
            const remote = remotePresets.find((p) => p.id === local.id);
            if (!remote) {
              if (local.synced) {
                idsToRemove.push(local.id);
              } else {
                const metadata = {};
                for (const sound of local.sounds) {
                  metadata[sound.id] = {
                    volume: sound.volume,
                    pan: sound.pan,
                    rate: sound.rate,
                  };
                }

                batch.collection("presets").create({
                  name: local.name,
                  user: pb.user.id,
                  sounds: local.sounds.map((s) => s.id),
                  metadata: metadata,
                });
                idsToUpdate.push(local.id);
              }
            }
          }

          if (idsToRemove.length > 0) {
            presets.value = presets.value.filter((p) => !idsToRemove.includes(p.id));
          }

          if (idsToUpdate.length > 0) {
            const response = await batch.send();

            for (const [key, entry] of response.entries()) {
              const localEntry = localPresets.find((p) => p.id === idsToUpdate[key]);
              localEntry.id = entry.body.id;
              localEntry.synced = true;
            }
          }
        } finally {
          isSyncing.value = false;
          syncPromise = null;
        }
      })();

      return syncPromise;
    };

    watch(
      () => pb.user,
      async () => {
        try {
          await performSync();
        } catch (err) {
          console.error("Failed to sync presets:", err);
          toast.error(`Failed to sync presets:\n${getErrorMessage(err)}`);
        }
      },
      { immediate: true },
    );

    return {
      presets,
      isSyncing,
      add,
      active,
      hide,
      hideAll,
      unhide,
      unhideAll,
      remove,
      removeHidden,
      savePlaying,
      play,
      migrate,
      sync: performSync,
    };
  },
  {
    persist: {
      pick: ["presets", "version"],
      afterHydrate(ctx) {
        if (ctx.store?.presets) {
          ctx.store.presets = ctx.store.presets.map((p) => new Preset(p)).filter((p) => !p.hidden);
        }
      },
    },
  },
);
