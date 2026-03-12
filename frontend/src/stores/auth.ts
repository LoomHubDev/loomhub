import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as usersApi from '@/api/users'
import { router } from '@/router'

interface User {
  id: string
  username: string
  email: string
  display_name: string
  bio: string
  is_admin: boolean
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('token'))

  const isAuthenticated = computed(() => !!user.value)
  const isAdmin = computed(() => user.value?.is_admin ?? false)

  async function login(login: string, password: string) {
    const res = await usersApi.login({ login, password })
    token.value = res.token
    localStorage.setItem('token', res.token)
    user.value = res
    router.push('/')
  }

  async function register(username: string, email: string, password: string) {
    const res = await usersApi.register({ username, email, password })
    token.value = res.token
    localStorage.setItem('token', res.token)
    user.value = res
    router.push('/')
  }

  async function fetchUser() {
    try {
      user.value = await usersApi.getCurrentUser()
    } catch {
      logout()
    }
  }

  function logout() {
    user.value = null
    token.value = null
    localStorage.removeItem('token')
    router.push('/login')
  }

  return { user, token, isAuthenticated, isAdmin, login, register, fetchUser, logout }
})
