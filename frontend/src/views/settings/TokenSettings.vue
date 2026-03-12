<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { listTokens, createToken, deleteToken, type AccessToken } from '@/api/tokens'

const tokens = ref<AccessToken[]>([])
const loading = ref(true)
const newName = ref('')
const newScopes = ref('read,write')
const creating = ref(false)
const newPlainToken = ref('')
const error = ref('')

onMounted(async () => {
  await load()
})

async function load() {
  loading.value = true
  try {
    tokens.value = await listTokens()
  } catch {
    error.value = 'Failed to load tokens'
  } finally {
    loading.value = false
  }
}

async function handleCreate() {
  if (!newName.value.trim()) return
  creating.value = true
  error.value = ''
  newPlainToken.value = ''
  try {
    const res = await createToken({ name: newName.value.trim(), scopes: newScopes.value })
    newPlainToken.value = res.token
    newName.value = ''
    await load()
  } catch {
    error.value = 'Failed to create token'
  } finally {
    creating.value = false
  }
}

async function handleDelete(id: string) {
  try {
    await deleteToken(id)
    await load()
  } catch {
    error.value = 'Failed to delete token'
  }
}

function formatDate(ts: string): string {
  if (!ts) return 'Never'
  return new Date(ts).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}
</script>

<template>
  <div class="max-w-xl mx-auto mt-8 px-4">
    <h1 class="text-xl font-semibold text-white mb-6">Access Tokens</h1>
    <p class="text-sm text-gray-500 mb-6">Personal access tokens for API and CLI authentication.</p>

    <!-- Create form -->
    <form @submit.prevent="handleCreate" class="border border-gray-800 rounded-lg p-4 mb-6 space-y-3">
      <h3 class="text-sm font-medium text-gray-300">Generate new token</h3>
      <div class="flex gap-3">
        <input
          v-model="newName"
          type="text"
          placeholder="Token name"
          class="flex-1 px-3 py-2 bg-gray-900 border border-gray-700 rounded-md text-white text-sm focus:outline-none focus:border-indigo-500"
        />
        <select
          v-model="newScopes"
          class="px-3 py-2 bg-gray-900 border border-gray-700 rounded-md text-white text-sm focus:outline-none focus:border-indigo-500"
        >
          <option value="read,write">Read & Write</option>
          <option value="read">Read only</option>
        </select>
      </div>
      <button
        type="submit"
        :disabled="creating || !newName.trim()"
        class="px-4 py-2 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white rounded-md text-sm font-medium transition-colors"
      >
        {{ creating ? 'Creating...' : 'Generate token' }}
      </button>
    </form>

    <!-- Newly created token display -->
    <div v-if="newPlainToken" class="border border-green-800 bg-green-900/20 rounded-lg p-4 mb-6">
      <p class="text-sm text-green-400 mb-2">Token created. Copy it now — it won't be shown again.</p>
      <code class="block text-sm bg-gray-900 text-white p-3 rounded font-mono break-all select-all">{{ newPlainToken }}</code>
    </div>

    <div v-if="error" class="text-red-400 text-sm mb-4">{{ error }}</div>

    <!-- Token list -->
    <div v-if="loading" class="text-gray-500 text-sm">Loading tokens...</div>
    <div v-else-if="tokens.length === 0" class="text-gray-500 text-sm">No tokens yet.</div>
    <div v-else class="border border-gray-800 rounded-lg divide-y divide-gray-800">
      <div v-for="t in tokens" :key="t.id" class="flex items-center justify-between px-4 py-3">
        <div>
          <p class="text-sm text-white font-medium">{{ t.name }}</p>
          <p class="text-xs text-gray-500">
            <span>{{ t.scopes }}</span>
            <span class="mx-1">&middot;</span>
            <span>Created {{ formatDate(t.created_at) }}</span>
            <span v-if="t.last_used_at" class="mx-1">&middot;</span>
            <span v-if="t.last_used_at">Last used {{ formatDate(t.last_used_at) }}</span>
          </p>
        </div>
        <button
          @click="handleDelete(t.id)"
          class="text-xs text-red-400 hover:text-red-300 transition-colors"
        >
          Delete
        </button>
      </div>
    </div>
  </div>
</template>
