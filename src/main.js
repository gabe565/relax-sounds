import Vue from 'vue';
import Vuex from 'vuex';
import VueRouter from 'vue-router';

import App from './App.vue';
import './registerServiceWorker';
import vuetify from './plugins/vuetify';
import storeData from './plugins/store/main';
import routes from './plugins/routes';

Vue.config.productionTip = false;

Vue.use(Vuex);
Vue.use(VueRouter);

const store = new Vuex.Store(storeData);

const router = new VueRouter({
  routes,
  mode: 'history',
});

new Vue({
  vuetify,
  store,
  router,
  render: (h) => h(App),
}).$mount('#app');
