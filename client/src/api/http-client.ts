import { Configuration } from "./configuration";
import { ServiceApi } from './api'

// Создаём конфигурацию
const config = new Configuration({
  basePath: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  // Дополнительные настройки:
  // middleware: [], // можно добавить middleware
  // apiKey: 'your-api-key' // если требуется авторизация
})

// Создаём экземпляр API
export const api = {
  service: new ServiceApi(config)
}
