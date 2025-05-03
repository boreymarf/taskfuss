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
    login(username: string, password: string): Promise<void> { // Добавили типы параметров
      return new Promise<void>((resolve) => { // Указали тип Promise<void>
        setTimeout(() => {
          // Добавим минимальную проверку пароля для примера
          if (password.length < 1) {
            throw new Error('Password cannot be empty')
          }

          this.user = { id: 1, username }
          this.token = 'fake-jwt-token'
          localStorage.setItem('token', this.token)
          resolve()
        }, 500)
      })
    },

    logout(): void { // Добавили тип возвращаемого значения
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
