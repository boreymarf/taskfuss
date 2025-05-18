import axios from 'axios'
import { defineStore } from 'pinia'
import type { LoginRequest, LoginResponse, RegisterRequest, RegisterResponse, User, ValidationError } from '../types/api'
import type { LoginForm, RegisterForm } from '../types/forms'
import { RegisterError } from '../types/stores'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    authToken: localStorage.getItem('token') || null as string | null
  }),

  actions: {
    async register(formData: RegisterForm): Promise<void> {
      const requestData: RegisterRequest = {
        username: formData.name,
        email: formData.email,
        password: formData.password
      }

      try {
        const responseData: RegisterResponse = await axios.post("/api/register", requestData)

        this.user = responseData.user
        this.authToken = responseData.authToken
        localStorage.setItem('token', responseData.authToken)
      } catch (error) {
        if (axios.isAxiosError(error)) {

          // If bad request, expect validation error
          if (error.response?.status === 400) {
            const serverError = error.response.data as ValidationError

            throw new RegisterError(
              serverError.message,
              serverError.code,
              serverError.details
            )
          }

          // Обработка других HTTP ошибок
          throw new Error(error.response?.data?.message || 'Request failed')
        }
      }
    },

    async login(formData: LoginForm): Promise<void> {

      const requestData: LoginRequest = {
        email: formData.email,
        password: formData.password
      }

      try {
        const responseData: LoginResponse = await axios.post("/api/login", requestData)

      } catch (error) {

      }
    },

    async logout(): Promise<void> {
      this.user = null
      this.authToken = null
      localStorage.removeItem('token')
    }
  },

  getters: {
    isAuthenticated(): boolean {
      return !!this.authToken
    }
  }
})
