<template>
  <div class="min-h-screen bg-gray-50 flex flex-col font-sans text-gray-900">
    <router-view />
    <NotificationContainer />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import NotificationContainer from '@/components/common/NotificationContainer.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

onMounted(async () => {
  const token = localStorage.getItem('auth_token')
  if (token) {
    try {
      await authStore.fetchProfile()
    } catch (err) {
      console.error('Failed to restore session:', err)
      localStorage.removeItem('auth_token')
    }
  }
})
</script>

<style>
/* Global resets if needed */
body {
  margin: 0;
  padding: 0;
}
</style>
