import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as loomsApi from '@/api/looms'

interface Loom {
  id: string
  owner: string
  name: string
  full_name: string
  description: string
  visibility: string
  default_stream: string
  checkpoint_count: number
  stream_count: number
  pin_count: number
  spin_count: number
  created_at: string
  updated_at: string
}

export const useLoomStore = defineStore('loom', () => {
  const loom = ref<Loom | null>(null)
  const activeStream = ref('main')
  const loading = ref(false)

  async function fetchLoom(owner: string, name: string) {
    loading.value = true
    try {
      loom.value = await loomsApi.getLoom(owner, name)
      activeStream.value = loom.value?.default_stream ?? 'main'
    } finally {
      loading.value = false
    }
  }

  function clear() {
    loom.value = null
    activeStream.value = 'main'
  }

  return { loom, activeStream, loading, fetchLoom, clear }
})
