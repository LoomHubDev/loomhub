<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useLoomStore } from '@/stores/loom'
import { listCheckpoints, listStreams, type Checkpoint } from '@/api/content'

const route = useRoute()
const router = useRouter()
const loomStore = useLoomStore()

const owner = computed(() => route.params.owner as string)
const loomName = computed(() => route.params.loom as string)

const checkpoints = ref<Checkpoint[]>([])
const total = ref(0)
const loading = ref(true)
const error = ref('')
const resolvedStream = ref('')

// Pagination — driven by query params
const perPage = ref(20)
const page = computed(() => {
  const p = parseInt(route.query.page as string)
  return isNaN(p) || p < 1 ? 1 : p
})
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / perPage.value)))

// Filters — driven by query params
const filterType = computed(() => (route.query.type as string) || '')
const filterAuthor = computed(() => (route.query.author as string) || '')
const filterPath = computed(() => (route.query.path as string) || '')

// Local filter inputs (so user can type before applying)
const inputType = ref('')
const inputAuthor = ref('')
const inputPath = ref('')
const showFilters = ref(false)

function syncInputsFromQuery() {
  inputType.value = filterType.value
  inputAuthor.value = filterAuthor.value
  inputPath.value = filterPath.value
  showFilters.value = !!(filterType.value || filterAuthor.value || filterPath.value)
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    if (!resolvedStream.value) {
      const streams = await listStreams(owner.value, loomName.value)
      resolvedStream.value = streams.length > 0 ? streams[0].id : ''
    }

    const res = await listCheckpoints(owner.value, loomName.value, {
      stream: resolvedStream.value,
      type: filterType.value || undefined,
      author: filterAuthor.value || undefined,
      path: filterPath.value || undefined,
      limit: perPage.value,
      offset: (page.value - 1) * perPage.value,
    })
    checkpoints.value = res.items ?? []
    total.value = res.total ?? 0
  } catch {
    error.value = 'Failed to load checkpoints'
  } finally {
    loading.value = false
  }
}

function applyFilters() {
  const q: Record<string, string> = {}
  if (inputType.value) q.type = inputType.value
  if (inputAuthor.value) q.author = inputAuthor.value
  if (inputPath.value) q.path = inputPath.value
  router.replace({ query: q })
}

function clearFilters() {
  inputType.value = ''
  inputAuthor.value = ''
  inputPath.value = ''
  router.replace({ query: {} })
}

function goToPage(p: number) {
  router.replace({ query: { ...route.query, page: p > 1 ? String(p) : undefined } })
}

const activeFilterCount = computed(() =>
  [filterType.value, filterAuthor.value, filterPath.value].filter(Boolean).length
)

// Watch route query changes to reload
watch(() => route.query, () => {
  syncInputsFromQuery()
  load()
}, { deep: true })

onMounted(() => {
  syncInputsFromQuery()
  load()
})

function formatTime(ts: string): string {
  if (!ts) return ''
  const d = new Date(ts)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  if (diff < 60_000) return 'just now'
  if (diff < 3600_000) return `${Math.floor(diff / 60_000)}m ago`
  if (diff < 86400_000) return `${Math.floor(diff / 3600_000)}h ago`
  if (diff < 604800_000) return `${Math.floor(diff / 86400_000)}d ago`
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

const typeColors: Record<string, string> = {
  create: 'text-green-400 bg-green-900/20',
  modify: 'text-yellow-400 bg-yellow-900/20',
  delete: 'text-red-400 bg-red-900/20',
  rename: 'text-blue-400 bg-blue-900/20',
}

// Page numbers to show
const visiblePages = computed(() => {
  const pages: number[] = []
  const tp = totalPages.value
  const cp = page.value
  if (tp <= 7) {
    for (let i = 1; i <= tp; i++) pages.push(i)
  } else {
    pages.push(1)
    if (cp > 3) pages.push(-1) // ellipsis
    for (let i = Math.max(2, cp - 1); i <= Math.min(tp - 1, cp + 1); i++) pages.push(i)
    if (cp < tp - 2) pages.push(-1)
    pages.push(tp)
  }
  return pages
})
</script>

<template>
  <div>
    <!-- Header -->
    <div class="flex items-center justify-between mb-4">
      <div class="flex items-center gap-3">
        <h2 class="text-lg font-semibold">Checkpoint Log</h2>
        <span class="text-xs text-gray-500">{{ total }} total</span>
      </div>
      <button
        @click="showFilters = !showFilters"
        class="flex items-center gap-1.5 px-3 py-1.5 text-xs rounded-md border transition-colors"
        :class="activeFilterCount > 0
          ? 'border-indigo-500 text-indigo-400 bg-indigo-500/10'
          : 'border-gray-700 text-gray-400 hover:text-gray-200 hover:border-gray-600'"
      >
        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
        </svg>
        Filter
        <span v-if="activeFilterCount" class="w-4 h-4 flex items-center justify-center bg-indigo-500 text-white rounded-full text-[10px] font-bold">{{ activeFilterCount }}</span>
      </button>
    </div>

    <!-- Filters -->
    <div v-if="showFilters" class="mb-4 p-4 border border-gray-800 rounded-lg bg-gray-900/50">
      <form @submit.prevent="applyFilters" class="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div>
          <label class="block text-xs text-gray-500 mb-1">Type</label>
          <select
            v-model="inputType"
            class="w-full px-2.5 py-1.5 bg-gray-900 border border-gray-700 rounded-md text-sm text-white focus:outline-none focus:border-indigo-500"
          >
            <option value="">All types</option>
            <option value="create">create</option>
            <option value="modify">modify</option>
            <option value="delete">delete</option>
            <option value="rename">rename</option>
          </select>
        </div>
        <div>
          <label class="block text-xs text-gray-500 mb-1">Author</label>
          <input
            v-model="inputAuthor"
            placeholder="Search author..."
            class="w-full px-2.5 py-1.5 bg-gray-900 border border-gray-700 rounded-md text-sm text-white focus:outline-none focus:border-indigo-500"
          />
        </div>
        <div>
          <label class="block text-xs text-gray-500 mb-1">Path</label>
          <input
            v-model="inputPath"
            placeholder="Search path..."
            class="w-full px-2.5 py-1.5 bg-gray-900 border border-gray-700 rounded-md text-sm text-white focus:outline-none focus:border-indigo-500"
          />
        </div>
      </form>
      <div class="flex gap-2 mt-3">
        <button
          @click="applyFilters"
          class="px-3 py-1.5 bg-indigo-600 hover:bg-indigo-500 text-white rounded-md text-xs font-medium transition-colors"
        >
          Apply
        </button>
        <button
          v-if="activeFilterCount > 0"
          @click="clearFilters"
          class="px-3 py-1.5 text-gray-400 hover:text-white text-xs transition-colors"
        >
          Clear filters
        </button>
      </div>
    </div>

    <!-- Loading / Error / Empty -->
    <div v-if="loading" class="text-gray-500 text-sm">Loading checkpoints...</div>
    <div v-else-if="error" class="text-red-400 text-sm">{{ error }}</div>
    <div v-else-if="checkpoints.length === 0 && activeFilterCount === 0" class="border border-gray-800 rounded-lg p-8 text-center">
      <p class="text-gray-500">No checkpoints yet.</p>
      <p class="text-gray-600 text-xs mt-1">Send content to this loom to create checkpoints.</p>
    </div>
    <div v-else-if="checkpoints.length === 0" class="border border-gray-800 rounded-lg p-8 text-center">
      <p class="text-gray-500">No results match your filters.</p>
      <button @click="clearFilters" class="text-indigo-400 hover:text-indigo-300 text-sm mt-2">Clear filters</button>
    </div>

    <!-- Log entries -->
    <template v-else>
      <div class="border border-gray-800 rounded-lg divide-y divide-gray-800">
        <div
          v-for="cp in checkpoints"
          :key="cp.id"
          class="px-4 py-3 hover:bg-gray-800/50 transition-colors"
        >
          <div class="flex items-start justify-between gap-3">
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <span class="font-mono text-indigo-400 text-xs shrink-0">#{{ cp.seq }}</span>
                <span
                  class="text-xs px-1.5 py-0.5 rounded font-medium shrink-0"
                  :class="typeColors[cp.type] || 'text-gray-400 bg-gray-800'"
                >
                  {{ cp.type }}
                </span>
                <span class="text-sm text-white truncate">{{ cp.path }}</span>
              </div>
              <div class="flex items-center gap-3 mt-1 text-xs text-gray-500">
                <span class="flex items-center gap-1">
                  <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                  {{ cp.author }}
                </span>
                <span v-if="cp.timestamp">{{ formatTime(cp.timestamp) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="flex items-center justify-between mt-4">
        <button
          :disabled="page <= 1"
          class="px-3 py-1.5 rounded-md text-sm border border-gray-700 text-gray-300 disabled:opacity-30 disabled:cursor-not-allowed hover:bg-gray-800 transition-colors"
          @click="goToPage(page - 1)"
        >
          Previous
        </button>

        <div class="flex items-center gap-1">
          <template v-for="(p, i) in visiblePages" :key="i">
            <span v-if="p === -1" class="px-1 text-gray-600">...</span>
            <button
              v-else
              @click="goToPage(p)"
              class="w-8 h-8 flex items-center justify-center rounded-md text-sm transition-colors"
              :class="p === page
                ? 'bg-indigo-600 text-white font-medium'
                : 'text-gray-400 hover:bg-gray-800 hover:text-white'"
            >
              {{ p }}
            </button>
          </template>
        </div>

        <button
          :disabled="page >= totalPages"
          class="px-3 py-1.5 rounded-md text-sm border border-gray-700 text-gray-300 disabled:opacity-30 disabled:cursor-not-allowed hover:bg-gray-800 transition-colors"
          @click="goToPage(page + 1)"
        >
          Next
        </button>
      </div>

      <!-- Per-page selector -->
      <div class="flex items-center justify-end mt-2 gap-2 text-xs text-gray-500">
        <span>Per page:</span>
        <select
          :value="perPage"
          @change="perPage = parseInt(($event.target as HTMLSelectElement).value); goToPage(1)"
          class="bg-gray-900 border border-gray-700 rounded px-1.5 py-0.5 text-gray-300 text-xs focus:outline-none"
        >
          <option :value="10">10</option>
          <option :value="20">20</option>
          <option :value="50">50</option>
          <option :value="100">100</option>
        </select>
      </div>
    </template>
  </div>
</template>
