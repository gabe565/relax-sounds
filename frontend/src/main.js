import { createApp } from "vue";
import { createPinia } from "pinia";
import piniaPluginPersistedstate from "pinia-plugin-persistedstate";
import Toast, { POSITION } from "vue-toastification";
import "vue-toastification/dist/index.css";
import vuetify from "./plugins/vuetify";
import router from "./plugins/router";
import App from "./App.vue";

import "./cast";
import "./plugins/buffer";

const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);

const toastOptions = {
  position: POSITION.BOTTOM_RIGHT,
  showCloseButtonOnHover: true,
  transition: "Vue-Toastification__fade",
};

createApp(App).use(pinia).use(vuetify).use(router).use(Toast, toastOptions).mount("#app");
