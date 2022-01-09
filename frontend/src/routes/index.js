import Vue from 'vue';
import VueRouter from 'vue-router';
import store from '../plugins/store/main';
import { decode } from '../util/shareUrl';
import SoundsPage from '../pages/SoundsPage.vue';
import PlaylistsPage from '../pages/PlaylistsPage.vue';
import NotFoundPage from '../pages/NotFoundPage.vue';

Vue.use(VueRouter);

export default new VueRouter({
  mode: 'history',
  routes: [
    {
      path: '/sounds',
      name: 'Sounds',
      component: SoundsPage,
      meta: {
        icon: 'fa-volume',
        showInNav: true,
      },
    },
    {
      path: '/playlists',
      name: 'Playlists',
      component: PlaylistsPage,
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
      component: NotFoundPage,
    },
  ],
});
