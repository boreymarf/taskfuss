import { Configuration } from "../generated/configuration";
import { AuthenticationApi, ProfileApi, ServiceApi, TasksApi } from '../generated/api'

// Создаём конфигурацию
const config = new Configuration({
  // basePath: import.meta.env.VITE_API_BASE_URL || 'http://localhost:4000',
  basePath: 'http://localhost:4000/api',
  // Дополнительные настройки:
  // middleware: [], // можно добавить middleware
  // apiKey: 'your-api-key' // если требуется авторизация
})

// Создаём экземпляр API
export const api = {
  auth: new AuthenticationApi(config),
  service: new ServiceApi(config),
  profile: new ProfileApi(config),
  tasks: new TasksApi(config),
}
