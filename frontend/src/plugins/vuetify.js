import { createVuetify } from "vuetify";
import { aliases, mdi } from "vuetify/iconsets/mdi-svg";
import colors from "vuetify/lib/util/colors";

import "vuetify/styles";

export default createVuetify({
  icons: {
    defaultSet: "mdi",
    aliases,
    sets: {
      mdi,
    },
  },
  theme: {
    defaultTheme: "dark",
    themes: {
      dark: {
        dark: true,
        colors: {
          background: "#150B29",
          primary: colors.deepOrange.base,
          accent: colors.deepPurple.base,
        },
      },
      light: {
        colors: {
          primary: colors.deepOrange.base,
          secondary: "#E4E4E4",
          accent: colors.deepPurple.base,
        },
      },
    },
  },
});
