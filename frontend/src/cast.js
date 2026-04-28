import { ref } from "vue";
import { usePlayer } from "@/plugins/store/player";
import { wait } from "@/util/helpers";

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

    await usePlayer().initializeCastApi();
    castEnabled.value = true;
  }
};
