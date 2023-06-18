import { createApp } from "vue";
import { createPinia } from "pinia";
import vuetify from "./plugins/vuetify";
import router from "./plugins/router";
import App from "./App.vue";

import "./cast";
import "./plugins/buffer";

import("./scss/fontawesome.scss");

createApp(App).use(createPinia()).use(vuetify).use(router).mount("#app");
