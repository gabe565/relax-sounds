import "./styles";
import { createPinia } from "pinia";
import piniaPluginPersistedstate from "pinia-plugin-persistedstate";
import { createApp } from "vue";
import App from "@/App.vue";
import "@/cast";
import router from "@/plugins/router";
import vuetify from "@/plugins/vuetify";

const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);

createApp(App).use(pinia).use(vuetify).use(router).mount("#app");
