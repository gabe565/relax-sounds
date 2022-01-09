import Vue from 'vue';
import VueRouter from 'vue-router';
import store from '../plugins/store/main';
import { decode } from '../util/shareUrl';
import Sounds from '../pages/Sounds.vue';
import Playlists from '../pages/Playlists.vue';
import NotFound from '../pages/NotFound.vue';

Vue.use(VueRouter);

export default new VueRouter({
  mode: 'history',
  routes: [
    {
      path: '/sounds',
      name: 'Sounds',
      component: Sounds,
      meta: {
        icon: 'fa-volume',
        showInNav: true,
      },
    },
    {
      path: '/playlists',
      name: 'Playlists',
      component: Playlists,
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
          const playlist = { ...decode(params), new: true };
          store.commit('playlists/add', { playlist });
        } catch (error) {
          redirectParams = {
            alert: {
              type: 'error',
              text: 'Could not import playlist. Please try again later.',
            },
          };
        }
        return { name: 'Playlists', params: redirectParams };
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
