<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useLoomStore } from '@/stores/loom'
import { listStreams, type Stream } from '@/api/content'

const route = useRoute()
const loomStore = useLoomStore()

const owner = computed(() => route.params.owner as string)
const loomName = computed(() => route.params.loom as string)

const streams = ref<Stream[]>([])
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    streams.value = await listStreams(owner.value, loomName.value)
  } catch {
    error.value = 'Failed to load streams'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-lg font-semibold">Streams</h2>
      <span class="text-xs text-gray-500">{{ streams.length }} stream{{ streams.length !== 1 ? 's' : '' }}</span>
    </div>

    <div v-if="loading" class="text-gray-500 text-sm">Loading streams...</div>
    <div v-else-if="error" class="text-red-400 text-sm">{{ error }}</div>
    <div v-else-if="streams.length === 0" class="border border-gray-800 rounded-lg p-8 text-center">
      <p class="text-gray-500">No streams yet.</p>
      <p class="text-gray-600 text-xs mt-1">Send data to this loom to create streams.</p>
    </div>
    <div v-else class="border border-gray-800 rounded-lg divide-y divide-gray-800">
      <RouterLink
        v-for="s in streams"
        :key="s.id"
        :to="`/${owner}/${loomName}/tree/${s.name}`"
        class="flex items-center justify-between px-4 py-3 hover:bg-gray-800/50 transition-colors"
      >
        <div class="flex items-center gap-3">
          <div class="w-2 h-2 rounded-full" :class="s.status === 'active' ? 'bg-green-500' : 'bg-gray-600'" />
          <span class="text-sm font-mono">{{ s.name }}</span>
          <span v-if="s.name === loomStore.loom?.default_stream" class="text-xs px-1.5 py-0.5 bg-indigo-900/50 text-indigo-400 rounded">default</span>
        </div>
        <span class="text-xs text-gray-500">seq {{ s.head_seq }}</span>
      </RouterLink>
    </div>
  </div>
</template>
