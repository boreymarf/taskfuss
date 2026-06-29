import { createApp } from 'vue'
import { createPinia } from 'pinia'
import '@/styles/tailwind.css'
import { VueQueryPlugin } from '@tanstack/vue-query'
import { client } from '@/api/generated/client.gen';

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(VueQueryPlugin)
app.use(createPinia())
app.use(router)
app.mount('#app')

client.setConfig({
  baseUrl: 'http://localhost:5000',
});
