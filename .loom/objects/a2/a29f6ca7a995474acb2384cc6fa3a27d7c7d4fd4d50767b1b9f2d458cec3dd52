<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const menuOpen = ref(false)
</script>

<template>
  <header class="border-b border-gray-800 bg-gray-950/80 backdrop-blur-sm sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-4 h-14 flex items-center justify-between">
      <div class="flex items-center gap-6">
        <RouterLink to="/" class="text-lg font-semibold tracking-tight text-white hover:text-indigo-400 transition-colors">
          LoomHub
        </RouterLink>
        <RouterLink to="/explore" class="hidden sm:inline text-sm text-gray-400 hover:text-gray-200 transition-colors">
          Explore
        </RouterLink>
      </div>

      <!-- Desktop nav -->
      <div class="hidden sm:flex items-center gap-4">
        <template v-if="auth.isAuthenticated">
          <RouterLink
            v-if="auth.user?.is_admin"
            to="/admin"
            class="text-sm text-gray-400 hover:text-gray-200 transition-colors"
          >
            Admin
          </RouterLink>
          <RouterLink
            to="/settings/profile"
            class="text-sm text-gray-400 hover:text-gray-200 transition-colors"
          >
            Settings
          </RouterLink>
          <button
            @click="auth.logout()"
            class="text-sm text-gray-400 hover:text-gray-200 transition-colors"
          >
            Sign out
          </button>
          <RouterLink :to="`/${auth.user?.username}`" class="w-8 h-8 rounded-full bg-indigo-600 flex items-center justify-center text-sm font-medium text-white">
            {{ auth.user?.username?.[0]?.toUpperCase() }}
          </RouterLink>
        </template>
        <template v-else>
          <RouterLink to="/login" class="text-sm text-gray-400 hover:text-gray-200 transition-colors">
            Sign in
          </RouterLink>
          <RouterLink
            to="/register"
            class="text-sm px-3 py-1.5 bg-indigo-600 hover:bg-indigo-500 text-white rounded-md transition-colors"
          >
            Sign up
          </RouterLink>
        </template>
      </div>

      <!-- Mobile menu button -->
      <button @click="menuOpen = !menuOpen" class="sm:hidden text-gray-400 hover:text-white p-1">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path v-if="!menuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/>
          <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
        </svg>
      </button>
    </div>

    <!-- Mobile menu -->
    <div v-if="menuOpen" class="sm:hidden border-t border-gray-800 px-4 py-3 space-y-2">
      <RouterLink to="/explore" @click="menuOpen = false" class="block text-sm text-gray-400 hover:text-gray-200 py-1">Explore</RouterLink>
      <template v-if="auth.isAuthenticated">
        <RouterLink :to="`/${auth.user?.username}`" @click="menuOpen = false" class="block text-sm text-gray-400 hover:text-gray-200 py-1">Profile</RouterLink>
        <RouterLink to="/settings/profile" @click="menuOpen = false" class="block text-sm text-gray-400 hover:text-gray-200 py-1">Settings</RouterLink>
        <RouterLink v-if="auth.user?.is_admin" to="/admin" @click="menuOpen = false" class="block text-sm text-gray-400 hover:text-gray-200 py-1">Admin</RouterLink>
        <button @click="auth.logout(); menuOpen = false" class="block text-sm text-gray-400 hover:text-gray-200 py-1">Sign out</button>
      </template>
      <template v-else>
        <RouterLink to="/login" @click="menuOpen = false" class="block text-sm text-gray-400 hover:text-gray-200 py-1">Sign in</RouterLink>
        <RouterLink to="/register" @click="menuOpen = false" class="block text-sm text-indigo-400 hover:text-indigo-300 py-1">Sign up</RouterLink>
      </template>
    </div>
  </header>
</template>
