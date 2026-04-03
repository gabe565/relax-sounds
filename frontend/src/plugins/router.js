import { nextTick } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import { toast } from "vue-sonner";
import MixerIcon from "~icons/material-symbols/instant-mix-rounded";
import LogoutIcon from "~icons/material-symbols/logout-rounded";
import PresetsIcon from "~icons/material-symbols/playlist-play-rounded";
import SoundsIcon from "~icons/material-symbols/sound-detection-loud-sound-rounded";
import { useAuth } from "@/composables/useAuth.js";
import { ApiPath } from "@/config/api";
import { getErrorMessage } from "@/plugins/pocketbase.js";
import { usePresetsStore } from "@/plugins/store/presets";
import { Preset } from "@/util/Preset";
import { wait } from "@/util/helpers";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/sounds",
      name: "Sounds",
      component: () => import("@/pages/SoundsPage.vue"),
      props: true,
      meta: {
        icon: SoundsIcon,
        showInNav: true,
      },
    },
    {
      path: "/mixer",
      name: "Mixer",
      component: () => import("@/pages/MixerPage.vue"),
      props: true,
      meta: {
        icon: MixerIcon,
        showInNav: true,
      },
    },
    {
      path: "/presets",
      name: "Presets",
      component: () => import("@/pages/PresetsPage.vue"),
      props: true,
      meta: {
        icon: PresetsIcon,
        showInNav: true,
      },
    },
    {
      path: "/login",
      name: "Login",
      component: () => import("@/pages/LoginPage.vue"),
      meta: {
        showInNav: false,
        guestOnly: true,
        hideLogin: true,
      },
    },
    {
      path: "/register",
      name: "Register",
      component: () => import("@/pages/LoginPage.vue"),
      props: { register: true },
      meta: {
        showInNav: false,
        guestOnly: true,
        hideLogin: true,
      },
    },
    {
      path: "/forgot-password",
      name: "Forgot Password",
      component: () => import("@/pages/ForgotPasswordPage.vue"),
      meta: {
        showInNav: false,
        guestOnly: true,
        hideLogin: true,
      },
    },
    {
      path: "/logout",
      name: "Logout",
      redirect: () => {
        useAuth().logout();
        return { name: "Sounds" };
      },
      meta: {
        icon: LogoutIcon,
        showInNav: false,
        authOnly: true,
      },
    },
    {
      path: "/import/:name/:songs",
      redirect: ({ params }) => {
        (async () => {
          try {
            await wait(500);
            const preset = new Preset();
            preset.encodedName = params.name;
            await preset.setEncodedShorthand(params.songs);
            await usePresetsStore().add({ preset });
            toast.success(`Imported ${preset.name}.`);
          } catch (err) {
            console.error(err);
            toast.error(`Failed to import preset:\n${getErrorMessage(err)}`);
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
      path: "/_",
      redirect() {
        window.location.href = ApiPath("/_");
      },
    },
    {
      path: "/:catchAll(.*)",
      name: "404 Not Found",
      component: () => import("@/pages/NotFoundPage.vue"),
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
