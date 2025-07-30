import { Configuration } from "../generated/configuration";
import { ServiceApi } from '../generated/api'

// Создаём конфигурацию
const config = new Configuration({
  basePath: import.meta.env.VITE_API_BASE_URL || 'http://localhost:4000',
  // Дополнительные настройки:
  // middleware: [], // можно добавить middleware
  // apiKey: 'your-api-key' // если требуется авторизация
})

// Создаём экземпляр API
export const api = {
  service: new ServiceApi(config)
}
