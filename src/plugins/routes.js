import SoundsPage from '../components/pages/SoundsPage.vue';
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
    path: '/',
    redirect: { name: 'Sounds' },
  },
  {
    path: '*',
    component: NotFoundPage,
  },
];
