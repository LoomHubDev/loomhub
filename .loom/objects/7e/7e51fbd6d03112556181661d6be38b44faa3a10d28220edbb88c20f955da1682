import { api } from './client'

export interface Activity {
  id: string
  actor_id: string
  loom_id: string | null
  type: string
  payload: string
  created_at: string
  actor: string
  loom_full_name: string | null
}

export function getActivityFeed() {
  return api<Activity[]>('/user/feed')
}

export function getUserActivity(username: string) {
  return api<Activity[]>(`/users/${username}/activity`)
}

export function getGlobalActivity() {
  return api<Activity[]>('/activity')
}
