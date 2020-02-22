import Vue from 'vue';
import Vuex from 'vuex';
import App from './App.vue';
import './registerServiceWorker';
import vuetify from './plugins/vuetify';
import storeData from './store/main';

Vue.config.productionTip = false;

Vue.use(Vuex);

const store = new Vuex.Store(storeData);

new Vue({
  vuetify,
  store,
  render: (h) => h(App),
}).$mount('#app');
