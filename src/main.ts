import { createApp } from "vue";
import App from "./App.vue";
import i18n from "./i18n.ts";
import { attachConsole } from "@tauri-apps/plugin-log";
import { info } from "@utils/logger";

// Styles
import "@fontsource-variable/inter";
import "@fontsource-variable/geist-mono";
import "./styles/main.scss";
import 'virtual:uno.css'

await attachConsole();

createApp(App).use(i18n).mount("#app");
info("The app is mounted.")