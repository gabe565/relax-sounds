import SoundsPage from '../components/pages/SoundsPage.vue';
import PlaylistsPage from '../components/pages/PlaylistsPage.vue';
import NotFoundPage from '../components/pages/NotFoundPage.vue';

export default [
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
    path: '/',
    redirect: { name: 'Sounds' },
  },
  {
    path: '*',
    component: NotFoundPage,
  },
];
