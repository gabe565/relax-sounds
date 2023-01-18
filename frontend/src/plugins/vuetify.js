import { createVuetify } from 'vuetify';
import { aliases, fa } from 'vuetify/iconsets/fa';
import colors from 'vuetify/lib/util/colors';
import '@fortawesome/fontawesome-pro/css/fontawesome.css';
import '@fortawesome/fontawesome-pro/css/light.css';
import '@fortawesome/fontawesome-pro/css/regular.css';
import '@fortawesome/fontawesome-pro/css/solid.css';

import 'vuetify/styles';

export default createVuetify({
  icons: {
    defaultSet: 'fa',
    aliases,
    sets: {
      fa,
    },
  },
  theme: {
    defaultTheme: 'dark',
    themes: {
      dark: {
        dark: true,
        colors: {
          background: '#150B29',
          primary: colors.deepOrange.base,
          accent: colors.deepPurple.base,
        },
      },
      light: {
        colors: {
          primary: colors.deepOrange.base,
          secondary: '#E4E4E4',
          accent: colors.deepPurple.base,
        },
      },
    },
  },
});
