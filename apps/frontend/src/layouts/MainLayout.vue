<template>
  <div class="flex h-screen bg-gray-50 overflow-hidden">
    <!-- Sidebar (Desktop) -->
    <aside class="hidden md:flex md:flex-shrink-0">
      <div class="flex flex-col w-64 border-r border-gray-200 bg-white">
        <div class="flex-1 flex flex-col pt-5 pb-4 overflow-y-auto">
          <div class="flex items-center flex-shrink-0 px-6 mb-8">
            <div class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center mr-3 shadow-sm">
              <span class="text-white font-bold">$</span>
            </div>
            <span class="text-xl font-bold text-gray-900 tracking-tight">gSpend</span>
          </div>
          <nav class="mt-2 flex-1 px-3 space-y-1">
            <router-link v-for="item in navigation" :key="item.name" :to="item.href" :class="[
              $route.path === item.href
                ? 'bg-primary-50 text-primary-600'
                : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900',
              'group flex items-center px-3 py-2.5 text-sm font-medium rounded-xl transition-all duration-200'
            ]">
              <component :is="item.icon" :class="[
                $route.path === item.href ? 'text-primary-600' : 'text-gray-400 group-hover:text-gray-500',
                'mr-3 flex-shrink-0 h-5 w-5'
              ]" aria-hidden="true" />
              {{ item.name }}
            </router-link>
          </nav>
        </div>
        <div class="flex-shrink-0 flex border-t border-gray-200 p-4">
          <div class="flex items-center w-full">
            <router-link to="/profile"
              class="flex items-center flex-1 group hover:bg-gray-50 rounded-lg p-2 -m-2 transition-colors">
              <div
                class="w-10 h-10 rounded-full bg-primary-100 flex items-center justify-center text-primary-700 font-bold border-2 border-white shadow-sm">
                {{ userInitials }}
              </div>
              <div class="ml-3 flex-1 overflow-hidden">
                <p class="text-sm font-semibold text-gray-900 truncate group-hover:text-primary-600">{{
                  authStore.user?.fullName || 'User' }}</p>
                <p class="text-xs text-gray-500 truncate">{{ authStore.user?.email || '' }}</p>
              </div>
            </router-link>
            <button @click="logout"
              class="ml-2 p-1.5 text-gray-400 hover:text-red-500 rounded-lg hover:bg-red-50 transition-colors"
              aria-label="Sign out">
              <LogOutIcon class="h-5 w-5" />
            </button>
          </div>
        </div>
      </div>
    </aside>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Mobile header -->
      <header class="md:hidden flex items-center justify-between px-4 py-3 bg-white border-b border-gray-200 shadow-sm">
        <div class="flex items-center">
          <div class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center mr-3 shadow-sm">
            <span class="text-white font-bold text-xs">$</span>
          </div>
          <span class="font-bold text-gray-900">gSpend</span>
        </div>
        <button @click="logout" class="text-gray-400" aria-label="Sign out">
          <LogOutIcon class="h-6 w-6" />
        </button>
      </header>

      <main class="flex-1 relative overflow-y-auto focus:outline-none">
        <!-- Email Verification Warning Banner -->
        <div v-if="showVerificationBanner" class="bg-amber-50 border-b border-amber-200">
          <div class="max-w-7xl mx-auto px-4 py-3 sm:px-6 lg:px-8">
            <div class="flex items-center justify-between flex-wrap gap-2">
              <div class="flex items-center">
                <span class="flex p-1 rounded-lg bg-amber-100">
                  <MailWarningIcon class="h-5 w-5 text-amber-600" />
                </span>
                <p class="ml-3 text-sm font-medium text-amber-800">
                  Please verify your email address to unlock all features.
                </p>
              </div>
              <div class="flex items-center gap-3">
                <button @click="resendVerification" :disabled="isResending"
                  class="text-sm font-semibold text-amber-800 hover:text-amber-900 underline disabled:opacity-50">
                  {{ isResending ? 'Sending...' : 'Resend verification email' }}
                </button>
                <button @click="dismissBanner" class="text-amber-600 hover:text-amber-800">
                  <XIcon class="h-5 w-5" />
                </button>
              </div>
            </div>
          </div>
        </div>
        <div class="py-6 sm:py-8 lg:py-10">
          <slot />
        </div>
      </main>

      <!-- Mobile Bottom Navigation -->
      <nav class="md:hidden bg-white border-t border-gray-200 flex justify-around items-center h-16 safe-area-bottom">
        <router-link v-for="item in navigation" :key="item.name" :to="item.href"
          class="flex flex-col items-center justify-center w-full h-full text-xs font-medium transition-colors duration-200"
          :class="[$route.path === item.href ? 'text-primary-600' : 'text-gray-400']">
          <component :is="item.icon" class="h-6 w-6 mb-1" />
          {{ item.name }}
        </router-link>
      </nav>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import {
  LayoutDashboardIcon,
  WalletIcon,
  BarChart3Icon,
  HistoryIcon,
  LogOutIcon,
  FileTextIcon,
  MailWarningIcon,
  XIcon
} from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notification'
import { authService } from '@/services/authService'

const router = useRouter()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()

const isResending = ref(false)
const bannerDismissed = ref(false)

const showVerificationBanner = computed(() => {
  return authStore.user &&
    authStore.user.emailVerified === false &&
    !bannerDismissed.value
})

const dismissBanner = () => {
  bannerDismissed.value = true
}

const resendVerification = async () => {
  if (!authStore.user?.email) return

  try {
    isResending.value = true
    await authService.resendVerification(authStore.user.email)
    notificationStore.success('Verification email sent! Check your inbox.')
  } catch (err) {
    notificationStore.error('Failed to send verification email. Please try again.')
  } finally {
    isResending.value = false
  }
}

const userInitials = computed(() => {
  if (!authStore.user?.fullName) return 'U'
  return authStore.user.fullName
    .split(' ')
    .map((n: string) => n[0])
    .join('')
    .toUpperCase()
    .substring(0, 2)
})

const navigation = [
  { name: 'Dashboard', href: '/', icon: LayoutDashboardIcon },
  { name: 'Income', href: '/income', icon: WalletIcon },
  { name: 'Budget', href: '/budget', icon: BarChart3Icon },
  { name: 'Transactions', href: '/transactions', icon: HistoryIcon },
  { name: 'Reports', href: '/reports', icon: FileTextIcon },
]

const logout = () => {
  authStore.logout()
  router.push('/login')
}
</script>
