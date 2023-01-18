import { nextTick } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';
import store from '../plugins/store/main';
import SoundsPage from '../pages/SoundsPage.vue';
import PresetsPage from '../pages/PresetsPage.vue';
import NotFoundPage from '../pages/NotFoundPage.vue';
import { Preset } from '../util/Preset';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/sounds',
      name: 'Sounds',
      component: SoundsPage,
      props: true,
      meta: {
        icon: 'fa-volume',
        showInNav: true,
      },
    },
    {
      path: '/presets',
      name: 'Presets',
      component: PresetsPage,
      props: true,
      meta: {
        icon: 'fa-list-music',
        showInNav: true,
      },
    },
    {
      path: '/import/:name/:songs',
      redirect: ({ params }) => {
        let redirectParams;
        try {
          const preset = new Preset({ new: true });
          preset.encodedName = params.name;
          preset.encodedShorthand = params.songs;
          store.commit('presets/add', { preset });
        } catch (error) {
          console.error(error);
          redirectParams = {
            alert: {
              type: 'error',
              text: 'Could not import preset. Please try again later.',
            },
          };
        }
        return { name: 'Presets', params: redirectParams };
      },
    },
    {
      path: '/',
      redirect: { name: 'Sounds' },
    },
    {
      path: '/:catchAll(.*)',
      name: '404 Not Found',
      component: NotFoundPage,
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
