import store from './plugins/store/main';

// eslint-disable-next-line no-underscore-dangle
// eslint-disable-next-line dot-notation
window['__onGCastApiAvailable'] = (isAvailable) => {
  if (isAvailable) {
    store.dispatch('player/initializeCastApi');
  }
};
