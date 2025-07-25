import axios from 'axios'
import { defineStore } from 'pinia'
import type { ApiGenericError, LoginRequest, LoginResponse, RegisterRequest, RegisterResponse, User, ValidationError } from '../types/api'
import type { LoginForm, RegisterForm } from '../types/forms'
import type { ApiResponse } from '../types/api/generic'
import { StoreGenericError, RegisterError } from '../types/stores'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    auth_token: localStorage.getItem('auth_token') || null as string | null
  }),

  actions: {
    async register(formData: RegisterForm): Promise<void> {
      const requestData: RegisterRequest = {
        username: formData.name,
        email: formData.email,
        password: formData.password
      }

      try {
        const response = await axios.post<RegisterResponse>("/api/register", requestData);
        const responseData = response.data;

        this.user = responseData.user
        this.auth_token = responseData.auth_token
        localStorage.setItem('auth_token', responseData.auth_token)
      } catch (error) {
        if (axios.isAxiosError(error)) {

          // If bad request, expect validation error
          // FIXME: DTO has changed
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
        const response = await axios.post<ApiResponse<LoginResponse>>("/api/login", requestData);
        const responseData = response.data.data;

        this.user = responseData.user
        this.auth_token = responseData.auth_token
        console.log(responseData)
        console.log(responseData.auth_token)
        localStorage.setItem('auth_token', responseData.auth_token)

      } catch (error) {

        // FIXME: DTO has changed
        if (axios.isAxiosError(error)) {
          const serverError = error.response?.data as ApiGenericError

          throw new StoreGenericError(
            serverError.code,
            serverError.message
          )

        }
      }
    },

    async logout(): Promise<void> {
      this.user = null
      this.auth_token = null
      localStorage.removeItem('auth_token')
    }
  },

  getters: {
    isAuthenticated(): boolean {
      return !!this.auth_token
    }
  }
})
