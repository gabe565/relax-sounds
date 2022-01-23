import Vue from 'vue';
import VueRouter from 'vue-router';
import store from '../plugins/store/main';
import { decode } from '../util/shareUrl';
import Sounds from '../pages/Sounds.vue';
import Presets from '../pages/Presets.vue';
import NotFound from '../pages/NotFound.vue';

Vue.use(VueRouter);

const router = new VueRouter({
  mode: 'history',
  routes: [
    {
      path: '/sounds',
      name: 'Sounds',
      component: Sounds,
      props: true,
      meta: {
        icon: 'fa-volume',
        showInNav: true,
      },
    },
    {
      path: '/presets',
      name: 'Presets',
      component: Presets,
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
          const preset = { ...decode(params), new: true };
          store.commit('presets/add', { preset });
        } catch (error) {
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
      path: '*',
      name: '404 Not Found',
      component: NotFound,
    },
  ],
});

const defaultTitle = document.title;
router.afterEach(async (to) => {
  // Use next tick to handle router history correctly
  // see: https://github.com/vuejs/vue-router/issues/914#issuecomment-384477609
  await Vue.nextTick();
  document.title = to.name ? `${to.name} - ${defaultTitle}` : defaultTitle;
});

export default router;