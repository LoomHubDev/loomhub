<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useLoomStore } from '@/stores/loom'
import { listEntities, listStreams, type Entity } from '@/api/content'

const route = useRoute()
const loomStore = useLoomStore()

const owner = computed(() => route.params.owner as string)
const loomName = computed(() => route.params.loom as string)
const currentPath = computed(() => {
  const p = route.params.treePath
  if (!p) return ''
  return Array.isArray(p) ? p.join('/') : p
})

const resolvedStream = ref('')
const allEntities = ref<Entity[]>([])
const loading = ref(true)
const error = ref('')

interface TreeEntry {
  name: string
  type: 'folder' | 'file'
  path: string
  entity?: Entity
}

const entries = computed<TreeEntry[]>(() => {
  const prefix = currentPath.value ? currentPath.value + '/' : ''
  const folders = new Map<string, boolean>()
  const files: TreeEntry[] = []

  for (const entity of allEntities.value) {
    if (!entity.path.startsWith(prefix) && prefix !== '') continue
    if (entity.path === prefix) continue

    const rest = entity.path.slice(prefix.length)
    const slashIdx = rest.indexOf('/')

    if (slashIdx !== -1) {
      // This entity is in a subdirectory
      const folderName = rest.slice(0, slashIdx)
      if (!folders.has(folderName)) {
        folders.set(folderName, true)
      }
    } else {
      // This entity is directly in the current directory
      files.push({
        name: rest,
        type: 'file',
        path: entity.path,
        entity,
      })
    }
  }

  const result: TreeEntry[] = []
  // Folders first, sorted
  for (const name of [...folders.keys()].sort()) {
    result.push({
      name,
      type: 'folder',
      path: prefix + name,
    })
  }
  // Then files, sorted
  files.sort((a, b) => a.name.localeCompare(b.name))
  result.push(...files)
  return result
})

const breadcrumbs = computed(() => {
  if (!currentPath.value) return []
  const parts = currentPath.value.split('/')
  return parts.map((part, i) => ({
    name: part,
    path: parts.slice(0, i + 1).join('/'),
  }))
})

async function load() {
  loading.value = true
  error.value = ''
  try {
    let streamFilter = ''
    const streams = await listStreams(owner.value, loomName.value)
    if (streams.length > 0) {
      streamFilter = streams[0].id
      resolvedStream.value = streams[0].name
    }
    allEntities.value = await listEntities(owner.value, loomName.value, streamFilter)
  } catch {
    error.value = 'Failed to load entities'
  } finally {
    loading.value = false
  }
}

onMounted(load)

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

function formatTime(ts: string): string {
  if (!ts) return ''
  const d = new Date(ts)
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function folderItemCount(folderPath: string): number {
  const prefix = folderPath + '/'
  return allEntities.value.filter(e => e.path.startsWith(prefix)).length
}
</script>

<template>
  <div>
    <!-- Header -->
    <div class="flex items-center justify-between mb-4">
      <div class="flex items-center gap-2">
        <h2 class="text-lg font-semibold">Entities</h2>
        <span class="text-xs px-2 py-0.5 bg-gray-800 rounded-full text-gray-400 font-mono">{{ resolvedStream || 'all' }}</span>
      </div>
      <span class="text-xs text-gray-500">{{ allEntities.length }} entit{{ allEntities.length !== 1 ? 'ies' : 'y' }}</span>
    </div>

    <!-- Breadcrumb -->
    <div v-if="currentPath" class="flex items-center gap-1 mb-4 text-sm">
      <RouterLink
        :to="`/${owner}/${loomName}/tree`"
        class="text-indigo-400 hover:text-indigo-300"
      >
        {{ loomName }}
      </RouterLink>
      <template v-for="crumb in breadcrumbs" :key="crumb.path">
        <span class="text-gray-600">/</span>
        <RouterLink
          v-if="crumb.path !== currentPath"
          :to="`/${owner}/${loomName}/tree/${crumb.path}`"
          class="text-indigo-400 hover:text-indigo-300"
        >
          {{ crumb.name }}
        </RouterLink>
        <span v-else class="text-white font-medium">{{ crumb.name }}</span>
      </template>
    </div>

    <div v-if="loading" class="text-gray-500 text-sm">Loading entities...</div>
    <div v-else-if="error" class="text-red-400 text-sm">{{ error }}</div>
    <div v-else-if="allEntities.length === 0" class="border border-gray-800 rounded-lg p-8 text-center">
      <p class="text-gray-500">No entities in this stream.</p>
      <p class="text-gray-600 text-xs mt-1">Send content to this loom to see entities here.</p>
    </div>
    <div v-else-if="entries.length === 0" class="border border-gray-800 rounded-lg p-8 text-center">
      <p class="text-gray-500">Empty folder.</p>
    </div>
    <div v-else class="border border-gray-800 rounded-lg divide-y divide-gray-800">
      <!-- Go up -->
      <RouterLink
        v-if="currentPath"
        :to="breadcrumbs.length > 1
          ? `/${owner}/${loomName}/tree/${breadcrumbs[breadcrumbs.length - 2].path}`
          : `/${owner}/${loomName}/tree`"
        class="flex items-center gap-3 px-4 py-2.5 hover:bg-gray-800/50 transition-colors text-gray-400"
      >
        <svg class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M11 17l-5-5m0 0l5-5m-5 5h12" />
        </svg>
        <span class="text-sm">..</span>
      </RouterLink>

      <!-- Folders -->
      <RouterLink
        v-for="entry in entries.filter(e => e.type === 'folder')"
        :key="entry.path"
        :to="`/${owner}/${loomName}/tree/${entry.path}`"
        class="flex items-center justify-between px-4 py-2.5 hover:bg-gray-800/50 transition-colors"
      >
        <div class="flex items-center gap-3 min-w-0">
          <svg class="w-4 h-4 text-indigo-400 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
          </svg>
          <span class="text-sm font-mono">{{ entry.name }}</span>
        </div>
        <span class="text-xs text-gray-600">{{ folderItemCount(entry.path) }} items</span>
      </RouterLink>

      <!-- Files -->
      <RouterLink
        v-for="entry in entries.filter(e => e.type === 'file')"
        :key="entry.path"
        :to="`/${owner}/${loomName}/blob/${entry.path}?stream=${resolvedStream}`"
        class="flex items-center justify-between px-4 py-2.5 hover:bg-gray-800/50 transition-colors"
      >
        <div class="flex items-center gap-3 min-w-0">
          <svg class="w-4 h-4 text-gray-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <span class="text-sm font-mono truncate">{{ entry.name }}</span>
        </div>
        <div class="flex items-center gap-4 text-xs text-gray-500 shrink-0 ml-4">
          <span v-if="entry.entity">{{ formatSize(entry.entity.size) }}</span>
          <span v-if="entry.entity?.mod_time">{{ formatTime(entry.entity.mod_time) }}</span>
        </div>
      </RouterLink>
    </div>
  </div>
</template>
