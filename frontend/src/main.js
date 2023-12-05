import { createApp } from "vue";
import { createPinia } from "pinia";
import piniaPluginPersistedstate from "pinia-plugin-persistedstate";
import Toast, { POSITION } from "vue-toastification";
import "vue-toastification/dist/index.css";
import vuetify from "./plugins/vuetify";
import router from "./plugins/router";
import App from "./App.vue";

import "./cast";

const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);

const toastOptions = {
  position: POSITION.TOP_RIGHT,
  containerClassName: "mt-16 pt-3 mt-sm-0 pt-sm-0",
  transition: "Vue-Toastification__fade",
};

createApp(App).use(pinia).use(vuetify).use(router).use(Toast, toastOptions).mount("#app");
