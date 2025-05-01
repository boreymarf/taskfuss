import { createApp } from 'vue'
const logger = pino(
  {
    level: process.env.NODE_ENV === 'production' ? 'info' : 'debug',
    formatters: {
      level: (label) => ({ level: label }),
      bindings: () => ({})
    },
    redact: {
      paths: ['req.headers.authorization', 'req.headers.cookie'],
      censor: '**REDACTED**'
    }
  },
  pino.multistream(streams)
)
import './style.css'
import App from './App.vue'
import router from './router'

createApp(App)
  .use(router)
  .mount('#app')
