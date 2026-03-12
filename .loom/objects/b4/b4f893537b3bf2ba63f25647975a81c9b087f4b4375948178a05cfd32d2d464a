<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { createWeaveRequest } from '@/api/weave-requests'

const route = useRoute()
const router = useRouter()
const owner = computed(() => route.params.owner as string)
const loomName = computed(() => route.params.loom as string)

const title = ref('')
const description = ref('')
const sourceStream = ref('')
const targetStream = ref('main')
const error = ref('')
const submitting = ref(false)

async function submit() {
  if (!title.value || !sourceStream.value || !targetStream.value) {
    error.value = 'Title, source stream, and target stream are required'
    return
  }
  submitting.value = true
  error.value = ''
  try {
    const wr = await createWeaveRequest(owner.value, loomName.value, {
      title: title.value,
      description: description.value,
      source_stream: sourceStream.value,
      target_stream: targetStream.value,
    })
    router.push(`/${owner.value}/${loomName.value}/weave-requests/${wr.number}`)
  } catch (e: any) {
    error.value = e.data?.error?.message || 'Failed to create weave request'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="max-w-2xl">
    <h2 class="text-lg font-semibold text-white mb-6">New Weave Request</h2>

    <div v-if="error" class="mb-4 p-3 bg-red-900/30 border border-red-800 rounded-md text-red-400 text-sm">
      {{ error }}
    </div>

    <form @submit.prevent="submit" class="space-y-4">
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block text-sm text-gray-400 mb-1">Source stream</label>
          <input
            v-model="sourceStream"
            placeholder="feature-x"
            class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md text-white text-sm focus:border-indigo-500 focus:outline-none"
          />
        </div>
        <div>
          <label class="block text-sm text-gray-400 mb-1">Target stream</label>
          <input
            v-model="targetStream"
            placeholder="main"
            class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md text-white text-sm focus:border-indigo-500 focus:outline-none"
          />
        </div>
      </div>

      <div>
        <label class="block text-sm text-gray-400 mb-1">Title</label>
        <input
          v-model="title"
          placeholder="What does this weave request do?"
          class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md text-white text-sm focus:border-indigo-500 focus:outline-none"
        />
      </div>

      <div>
        <label class="block text-sm text-gray-400 mb-1">Description</label>
        <textarea
          v-model="description"
          rows="6"
          placeholder="Describe the changes..."
          class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md text-white text-sm focus:border-indigo-500 focus:outline-none resize-y"
        />
      </div>

      <div class="flex gap-3">
        <button
          type="submit"
          :disabled="submitting"
          class="px-4 py-2 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white text-sm rounded-md transition-colors"
        >
          {{ submitting ? 'Creating...' : 'Create Weave Request' }}
        </button>
        <RouterLink
          :to="`/${owner}/${loomName}/weave-requests`"
          class="px-4 py-2 text-gray-400 hover:text-gray-200 text-sm"
        >
          Cancel
        </RouterLink>
      </div>
    </form>
  </div>
</template>
