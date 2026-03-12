<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { listWeaveRequests, type WeaveRequest } from '@/api/weave-requests'
import { useAuthStore } from '@/stores/auth'
import { timeAgo } from '@/utils/format'

const route = useRoute()
const auth = useAuthStore()
const owner = computed(() => route.params.owner as string)
const loomName = computed(() => route.params.loom as string)

const wrs = ref<WeaveRequest[]>([])
const total = ref(0)
const filter = ref<string>('open')
const loading = ref(true)

async function load() {
  loading.value = true
  try {
    const res = await listWeaveRequests(owner.value, loomName.value, filter.value || undefined)
    wrs.value = res.items || []
    total.value = res.total
  } finally {
    loading.value = false
  }
}

onMounted(load)
watch(filter, load)

function statusColor(status: string) {
  if (status === 'open') return 'text-green-400'
  if (status === 'woven') return 'text-purple-400'
  return 'text-gray-500'
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <div class="flex gap-2">
        <button
          v-for="s in ['open', 'woven', 'closed']"
          :key="s"
          @click="filter = s"
          class="px-3 py-1 text-sm rounded-md transition-colors capitalize"
          :class="filter === s ? 'bg-gray-700 text-white' : 'text-gray-400 hover:text-gray-200'"
        >
          {{ s }}
        </button>
      </div>
      <RouterLink
        v-if="auth.user"
        :to="`/${owner}/${loomName}/weave-requests/new`"
        class="px-4 py-2 bg-indigo-600 hover:bg-indigo-500 text-white text-sm rounded-md transition-colors"
      >
        New Weave Request
      </RouterLink>
    </div>

    <div v-if="loading" class="text-gray-500 text-sm">Loading...</div>
    <div v-else-if="wrs.length === 0" class="text-gray-500 text-sm py-8 text-center">
      No {{ filter }} weave requests
    </div>
    <div v-else class="space-y-1">
      <RouterLink
        v-for="wr in wrs"
        :key="wr.id"
        :to="`/${owner}/${loomName}/weave-requests/${wr.number}`"
        class="block p-4 rounded-lg bg-gray-900 border border-gray-800 hover:border-gray-700 transition-colors"
      >
        <div class="flex items-start gap-3">
          <span :class="statusColor(wr.status)" class="mt-0.5 text-lg leading-none">&#9679;</span>
          <div class="flex-1 min-w-0">
            <div class="flex items-baseline gap-2">
              <span class="font-medium text-white">{{ wr.title }}</span>
              <span class="text-gray-600 text-sm">#{{ wr.number }}</span>
            </div>
            <div class="text-xs text-gray-500 mt-1">
              {{ wr.author }} opened {{ timeAgo(wr.created_at) }}
              &middot; {{ wr.source_stream }} &rarr; {{ wr.target_stream }}
            </div>
          </div>
        </div>
      </RouterLink>
    </div>
  </div>
</template>
