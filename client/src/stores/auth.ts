import axios from 'axios'
import { defineStore } from 'pinia'

interface User {
  id: number
  username: string
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    token: localStorage.getItem('token') || null as string | null
  }),

  actions: {
    async login(username: string, password: string): Promise<void> { // Добавили типы параметров
      return new Promise<void>(async (resolve) => { // Указали тип Promise<void>

        try {

          const response = await axios.post('/api/login/', {
            name: username,
            password: password
          })

        } catch (error) {

        }

        resolve()
      })
    },

    async register(): Promise<void> {

    },

    async logout(): Promise<void> { // Добавили тип возвращаемого значения
      this.user = null
      this.token = null
      localStorage.removeItem('token')
    }
  },

  getters: {
    isAuthenticated(): boolean { // Добавили тип возвращаемого значения
      return !!this.token
    }
  }
})
