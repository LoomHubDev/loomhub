<script setup lang="ts">
import { ref, onMounted, watch, shallowRef } from 'vue'
import { useRoute } from 'vue-router'
import { getUser } from '@/api/users'
import { getOrg } from '@/api/orgs'

import ProfileView from './Profile.vue'
import OrgProfileView from './OrgProfile.vue'

const route = useRoute()
const ownerType = ref<'user' | 'org' | null>(null)
const loading = ref(true)
const notFound = ref(false)

async function detect() {
  loading.value = true
  notFound.value = false
  ownerType.value = null
  const owner = route.params.owner as string

  try {
    await getUser(owner)
    ownerType.value = 'user'
  } catch {
    try {
      await getOrg(owner)
      ownerType.value = 'org'
    } catch {
      notFound.value = true
    }
  } finally {
    loading.value = false
  }
}

onMounted(detect)
watch(() => route.params.owner, detect)
</script>

<template>
  <div v-if="loading" class="max-w-5xl mx-auto mt-8 px-4">
    <div class="text-gray-500 text-sm">Loading...</div>
  </div>
  <div v-else-if="notFound" class="text-center py-16">
    <p class="text-gray-500">Not found.</p>
  </div>
  <ProfileView v-else-if="ownerType === 'user'" />
  <OrgProfileView v-else-if="ownerType === 'org'" />
</template>
