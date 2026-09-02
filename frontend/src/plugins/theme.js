import { kebabCase } from "lodash-es";

// Material Theme Builder export, seed #BB86FC. Paste new light/dark schemes here.
const schemes = {
  light: {
    primary: "#6C538B",
    onPrimary: "#FFFFFF",
    primaryContainer: "#EEDBFF",
    onPrimaryContainer: "#543B72",
    secondary: "#655A6F",
    onSecondary: "#FFFFFF",
    secondaryContainer: "#ECDDF7",
    onSecondaryContainer: "#4D4357",
    tertiary: "#805159",
    onTertiary: "#FFFFFF",
    tertiaryContainer: "#FFD9DE",
    onTertiaryContainer: "#653B42",
    error: "#BA1A1A",
    onError: "#FFFFFF",
    errorContainer: "#FFDAD6",
    onErrorContainer: "#93000A",
    background: "#FFF7FF",
    onBackground: "#1D1A20",
    surface: "#FFF7FF",
    onSurface: "#1D1A20",
    surfaceVariant: "#E8E0EB",
    onSurfaceVariant: "#4A454E",
    surfaceTint: "#6C538B",
    surfaceDim: "#DFD8E0",
    surfaceBright: "#FFF7FF",
    surfaceContainerLowest: "#FFFFFF",
    surfaceContainerLow: "#F9F1F9",
    surfaceContainer: "#F3EBF3",
    surfaceContainerHigh: "#EDE6EE",
    surfaceContainerHighest: "#E8E0E8",
    outline: "#7B757F",
    outlineVariant: "#CCC4CF",
    inverseSurface: "#332F35",
    inverseOnSurface: "#F6EEF6",
    inversePrimary: "#D8BAFA",
    shadow: "#000000",
    scrim: "#000000",
  },
  dark: {
    primary: "#D8BAFA",
    onPrimary: "#3C245A",
    primaryContainer: "#543B72",
    onPrimaryContainer: "#EEDBFF",
    secondary: "#CFC1DA",
    onSecondary: "#362D40",
    secondaryContainer: "#4D4357",
    onSecondaryContainer: "#ECDDF7",
    tertiary: "#F2B7C0",
    onTertiary: "#4B252C",
    tertiaryContainer: "#653B42",
    onTertiaryContainer: "#FFD9DE",
    error: "#FFB4AB",
    onError: "#690005",
    errorContainer: "#93000A",
    onErrorContainer: "#FFDAD6",
    background: "#151218",
    onBackground: "#E8E0E8",
    surface: "#151218",
    onSurface: "#E8E0E8",
    surfaceVariant: "#4A454E",
    onSurfaceVariant: "#CCC4CF",
    surfaceTint: "#D8BAFA",
    surfaceDim: "#151218",
    surfaceBright: "#3C383E",
    surfaceContainerLowest: "#100D12",
    surfaceContainerLow: "#1D1A20",
    surfaceContainer: "#221E24",
    surfaceContainerHigh: "#2C292F",
    surfaceContainerHighest: "#373339",
    outline: "#958E98",
    outlineVariant: "#4A454E",
    inverseSurface: "#E8E0E8",
    inverseOnSurface: "#332F35",
    inversePrimary: "#6C538B",
    shadow: "#000000",
    scrim: "#000000",
  },
};

const status = {
  light: {
    success: "#3F6B3F",
    onSuccess: "#FFFFFF",
    info: "#37618E",
    onInfo: "#FFFFFF",
    warning: "#7D5700",
    onWarning: "#FFFFFF",
  },
  dark: {
    success: "#A5D3A0",
    onSuccess: "#0C3910",
    info: "#A2C9FE",
    onInfo: "#00325B",
    warning: "#F2BF48",
    onWarning: "#412D00",
  },
};

const vuetifySurfaceAliases = (scheme) => ({
  "surface-light": scheme.surfaceContainerHigh,
  "primary-darken-1": scheme.onPrimaryContainer,
  "secondary-darken-1": scheme.onSecondaryContainer,
});

const colors = (name) =>
  Object.fromEntries(
    Object.entries({ ...schemes[name], ...status[name] }).map(([key, value]) => [
      kebabCase(key),
      value,
    ]),
  );

const variables = (name) => ({
  "border-color": schemes[name].outlineVariant,
  "border-opacity": 1,
  "shadow-color": schemes[name].shadow,
  "theme-kbd": schemes[name].surfaceContainerHighest,
  "theme-on-kbd": schemes[name].onSurface,
  "theme-code": schemes[name].surfaceContainer,
  "theme-on-code": schemes[name].onSurfaceVariant,
});

export const themes = {
  light: {
    dark: false,
    colors: { ...colors("light"), ...vuetifySurfaceAliases(schemes.light) },
    variables: variables("light"),
  },
  dark: {
    dark: true,
    colors: { ...colors("dark"), ...vuetifySurfaceAliases(schemes.dark) },
    variables: variables("dark"),
  },
};

export const colorKeys = Object.keys(themes.dark.colors);
