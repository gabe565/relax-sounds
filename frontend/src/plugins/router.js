import { nextTick } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import { Preset } from "../util/Preset";
import { usePresetsStore } from "./store/presets";
import SoundsIcon from "~icons/solar/soundwave-bold";
import PresetsIcon from "~icons/solar/playlist-bold";

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
        let redirectParams;
        try {
          const preset = new Preset({ new: true });
          preset.encodedName = params.name;
          preset.encodedShorthand = params.songs;
          usePresetsStore().add({ preset });
        } catch (error) {
          console.error(error);
          redirectParams = {
            alert: {
              type: "error",
              text: "Could not import preset. Please try again later.",
            },
          };
        }
        return { name: "Presets", params: redirectParams };
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
