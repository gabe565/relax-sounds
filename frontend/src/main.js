import { createPinia } from "pinia";
import piniaPluginPersistedstate from "pinia-plugin-persistedstate";
import { createApp } from "vue";
import Toast, { POSITION } from "vue-toastification";
import "vue-toastification/dist/index.css";
import App from "./App.vue";
import "./cast";
import router from "./plugins/router";
import vuetify from "./plugins/vuetify";

const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);

const toastOptions = {
  position: POSITION.TOP_RIGHT,
  containerClassName: "mt-16 pt-3",
  transition: "Vue-Toastification__fade",
};

createApp(App).use(pinia).use(vuetify).use(router).use(Toast, toastOptions).mount("#app");
