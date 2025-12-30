<template>
  <div class="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
    <div class="sm:mx-auto sm:w-full sm:max-w-md">
      <div class="mx-auto flex items-center justify-center w-12 h-12 bg-primary-600 rounded-2xl shadow-lg shadow-primary-200">
        <span class="text-white text-2xl font-bold">$</span>
      </div>
      <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900 tracking-tight">
        Sign in to your account
      </h2>
      <p class="mt-2 text-center text-sm text-gray-600">
        Experience the gSpend prototype journey
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md px-4">
      <div class="bg-white py-10 px-6 shadow-xl shadow-gray-200/50 rounded-3xl sm:px-10 border border-gray-100">
        <div v-if="error" class="mb-6 p-4 bg-red-50 border border-red-100 text-red-600 text-sm rounded-2xl flex items-center">
          <svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
          {{ error }}
        </div>
        <form @submit.prevent="handleLogin" class="space-y-6">
          <div>
            <label for="email" class="block text-sm font-bold text-gray-700 ml-1 mb-1.5">
              Email address
            </label>
            <div class="mt-1">
              <input 
                id="email" 
                v-model="email" 
                type="email" 
                placeholder="demo@gspend.app"
                required 
                class="appearance-none block w-full px-4 py-3 border border-gray-200 rounded-2xl shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all" 
              />
            </div>
          </div>

          <div>
            <label for="password" class="block text-sm font-bold text-gray-700 ml-1 mb-1.5">
              Password
            </label>
            <div class="mt-1">
              <input 
                id="password" 
                v-model="password" 
                type="password" 
                placeholder="••••••••"
                required 
                class="appearance-none block w-full px-4 py-3 border border-gray-200 rounded-2xl shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all" 
              />
            </div>
          </div>

          <div class="flex items-center justify-between px-1">
            <div class="flex items-center">
              <input id="remember-me" name="remember-me" type="checkbox" class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded-lg" />
              <label for="remember-me" class="ml-2 block text-sm font-medium text-gray-600">
                Remember me
              </label>
            </div>

            <div class="text-sm">
              <a href="#" class="font-bold text-primary-600 hover:text-primary-500">
                Forgot password?
              </a>
            </div>
          </div>

          <div>
            <button 
              type="submit" 
              :disabled="isLoading"
              class="w-full flex justify-center py-4 px-4 border border-transparent rounded-2xl shadow-lg shadow-primary-200 text-sm font-bold text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="isLoading">Signing in...</span>
              <span v-else>Sign In</span>
            </button>
          </div>
        </form>

        <div class="mt-8">
          <div class="relative">
            <div class="absolute inset-0 flex items-center">
              <div class="w-full border-t border-gray-100"></div>
            </div>
            <div class="relative flex justify-center text-sm">
              <span class="px-2 bg-white text-gray-400">
                New to gSpend?
              </span>
            </div>
          </div>

          <div class="mt-6">
            <router-link 
              to="/register" 
              class="w-full flex justify-center py-3.5 px-4 border border-gray-200 rounded-2xl shadow-sm text-sm font-bold text-gray-700 bg-white hover:bg-gray-50 transition-all"
            >
              Create an account
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { authService } from '@/services/authService'
import { useAuthStore } from '@/stores/auth'
import { getOperationError } from '@/utils/errorHandling'

const email = ref('')
const password = ref('')
const error = ref('')
const isLoading = ref(false)
const router = useRouter()
const authStore = useAuthStore()

const handleLogin = async () => {
  try {
    isLoading.value = true
    error.value = ''
    
    const response = await authService.login({
      email: email.value,
      password: password.value
    })
    
    authStore.setSession(response.accessToken, response.user)
    router.push('/')
  } catch (err: any) {
    console.error('Login failed:', err)
    const friendlyError = getOperationError('auth', 'login', err)
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
