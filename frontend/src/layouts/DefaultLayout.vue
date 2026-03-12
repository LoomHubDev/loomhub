<script setup lang="ts">
import AppNav from '@/components/nav/AppNav.vue'
</script>

<template>
  <div class="min-h-screen flex flex-col">
    <AppNav />
    <main class="flex-1">
      <RouterView />
    </main>
  </div>
</template>
