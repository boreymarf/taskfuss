import axios from 'axios'
import { defineStore } from 'pinia'
import type { DtoRegisterRequest, DtoUser, DtoRegisterResponse, ApiError, DtoLoginRequest, DtoLoginResponse, DtoProfileResponse } from '../api/generated'
import type { LoginForm, RegisterForm } from '../types/forms'
import { api } from '../api/client/http'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as DtoUser | null,
    auth_token: localStorage.getItem('auth_token') || null as string | null,
    isLoading: false,
    error: null as null | Error
  }),

  actions: {
    async register(formData: RegisterForm): Promise<void> {
      const requestData: DtoRegisterRequest = {
        username: formData.name,
        email: formData.email,
        password: formData.password
      }

      try {
        const response = await api.auth.authRegisterPost({
          registerRequest: requestData
        });

        const dtoRegisterResponse = response.data as DtoRegisterResponse;

        if (!dtoRegisterResponse) {
          throw new Error('Invalid response data: missing response payload');
        }

        if (!dtoRegisterResponse.user || !dtoRegisterResponse.auth_token) {
          throw new Error('Invalid response data: missing user or auth_token');
        }

        this.user = dtoRegisterResponse.user
        this.auth_token = dtoRegisterResponse.auth_token
        localStorage.setItem('auth_token', dtoRegisterResponse.auth_token)

      } catch (error) {

        if (!axios.isAxiosError(error)) {
          throw error;
        }

        if (!error.response?.data) {
          throw new Error('Unknown error');
        }

        const { status, data } = error.response;
        const apiError = data as ApiError;

        if (!apiError.code) {
          throw new Error(apiError.message || 'Unknown error');
        }

        switch (`${status}_${apiError.code}`) {
          case "400_BAD_REQUEST":
            // TODO: Handle later
            console.log("BAD_REQUEST")
            break;
          case "400_DUPLICATE_USER":
            // TODO: Handle later
            console.log("DUBLICATE_USER")
            break;
          default:
            throw new Error(apiError.message || 'Request failed');
        }
      }
    },

    async login(formData: LoginForm): Promise<void> {

      const requestData: DtoLoginRequest = {
        email: formData.email,
        password: formData.password
      }

      try {
        const response = await api.auth.authLoginPost({
          loginRequest: requestData
        });

        const dtoLoginResponse = response.data as DtoLoginResponse;

        if (!dtoLoginResponse) {
          throw new Error('Invalid response data: missing response payload');
        }

        if (!dtoLoginResponse.user || !dtoLoginResponse.auth_token) {
          throw new Error('Invalid response data: missing user or auth_token');
        }

        this.user = dtoLoginResponse.user
        this.auth_token = dtoLoginResponse.auth_token
        localStorage.setItem('auth_token', dtoLoginResponse.auth_token)

      } catch (error) {

        if (!axios.isAxiosError(error)) {
          throw error;
        }

        if (!error.response?.data) {
          throw new Error('Unknown error');
        }

        const { status, data } = error.response;
        const apiError = data as ApiError;

        if (!apiError.code) {
          throw new Error(apiError.message || 'Unknown error');
        }

        switch (`${status}_${apiError.code}`) {
          case "400_BAD_REQUEST":
            // TODO: Handle later
            console.log("BAD_REQUEST")
            break;
          default:
            throw new Error(apiError.message || 'Request failed');
        }
      }
    },

    async logout(): Promise<void> {
      this.user = null
      this.auth_token = null
      localStorage.removeItem('auth_token')
    },

    async fetchUser(): Promise<void> {
      this.isLoading = true

      try {

        if (!this.auth_token) {
          this.error = new Error("No auth token found!")
          return
        }

        const response = await api.profile.profileGet({ authorization: `Bearer ${this.auth_token}` })
        console.log(response)

        const dtoProfileResponse = response.data as DtoProfileResponse;

        if (!dtoProfileResponse) {
          throw new Error('Invalid response data: missing response payload');
        }

        if (!dtoProfileResponse.user) {
          throw new Error('Invalid response data: missing user');
        }

        this.user = dtoProfileResponse.user

      } catch (error) {

        if (!axios.isAxiosError(error)) {
          throw error;
        }

        if (!error.response?.data) {
          throw new Error('Unknown error');
        }

        const { status, data } = error.response;
        const apiError = data as ApiError;

        if (!apiError.code) {
          throw new Error(apiError.message || 'Unknown error');
        }

        switch (`${status}_${apiError.code}`) {
          case "400_BAD_REQUEST":
            // TODO: Handle later
            console.log("BAD_REQUEST")
            break;
          default:
            throw new Error(apiError.message || 'Request failed');
        }
      }
    }
  },

  getters: {
    isAuthenticated(): boolean {
      return !!this.auth_token
    }
  }
})
