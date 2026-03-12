<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

if (!auth.user?.is_admin) {
  router.push('/')
}

const tabs = [
  { name: 'Dashboard', to: '/admin' },
  { name: 'Users', to: '/admin/users' },
  { name: 'Looms', to: '/admin/looms' },
]
</script>

<template>
  <div class="max-w-6xl mx-auto mt-6 px-4">
    <h1 class="text-xl font-semibold text-white mb-4">Admin</h1>
    <nav class="flex gap-1 border-b border-gray-800 mb-6">
      <RouterLink
        v-for="tab in tabs"
        :key="tab.name"
        :to="tab.to"
        class="px-4 py-2 text-sm rounded-t-md transition-colors"
        :class="route.path === tab.to
          ? 'text-white bg-gray-800 border-b-2 border-indigo-500'
          : 'text-gray-400 hover:text-gray-200'"
      >
        {{ tab.name }}
      </RouterLink>
    </nav>
    <RouterView />
  </div>
</template>
