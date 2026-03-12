import { api } from './client'

export interface Stream {
  id: string
  name: string
  head_seq: number
  status: string
}

export interface Entity {
  id: string
  space_id: string
  path: string
  object_ref: string
  size: number
  mod_time: string
}

export interface Checkpoint {
  id: string
  seq: number
  stream_id: string
  space_id: string
  entity_id: string
  type: string
  path: string
  object_ref: string
  parent_seq: number
  author: string
  timestamp: string
  meta: string
}

export function listStreams(owner: string, loom: string): Promise<Stream[]> {
  return api(`/looms/${owner}/${loom}/streams`)
}

export function listEntities(owner: string, loom: string, stream?: string): Promise<Entity[]> {
  const params = stream ? `?stream=${stream}` : ''
  return api(`/looms/${owner}/${loom}/entities${params}`)
}

export interface DiffEntry {
  path: string
  status: 'added' | 'modified' | 'deleted'
  source_ref?: string
  target_ref?: string
  source_content?: string
  target_content?: string
}

export function getEntityContent(owner: string, loom: string, ref: string): Promise<string> {
  return api(`/looms/${owner}/${loom}/objects/${ref}`, { responseType: 'text' })
}

export function compareStreams(owner: string, loom: string, source: string, target: string): Promise<DiffEntry[]> {
  return api(`/looms/${owner}/${loom}/compare?source=${encodeURIComponent(source)}&target=${encodeURIComponent(target)}`)
}

export function listCheckpoints(owner: string, loom: string, opts?: {
  stream?: string
  type?: string
  author?: string
  path?: string
  limit?: number
  offset?: number
}): Promise<{ items: Checkpoint[]; total: number }> {
  const q = new URLSearchParams()
  if (opts?.stream) q.set('stream', opts.stream)
  if (opts?.type) q.set('type', opts.type)
  if (opts?.author) q.set('author', opts.author)
  if (opts?.path) q.set('path', opts.path)
  if (opts?.limit) q.set('limit', String(opts.limit))
  if (opts?.offset) q.set('offset', String(opts.offset))
  const qs = q.toString()
  return api(`/looms/${owner}/${loom}/checkpoints${qs ? '?' + qs : ''}`)
}
