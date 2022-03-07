import store from './plugins/store/main';
import { wait } from './util/helpers';

window.__onGCastApiAvailable = async (isAvailable) => {
  if (isAvailable) {
    // Workaround for __onGCastApiAvailable called before window.cast is set
    let waitMs = 100;
    while (!window.cast) {
      console.warn(`Cast is undefined. Retrying setup in ${waitMs}ms.`);
      // eslint-disable-next-line no-await-in-loop
      await wait(waitMs);
      waitMs *= 2;
    }

    await store.dispatch('player/initializeCastApi');
  }
};
