import { defineStore } from "pinia";
import Fuse from "fuse.js";
import { computed, reactive, ref, watch } from "vue";
import { usePlayerStore } from "./player";

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
    page: 1,
  });

  const filteredSounds = computed(() => {
    let result;
    if (filters.value.word) {
      result = reactive(fuse.search(filters.value.word).map((e) => e.item));
    } else {
      result = usePlayerStore().sounds;
    }
    return result;
  });

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
  };
});
