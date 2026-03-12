<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { updateCurrentUser } from '@/api/users'

const auth = useAuthStore()
const displayName = ref(auth.user?.display_name ?? '')
const bio = ref(auth.user?.bio ?? '')
const saving = ref(false)
const saved = ref(false)

async function save() {
  saving.value = true
  saved.value = false
  try {
    await updateCurrentUser({ display_name: displayName.value, bio: bio.value })
    await auth.fetchUser()
    saved.value = true
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="max-w-xl mx-auto mt-8 px-4">
    <h1 class="text-xl font-semibold text-white mb-6">Profile settings</h1>

    <form @submit.prevent="save" class="space-y-4">
      <div>
        <label class="block text-sm text-gray-400 mb-1">Display name</label>
        <input
          v-model="displayName"
          type="text"
          class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md text-white focus:outline-none focus:border-indigo-500"
        />
      </div>

      <div>
        <label class="block text-sm text-gray-400 mb-1">Bio</label>
        <textarea
          v-model="bio"
          rows="3"
          class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md text-white focus:outline-none focus:border-indigo-500"
        ></textarea>
      </div>

      <div class="flex items-center gap-3">
        <button
          type="submit"
          :disabled="saving"
          class="px-4 py-2 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white rounded-md text-sm font-medium transition-colors"
        >
          {{ saving ? 'Saving...' : 'Save' }}
        </button>
        <span v-if="saved" class="text-sm text-green-400">Saved!</span>
      </div>
    </form>
  </div>
</template>
