<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const username = ref('')
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function submit() {
  error.value = ''
  loading.value = true
  try {
    await auth.register(username.value, email.value, password.value)
  } catch (e: any) {
    error.value = e?.data?.error?.message || 'Registration failed'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="max-w-sm mx-auto mt-24 px-4">
    <h1 class="text-2xl font-semibold text-white mb-8 text-center">Create your account</h1>

    <form @submit.prevent="submit" class="space-y-4">
      <div v-if="error" class="p-3 bg-red-900/30 border border-red-800 rounded-md text-sm text-red-400">
        {{ error }}
      </div>

      <div>
        <label class="block text-sm text-gray-400 mb-1">Username</label>
        <input
          v-model="username"
          type="text"
          required
          placeholder="lowercase, alphanumeric, hyphens"
          class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md text-white placeholder-gray-500 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500"
        />
      </div>

      <div>
        <label class="block text-sm text-gray-400 mb-1">Email</label>
        <input
          v-model="email"
          type="email"
          required
          class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md text-white placeholder-gray-500 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500"
        />
      </div>

      <div>
        <label class="block text-sm text-gray-400 mb-1">Password</label>
        <input
          v-model="password"
          type="password"
          required
          minlength="8"
          class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md text-white placeholder-gray-500 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500"
        />
        <p class="text-xs text-gray-500 mt-1">At least 8 characters</p>
      </div>

      <button
        type="submit"
        :disabled="loading"
        class="w-full py-2 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white rounded-md font-medium transition-colors"
      >
        {{ loading ? 'Creating account...' : 'Create account' }}
      </button>
    </form>

    <p class="mt-6 text-center text-sm text-gray-500">
      Already have an account?
      <RouterLink to="/login" class="text-indigo-400 hover:text-indigo-300">Sign in</RouterLink>
    </p>
  </div>
</template>
