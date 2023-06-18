import { defineStore } from "pinia";
import { ref } from "vue";

export const useAlertStore = defineStore("alert", () => {
  const show = ref(false);
  const message = ref("");
  const type = ref("");

  return { show, message, type };
});
