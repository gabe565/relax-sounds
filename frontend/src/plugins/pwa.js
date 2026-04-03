import { ref } from "vue";
import { toast } from "vue-sonner";

const updateSW = ref(undefined);
const intervalMS = 60 * 60 * 1000;

export const registerSW = async () => {
  try {
    const { registerSW } = await import("virtual:pwa-register");
    updateSW.value = registerSW({
      immediate: true,
      onOfflineReady() {
        toast.info("App ready to work offline.", { duration: 3000 });
      },
      onNeedRefresh() {
        toast.info("New content available, click on reload button to update.", {
          duration: Infinity,
          action: {
            label: "Reload",
            onClick: () => updateSW.value(true),
          },
        });
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
