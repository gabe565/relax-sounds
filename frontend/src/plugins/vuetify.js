import { createVuetify } from "vuetify";
import { aliases, mdi } from "vuetify/iconsets/mdi-svg";
import colors from "vuetify/util/colors";

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
          cardBackground: "#271A40",
          newPresetCardBackground: colors.deepPurple.darken2,
          primary: colors.deepOrange.base,
          accent: colors.deepPurple.base,
        },
      },
      light: {
        colors: {
          cardBackground: "#F4F4F4",
          newPresetCardBackground: colors.deepPurple.lighten1,
          primary: colors.deepOrange.base,
          secondary: "#E4E4E4",
          accent: colors.deepPurple.base,
        },
      },
    },
  },
});
