<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { api } from '@/api/client'

const route = useRoute()
const owner = computed(() => route.params.owner as string)
const loomName = computed(() => route.params.loom as string)

interface Webhook {
  id: string
  url: string
  events: string
  active: boolean
  created_at: string
}

const webhooks = ref<Webhook[]>([])
const loading = ref(true)
const showForm = ref(false)
const newURL = ref('')
const newEvents = ref('send')

async function load() {
  loading.value = true
  try {
    webhooks.value = await api(`/looms/${owner.value}/${loomName.value}/webhooks`)
  } finally {
    loading.value = false
  }
}

onMounted(load)

async function create() {
  if (!newURL.value) return
  await api(`/looms/${owner.value}/${loomName.value}/webhooks`, {
    method: 'POST',
    body: { url: newURL.value, events: newEvents.value },
  })
  newURL.value = ''
  showForm.value = false
  await load()
}

async function toggle(wh: Webhook) {
  await api(`/looms/${owner.value}/${loomName.value}/webhooks/${wh.id}`, {
    method: 'PATCH',
    body: { active: !wh.active },
  })
  await load()
}

async function remove(wh: Webhook) {
  if (!confirm('Delete this webhook?')) return
  await api(`/looms/${owner.value}/${loomName.value}/webhooks/${wh.id}`, { method: 'DELETE' })
  await load()
}
</script>

<template>
  <div class="max-w-2xl">
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-lg font-medium text-white">Webhooks</h2>
      <button @click="showForm = !showForm" class="text-sm px-3 py-1.5 bg-indigo-600 hover:bg-indigo-500 text-white rounded-md transition-colors">
        {{ showForm ? 'Cancel' : 'Add webhook' }}
      </button>
    </div>

    <form v-if="showForm" @submit.prevent="create" class="mb-6 p-4 bg-gray-900 border border-gray-800 rounded-lg space-y-3">
      <div>
        <label class="block text-sm text-gray-400 mb-1">Payload URL</label>
        <input v-model="newURL" placeholder="https://example.com/hook" class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white text-sm focus:border-indigo-500 focus:outline-none" />
      </div>
      <div>
        <label class="block text-sm text-gray-400 mb-1">Events (comma-separated)</label>
        <input v-model="newEvents" placeholder="send,weave,review" class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white text-sm focus:border-indigo-500 focus:outline-none" />
      </div>
      <button type="submit" class="px-4 py-2 bg-indigo-600 hover:bg-indigo-500 text-white text-sm rounded-md transition-colors">Create</button>
    </form>

    <div v-if="loading" class="text-gray-500 text-sm">Loading...</div>
    <div v-else-if="webhooks.length === 0" class="text-gray-500 text-sm">No webhooks configured.</div>
    <div v-else class="space-y-3">
      <div v-for="wh in webhooks" :key="wh.id" class="p-4 bg-gray-900 border border-gray-800 rounded-lg flex items-center justify-between">
        <div>
          <p class="text-sm text-white font-mono">{{ wh.url }}</p>
          <p class="text-xs text-gray-500 mt-1">Events: {{ wh.events }}</p>
        </div>
        <div class="flex items-center gap-3">
          <button @click="toggle(wh)" class="text-xs px-2 py-1 rounded" :class="wh.active ? 'bg-green-900/30 text-green-400' : 'bg-gray-800 text-gray-500'">
            {{ wh.active ? 'Active' : 'Inactive' }}
          </button>
          <button @click="remove(wh)" class="text-xs text-red-400 hover:text-red-300">Delete</button>
        </div>
      </div>
    </div>
  </div>
</template>
