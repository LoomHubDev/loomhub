import { api } from './client'

export interface Collaborator {
  user_id: string
  username: string
  display_name: string
  avatar_url: string
  permission: 'read' | 'write' | 'admin'
  created_at: string
}

export function listCollaborators(owner: string, loom: string): Promise<Collaborator[]> {
  return api(`/looms/${owner}/${loom}/collaborators`)
}

export function addCollaborator(owner: string, loom: string, username: string, permission: string): Promise<void> {
  return api(`/looms/${owner}/${loom}/collaborators`, {
    method: 'POST',
    body: { username, permission },
  })
}

export function removeCollaborator(owner: string, loom: string, username: string): Promise<void> {
  return api(`/looms/${owner}/${loom}/collaborators/${username}`, {
    method: 'DELETE',
  })
}
