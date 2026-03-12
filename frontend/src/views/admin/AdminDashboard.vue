<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/api/client'

const stats = ref<Record<string, number>>({})
const loading = ref(true)

onMounted(async () => {
  try {
    stats.value = await api('/admin/stats')
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div v-if="loading" class="text-gray-500 text-sm">Loading...</div>
  <div v-else class="grid grid-cols-2 md:grid-cols-4 gap-4">
    <div v-for="(value, key) in stats" :key="key" class="p-4 bg-gray-900 border border-gray-800 rounded-lg">
      <p class="text-2xl font-bold text-white">{{ value }}</p>
      <p class="text-sm text-gray-500 capitalize">{{ String(key).replace('_', ' ') }}</p>
    </div>
  </div>
</template>
