import Vue from 'vue';
import Vuetify from 'vuetify/lib';
import '@fortawesome/fontawesome-pro/css/all.css';

Vue.use(Vuetify);

export default new Vuetify({
  icons: {
    iconfont: 'fa',
    values: {
      clear: 'fal fa-times',
      close: 'fal fa-times',
      menu: 'fal fa-bars',
    },
  },
});
