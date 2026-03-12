import { api } from './client'

export function register(data: { username: string; email: string; password: string; display_name?: string }) {
  return api('/users', { method: 'POST', body: data })
}

export function login(data: { login: string; password: string }) {
  return api('/auth/login', { method: 'POST', body: data })
}

export function getCurrentUser() {
  return api('/user')
}

export function updateCurrentUser(data: { display_name?: string; bio?: string }) {
  return api('/user', { method: 'PATCH', body: data })
}

export function getUser(username: string) {
  return api(`/users/${username}`)
}

export function getUserLooms(username: string) {
  return api(`/users/${username}/looms`)
}
