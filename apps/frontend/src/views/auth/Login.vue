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
              class="w-full flex justify-center py-4 px-4 border border-transparent rounded-2xl shadow-lg shadow-primary-200 text-sm font-bold text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-all active:scale-95"
            >
              Sign In (Demo)
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

const email = ref('')
const password = ref('')
const router = useRouter()

const handleLogin = () => {
  localStorage.setItem('auth_token', 'mock_session_' + Date.now())
  localStorage.setItem('user_profile', JSON.stringify({
    id: 'mock-user-1',
    email: email.value || 'demo@gspend.app',
    fullName: 'Demo Family',
    familySize: 5
  }))
  router.push('/')
}
</script>
