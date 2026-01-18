import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User } from '@/types'

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const token = ref<string>(localStorage.getItem('token') || '')

  function setUser(userData: User) {
    user.value = userData
  }

  function setToken(tokenValue: string) {
    token.value = tokenValue
    localStorage.setItem('token', tokenValue)
  }

  function logout() {
    user.value = null
    token.value = ''
    localStorage.removeItem('token')
  }

  function isLoggedIn(): boolean {
    return !!token.value
  }

  return {
    user,
    token,
    setUser,
    setToken,
    logout,
    isLoggedIn
  }
})
