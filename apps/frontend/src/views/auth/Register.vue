<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-100 px-4">
    <div class="w-full max-w-md p-8 bg-white rounded-3xl shadow-xl">
      <h1 class="text-3xl font-bold text-center text-primary-600 mb-8">Create Account</h1>
      
      <div v-if="error" class="mb-6 p-4 bg-red-50 border border-red-100 text-red-600 text-sm rounded-2xl flex items-center">
        {{ error }}
      </div>

      <form @submit.prevent="handleRegister" class="space-y-6">
        <div>
          <label class="block text-sm font-medium text-gray-700">Full Name</label>
          <input type="text" v-model="name" class="w-full mt-1 p-3 border border-gray-300 rounded-2xl focus:ring-primary-500 focus:border-primary-500" placeholder="John Doe" required />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700">Email</label>
          <input type="email" v-model="email" class="w-full mt-1 p-3 border border-gray-300 rounded-2xl focus:ring-primary-500 focus:border-primary-500" placeholder="demo@example.com" required />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700">Password</label>
          <input type="password" v-model="password" class="w-full mt-1 p-3 border border-gray-300 rounded-2xl focus:ring-primary-500 focus:border-primary-500" placeholder="••••••••" required />
        </div>
        <button 
          type="submit" 
          :disabled="isLoading"
          class="w-full py-3 px-4 bg-primary-600 text-white font-semibold rounded-2xl hover:bg-primary-700 transition duration-200 disabled:opacity-50"
        >
          <span v-if="isLoading">Creating account...</span>
          <span v-else>Register</span>
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
import { authService } from '@/services/authService'
import { useAuthStore } from '@/stores/auth'
import { getOperationError } from '@/utils/errorHandling'

const name = ref('')
const email = ref('')
const password = ref('')
const error = ref('')
const isLoading = ref(false)
const router = useRouter()
const authStore = useAuthStore()

const handleRegister = async () => {
  try {
    isLoading.value = true
    error.value = ''
    
    const response = await authService.register({
      email: email.value,
      password: password.value,
      fullName: name.value,
      familySize: 0  // Default to 0, user can update in profile later
    })
    
    authStore.setSession(response.accessToken, response.user)
    router.push('/')
  } catch (err: any) {
    console.error('Registration failed:', err)
    const friendlyError = getOperationError('auth', 'register', err)
    error.value = friendlyError.message
    
    // Add suggestion if available
    if (friendlyError.suggestion) {
      error.value += ` ${friendlyError.suggestion}`
    }
  } finally {
    isLoading.value = false
  }
}
</script>
