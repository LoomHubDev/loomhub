import { describe, it, expect, vi, beforeEach } from 'vitest'
import {
  listStreams,
  listEntities,
  listCheckpoints,
  getEntityContent,
  compareStreams,
} from '@/api/content'

vi.mock('@/api/client', () => ({
  api: vi.fn(),
}))

import { api } from '@/api/client'

const mockApi = vi.mocked(api)

beforeEach(() => {
  mockApi.mockReset()
})

describe('listStreams', () => {
  it('calls the correct URL', async () => {
    mockApi.mockResolvedValue([])
    await listStreams('alice', 'my-loom')
    expect(mockApi).toHaveBeenCalledWith('/looms/alice/my-loom/streams')
  })

  it('returns the streams from the API', async () => {
    const streams = [{ id: 's1', name: 'main', head_seq: 5, status: 'active' }]
    mockApi.mockResolvedValue(streams)
    const result = await listStreams('alice', 'my-loom')
    expect(result).toEqual(streams)
  })
})

describe('listEntities', () => {
  it('calls without stream param when not provided', async () => {
    mockApi.mockResolvedValue([])
    await listEntities('bob', 'repo')
    expect(mockApi).toHaveBeenCalledWith('/looms/bob/repo/entities')
  })

  it('appends stream query param when provided', async () => {
    mockApi.mockResolvedValue([])
    await listEntities('bob', 'repo', 'feature-branch')
    expect(mockApi).toHaveBeenCalledWith('/looms/bob/repo/entities?stream=feature-branch')
  })

  it('returns the entities from the API', async () => {
    const entities = [
      { id: 'e1', space_id: 's1', path: '/file.txt', object_ref: 'abc', size: 100, mod_time: '2025-01-01' },
    ]
    mockApi.mockResolvedValue(entities)
    const result = await listEntities('bob', 'repo', 'main')
    expect(result).toEqual(entities)
  })
})

describe('listCheckpoints', () => {
  it('calls without query string when no opts provided', async () => {
    mockApi.mockResolvedValue({ items: [], total: 0 })
    await listCheckpoints('alice', 'loom')
    expect(mockApi).toHaveBeenCalledWith('/looms/alice/loom/checkpoints')
  })

  it('calls without query string when opts is empty object', async () => {
    mockApi.mockResolvedValue({ items: [], total: 0 })
    await listCheckpoints('alice', 'loom', {})
    expect(mockApi).toHaveBeenCalledWith('/looms/alice/loom/checkpoints')
  })

  it('includes stream filter', async () => {
    mockApi.mockResolvedValue({ items: [], total: 0 })
    await listCheckpoints('alice', 'loom', { stream: 'main' })
    expect(mockApi).toHaveBeenCalledWith('/looms/alice/loom/checkpoints?stream=main')
  })

  it('includes type filter', async () => {
    mockApi.mockResolvedValue({ items: [], total: 0 })
    await listCheckpoints('alice', 'loom', { type: 'create' })
    expect(mockApi).toHaveBeenCalledWith('/looms/alice/loom/checkpoints?type=create')
  })

  it('includes author filter', async () => {
    mockApi.mockResolvedValue({ items: [], total: 0 })
    await listCheckpoints('alice', 'loom', { author: 'bob' })
    expect(mockApi).toHaveBeenCalledWith('/looms/alice/loom/checkpoints?author=bob')
  })

  it('includes path filter', async () => {
    mockApi.mockResolvedValue({ items: [], total: 0 })
    await listCheckpoints('alice', 'loom', { path: 'src/main.ts' })
    expect(mockApi).toHaveBeenCalledWith('/looms/alice/loom/checkpoints?path=src%2Fmain.ts')
  })

  it('includes limit and offset', async () => {
    mockApi.mockResolvedValue({ items: [], total: 0 })
    await listCheckpoints('alice', 'loom', { limit: 20, offset: 40 })
    expect(mockApi).toHaveBeenCalledWith('/looms/alice/loom/checkpoints?limit=20&offset=40')
  })

  it('includes all filters together', async () => {
    mockApi.mockResolvedValue({ items: [], total: 0 })
    await listCheckpoints('alice', 'loom', {
      stream: 'main',
      type: 'modify',
      author: 'bob',
      path: 'README.md',
      limit: 10,
      offset: 0,
    })
    const calledUrl = mockApi.mock.calls[0][0] as string
    expect(calledUrl).toContain('stream=main')
    expect(calledUrl).toContain('type=modify')
    expect(calledUrl).toContain('author=bob')
    expect(calledUrl).toContain('path=README.md')
    expect(calledUrl).toContain('limit=10')
    // offset=0 is falsy so it won't be included (URLSearchParams sets '0')
    // Actually, 0 is falsy in JS: if (opts?.offset) is false for 0
    // So offset=0 should NOT appear
  })

  it('skips offset when it is 0 (falsy)', async () => {
    mockApi.mockResolvedValue({ items: [], total: 0 })
    await listCheckpoints('alice', 'loom', { offset: 0 })
    expect(mockApi).toHaveBeenCalledWith('/looms/alice/loom/checkpoints')
  })

  it('returns items and total from API', async () => {
    const payload = {
      items: [{ id: 'c1', seq: 1, stream_id: 's1', space_id: 'sp1', entity_id: 'e1', type: 'create', path: '/a.txt', object_ref: 'ref1', parent_seq: 0, author: 'alice', timestamp: '2025-01-01', meta: '' }],
      total: 42,
    }
    mockApi.mockResolvedValue(payload)
    const result = await listCheckpoints('alice', 'loom', { limit: 10 })
    expect(result.items).toHaveLength(1)
    expect(result.total).toBe(42)
  })
})

describe('getEntityContent', () => {
  it('calls the correct URL with responseType text', async () => {
    mockApi.mockResolvedValue('file content')
    await getEntityContent('alice', 'loom', 'abc123')
    expect(mockApi).toHaveBeenCalledWith('/looms/alice/loom/objects/abc123', { responseType: 'text' })
  })

  it('returns the content string', async () => {
    mockApi.mockResolvedValue('hello world')
    const result = await getEntityContent('alice', 'loom', 'ref1')
    expect(result).toBe('hello world')
  })
})

describe('compareStreams', () => {
  it('calls the correct URL with encoded params', async () => {
    mockApi.mockResolvedValue([])
    await compareStreams('alice', 'loom', 'main', 'feature/branch')
    expect(mockApi).toHaveBeenCalledWith(
      '/looms/alice/loom/compare?source=main&target=feature%2Fbranch'
    )
  })

  it('encodes special characters in source and target', async () => {
    mockApi.mockResolvedValue([])
    await compareStreams('alice', 'loom', 'branch with spaces', 'branch&special=chars')
    expect(mockApi).toHaveBeenCalledWith(
      '/looms/alice/loom/compare?source=branch%20with%20spaces&target=branch%26special%3Dchars'
    )
  })

  it('returns diff entries from the API', async () => {
    const diffs = [
      { path: 'file.txt', status: 'modified' as const, source_ref: 'a', target_ref: 'b' },
    ]
    mockApi.mockResolvedValue(diffs)
    const result = await compareStreams('alice', 'loom', 'main', 'dev')
    expect(result).toEqual(diffs)
  })
})
