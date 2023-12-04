import { createApp } from "vue";
import { createPinia } from "pinia";
import piniaPluginPersistedstate from "pinia-plugin-persistedstate";
import vuetify from "./plugins/vuetify";
import router from "./plugins/router";
import App from "./App.vue";

import "./cast";
import "./plugins/buffer";

const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);

createApp(App).use(pinia).use(vuetify).use(router).mount("#app");
