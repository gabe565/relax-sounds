import { nextTick } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import { Preset } from "../util/Preset";
import { usePresetsStore } from "./store/presets";
import SoundsIcon from "~icons/material-symbols/sound-detection-loud-sound-rounded";
import PresetsIcon from "~icons/material-symbols/playlist-play-rounded";
import MixerIcon from "~icons/material-symbols/instant-mix-rounded";
import { useToast } from "vue-toastification";
import { wait } from "../util/helpers";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/sounds",
      name: "Sounds",
      component: () => import("../pages/SoundsPage.vue"),
      props: true,
      meta: {
        icon: SoundsIcon,
        showInNav: true,
      },
    },
    {
      path: "/mixer",
      name: "Mixer",
      component: () => import("../pages/MixerPage.vue"),
      props: true,
      meta: {
        icon: MixerIcon,
        showInNav: true,
      },
    },
    {
      path: "/presets",
      name: "Presets",
      component: () => import("../pages/PresetsPage.vue"),
      props: true,
      meta: {
        icon: PresetsIcon,
        showInNav: true,
      },
    },
    {
      path: "/import/:name/:songs",
      redirect: ({ params }) => {
        (async () => {
          const toast = useToast();
          try {
            await wait(500);
            const preset = new Preset({ new: true });
            preset.encodedName = params.name;
            await preset.setEncodedShorthand(params.songs);
            usePresetsStore().add({ preset });
            toast.success(`Imported ${preset.name}.`);
          } catch (error) {
            console.error(error);
            toast.error(`Failed to import preset:\n${error}`);
          }
        })();
        return { name: "Presets" };
      },
    },
    {
      path: "/",
      redirect: { name: "Sounds" },
    },
    {
      path: "/:catchAll(.*)",
      name: "404 Not Found",
      component: () => import("../pages/NotFoundPage.vue"),
    },
  ],
});

const defaultTitle = document.title;
router.afterEach(async (to) => {
  // Use next tick to handle router history correctly
  // see: https://github.com/vuejs/vue-router/issues/914#issuecomment-384477609
  await nextTick();
  document.title = to.name ? `${to.name} - ${defaultTitle}` : defaultTitle;
});

export default router;
