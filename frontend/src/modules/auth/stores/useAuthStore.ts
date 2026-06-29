import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(null)
  const tokenType = ref<string | null>(null)

  function initialize() {
    const savedToken = localStorage.getItem('access_token')
    const savedTokenType = localStorage.getItem('token_type')

    if (savedToken && savedTokenType) {
      token.value = savedToken
      tokenType.value = savedTokenType
    }
  }

  function setAuth(accessToken: string, type: string = 'bearer') {
    token.value = accessToken
    tokenType.value = type

    localStorage.setItem('access_token', accessToken)
    localStorage.setItem('token_type', type)
  }

  function logout() {
    token.value = null
    tokenType.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('token_type')
  }

  const isAuthenticated = computed(() => !!token.value)

  const authHeader = computed(() => {
    if (!token.value || !tokenType.value) return null
    return `${tokenType.value} ${token.value}`
  })

  // Get token from storage
  initialize()

  return {
    token: computed(() => token.value), // readonly
    tokenType: computed(() => tokenType.value), // readonly
    isAuthenticated,
    authHeader,
    setAuth,
    logout,
  }
})
