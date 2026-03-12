<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/api/client'
import { timeAgo } from '@/utils/format'

interface User {
  id: string
  username: string
  email: string
  display_name: string
  is_admin: boolean
  created_at: string
}

const users = ref<User[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    users.value = await api('/admin/users')
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div v-if="loading" class="text-gray-500 text-sm">Loading...</div>
  <div v-else>
    <table class="w-full text-sm">
      <thead>
        <tr class="text-left text-gray-500 border-b border-gray-800">
          <th class="pb-2 font-medium">Username</th>
          <th class="pb-2 font-medium">Email</th>
          <th class="pb-2 font-medium">Role</th>
          <th class="pb-2 font-medium">Joined</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="u in users" :key="u.id" class="border-b border-gray-800/50">
          <td class="py-2">
            <RouterLink :to="`/${u.username}`" class="text-indigo-400 hover:text-indigo-300">{{ u.username }}</RouterLink>
          </td>
          <td class="py-2 text-gray-400">{{ u.email }}</td>
          <td class="py-2">
            <span v-if="u.is_admin" class="text-xs px-2 py-0.5 bg-amber-900/30 text-amber-400 rounded-full">Admin</span>
            <span v-else class="text-xs text-gray-500">User</span>
          </td>
          <td class="py-2 text-gray-500">{{ timeAgo(u.created_at) }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
