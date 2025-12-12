import { createApp } from "vue";
import "./index.css";
import App from "./App.vue";
import router from "./router";
import { createPinia } from "pinia";
import i18n from "./lib/i18n";

export const pinia = createPinia();

createApp(App).use(router).use(pinia).use(i18n).mount("#app");
