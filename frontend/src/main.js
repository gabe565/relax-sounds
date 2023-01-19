import { createApp } from "vue";
import vuetify from "./plugins/vuetify";
import store from "./plugins/store/main";
import router from "./plugins/router";
import App from "./App.vue";

import "./cast";
import "./plugins/buffer";

import("./scss/fontawesome.scss");

createApp(App).use(vuetify).use(store).use(router).mount("#app");
