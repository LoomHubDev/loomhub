import { api } from './client'

export interface AccessToken {
  id: string
  name: string
  token?: string
  scopes: string
  expires_at: string
  last_used_at: string
  created_at: string
}

export function listTokens(): Promise<AccessToken[]> {
  return api('/user/tokens')
}

export function createToken(data: { name: string; scopes?: string }): Promise<AccessToken & { token: string }> {
  return api('/user/tokens', { method: 'POST', body: data })
}

export function deleteToken(id: string): Promise<void> {
  return api(`/user/tokens/${id}`, { method: 'DELETE' })
}
