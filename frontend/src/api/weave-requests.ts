import { api } from './client'

export interface WeaveRequest {
  id: string
  loom_id: string
  number: number
  author_id: string
  title: string
  description: string
  source_stream: string
  target_stream: string
  status: 'open' | 'woven' | 'closed'
  woven_by: string | null
  woven_at: string | null
  closed_at: string | null
  created_at: string
  updated_at: string
  author: string
}

export interface Comment {
  id: string
  wr_id: string
  author_id: string
  body: string
  space_id: string | null
  entity_path: string | null
  line_number: number | null
  is_system: boolean
  event_type: string | null
  created_at: string
  updated_at: string
  author: string
}

export interface Review {
  id: string
  wr_id: string
  reviewer_id: string
  status: 'approved' | 'changes_requested' | 'commented'
  body: string
  created_at: string
  reviewer: string
}

export function listWeaveRequests(owner: string, loom: string, status?: string) {
  const params = status ? `?status=${status}` : ''
  return api<{ items: WeaveRequest[]; total: number }>(`/looms/${owner}/${loom}/weave-requests${params}`)
}

export function getWeaveRequest(owner: string, loom: string, number: number) {
  return api<WeaveRequest>(`/looms/${owner}/${loom}/weave-requests/${number}`)
}

export function createWeaveRequest(owner: string, loom: string, data: {
  title: string
  description?: string
  source_stream: string
  target_stream: string
}) {
  return api<WeaveRequest>(`/looms/${owner}/${loom}/weave-requests`, { method: 'POST', body: data })
}

export function updateWeaveRequest(owner: string, loom: string, number: number, data: {
  title?: string
  description?: string
}) {
  return api<WeaveRequest>(`/looms/${owner}/${loom}/weave-requests/${number}`, { method: 'PATCH', body: data })
}

export function weaveWR(owner: string, loom: string, number: number) {
  return api<WeaveRequest>(`/looms/${owner}/${loom}/weave-requests/${number}/weave`, { method: 'POST' })
}

export function closeWR(owner: string, loom: string, number: number) {
  return api<WeaveRequest>(`/looms/${owner}/${loom}/weave-requests/${number}/close`, { method: 'POST' })
}

export function listComments(owner: string, loom: string, number: number) {
  return api<Comment[]>(`/looms/${owner}/${loom}/weave-requests/${number}/comments`)
}

export function createComment(owner: string, loom: string, number: number, body: string) {
  return api<Comment>(`/looms/${owner}/${loom}/weave-requests/${number}/comments`, { method: 'POST', body: { body } })
}

export function listReviews(owner: string, loom: string, number: number) {
  return api<Review[]>(`/looms/${owner}/${loom}/weave-requests/${number}/reviews`)
}

export function createReview(owner: string, loom: string, number: number, data: {
  status: 'approved' | 'changes_requested' | 'commented'
  body: string
}) {
  return api<Review>(`/looms/${owner}/${loom}/weave-requests/${number}/reviews`, { method: 'POST', body: data })
}

export function pinLoom(owner: string, loom: string) {
  return api(`/looms/${owner}/${loom}/pin`, { method: 'PUT' })
}

export function unpinLoom(owner: string, loom: string) {
  return api(`/looms/${owner}/${loom}/pin`, { method: 'DELETE' })
}
