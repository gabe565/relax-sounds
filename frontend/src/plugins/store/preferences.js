import { defineStore } from "pinia";
import { ref } from "vue";

export const usePreferencesStore = defineStore(
  "preferences",
  () => {
    const shrinkLeftPanel = ref(false);

    return { shrinkLeftPanel };
  },
  {
    persist: true,
  },
);
