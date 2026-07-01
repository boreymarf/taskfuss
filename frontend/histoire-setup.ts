import './src/styles/tailwind.css'

import type { App } from 'vue'
import { VueQueryPlugin } from '@tanstack/vue-query'
import { createPinia } from 'pinia'
import { client } from './src/api/generated/client.gen';

export const setupVue3 = ({ app }: { app: App }) => {
  app.use(VueQueryPlugin)
  app.use(createPinia())


  client.setConfig({
    baseUrl: 'http://localhost:5000',
  });
}
