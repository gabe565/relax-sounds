/* eslint-disable no-unused-vars */
import { getSounds } from '../data/sounds';
import { getTags } from '../data/tags';

export default {
  name: 'useRegisterSW',
  data() {
    return {
      updateSW: undefined,
      offlineReady: false,
      needRefresh: false,
    };
  },
  async mounted() {
    try {
      const { registerSW } = await import('virtual:pwa-register');
      const vm = this;
      this.updateSW = registerSW({
        immediate: true,
        onOfflineReady() {
          vm.offlineReady = true;
          vm.onOfflineReadyFn();
        },
        onNeedRefresh() {
          vm.needRefresh = true;
          vm.onNeedRefreshFn();
        },
        onRegistered(swRegistration) {
          if (swRegistration) {
            vm.handleSWManualUpdates(swRegistration);
          }
        },
        onRegisterError(e) {
          vm.handleSWRegisterError(e);
        },
      });
    } catch {
      console.log('PWA disabled.');
    }
  },
  methods: {
    async closePromptUpdateSW() {
      this.offlineReady = false;
      this.needRefresh = false;
    },
    async onOfflineReadyFn() {
      try {
        await getSounds(true);
        await getTags(true);
      } catch (error) {
        console.error(error);
      }
      console.log('onOfflineReady');
    },
    onNeedRefreshFn() {
      console.log('onNeedRefresh');
    },
    updateServiceWorker() {
      if (this.updateSW) {
        this.updateSW(true);
      }
    },
    handleSWManualUpdates(swRegistration) {},
    handleSWRegisterError(error) {
      console.error(error);
    },
  },
};
