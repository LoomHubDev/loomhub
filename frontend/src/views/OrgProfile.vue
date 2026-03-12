<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getOrg, listOrgMembers, getOrgLooms, addOrgMember, removeOrgMember, type Organization, type OrgMember } from '@/api/orgs'
import { useAuthStore } from '@/stores/auth'
import { timeAgo } from '@/utils/format'

const route = useRoute()
const auth = useAuthStore()

const org = ref<Organization | null>(null)
const members = ref<OrgMember[]>([])
const looms = ref<any[]>([])
const loading = ref(true)
const notFound = ref(false)

const newMember = ref('')
const newRole = ref('member')
const addError = ref('')
const adding = ref(false)

const isOwner = ref(false)

async function load() {
  loading.value = true
  notFound.value = false
  const name = route.params.owner as string
  try {
    const [orgData, memberData, loomData] = await Promise.all([
      getOrg(name),
      listOrgMembers(name),
      getOrgLooms(name),
    ])
    org.value = orgData
    members.value = memberData || []
    looms.value = loomData.items || []
    isOwner.value = members.value.some(
      m => m.username === auth.user?.username && m.role === 'owner'
    )
  } catch {
    notFound.value = true
  } finally {
    loading.value = false
  }
}

async function handleAddMember() {
  if (!newMember.value.trim()) return
  addError.value = ''
  adding.value = true
  try {
    await addOrgMember(org.value!.name, newMember.value.trim(), newRole.value)
    newMember.value = ''
    newRole.value = 'member'
    const memberData = await listOrgMembers(org.value!.name)
    members.value = memberData || []
  } catch (e: any) {
    addError.value = e?.data?.message || 'Failed to add member'
  } finally {
    adding.value = false
  }
}

async function handleRemoveMember(username: string) {
  try {
    await removeOrgMember(org.value!.name, username)
    members.value = members.value.filter(m => m.username !== username)
  } catch {
    // ignore
  }
}

onMounted(load)
watch(() => route.params.owner, load)
</script>

<template>
  <div class="max-w-5xl mx-auto mt-8 px-4">
    <div v-if="loading" class="text-gray-500 text-sm">Loading...</div>

    <div v-else-if="notFound" class="text-center py-16">
      <p class="text-gray-500">Organization not found.</p>
    </div>

    <div v-else-if="org">
      <!-- Header -->
      <div class="mb-8 flex items-center gap-4">
        <div class="w-16 h-16 rounded-lg bg-indigo-600/20 border border-indigo-500/30 flex items-center justify-center text-2xl font-semibold text-indigo-400">
          {{ org.name[0].toUpperCase() }}
        </div>
        <div>
          <h1 class="text-xl font-semibold text-white">{{ org.display_name || org.name }}</h1>
          <p class="text-sm text-gray-500">@{{ org.name }}</p>
          <p v-if="org.description" class="text-sm text-gray-400 mt-1">{{ org.description }}</p>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Looms -->
        <div class="lg:col-span-2">
          <h2 class="text-lg font-medium text-white mb-4">Looms</h2>
          <div v-if="looms.length === 0" class="text-gray-500 text-sm">No looms yet.</div>
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

        <!-- Members -->
        <div>
          <h2 class="text-lg font-medium text-white mb-4">Members</h2>
          <div class="space-y-2">
            <div
              v-for="m in members"
              :key="m.user_id"
              class="flex items-center justify-between py-2"
            >
              <RouterLink :to="`/${m.username}`" class="flex items-center gap-2">
                <div class="w-7 h-7 rounded-full bg-gray-700 flex items-center justify-center text-xs font-medium text-white">
                  {{ m.username[0].toUpperCase() }}
                </div>
                <span class="text-sm text-gray-300 hover:text-white">{{ m.username }}</span>
              </RouterLink>
              <div class="flex items-center gap-2">
                <span class="text-xs text-gray-600 capitalize">{{ m.role }}</span>
                <button
                  v-if="isOwner && m.role !== 'owner'"
                  @click="handleRemoveMember(m.username)"
                  class="text-xs text-red-500 hover:text-red-400"
                >
                  Remove
                </button>
              </div>
            </div>
          </div>

          <!-- Add member form -->
          <div v-if="isOwner" class="mt-6 pt-4 border-t border-gray-800">
            <h3 class="text-sm font-medium text-gray-300 mb-3">Add member</h3>
            <form @submit.prevent="handleAddMember" class="space-y-2">
              <input
                v-model="newMember"
                type="text"
                placeholder="Username"
                class="w-full px-3 py-1.5 bg-gray-900 border border-gray-700 rounded-md text-sm text-white placeholder-gray-600 focus:border-indigo-500 focus:outline-none"
              />
              <select
                v-model="newRole"
                class="w-full px-3 py-1.5 bg-gray-900 border border-gray-700 rounded-md text-sm text-white focus:border-indigo-500 focus:outline-none"
              >
                <option value="member">Member</option>
                <option value="admin">Admin</option>
              </select>
              <div v-if="addError" class="text-xs text-red-400">{{ addError }}</div>
              <button
                type="submit"
                :disabled="adding"
                class="w-full py-1.5 text-sm bg-gray-800 hover:bg-gray-700 disabled:opacity-50 text-white rounded-md transition-colors"
              >
                {{ adding ? 'Adding...' : 'Add' }}
              </button>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
