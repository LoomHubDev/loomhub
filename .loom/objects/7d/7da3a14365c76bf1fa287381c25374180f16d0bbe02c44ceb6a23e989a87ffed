<script setup lang="ts">
import { computed } from 'vue'
import { useLoomStore } from '@/stores/loom'
import { useAuthStore } from '@/stores/auth'
import { timeAgo } from '@/utils/format'

const loomStore = useLoomStore()
const auth = useAuthStore()
const loom = computed(() => loomStore.loom)
</script>

<template>
  <div v-if="loom">
    <div class="border border-gray-800 rounded-lg p-6">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-lg font-medium text-white">About</h2>
        <RouterLink
          v-if="auth.isAuthenticated && auth.user?.username === loom.owner"
          :to="`/${loom.full_name}/settings`"
          class="text-xs text-gray-500 hover:text-gray-300"
        >
          Settings
        </RouterLink>
      </div>

      <p v-if="loom.description" class="text-gray-400 mb-4">{{ loom.description }}</p>
      <p v-else class="text-gray-600 italic mb-4">No description provided.</p>

      <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 text-center">
        <div class="p-3 bg-gray-900 rounded-md">
          <div class="text-lg font-semibold text-white">{{ loom.checkpoint_count }}</div>
          <div class="text-xs text-gray-500">Checkpoints</div>
        </div>
        <div class="p-3 bg-gray-900 rounded-md">
          <div class="text-lg font-semibold text-white">{{ loom.stream_count }}</div>
          <div class="text-xs text-gray-500">Streams</div>
        </div>
        <div class="p-3 bg-gray-900 rounded-md">
          <div class="text-lg font-semibold text-white">{{ loom.pin_count }}</div>
          <div class="text-xs text-gray-500">Pins</div>
        </div>
        <div class="p-3 bg-gray-900 rounded-md">
          <div class="text-lg font-semibold text-white">{{ loom.spin_count }}</div>
          <div class="text-xs text-gray-500">Spins</div>
        </div>
      </div>

      <div class="mt-6 pt-4 border-t border-gray-800">
        <h3 class="text-sm font-medium text-gray-400 mb-2">Quick setup</h3>
        <code class="block p-3 bg-gray-900 rounded-md text-sm text-gray-300 overflow-x-auto">
          loom hub add origin https://loomhub.dev/{{ loom.full_name }}
        </code>
      </div>
    </div>
  </div>
</template>
