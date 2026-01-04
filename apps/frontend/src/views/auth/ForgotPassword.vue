<template>
  <div class="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
    <div class="sm:mx-auto sm:w-full sm:max-w-md">
      <div class="mx-auto flex items-center justify-center w-12 h-12 bg-primary-600 rounded-2xl shadow-lg shadow-primary-200">
        <span class="text-white text-2xl font-bold">$</span>
      </div>
      <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900 tracking-tight">
        Reset your password
      </h2>
      <p class="mt-2 text-center text-sm text-gray-600">
        Enter your email and we'll send you a reset link
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md px-4">
      <div class="bg-white py-10 px-6 shadow-xl shadow-gray-200/50 rounded-3xl sm:px-10 border border-gray-100">
        <!-- Success Message -->
        <div v-if="success" class="mb-6 p-4 bg-green-50 border border-green-100 text-green-700 text-sm rounded-2xl flex items-center">
          <svg class="w-5 h-5 mr-2 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
          </svg>
          <span>{{ successMessage }}</span>
        </div>

        <!-- Error Message -->
        <div v-if="error" class="mb-6 p-4 bg-red-50 border border-red-100 text-red-600 text-sm rounded-2xl flex items-center">
          <svg class="w-5 h-5 mr-2 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
          <span>{{ error }}</span>
        </div>

        <form v-if="!success" @submit.prevent="handleSubmit" class="space-y-6">
          <div>
            <label for="email" class="block text-sm font-bold text-gray-700 ml-1 mb-1.5">
              Email address
            </label>
            <div class="mt-1">
              <input 
                id="email" 
                v-model="email" 
                type="email" 
                placeholder="your@email.com"
                required 
                class="appearance-none block w-full px-4 py-3 border border-gray-200 rounded-2xl shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all" 
              />
            </div>
          </div>

          <div>
            <button 
              type="submit" 
              :disabled="isLoading"
              class="w-full flex justify-center py-4 px-4 border border-transparent rounded-2xl shadow-lg shadow-primary-200 text-sm font-bold text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="isLoading">Sending...</span>
              <span v-else>Send Reset Link</span>
            </button>
          </div>
        </form>

        <div class="mt-6 text-center">
          <router-link to="/login" class="font-bold text-primary-600 hover:text-primary-500 text-sm">
            ← Back to Sign In
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { authService } from '@/services/authService'

const email = ref('')
const error = ref('')
const success = ref(false)
const successMessage = ref('')
const isLoading = ref(false)

const handleSubmit = async () => {
  try {
    isLoading.value = true
    error.value = ''
    
    const response = await authService.forgotPassword(email.value)
    success.value = true
    successMessage.value = response.message || 'If your email is registered, a password reset link has been sent.'
  } catch (err: any) {
    console.error('Forgot password failed:', err)
    error.value = err.response?.data?.message || 'Failed to send reset link. Please try again.'
  } finally {
    isLoading.value = false
  }
}
</script>
