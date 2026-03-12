<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { createOrg } from '@/api/orgs'

const router = useRouter()

const name = ref('')
const displayName = ref('')
const description = ref('')
const error = ref('')
const submitting = ref(false)

async function submit() {
  error.value = ''
  if (!name.value.trim()) {
    error.value = 'Name is required'
    return
  }
  submitting.value = true
  try {
    const org = await createOrg({
      name: name.value.trim(),
      display_name: displayName.value.trim() || undefined,
      description: description.value.trim() || undefined,
    })
    router.push(`/${org.name}`)
  } catch (e: any) {
    error.value = e?.data?.message || 'Failed to create organization'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="max-w-lg mx-auto mt-16 px-4">
    <h1 class="text-2xl font-semibold text-white mb-6">Create a new organization</h1>

    <form @submit.prevent="submit" class="space-y-5">
      <div>
        <label class="block text-sm font-medium text-gray-300 mb-1.5">Organization name *</label>
        <input
          v-model="name"
          type="text"
          placeholder="my-org"
          pattern="[a-zA-Z0-9_-]+"
          class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md text-white placeholder-gray-600 focus:border-indigo-500 focus:outline-none"
        />
        <p class="mt-1 text-xs text-gray-600">Letters, numbers, hyphens, and underscores only.</p>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-300 mb-1.5">Display name</label>
        <input
          v-model="displayName"
          type="text"
          placeholder="My Organization"
          class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md text-white placeholder-gray-600 focus:border-indigo-500 focus:outline-none"
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-300 mb-1.5">Description</label>
        <textarea
          v-model="description"
          rows="3"
          placeholder="What does this organization do?"
          class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md text-white placeholder-gray-600 focus:border-indigo-500 focus:outline-none resize-none"
        />
      </div>

      <div v-if="error" class="text-sm text-red-400">{{ error }}</div>

      <button
        type="submit"
        :disabled="submitting"
        class="w-full py-2 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white rounded-md font-medium transition-colors"
      >
        {{ submitting ? 'Creating...' : 'Create organization' }}
      </button>
    </form>
  </div>
</template>
