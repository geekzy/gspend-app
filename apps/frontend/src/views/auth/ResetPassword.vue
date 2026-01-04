<template>
  <div class="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
    <div class="sm:mx-auto sm:w-full sm:max-w-md">
      <div class="mx-auto flex items-center justify-center w-12 h-12 bg-primary-600 rounded-2xl shadow-lg shadow-primary-200">
        <span class="text-white text-2xl font-bold">$</span>
      </div>
      <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900 tracking-tight">
        Create new password
      </h2>
      <p class="mt-2 text-center text-sm text-gray-600">
        Enter your new password below
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md px-4">
      <div class="bg-white py-10 px-6 shadow-xl shadow-gray-200/50 rounded-3xl sm:px-10 border border-gray-100">
        <!-- Success Message -->
        <div v-if="success" class="text-center">
          <div class="mx-auto flex items-center justify-center w-16 h-16 bg-green-100 rounded-full mb-4">
            <svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h3 class="text-lg font-bold text-gray-900 mb-2">Password Reset Successful!</h3>
          <p class="text-gray-600 text-sm mb-6">{{ successMessage }}</p>
          <router-link 
            to="/login" 
            class="inline-flex justify-center py-3 px-6 border border-transparent rounded-2xl shadow-lg shadow-primary-200 text-sm font-bold text-white bg-primary-600 hover:bg-primary-700 transition-all"
          >
            Sign In Now
          </router-link>
        </div>

        <!-- Error Message -->
        <div v-if="error && !success" class="mb-6 p-4 bg-red-50 border border-red-100 text-red-600 text-sm rounded-2xl flex items-center">
          <svg class="w-5 h-5 mr-2 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
          <span>{{ error }}</span>
        </div>

        <form v-if="!success" @submit.prevent="handleSubmit" class="space-y-6">
          <div>
            <label for="password" class="block text-sm font-bold text-gray-700 ml-1 mb-1.5">
              New Password
            </label>
            <div class="mt-1">
              <input 
                id="password" 
                v-model="password" 
                type="password" 
                placeholder="••••••••"
                required 
                minlength="8"
                class="appearance-none block w-full px-4 py-3 border border-gray-200 rounded-2xl shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all" 
              />
            </div>
            <p class="mt-1.5 text-xs text-gray-500 ml-1">
              At least 8 characters with uppercase, lowercase, and number
            </p>
          </div>

          <div>
            <label for="confirmPassword" class="block text-sm font-bold text-gray-700 ml-1 mb-1.5">
              Confirm Password
            </label>
            <div class="mt-1">
              <input 
                id="confirmPassword" 
                v-model="confirmPassword" 
                type="password" 
                placeholder="••••••••"
                required 
                class="appearance-none block w-full px-4 py-3 border border-gray-200 rounded-2xl shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all" 
              />
            </div>
          </div>

          <div>
            <button 
              type="submit" 
              :disabled="isLoading || password !== confirmPassword"
              class="w-full flex justify-center py-4 px-4 border border-transparent rounded-2xl shadow-lg shadow-primary-200 text-sm font-bold text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="isLoading">Resetting...</span>
              <span v-else>Reset Password</span>
            </button>
          </div>
        </form>

        <div v-if="!success" class="mt-6 text-center">
          <router-link to="/login" class="font-bold text-primary-600 hover:text-primary-500 text-sm">
            ← Back to Sign In
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { authService } from '@/services/authService'

const route = useRoute()

const password = ref('')
const confirmPassword = ref('')
const error = ref('')
const success = ref(false)
const successMessage = ref('')
const isLoading = ref(false)
const token = ref('')

onMounted(() => {
  // Get token from query params
  token.value = route.query.token as string || ''
  if (!token.value) {
    error.value = 'Invalid or missing reset token. Please request a new password reset.'
  }
})

const handleSubmit = async () => {
  if (password.value !== confirmPassword.value) {
    error.value = 'Passwords do not match'
    return
  }

  try {
    isLoading.value = true
    error.value = ''
    
    const response = await authService.resetPassword(token.value, password.value)
    success.value = true
    successMessage.value = response.message || 'Your password has been reset successfully!'
  } catch (err: any) {
    console.error('Reset password failed:', err)
    error.value = err.response?.data?.message || 'Failed to reset password. The link may be expired or invalid.'
  } finally {
    isLoading.value = false
  }
}
</script>
