import { api } from './client'

export function createLoom(data: { name: string; description?: string; visibility?: string; owner?: string }) {
  return api('/looms', { method: 'POST', body: data })
}

export function getLoom(owner: string, loom: string) {
  return api(`/looms/${owner}/${loom}`)
}

export function updateLoom(owner: string, loom: string, data: { description?: string; visibility?: string }) {
  return api(`/looms/${owner}/${loom}`, { method: 'PATCH', body: data })
}

export function deleteLoom(owner: string, loom: string) {
  return api(`/looms/${owner}/${loom}`, { method: 'DELETE' })
}

export function getCurrentUserLooms() {
  return api('/user/looms')
}
