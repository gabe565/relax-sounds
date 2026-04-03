import { defineStore } from "pinia";
import { ref } from "vue";

export const Theme = {
  auto: "auto",
  dark: "dark",
  light: "light",
};

export const usePreferences = defineStore(
  "preferences",
  () => {
    const theme = ref(Theme.auto);
    const shrinkLeftPanel = ref(false);

    return { theme, shrinkLeftPanel };
  },
  {
    persist: true,
  },
);
