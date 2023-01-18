import { createStore } from "vuex";
import filters from "./filters";
import presets from "./presets";
import player from "./player";

export default createStore({
  modules: {
    filters,
    presets,
    player,
  },
});
