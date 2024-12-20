import { usePlayerStore } from "./plugins/store/player";
import { wait } from "./util/helpers";
import { ref } from "vue";

export const castEnabled = ref(false);

window.__onGCastApiAvailable = async (isAvailable) => {
  if (isAvailable) {
    // Workaround for __onGCastApiAvailable called before window.cast is set
    let waitMs = 100;
    while (!window.cast) {
      console.warn(`Cast is undefined. Retrying setup in ${waitMs}ms.`);
      await wait(waitMs);
      waitMs *= 2;
    }

    castEnabled.value = true;
    await usePlayerStore().initializeCastApi();
  }
};
