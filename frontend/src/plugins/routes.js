import Vue from 'vue';
import VueRouter from 'vue-router';
import store from './store/main';
import { decode } from '../util/shareUrl';
import SoundsPage from '../components/pages/SoundsPage.vue';
import PlaylistsPage from '../components/pages/PlaylistsPage.vue';
import NotFoundPage from '../components/pages/NotFoundPage.vue';

Vue.use(VueRouter);

export default new VueRouter({
  mode: 'history',
  routes: [
    {
      path: '/sounds',
      name: 'Sounds',
      component: SoundsPage,
      meta: {
        icon: 'fa-speaker',
      },
    },
    {
      path: '/playlists',
      name: 'Playlists',
      component: PlaylistsPage,
      meta: {
        icon: 'fa-list-music',
      },
    },
    {
      path: '/import/:name/:songs',
      redirect: ({ params }) => {
        const playlist = { ...decode(params), new: true };
        store.commit('playlists/add', { playlist });
        return { name: 'Playlists' };
      },
    },
    {
      path: '/',
      redirect: { name: 'Sounds' },
    },
    {
      path: '*',
      component: NotFoundPage,
    },
  ],
});
