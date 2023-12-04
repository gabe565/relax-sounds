import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { SoundState } from "../../util/Sound";
import { Preset } from "../../util/Preset";
import { usePlayerStore } from "./player";

let stateVersion = 0;
const version = 3;

const saveState = (presets) => {
  const state = { version, presets };
  localStorage.setItem("presets", JSON.stringify(state));
};

const loadState = () => {
  let state = JSON.parse(localStorage.getItem("presets"));

  if (!state) {
    // Playlist to preset migration
    const playlists = JSON.parse(localStorage.getItem("playlists"));
    if (playlists) {
      state = { version, presets: playlists.playlists };
      localStorage.setItem("presets", JSON.stringify(state));
      localStorage.removeItem("playlists");
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

    stateVersion = state.version;
    if (dirty) {
      saveState(state);
    }

    return state.presets.map((preset) => new Preset(preset)).filter((preset) => !preset.hidden);
  }
  return [];
};

export const usePresetsStore = defineStore("presets", () => {
  const presets = ref(loadState());

  const add = ({ preset, playing = true }) => {
    for (const sound of preset.sounds) {
      if (typeof sound.id === "number") {
        sound.id = sound.id.toString();
      }
    }
    presets.value.push(new Preset(preset));
    if (playing) {
      usePlayerStore().currentName = preset.name;
    }
    saveState(presets.value);
  };

  const active = computed(() => {
    return presets.value.filter((preset) => !preset.hidden);
  });

  const hide = ({ preset }) => {
    preset.hidden = true;
    saveState(presets.value);
  };

  const hideAll = () => {
    for (const preset of presets.value) {
      preset.hidden = true;
    }
    saveState(presets.value);
  };

  const unhide = ({ preset }) => {
    delete preset.hidden;
    saveState(presets.value);
  };

  const unhideAll = () => {
    for (const preset of presets.value) {
      delete preset.hidden;
    }
    saveState(presets.value);
  };

  const remove = ({ preset }) => {
    const index = presets.value.indexOf(preset);
    presets.value.splice(index, 1);
    saveState(presets.value);
  };

  const removeHidden = () => {
    presets.value = presets.value.filter((preset) => !preset.hidden);
    saveState(presets.value);
  };

  const savePlaying = ({ name }) => {
    const sounds = usePlayerStore().soundsPlaying.map((sound) => ({
      id: sound.id,
      volume: sound.volume,
      rate: sound.rate,
      pan: sound.pan,
    }));

    add({
      preset: {
        name,
        sounds,
        new: true,
      },
    });
  };

  const play = async ({ preset }) => {
    const playerStore = usePlayerStore();
    if (playerStore.state !== SoundState.STOPPED) {
      playerStore.stopAll({ fade: 0, local: true });
    }
    await Promise.all(
      preset.sounds.map((savedSound) => {
        const sound = playerStore.soundById(savedSound.id);
        sound.volume = savedSound.volume || 1;
        sound.rate = savedSound.rate || 1;
        sound.pan = savedSound.pan || 0;
        const fade = playerStore.state === SoundState.STOPPED ? 500 : false;
        return playerStore.playStop({ sound, fade, local: true });
      }),
    );
    playerStore.currentName = preset.name;
    if (preset.new) {
      preset.new = false;
      saveState(presets.value);
    }
    playerStore.updateCast();
  };

  const migrate = async () => {
    if (stateVersion === 2) {
      await Promise.all(presets.value.map((preset) => preset.migrate()));
      saveState(presets.value);
    }
  };

  return {
    presets,
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
  };
});
