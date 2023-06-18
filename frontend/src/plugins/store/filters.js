import { defineStore } from "pinia";
import Fuse from "fuse.js";
import { computed, reactive, ref, watch } from "vue";
import { usePlayerStore } from "./player";

const PER_PAGE = 48;

export const fuse = new Fuse([], {
  shouldSort: true,
  threshold: 0.3,
  location: 0,
  distance: 100,
  maxPatternLength: 32,
  minMatchCharLength: 1,
  keys: ["name", "tags"],
});

export const useFiltersStore = defineStore("filters", () => {
  const filters = ref({
    word: "",
    playing: false,
    page: 1,
  });

  const filteredSounds = computed(() => {
    let result;
    if (filters.value.word) {
      result = reactive(fuse.search(filters.value.word).map((e) => e.item));
    } else {
      result = usePlayerStore().sounds;
    }
    if (filters.value.playing) {
      result = result.filter((e) => !e.isStopped);
    }
    return result;
  });

  const sounds = computed(() => {
    const result = filteredSounds.value;
    const offset = PER_PAGE * (filters.value.page - 1);
    return result?.slice(offset, offset + PER_PAGE);
  });

  const pages = computed(() => Math.max(Math.ceil(filteredSounds.value.length / PER_PAGE), 1));

  const updateSounds = (val = null) => {
    if (val === null) {
      val = usePlayerStore().sounds;
    }
    fuse.setCollection(val);
  };
  watch(() => usePlayerStore().sounds, updateSounds);
  updateSounds();

  return {
    filters,
    filteredSounds,
    sounds,
    pages,
  };
});
