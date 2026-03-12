<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { api } from '@/api/client'
import { timeAgo } from '@/utils/format'

interface Loom {
  id: string
  name: string
  full_name: string
  description: string
  pin_count: number
  checkpoint_count: number
  updated_at: string
}

const looms = ref<Loom[]>([])
const total = ref(0)
const sort = ref('recent')
const search = ref('')
const loading = ref(true)

async function load() {
  loading.value = true
  try {
    if (search.value) {
      const res = await api<Loom[]>(`/search/looms?q=${encodeURIComponent(search.value)}`)
      looms.value = res || []
      total.value = looms.value.length
    } else {
      const res = await api<{ items: Loom[]; total: number }>(`/explore?sort=${sort.value}`)
      looms.value = res.items || []
      total.value = res.total
    }
  } finally {
    loading.value = false
  }
}

onMounted(load)
watch(sort, load)

let searchTimer: ReturnType<typeof setTimeout>
function onSearch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(load, 300)
}
</script>

<template>
  <div class="max-w-5xl mx-auto mt-8 px-4">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-xl font-semibold text-white">Explore</h1>
      <div class="flex gap-2">
        <button
          v-for="s in [{ key: 'recent', label: 'Recent' }, { key: 'pins', label: 'Most pinned' }, { key: 'created', label: 'Newest' }]"
          :key="s.key"
          @click="sort = s.key; search = ''"
          class="px-3 py-1 text-sm rounded-md transition-colors"
          :class="sort === s.key && !search ? 'bg-gray-700 text-white' : 'text-gray-400 hover:text-gray-200'"
        >
          {{ s.label }}
        </button>
      </div>
    </div>

    <div class="mb-6">
      <input
        v-model="search"
        @input="onSearch"
        placeholder="Search looms..."
        class="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded-md text-white text-sm focus:border-indigo-500 focus:outline-none"
      />
    </div>

    <div v-if="loading" class="text-gray-500 text-sm">Loading...</div>
    <div v-else-if="looms.length === 0" class="text-gray-500 text-sm text-center py-8">
      {{ search ? 'No looms match your search.' : 'No public looms yet.' }}
    </div>
    <div v-else>
      <p class="text-xs text-gray-600 mb-3">{{ total }} looms</p>
      <div class="space-y-3">
        <RouterLink
          v-for="l in looms"
          :key="l.id"
          :to="`/${l.full_name}`"
          class="block p-4 border border-gray-800 rounded-lg hover:border-gray-600 transition-colors"
        >
          <div class="flex items-center justify-between">
            <div>
              <h3 class="font-medium text-indigo-400">{{ l.full_name }}</h3>
              <p v-if="l.description" class="text-sm text-gray-500 mt-1">{{ l.description }}</p>
            </div>
            <div class="flex items-center gap-4 text-xs text-gray-500">
              <span v-if="l.pin_count">{{ l.pin_count }} pins</span>
              <span>{{ timeAgo(l.updated_at) }}</span>
            </div>
          </div>
        </RouterLink>
      </div>
    </div>
  </div>
</template>
