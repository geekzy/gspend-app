<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-100 px-4">
    <div class="w-full max-w-md p-8 bg-white rounded-2xl shadow-xl">
      <h1 class="text-3xl font-bold text-center text-primary-600 mb-8">Create Account</h1>
      <form @submit.prevent="handleRegister" class="space-y-6">
        <div>
          <label class="block text-sm font-medium text-gray-700">Full Name</label>
          <input type="text" v-model="name" class="w-full mt-1 p-3 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500" placeholder="John Doe" required />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700">Email</label>
          <input type="email" v-model="email" class="w-full mt-1 p-3 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500" placeholder="demo@example.com" required />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700">Password</label>
          <input type="password" v-model="password" class="w-full mt-1 p-3 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500" placeholder="••••••••" required />
        </div>
        <button type="submit" class="w-full py-3 px-4 bg-primary-600 text-white font-semibold rounded-lg hover:bg-primary-700 transition duration-200">
          Register (Demo)
        </button>
      </form>
      <p class="mt-6 text-center text-sm text-gray-600">
        Already have an account? <router-link to="/login" class="text-primary-600 font-medium hover:underline">Login</router-link>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const name = ref('')
const email = ref('')
const password = ref('')
const router = useRouter()

const handleRegister = () => {
  // Mock register
  localStorage.setItem('auth_token', 'mock_session_' + Date.now())
  localStorage.setItem('user_profile', JSON.stringify({
    id: 'mock-user-1',
    email: email.value || 'demo@example.com',
    fullName: name.value || 'Demo Family',
    familySize: 5
  }))
  router.push('/')
}
</script>
