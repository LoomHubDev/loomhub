<script setup lang="ts">
import { onMounted, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useLoomStore } from '@/stores/loom'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const loomStore = useLoomStore()
const auth = useAuthStore()

const owner = computed(() => route.params.owner as string)
const loomName = computed(() => route.params.loom as string)

async function load() {
  await loomStore.fetchLoom(owner.value, loomName.value)
}

onMounted(load)
watch([owner, loomName], load)

const isOwner = computed(() => auth.isAuthenticated && auth.user?.username === owner.value)

const tabs = computed(() => {
  const t = [
    { name: 'Overview', to: `/${owner.value}/${loomName.value}` },
    { name: 'Tree', to: `/${owner.value}/${loomName.value}/tree` },
    { name: 'Log', to: `/${owner.value}/${loomName.value}/log` },
    { name: 'Streams', to: `/${owner.value}/${loomName.value}/streams` },
    { name: 'Weave Requests', to: `/${owner.value}/${loomName.value}/weave-requests` },
  ]
  if (isOwner.value) {
    t.push({ name: 'Settings', to: `/${owner.value}/${loomName.value}/settings` })
  }
  return t
})
</script>

<template>
  <div class="max-w-6xl mx-auto px-4 mt-6">
    <!-- Loom header -->
    <div class="mb-6" v-if="loomStore.loom">
      <h1 class="text-xl">
        <RouterLink :to="`/${owner}`" class="text-indigo-400 hover:text-indigo-300">{{ owner }}</RouterLink>
        <span class="text-gray-600 mx-1">/</span>
        <span class="font-semibold text-white">{{ loomName }}</span>
        <span v-if="loomStore.loom.visibility === 'private'" class="ml-2 text-xs px-2 py-0.5 border border-gray-700 rounded-full text-gray-400">Private</span>
      </h1>
      <p v-if="loomStore.loom.description" class="text-sm text-gray-500 mt-1">{{ loomStore.loom.description }}</p>

      <div class="flex gap-4 mt-3 text-xs text-gray-500">
        <span>{{ loomStore.loom.checkpoint_count }} checkpoints</span>
        <span>{{ loomStore.loom.stream_count }} streams</span>
        <span>{{ loomStore.loom.pin_count }} pins</span>
      </div>
    </div>

    <!-- Tabs -->
    <nav class="flex gap-1 border-b border-gray-800 mb-6">
      <RouterLink
        v-for="tab in tabs"
        :key="tab.name"
        :to="tab.to"
        class="px-4 py-2 text-sm rounded-t-md transition-colors"
        :class="$route.path === tab.to
          || (tab.to.endsWith('/tree') && ($route.path.includes('/tree') || $route.path.includes('/blob/')))
          || (tab.to.endsWith('weave-requests') && $route.path.includes('/weave-requests'))
          || (tab.to.endsWith('/settings') && $route.path.includes('/settings'))
          ? 'text-white bg-gray-800 border-b-2 border-indigo-500'
          : 'text-gray-400 hover:text-gray-200'"
      >
        {{ tab.name }}
      </RouterLink>
    </nav>

    <div v-if="loomStore.loading" class="text-gray-500 text-sm">Loading...</div>
    <RouterView v-else />
  </div>
</template>
