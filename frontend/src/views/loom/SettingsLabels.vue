<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { api } from '@/api/client'

const route = useRoute()
const owner = computed(() => route.params.owner as string)
const loomName = computed(() => route.params.loom as string)

interface Label {
  id: string
  name: string
  color: string
  description: string
}

const labels = ref<Label[]>([])
const loading = ref(true)
const showForm = ref(false)
const newName = ref('')
const newColor = ref('#888888')
const newDesc = ref('')

async function load() {
  loading.value = true
  try {
    labels.value = await api(`/looms/${owner.value}/${loomName.value}/labels`)
  } finally {
    loading.value = false
  }
}

onMounted(load)

async function create() {
  if (!newName.value) return
  await api(`/looms/${owner.value}/${loomName.value}/labels`, {
    method: 'POST',
    body: { name: newName.value, color: newColor.value, description: newDesc.value },
  })
  newName.value = ''
  newColor.value = '#888888'
  newDesc.value = ''
  showForm.value = false
  await load()
}

async function remove(label: Label) {
  if (!confirm(`Delete label "${label.name}"?`)) return
  await api(`/looms/${owner.value}/${loomName.value}/labels/${label.id}`, { method: 'DELETE' })
  await load()
}
</script>

<template>
  <div class="max-w-2xl">
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-lg font-medium text-white">Labels</h2>
      <button @click="showForm = !showForm" class="text-sm px-3 py-1.5 bg-indigo-600 hover:bg-indigo-500 text-white rounded-md transition-colors">
        {{ showForm ? 'Cancel' : 'New label' }}
      </button>
    </div>

    <form v-if="showForm" @submit.prevent="create" class="mb-6 p-4 bg-gray-900 border border-gray-800 rounded-lg space-y-3">
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label class="block text-sm text-gray-400 mb-1">Name</label>
          <input v-model="newName" placeholder="bug" class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white text-sm focus:border-indigo-500 focus:outline-none" />
        </div>
        <div>
          <label class="block text-sm text-gray-400 mb-1">Color</label>
          <div class="flex gap-2">
            <input v-model="newColor" type="color" class="w-10 h-10 rounded cursor-pointer bg-transparent border-0" />
            <input v-model="newColor" class="flex-1 px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white text-sm font-mono focus:border-indigo-500 focus:outline-none" />
          </div>
        </div>
      </div>
      <div>
        <label class="block text-sm text-gray-400 mb-1">Description</label>
        <input v-model="newDesc" placeholder="Something isn't working" class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white text-sm focus:border-indigo-500 focus:outline-none" />
      </div>
      <button type="submit" class="px-4 py-2 bg-indigo-600 hover:bg-indigo-500 text-white text-sm rounded-md transition-colors">Create label</button>
    </form>

    <div v-if="loading" class="text-gray-500 text-sm">Loading...</div>
    <div v-else-if="labels.length === 0" class="text-gray-500 text-sm">No labels yet.</div>
    <div v-else class="space-y-2">
      <div v-for="label in labels" :key="label.id" class="flex items-center justify-between p-3 bg-gray-900 border border-gray-800 rounded-lg">
        <div class="flex items-center gap-3">
          <span class="w-4 h-4 rounded-full" :style="{ backgroundColor: label.color }" />
          <span class="text-sm text-white font-medium">{{ label.name }}</span>
          <span v-if="label.description" class="text-xs text-gray-500">{{ label.description }}</span>
        </div>
        <button @click="remove(label)" class="text-xs text-red-400 hover:text-red-300">Delete</button>
      </div>
    </div>
  </div>
</template>
