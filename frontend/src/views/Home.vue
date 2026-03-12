<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { getCurrentUserLooms } from '@/api/looms'
import { getActivityFeed, type Activity } from '@/api/activities'
import { timeAgo } from '@/utils/format'

const auth = useAuthStore()

interface Loom {
  id: string
  name: string
  full_name: string
  description: string
  visibility: string
  pin_count: number
  updated_at: string
}

const looms = ref<Loom[]>([])
const activities = ref<Activity[]>([])
const loading = ref(false)

onMounted(async () => {
  if (auth.isAuthenticated) {
    loading.value = true
    try {
      const [loomRes, feedRes] = await Promise.all([
        getCurrentUserLooms(),
        getActivityFeed(),
      ])
      looms.value = loomRes.items || []
      activities.value = feedRes || []
    } finally {
      loading.value = false
    }
  }
})

function activityText(a: Activity): string {
  const payload = JSON.parse(a.payload || '{}')
  switch (a.type) {
    case 'create_loom': return `created loom ${a.loom_full_name}`
    case 'create_wr': return `opened WR #${payload.wr_number} on ${a.loom_full_name}`
    case 'weave': return `wove WR #${payload.wr_number} on ${a.loom_full_name}`
    case 'review': return `reviewed WR #${payload.wr_number} on ${a.loom_full_name}`
    case 'pin': return `pinned ${a.loom_full_name}`
    case 'create_org': return `created org ${payload.org}`
    default: return a.type
  }
}
</script>

<template>
  <!-- Guest landing -->
  <div v-if="!auth.isAuthenticated" class="max-w-2xl mx-auto mt-32 px-4 text-center">
    <h1 class="text-4xl font-bold text-white mb-4">LoomHub</h1>
    <p class="text-lg text-gray-400 mb-8">
      The hosting platform for Loom — the versioning system built for 2026 and beyond.
    </p>
    <div class="flex gap-4 justify-center">
      <RouterLink
        to="/register"
        class="px-6 py-2.5 bg-indigo-600 hover:bg-indigo-500 text-white rounded-md font-medium transition-colors"
      >
        Get started
      </RouterLink>
      <RouterLink
        to="/explore"
        class="px-6 py-2.5 border border-gray-700 hover:border-gray-500 text-gray-300 rounded-md font-medium transition-colors"
      >
        Explore looms
      </RouterLink>
    </div>
  </div>

  <!-- Dashboard -->
  <div v-else class="max-w-5xl mx-auto mt-8 px-4">
    <div v-if="loading" class="text-gray-500 text-sm">Loading...</div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Looms column -->
      <div class="lg:col-span-2">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-white">Your Looms</h2>
        </div>

        <div v-if="looms.length === 0" class="text-center py-12">
          <p class="text-gray-500 mb-4">You don't have any looms yet.</p>
          <p class="text-sm text-gray-600">
            Create one from the CLI: <code class="text-indigo-400">loom hub add origin https://loomhub.dev/{{ auth.user?.username }}/my-app</code>
          </p>
        </div>

        <div v-else class="space-y-3">
          <RouterLink
            v-for="l in looms"
            :key="l.id"
            :to="`/${l.full_name}`"
            class="block p-4 border border-gray-800 rounded-lg hover:border-gray-600 transition-colors"
          >
            <div class="flex items-center justify-between">
              <div>
                <h3 class="font-medium text-indigo-400">{{ l.full_name }}</h3>
                <p v-if="l.description" class="text-sm text-gray-500 mt-1">{{ l.description }}</p>
              </div>
              <div class="flex items-center gap-4 text-xs text-gray-500">
                <span v-if="l.visibility === 'private'" class="px-2 py-0.5 border border-gray-700 rounded-full">Private</span>
                <span>{{ timeAgo(l.updated_at) }}</span>
              </div>
            </div>
          </RouterLink>
        </div>
      </div>

      <!-- Activity feed column -->
      <div>
        <h2 class="text-lg font-semibold text-white mb-4">Activity</h2>
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
</template>
