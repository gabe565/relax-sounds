import { createVuetify } from "vuetify";
import "vuetify/styles";

export default createVuetify({
  display: {
    mobileBreakpoint: "md",
    thresholds: {
      md: 760,
    },
  },
  theme: {
    defaultTheme: "dark",
    themes: {
      dark: {
        dark: true,
        colors: {
          background: "#160E27",
          surface: "#251842",
          cardBackground: "#332457",
          primary: "#BB86FC",
          secondary: "#FFB74D",
          accent: "#7C4DFF",
        },
      },
      light: {
        colors: {
          background: "#F5F1FA",
          cardBackground: "#FFFFFF",
          primary: "#7C4DFF",
          secondary: "#FB8C00",
          accent: "#9C27B0",
        },
      },
    },
  },
});
