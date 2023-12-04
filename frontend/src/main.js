import { createApp } from "vue";
import { createPinia } from "pinia";
import piniaPluginPersistedstate from "pinia-plugin-persistedstate";
import Vue3Toastify, { toast } from "vue3-toastify";
import "vue3-toastify/dist/index.css";
import vuetify from "./plugins/vuetify";
import router from "./plugins/router";
import App from "./App.vue";

import "./cast";
import "./plugins/buffer";

const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);

const toastOptions = {
  position: toast.POSITION.BOTTOM_RIGHT,
};

createApp(App).use(pinia).use(vuetify).use(router).use(Vue3Toastify, toastOptions).mount("#app");
