<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/api/client'
import { timeAgo } from '@/utils/format'

interface Loom {
  id: string
  full_name: string
  owner: string
  name: string
  visibility: string
  pin_count: number
  checkpoint_count: number
  updated_at: string
}

const looms = ref<Loom[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    looms.value = await api('/admin/looms')
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div v-if="loading" class="text-gray-500 text-sm">Loading...</div>
  <div v-else>
    <table class="w-full text-sm">
      <thead>
        <tr class="text-left text-gray-500 border-b border-gray-800">
          <th class="pb-2 font-medium">Loom</th>
          <th class="pb-2 font-medium">Visibility</th>
          <th class="pb-2 font-medium">Pins</th>
          <th class="pb-2 font-medium">Last updated</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="l in looms" :key="l.id" class="border-b border-gray-800/50">
          <td class="py-2">
            <RouterLink :to="`/${l.full_name}`" class="text-indigo-400 hover:text-indigo-300">{{ l.full_name }}</RouterLink>
          </td>
          <td class="py-2">
            <span class="text-xs px-2 py-0.5 border border-gray-700 rounded-full text-gray-400">{{ l.visibility }}</span>
          </td>
          <td class="py-2 text-gray-400">{{ l.pin_count }}</td>
          <td class="py-2 text-gray-500">{{ timeAgo(l.updated_at) }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
