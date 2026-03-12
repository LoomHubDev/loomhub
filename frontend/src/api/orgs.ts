import { api } from './client'

export interface Organization {
  id: string
  name: string
  display_name: string
  description: string
  avatar_url: string
  created_at: string
  updated_at: string
}

export interface OrgMember {
  org_id: string
  user_id: string
  role: 'owner' | 'admin' | 'member'
  created_at: string
  username: string
}

export function createOrg(data: { name: string; display_name?: string; description?: string }) {
  return api<Organization>('/orgs', { method: 'POST', body: data })
}

export function getOrg(name: string) {
  return api<Organization>(`/orgs/${name}`)
}

export function updateOrg(name: string, data: { display_name?: string; description?: string; avatar_url?: string }) {
  return api<Organization>(`/orgs/${name}`, { method: 'PATCH', body: data })
}

export function listOrgMembers(name: string) {
  return api<OrgMember[]>(`/orgs/${name}/members`)
}

export function addOrgMember(orgName: string, username: string, role: string) {
  return api(`/orgs/${orgName}/members`, { method: 'POST', body: { username, role } })
}

export function removeOrgMember(orgName: string, username: string) {
  return api(`/orgs/${orgName}/members/${username}`, { method: 'DELETE' })
}

export function getUserOrgs() {
  return api<Organization[]>('/user/orgs')
}

export function getOrgLooms(name: string) {
  return api<{ items: any[]; total: number }>(`/orgs/${name}/looms`)
}
