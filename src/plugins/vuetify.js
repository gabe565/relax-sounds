import Vue from 'vue';
import Vuetify from 'vuetify/lib';
import colors from 'vuetify/lib/util/colors';
import '@fortawesome/fontawesome-pro/css/fontawesome.css';
import '@fortawesome/fontawesome-pro/css/light.css';
import '@fortawesome/fontawesome-pro/css/solid.css';

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
  theme: {
    dark: true,
    themes: {
      dark: {
        primary: colors.deepOrange,
        accent: colors.deepPurple,
      },
      light: {
        primary: colors.deepOrange,
        secondary: '#E4E4E4',
        accent: colors.deepPurple,
      },
    },
  },
});
