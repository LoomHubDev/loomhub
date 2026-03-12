<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useLoomStore } from '@/stores/loom'
import { updateLoom, deleteLoom } from '@/api/looms'

const router = useRouter()
const loomStore = useLoomStore()
const loom = computed(() => loomStore.loom)

const description = ref(loom.value?.description ?? '')
const visibility = ref(loom.value?.visibility ?? 'public')
const saving = ref(false)
const deleting = ref(false)

async function save() {
  if (!loom.value) return
  saving.value = true
  try {
    await updateLoom(loom.value.owner, loom.value.name, {
      description: description.value,
      visibility: visibility.value,
    })
    await loomStore.fetchLoom(loom.value.owner, loom.value.name)
  } finally {
    saving.value = false
  }
}

async function confirmDelete() {
  if (!loom.value) return
  if (!confirm(`Delete ${loom.value.full_name}? This cannot be undone.`)) return
  deleting.value = true
  try {
    await deleteLoom(loom.value.owner, loom.value.name)
    router.push('/')
  } finally {
    deleting.value = false
  }
}
</script>

<template>
  <div v-if="loom" class="max-w-xl">
    <h2 class="text-lg font-medium text-white mb-6">Settings</h2>

    <!-- Settings sub-nav -->
    <div class="flex gap-4 mb-6 text-sm">
      <RouterLink :to="`/${loom.owner}/${loom.name}/settings`" class="text-white font-medium">General</RouterLink>
      <RouterLink :to="`/${loom.owner}/${loom.name}/settings/collaborators`" class="text-gray-400 hover:text-gray-200">Collaborators</RouterLink>
      <RouterLink :to="`/${loom.owner}/${loom.name}/settings/labels`" class="text-gray-400 hover:text-gray-200">Labels</RouterLink>
      <RouterLink :to="`/${loom.owner}/${loom.name}/settings/webhooks`" class="text-gray-400 hover:text-gray-200">Webhooks</RouterLink>
    </div>

    <form @submit.prevent="save" class="space-y-4 mb-12">
      <div>
        <label class="block text-sm text-gray-400 mb-1">Description</label>
        <input
          v-model="description"
          type="text"
          class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md text-white focus:outline-none focus:border-indigo-500"
        />
      </div>

      <div>
        <label class="block text-sm text-gray-400 mb-1">Visibility</label>
        <select
          v-model="visibility"
          class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md text-white focus:outline-none focus:border-indigo-500"
        >
          <option value="public">Public</option>
          <option value="private">Private</option>
        </select>
      </div>

      <button
        type="submit"
        :disabled="saving"
        class="px-4 py-2 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white rounded-md text-sm font-medium transition-colors"
      >
        {{ saving ? 'Saving...' : 'Save changes' }}
      </button>
    </form>

    <div class="border-t border-red-900/50 pt-6">
      <h3 class="text-sm font-medium text-red-400 mb-2">Danger zone</h3>
      <button
        @click="confirmDelete"
        :disabled="deleting"
        class="px-4 py-2 border border-red-800 text-red-400 hover:bg-red-900/20 rounded-md text-sm transition-colors"
      >
        {{ deleting ? 'Deleting...' : 'Delete this loom' }}
      </button>
    </div>
  </div>
</template>
