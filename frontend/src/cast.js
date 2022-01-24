import store from './plugins/store/main';

window.__onGCastApiAvailable = (isAvailable) => {
  if (isAvailable) {
    store.dispatch('player/initializeCastApi');
  }
};
