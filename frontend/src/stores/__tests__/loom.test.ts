import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useLoomStore } from '@/stores/loom'

vi.mock('@/api/looms', () => ({
  getLoom: vi.fn(),
}))

import * as loomsApi from '@/api/looms'

const mockGetLoom = vi.mocked(loomsApi.getLoom)

beforeEach(() => {
  setActivePinia(createPinia())
  vi.clearAllMocks()
})

describe('useLoomStore', () => {
  const fakeLoom = {
    id: 'loom1',
    owner: 'alice',
    name: 'my-project',
    full_name: 'alice/my-project',
    description: 'A test loom',
    visibility: 'public',
    default_stream: 'develop',
    checkpoint_count: 42,
    stream_count: 3,
    pin_count: 5,
    spin_count: 2,
    created_at: '2025-01-01',
    updated_at: '2025-06-01',
  }

  describe('initial state', () => {
    it('has null loom initially', () => {
      const store = useLoomStore()
      expect(store.loom).toBeNull()
    })

    it('defaults activeStream to main', () => {
      const store = useLoomStore()
      expect(store.activeStream).toBe('main')
    })

    it('is not loading initially', () => {
      const store = useLoomStore()
      expect(store.loading).toBe(false)
    })
  })

  describe('fetchLoom', () => {
    it('loads loom data from the API', async () => {
      mockGetLoom.mockResolvedValue(fakeLoom)

      const store = useLoomStore()
      await store.fetchLoom('alice', 'my-project')

      expect(mockGetLoom).toHaveBeenCalledWith('alice', 'my-project')
      expect(store.loom).toEqual(fakeLoom)
    })

    it('sets activeStream to the loom default_stream', async () => {
      mockGetLoom.mockResolvedValue(fakeLoom)

      const store = useLoomStore()
      await store.fetchLoom('alice', 'my-project')

      expect(store.activeStream).toBe('develop')
    })

    it('keeps empty string when default_stream is empty (not nullish)', async () => {
      mockGetLoom.mockResolvedValue({ ...fakeLoom, default_stream: '' })

      const store = useLoomStore()
      await store.fetchLoom('alice', 'my-project')

      // ?? only triggers for null/undefined, not empty string
      expect(store.activeStream).toBe('')
    })

    it('falls back to main when loom is null (API returns null)', async () => {
      mockGetLoom.mockResolvedValue(null)

      const store = useLoomStore()
      await store.fetchLoom('alice', 'my-project')

      // loom.value?.default_stream is undefined, so ?? 'main' applies
      expect(store.activeStream).toBe('main')
    })

    it('sets loading to true during fetch and false after', async () => {
      let resolveApi: (value: unknown) => void
      mockGetLoom.mockImplementation(() => new Promise((resolve) => { resolveApi = resolve }))

      const store = useLoomStore()
      const promise = store.fetchLoom('alice', 'my-project')

      expect(store.loading).toBe(true)

      resolveApi!(fakeLoom)
      await promise

      expect(store.loading).toBe(false)
    })

    it('sets loading to false even when API throws', async () => {
      mockGetLoom.mockRejectedValue(new Error('Network error'))

      const store = useLoomStore()

      await expect(store.fetchLoom('alice', 'my-project')).rejects.toThrow('Network error')
      expect(store.loading).toBe(false)
    })
  })

  describe('clear', () => {
    it('resets loom to null', async () => {
      mockGetLoom.mockResolvedValue(fakeLoom)

      const store = useLoomStore()
      await store.fetchLoom('alice', 'my-project')
      expect(store.loom).not.toBeNull()

      store.clear()
      expect(store.loom).toBeNull()
    })

    it('resets activeStream to main', async () => {
      mockGetLoom.mockResolvedValue(fakeLoom)

      const store = useLoomStore()
      await store.fetchLoom('alice', 'my-project')
      expect(store.activeStream).toBe('develop')

      store.clear()
      expect(store.activeStream).toBe('main')
    })
  })
})
