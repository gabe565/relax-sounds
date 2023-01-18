import { createApp } from "vue";
import vuetify from "./plugins/vuetify";
import store from "./plugins/store/main";
import router from "./router";
import App from "./App.vue";

import "./cast";
import "./plugins/buffer";

createApp(App).use(vuetify).use(store).use(router).mount("#app");
