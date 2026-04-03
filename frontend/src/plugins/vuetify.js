import { createVuetify } from "vuetify";

export default createVuetify({
  display: {
    mobileBreakpoint: "md",
    thresholds: {
      xs: 0,
      sm: 600,
      md: 840,
      lg: 1145,
      xl: 1545,
      xxl: 2138,
    },
  },
  theme: {
    defaultTheme: "dark",
    utilities: false,
    themes: {
      dark: {
        dark: true,
        colors: {
          background: "#160E27",
          surface: "#251842",
          "card-background": "#332457",
          primary: "#BB86FC",
          secondary: "#FFB74D",
          accent: "#7C4DFF",
        },
      },
      light: {
        colors: {
          background: "#F5F1FA",
          "card-background": "#FFF",
          primary: "#7C4DFF",
          secondary: "#FB8C00",
          accent: "#7C4DFF",
        },
      },
    },
  },
});
