<template>
  <div class="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
    <div class="sm:mx-auto sm:w-full sm:max-w-md">
      <div class="mx-auto flex items-center justify-center w-12 h-12 bg-primary-600 rounded-2xl shadow-lg shadow-primary-200">
        <span class="text-white text-2xl font-bold">$</span>
      </div>
      <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900 tracking-tight">
        Email Verification
      </h2>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md px-4">
      <div class="bg-white py-10 px-6 shadow-xl shadow-gray-200/50 rounded-3xl sm:px-10 border border-gray-100">
        <!-- Loading State -->
        <div v-if="isLoading" class="text-center">
          <div class="animate-spin mx-auto w-12 h-12 border-4 border-primary-200 border-t-primary-600 rounded-full mb-4"></div>
          <p class="text-gray-600">Verifying your email...</p>
        </div>

        <!-- Success State -->
        <div v-else-if="success" class="text-center">
          <div class="mx-auto flex items-center justify-center w-16 h-16 bg-green-100 rounded-full mb-4">
            <svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h3 class="text-lg font-bold text-gray-900 mb-2">Email Verified!</h3>
          <p class="text-gray-600 text-sm mb-6">{{ successMessage }}</p>
          <router-link 
            to="/login" 
            class="inline-flex justify-center py-3 px-6 border border-transparent rounded-2xl shadow-lg shadow-primary-200 text-sm font-bold text-white bg-primary-600 hover:bg-primary-700 transition-all"
          >
            Sign In Now
          </router-link>
        </div>

        <!-- Error State -->
        <div v-else class="text-center">
          <div class="mx-auto flex items-center justify-center w-16 h-16 bg-red-100 rounded-full mb-4">
            <svg class="w-8 h-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </div>
          <h3 class="text-lg font-bold text-gray-900 mb-2">Verification Failed</h3>
          <p class="text-gray-600 text-sm mb-6">{{ error }}</p>
          
          <div class="space-y-3">
            <button 
              @click="resendVerification"
              :disabled="isResending"
              class="w-full flex justify-center py-3 px-6 border border-transparent rounded-2xl shadow-lg shadow-primary-200 text-sm font-bold text-white bg-primary-600 hover:bg-primary-700 transition-all disabled:opacity-50"
            >
              <span v-if="isResending">Sending...</span>
              <span v-else>Resend Verification Email</span>
            </button>
            <router-link 
              to="/login" 
              class="block w-full py-3 px-6 border border-gray-200 rounded-2xl text-sm font-bold text-gray-700 bg-white hover:bg-gray-50 transition-all text-center"
            >
              Back to Sign In
            </router-link>
          </div>
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

const isLoading = ref(true)
const isResending = ref(false)
const success = ref(false)
const successMessage = ref('')
const error = ref('')

onMounted(async () => {
  const token = route.query.token as string
  
  if (!token) {
    isLoading.value = false
    error.value = 'Invalid or missing verification token.'
    return
  }

  try {
    const response = await authService.verifyEmail(token)
    success.value = true
    successMessage.value = response.message || 'Your email has been verified successfully!'
  } catch (err: any) {
    console.error('Email verification failed:', err)
    error.value = err.response?.data?.message || 'The verification link is invalid or has expired.'
  } finally {
    isLoading.value = false
  }
})

const resendVerification = async () => {
  const userEmail = prompt('Enter your email address to resend verification:')
  if (!userEmail) return

  try {
    isResending.value = true
    await authService.resendVerification(userEmail)
    alert('If your email is registered, a new verification link has been sent.')
  } catch (err) {
    console.error('Resend verification failed:', err)
    alert('Failed to resend verification email. Please try again.')
  } finally {
    isResending.value = false
  }
}
</script>
