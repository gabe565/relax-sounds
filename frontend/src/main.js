import Vue from 'vue';
import './registerServiceWorker';
import vuetify from './plugins/vuetify';
import store from './plugins/store/main';
import router from './plugins/routes';
import App from './App.vue';

Vue.config.productionTip = false;

Vue.config.ignoredElements = ['google-cast-launcher'];

window.vue = new Vue({
  vuetify,
  store,
  router,
  render: (h) => h(App),
}).$mount('#app');

require('./cast');
