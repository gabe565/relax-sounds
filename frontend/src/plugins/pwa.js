import { ref } from "vue";
import { useToast } from "vue-toastification";
import UpdateToast from "@/components/UpdateToast.vue";

const toast = useToast();
const updateSW = ref(undefined);
const intervalMS = 60 * 60 * 1000;

export const registerSW = async () => {
  try {
    const { registerSW } = await import("virtual:pwa-register");
    updateSW.value = registerSW({
      immediate: true,
      onOfflineReady() {
        toast.info("App ready to work offline.", { timeout: 3000 });
      },
      onNeedRefresh() {
        toast.info(
          {
            component: UpdateToast,
            listeners: { refresh: () => updateSW.value(true) },
          },
          { timeout: false, closeOnClick: false },
        );
      },
      onRegistered(swRegistration) {
        if (swRegistration) {
          setInterval(() => swRegistration.update(), intervalMS);
        }
      },
      onRegisterError(error) {
        console.error(error);
      },
    });
  } catch {
    console.log("PWA disabled.");
  }
};
