<template>
  <MainLayout>
    <div class="px-4 sm:px-6 lg:px-8 max-w-4xl mx-auto w-full">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
          Profile Settings
        </h1>
        <p class="mt-1 text-sm text-gray-500">
          Manage your family account information and settings.
        </p>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="flex items-center justify-center min-h-[400px]">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-xl p-6 mb-6">
        <div class="flex">
          <AlertCircleIcon class="h-5 w-5 text-red-400" />
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">Error loading profile</h3>
            <p class="mt-1 text-sm text-red-700">{{ error }}</p>
            <button 
              @click="loadProfile" 
              class="mt-3 text-sm font-medium text-red-800 hover:text-red-900 underline"
            >
              Try again
            </button>
          </div>
        </div>
      </div>

      <!-- Profile Content -->
      <div v-else class="space-y-6">
        <!-- Profile Information Card -->
        <div class="bg-white shadow-sm rounded-2xl border border-gray-100 overflow-hidden">
          <div class="px-6 py-5 border-b border-gray-100">
            <h3 class="text-lg font-semibold text-gray-900">Family Information</h3>
            <p class="mt-1 text-sm text-gray-500">Your current family account details</p>
          </div>
          
          <div class="px-6 py-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- Profile Avatar and Basic Info -->
              <div class="flex items-start space-x-4">
                <div class="w-16 h-16 rounded-full bg-primary-100 flex items-center justify-center text-primary-700 font-bold text-xl border-2 border-white shadow-sm">
                  {{ userInitials }}
                </div>
                <div class="flex-1">
                  <h4 class="text-lg font-semibold text-gray-900">{{ profileData?.fullName || 'Family Name' }}</h4>
                  <p class="text-sm text-gray-500">{{ profileData?.email || 'email@example.com' }}</p>
                  <div class="mt-2 flex items-center text-sm text-gray-500">
                    <CalendarIcon class="h-4 w-4 mr-1" />
                    Member since {{ formatDate(profileData?.createdAt) }}
                  </div>
                </div>
              </div>

              <!-- Family Details -->
              <div class="space-y-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Family Name</label>
                  <div class="text-sm text-gray-900 bg-gray-50 px-3 py-2 rounded-lg border">
                    {{ profileData?.fullName || 'Not set' }}
                  </div>
                </div>
                
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Family Size</label>
                  <div class="text-sm text-gray-900 bg-gray-50 px-3 py-2 rounded-lg border">
                    {{ familySizeText }}
                  </div>
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Email Address</label>
                  <div class="text-sm text-gray-900 bg-gray-50 px-3 py-2 rounded-lg border">
                    {{ profileData?.email || 'Not set' }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Account Details Card -->
        <div class="bg-white shadow-sm rounded-2xl border border-gray-100 overflow-hidden">
          <div class="px-6 py-5 border-b border-gray-100">
            <h3 class="text-lg font-semibold text-gray-900">Account Details</h3>
            <p class="mt-1 text-sm text-gray-500">Account creation and last update information</p>
          </div>
          
          <div class="px-6 py-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Account Created</label>
                <div class="flex items-center text-sm text-gray-900">
                  <CalendarIcon class="h-4 w-4 mr-2 text-gray-400" />
                  {{ formatDate(profileData?.createdAt) }}
                </div>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Last Updated</label>
                <div class="flex items-center text-sm text-gray-900">
                  <ClockIcon class="h-4 w-4 mr-2 text-gray-400" />
                  {{ formatDate(profileData?.updatedAt) }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="bg-white shadow-sm rounded-2xl border border-gray-100 overflow-hidden">
          <div class="px-6 py-5 border-b border-gray-100">
            <h3 class="text-lg font-semibold text-gray-900">Account Management</h3>
            <p class="mt-1 text-sm text-gray-500">Update your profile information and security settings</p>
          </div>
          
          <div class="px-6 py-6">
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
              <!-- Edit Profile Button -->
              <button 
                @click="navigateToEditProfile"
                class="flex items-center justify-center px-4 py-3 border border-gray-300 rounded-xl shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-all"
              >
                <UserIcon class="h-5 w-5 mr-2" />
                Edit Profile
              </button>

              <!-- Change Password Button -->
              <button 
                @click="navigateToChangePassword"
                class="flex items-center justify-center px-4 py-3 border border-gray-300 rounded-xl shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-all"
              >
                <KeyIcon class="h-5 w-5 mr-2" />
                Change Password
              </button>

              <!-- Refresh Profile Button -->
              <button 
                @click="loadProfile"
                :disabled="isLoading"
                class="flex items-center justify-center px-4 py-3 border border-gray-300 rounded-xl shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <RefreshCwIcon class="h-5 w-5 mr-2" :class="{ 'animate-spin': isLoading }" />
                Refresh
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import { useAuthStore } from '@/stores/auth'
import { 
  UserIcon, 
  KeyIcon, 
  CalendarIcon, 
  ClockIcon, 
  RefreshCwIcon,
  AlertCircleIcon
} from 'lucide-vue-next'

interface UserProfile {
  id: string
  email: string
  fullName: string
  familySize: number
  createdAt: string
  updatedAt: string
}

const router = useRouter()
const authStore = useAuthStore()

const profileData = ref<UserProfile | null>(null)
const isLoading = ref(true)
const error = ref<string | null>(null)

const userInitials = computed(() => {
  if (!profileData.value?.fullName) return 'U'
  return profileData.value.fullName
    .split(' ')
    .map((n: string) => n[0])
    .join('')
    .toUpperCase()
    .substring(0, 2)
})

const familySizeText = computed(() => {
  if (!profileData.value) return 'Not set'
  const size = profileData.value.familySize
  if (size === 0) return 'No children'
  if (size === 1) return '1 child'
  return `${size} children`
})

const loadProfile = async () => {
  try {
    isLoading.value = true
    error.value = null
    
    // Try to fetch fresh profile data
    const userData = await authStore.fetchProfile()
    profileData.value = userData
  } catch (err) {
    console.error('Failed to fetch profile:', err)
    error.value = 'Failed to load profile information. Please try again.'
    
    // Fallback to cached user data if available
    if (authStore.user) {
      profileData.value = authStore.user
    }
  } finally {
    isLoading.value = false
  }
}

const navigateToEditProfile = () => {
  // TODO: Navigate to profile edit page when implemented
  router.push('/profile/edit')
}

const navigateToChangePassword = () => {
  // TODO: Navigate to password change page when implemented
  router.push('/profile/password')
}

const formatDate = (dateStr: string | undefined) => {
  if (!dateStr) return 'Unknown'
  
  try {
    const date = new Date(dateStr)
    return date.toLocaleDateString(undefined, { 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric' 
    })
  } catch {
    return 'Invalid date'
  }
}

onMounted(() => {
  // Initialize with cached data if available
  if (authStore.user) {
    profileData.value = authStore.user
  }
  
  // Then fetch fresh data
  loadProfile()
})
</script>