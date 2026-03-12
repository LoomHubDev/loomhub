<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useLoomStore } from '@/stores/loom'
import { listCollaborators, addCollaborator, removeCollaborator, type Collaborator } from '@/api/collaborators'

const loomStore = useLoomStore()
const loom = computed(() => loomStore.loom)

const collaborators = ref<Collaborator[]>([])
const loading = ref(true)
const error = ref('')

const newUsername = ref('')
const newPermission = ref('write')
const adding = ref(false)
const addError = ref('')

async function load() {
  if (!loom.value) return
  loading.value = true
  error.value = ''
  try {
    collaborators.value = await listCollaborators(loom.value.owner, loom.value.name)
  } catch {
    error.value = 'Failed to load collaborators'
  } finally {
    loading.value = false
  }
}

async function handleAdd() {
  if (!loom.value || !newUsername.value.trim()) return
  adding.value = true
  addError.value = ''
  try {
    await addCollaborator(loom.value.owner, loom.value.name, newUsername.value.trim(), newPermission.value)
    newUsername.value = ''
    newPermission.value = 'write'
    await load()
  } catch (e: any) {
    addError.value = e?.data?.error?.message || 'Failed to add collaborator'
  } finally {
    adding.value = false
  }
}

async function handleRemove(username: string) {
  if (!loom.value) return
  if (!confirm(`Remove ${username} as collaborator?`)) return
  try {
    await removeCollaborator(loom.value.owner, loom.value.name, username)
    await load()
  } catch {
    error.value = 'Failed to remove collaborator'
  }
}

onMounted(load)
</script>

<template>
  <div v-if="loom" class="max-w-xl">
    <h2 class="text-lg font-medium text-white mb-6">Settings</h2>

    <!-- Settings sub-nav -->
    <div class="flex gap-4 mb-6 text-sm">
      <RouterLink :to="`/${loom.owner}/${loom.name}/settings`" class="text-gray-400 hover:text-gray-200">General</RouterLink>
      <RouterLink :to="`/${loom.owner}/${loom.name}/settings/collaborators`" class="text-white font-medium">Collaborators</RouterLink>
      <RouterLink :to="`/${loom.owner}/${loom.name}/settings/labels`" class="text-gray-400 hover:text-gray-200">Labels</RouterLink>
      <RouterLink :to="`/${loom.owner}/${loom.name}/settings/webhooks`" class="text-gray-400 hover:text-gray-200">Webhooks</RouterLink>
    </div>

    <!-- Add collaborator -->
    <form @submit.prevent="handleAdd" class="flex gap-2 mb-6">
      <input
        v-model="newUsername"
        placeholder="Username"
        class="flex-1 px-3 py-2 bg-gray-900 border border-gray-700 rounded-md text-white text-sm focus:outline-none focus:border-indigo-500"
      />
      <select
        v-model="newPermission"
        class="px-3 py-2 bg-gray-900 border border-gray-700 rounded-md text-white text-sm focus:outline-none focus:border-indigo-500"
      >
        <option value="read">Read</option>
        <option value="write">Write</option>
        <option value="admin">Admin</option>
      </select>
      <button
        type="submit"
        :disabled="adding || !newUsername.trim()"
        class="px-4 py-2 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white rounded-md text-sm font-medium transition-colors"
      >
        {{ adding ? 'Adding...' : 'Add' }}
      </button>
    </form>
    <p v-if="addError" class="text-red-400 text-sm mb-4">{{ addError }}</p>

    <!-- Collaborator list -->
    <div v-if="loading" class="text-gray-500 text-sm">Loading...</div>
    <div v-else-if="error" class="text-red-400 text-sm">{{ error }}</div>
    <div v-else-if="collaborators.length === 0" class="border border-gray-800 rounded-lg p-6 text-center">
      <p class="text-gray-500 text-sm">No collaborators yet.</p>
      <p class="text-gray-600 text-xs mt-1">Add users to give them access to this loom.</p>
    </div>
    <div v-else class="border border-gray-800 rounded-lg divide-y divide-gray-800">
      <div
        v-for="c in collaborators"
        :key="c.user_id"
        class="flex items-center justify-between px-4 py-3"
      >
        <div class="flex items-center gap-3 min-w-0">
          <div class="w-8 h-8 rounded-full bg-gray-700 flex items-center justify-center text-xs text-gray-300 font-medium shrink-0">
            {{ c.username[0].toUpperCase() }}
          </div>
          <div class="min-w-0">
            <RouterLink :to="`/${c.username}`" class="text-sm font-medium text-white hover:text-indigo-400">
              {{ c.username }}
            </RouterLink>
            <span v-if="c.display_name" class="text-xs text-gray-500 ml-1">{{ c.display_name }}</span>
          </div>
        </div>
        <div class="flex items-center gap-3 shrink-0 ml-4">
          <span
            class="text-xs px-2 py-0.5 rounded-full"
            :class="{
              'bg-green-900/30 text-green-400': c.permission === 'admin',
              'bg-blue-900/30 text-blue-400': c.permission === 'write',
              'bg-gray-800 text-gray-400': c.permission === 'read',
            }"
          >
            {{ c.permission }}
          </span>
          <button
            @click="handleRemove(c.username)"
            class="text-gray-500 hover:text-red-400 text-sm transition-colors"
          >
            Remove
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
