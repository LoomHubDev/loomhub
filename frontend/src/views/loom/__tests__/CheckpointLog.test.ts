import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import CheckpointLog from '@/views/loom/CheckpointLog.vue'

// Mock the content API
vi.mock('@/api/content', () => ({
  listCheckpoints: vi.fn(),
  listStreams: vi.fn(),
}))

import { listCheckpoints, listStreams } from '@/api/content'

const mockListCheckpoints = vi.mocked(listCheckpoints)
const mockListStreams = vi.mocked(listStreams)

function createTestRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      {
        path: '/:owner/:loom/log',
        name: 'loom-log',
        component: CheckpointLog,
      },
    ],
  })
}

function makeCheckpoint(overrides: Record<string, unknown> = {}) {
  return {
    id: 'cp-1',
    seq: 1,
    stream_id: 's1',
    space_id: 'sp1',
    entity_id: 'e1',
    type: 'create',
    path: '/file.txt',
    object_ref: 'ref1',
    parent_seq: 0,
    author: 'alice',
    timestamp: '2025-06-15T12:00:00Z',
    meta: '',
    ...overrides,
  }
}

async function mountComponent(query: Record<string, string> = {}) {
  const router = createTestRouter()
  const pinia = createPinia()
  setActivePinia(pinia)

  // Build path with query
  const queryString = Object.entries(query)
    .map(([k, v]) => `${k}=${encodeURIComponent(v)}`)
    .join('&')
  const path = `/alice/my-loom/log${queryString ? '?' + queryString : ''}`

  router.push(path)
  await router.isReady()

  const wrapper = mount(CheckpointLog, {
    global: {
      plugins: [router, pinia],
    },
  })

  return { wrapper, router }
}

beforeEach(() => {
  vi.clearAllMocks()
  mockListStreams.mockResolvedValue([
    { id: 'stream-1', name: 'main', head_seq: 10, status: 'active' },
  ])
})

describe('CheckpointLog', () => {
  describe('loading state', () => {
    it('shows loading text initially', async () => {
      // Don't resolve the API call yet
      mockListCheckpoints.mockImplementation(() => new Promise(() => {}))

      const { wrapper } = await mountComponent()

      expect(wrapper.text()).toContain('Loading checkpoints...')
    })
  })

  describe('rendering checkpoint entries', () => {
    it('renders checkpoint entries after load', async () => {
      const checkpoints = [
        makeCheckpoint({ id: 'cp-1', seq: 1, path: '/src/main.ts', author: 'alice', type: 'create' }),
        makeCheckpoint({ id: 'cp-2', seq: 2, path: '/README.md', author: 'bob', type: 'modify' }),
      ]
      mockListCheckpoints.mockResolvedValue({ items: checkpoints, total: 2 })

      const { wrapper } = await mountComponent()
      await flushPromises()

      expect(wrapper.text()).toContain('#1')
      expect(wrapper.text()).toContain('#2')
      expect(wrapper.text()).toContain('/src/main.ts')
      expect(wrapper.text()).toContain('/README.md')
      expect(wrapper.text()).toContain('alice')
      expect(wrapper.text()).toContain('bob')
      expect(wrapper.text()).toContain('create')
      expect(wrapper.text()).toContain('modify')
    })

    it('shows total count in header', async () => {
      mockListCheckpoints.mockResolvedValue({ items: [makeCheckpoint()], total: 42 })

      const { wrapper } = await mountComponent()
      await flushPromises()

      expect(wrapper.text()).toContain('42 total')
    })

    it('shows empty state when no checkpoints exist', async () => {
      mockListCheckpoints.mockResolvedValue({ items: [], total: 0 })

      const { wrapper } = await mountComponent()
      await flushPromises()

      expect(wrapper.text()).toContain('No checkpoints yet.')
    })
  })

  describe('filter panel', () => {
    it('toggles filter visibility when button is clicked', async () => {
      mockListCheckpoints.mockResolvedValue({ items: [makeCheckpoint()], total: 1 })

      const { wrapper } = await mountComponent()
      await flushPromises()

      // Filter panel should be hidden by default -- look for the Apply button which is unique to the filter panel
      expect(wrapper.findAll('button').some((b) => b.text() === 'Apply')).toBe(false)

      // Click the Filter button
      const filterBtn = wrapper.findAll('button').find((b) => b.text().includes('Filter'))
      expect(filterBtn).toBeDefined()
      await filterBtn!.trigger('click')

      // Now the filter panel with Apply button and labels should be visible
      expect(wrapper.findAll('button').some((b) => b.text() === 'Apply')).toBe(true)
      expect(wrapper.text()).toContain('Type')
      expect(wrapper.text()).toContain('Author')
      expect(wrapper.text()).toContain('Path')
    })
  })

  describe('pagination', () => {
    it('shows pagination controls when total > perPage', async () => {
      const items = Array.from({ length: 20 }, (_, i) =>
        makeCheckpoint({ id: `cp-${i}`, seq: i + 1 })
      )
      mockListCheckpoints.mockResolvedValue({ items, total: 45 })

      const { wrapper } = await mountComponent()
      await flushPromises()

      expect(wrapper.text()).toContain('Previous')
      expect(wrapper.text()).toContain('Next')
    })

    it('does not show pagination when total fits on one page', async () => {
      mockListCheckpoints.mockResolvedValue({
        items: [makeCheckpoint()],
        total: 1,
      })

      const { wrapper } = await mountComponent()
      await flushPromises()

      expect(wrapper.text()).not.toContain('Previous')
      expect(wrapper.text()).not.toContain('Next')
    })

    it('disables Previous button on first page', async () => {
      const items = Array.from({ length: 20 }, (_, i) =>
        makeCheckpoint({ id: `cp-${i}`, seq: i + 1 })
      )
      mockListCheckpoints.mockResolvedValue({ items, total: 45 })

      const { wrapper } = await mountComponent()
      await flushPromises()

      const prevBtn = wrapper.findAll('button').find((b) => b.text() === 'Previous')
      expect(prevBtn).toBeDefined()
      expect(prevBtn!.attributes('disabled')).toBeDefined()
    })

    it('passes limit and offset to the API', async () => {
      mockListCheckpoints.mockResolvedValue({ items: [], total: 0 })

      await mountComponent()
      await flushPromises()

      expect(mockListCheckpoints).toHaveBeenCalledWith(
        'alice',
        'my-loom',
        expect.objectContaining({
          limit: 20,
          offset: 0,
        })
      )
    })
  })

  describe('error state', () => {
    it('shows error message when API fails', async () => {
      mockListStreams.mockResolvedValue([
        { id: 'stream-1', name: 'main', head_seq: 10, status: 'active' },
      ])
      mockListCheckpoints.mockRejectedValue(new Error('Server error'))

      const { wrapper } = await mountComponent()
      await flushPromises()

      expect(wrapper.text()).toContain('Failed to load checkpoints')
    })
  })
})
