<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getUser, getUserLooms } from '@/api/users'
import { getUserActivity, type Activity } from '@/api/activities'
import { timeAgo } from '@/utils/format'

const route = useRoute()

interface User {
  id: string
  username: string
  display_name: string
  bio: string
  created_at: string
}
interface Loom {
  id: string
  name: string
  full_name: string
  description: string
  pin_count: number
  updated_at: string
}

const user = ref<User | null>(null)
const looms = ref<Loom[]>([])
const activities = ref<Activity[]>([])
const loading = ref(true)
const notFound = ref(false)

async function load() {
  loading.value = true
  notFound.value = false
  const owner = route.params.owner as string
  try {
    const [userData, loomRes, activityRes] = await Promise.all([
      getUser(owner),
      getUserLooms(owner),
      getUserActivity(owner),
    ])
    user.value = userData
    looms.value = loomRes.items || []
    activities.value = activityRes || []
  } catch {
    notFound.value = true
  } finally {
    loading.value = false
  }
}

onMounted(load)
watch(() => route.params.owner, load)

function activityText(a: Activity): string {
  const payload = JSON.parse(a.payload || '{}')
  switch (a.type) {
    case 'create_loom': return `Created loom ${a.loom_full_name}`
    case 'create_wr': return `Opened WR #${payload.wr_number} on ${a.loom_full_name}`
    case 'weave': return `Wove WR #${payload.wr_number} on ${a.loom_full_name}`
    case 'review': return `Reviewed WR #${payload.wr_number} on ${a.loom_full_name}`
    case 'pin': return `Pinned ${a.loom_full_name}`
    case 'create_org': return `Created org ${payload.org}`
    default: return a.type
  }
}
</script>

<template>
  <div class="max-w-5xl mx-auto mt-8 px-4">
    <div v-if="loading" class="text-gray-500 text-sm">Loading...</div>

    <div v-else-if="notFound" class="text-center py-16">
      <p class="text-gray-500">User not found.</p>
    </div>

    <div v-else-if="user">
      <div class="mb-8">
        <div class="flex items-center gap-4">
          <div class="w-16 h-16 rounded-full bg-indigo-600 flex items-center justify-center text-2xl font-semibold text-white">
            {{ user.username[0].toUpperCase() }}
          </div>
          <div>
            <h1 class="text-xl font-semibold text-white">{{ user.display_name || user.username }}</h1>
            <p class="text-sm text-gray-500">@{{ user.username }}</p>
          </div>
        </div>
        <p v-if="user.bio" class="mt-4 text-gray-400">{{ user.bio }}</p>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Looms -->
        <div class="lg:col-span-2">
          <h2 class="text-lg font-medium text-white mb-4">Looms</h2>
          <div v-if="looms.length === 0" class="text-gray-500 text-sm">No public looms yet.</div>
          <div v-else class="space-y-3">
            <RouterLink
              v-for="l in looms"
              :key="l.id"
              :to="`/${l.full_name}`"
              class="block p-4 border border-gray-800 rounded-lg hover:border-gray-600 transition-colors"
            >
              <div class="flex items-center justify-between">
                <div>
                  <h3 class="font-medium text-indigo-400">{{ l.name }}</h3>
                  <p v-if="l.description" class="text-sm text-gray-500 mt-1">{{ l.description }}</p>
                </div>
                <span class="text-xs text-gray-500">{{ timeAgo(l.updated_at) }}</span>
              </div>
            </RouterLink>
          </div>
        </div>

        <!-- Activity -->
        <div>
          <h2 class="text-lg font-medium text-white mb-4">Activity</h2>
          <div v-if="activities.length === 0" class="text-sm text-gray-600">No recent activity</div>
          <div v-else class="space-y-3">
            <div
              v-for="a in activities"
              :key="a.id"
              class="text-sm border-l-2 border-gray-800 pl-3 py-1"
            >
              <p class="text-gray-300">{{ activityText(a) }}</p>
              <p class="text-xs text-gray-600 mt-0.5">{{ timeAgo(a.created_at) }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
