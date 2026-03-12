import { describe, it, expect, vi, beforeEach } from 'vitest'
import {
  listCollaborators,
  addCollaborator,
  removeCollaborator,
} from '@/api/collaborators'

vi.mock('@/api/client', () => ({
  api: vi.fn(),
}))

import { api } from '@/api/client'

const mockApi = vi.mocked(api)

beforeEach(() => {
  mockApi.mockReset()
})

describe('listCollaborators', () => {
  it('calls the correct URL', async () => {
    mockApi.mockResolvedValue([])
    await listCollaborators('alice', 'my-loom')
    expect(mockApi).toHaveBeenCalledWith('/looms/alice/my-loom/collaborators')
  })

  it('returns collaborators from the API', async () => {
    const collabs = [
      {
        user_id: 'u1',
        username: 'bob',
        display_name: 'Bob',
        avatar_url: '',
        permission: 'write' as const,
        created_at: '2025-01-01',
      },
    ]
    mockApi.mockResolvedValue(collabs)
    const result = await listCollaborators('alice', 'my-loom')
    expect(result).toEqual(collabs)
    expect(result).toHaveLength(1)
    expect(result[0].username).toBe('bob')
  })
})

describe('addCollaborator', () => {
  it('sends a POST with username and permission', async () => {
    mockApi.mockResolvedValue(undefined)
    await addCollaborator('alice', 'my-loom', 'bob', 'write')
    expect(mockApi).toHaveBeenCalledWith('/looms/alice/my-loom/collaborators', {
      method: 'POST',
      body: { username: 'bob', permission: 'write' },
    })
  })

  it('sends correct data for admin permission', async () => {
    mockApi.mockResolvedValue(undefined)
    await addCollaborator('alice', 'my-loom', 'charlie', 'admin')
    expect(mockApi).toHaveBeenCalledWith('/looms/alice/my-loom/collaborators', {
      method: 'POST',
      body: { username: 'charlie', permission: 'admin' },
    })
  })
})

describe('removeCollaborator', () => {
  it('sends a DELETE to the correct URL', async () => {
    mockApi.mockResolvedValue(undefined)
    await removeCollaborator('alice', 'my-loom', 'bob')
    expect(mockApi).toHaveBeenCalledWith('/looms/alice/my-loom/collaborators/bob', {
      method: 'DELETE',
    })
  })

  it('uses the username in the URL path', async () => {
    mockApi.mockResolvedValue(undefined)
    await removeCollaborator('org', 'project', 'user-to-remove')
    expect(mockApi).toHaveBeenCalledWith('/looms/org/project/collaborators/user-to-remove', {
      method: 'DELETE',
    })
  })
})
