<template>
  <MainLayout>
    <div class="px-4 sm:px-6 lg:px-8 max-w-2xl mx-auto w-full">
      <!-- Header -->
      <div class="mb-8">
        <div class="flex items-center mb-4">
          <button 
            @click="goBack"
            class="mr-4 p-2 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 transition-all"
          >
            <ArrowLeftIcon class="w-5 h-5" />
          </button>
          <div>
            <h1 class="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
              Edit Profile
            </h1>
            <p class="mt-1 text-sm text-gray-500">
              Update your family information and account settings.
            </p>
          </div>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="flex items-center justify-center min-h-[400px]">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>

      <!-- Error Display -->
      <div v-else-if="errors.length > 0" class="mb-6 p-4 bg-red-50 border border-red-200 rounded-xl">
        <div class="flex items-center mb-2">
          <AlertCircleIcon class="w-5 h-5 text-red-500 mr-2" />
          <span class="text-sm font-bold text-red-700">Please fix the following errors:</span>
        </div>
        <ul class="text-sm text-red-600 space-y-1">
          <li v-for="error in errors" :key="error">• {{ error }}</li>
        </ul>
      </div>

      <!-- Success Message -->
      <div v-if="successMessage" class="mb-6 p-4 bg-green-50 border border-green-200 rounded-xl">
        <div class="flex items-center">
          <CheckCircleIcon class="w-5 h-5 text-green-500 mr-2" />
          <span class="text-sm font-medium text-green-700">{{ successMessage }}</span>
        </div>
      </div>

      <!-- Edit Form -->
      <div v-else class="bg-white shadow-sm rounded-2xl border border-gray-100 overflow-hidden">
        <div class="px-6 py-5 border-b border-gray-100">
          <h3 class="text-lg font-semibold text-gray-900">Family Information</h3>
          <p class="mt-1 text-sm text-gray-500">Update your family details and account information</p>
        </div>
        
        <form @submit.prevent="handleSave" class="px-6 py-6 space-y-6">
          <!-- Family Name -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Family Name <span class="text-red-500">*</span>
            </label>
            <input 
              v-model="editForm.fullName" 
              type="text" 
              :class="[
                'w-full px-4 py-3 border rounded-xl focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none transition-all',
                fieldErrors.fullName ? 'border-red-300 bg-red-50' : 'border-gray-200'
              ]"
              placeholder="Enter your family name" 
              required 
            />
            <p v-if="fieldErrors.fullName" class="mt-1 text-xs text-red-500">{{ fieldErrors.fullName }}</p>
            <p class="mt-1 text-xs text-gray-500">This will be displayed as your family account name</p>
          </div>

          <!-- Email Address -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Email Address <span class="text-red-500">*</span>
            </label>
            <input 
              v-model="editForm.email" 
              type="email" 
              :class="[
                'w-full px-4 py-3 border rounded-xl focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none transition-all',
                fieldErrors.email ? 'border-red-300 bg-red-50' : 'border-gray-200'
              ]"
              placeholder="Enter your email address" 
              required 
            />
            <p v-if="fieldErrors.email" class="mt-1 text-xs text-red-500">{{ fieldErrors.email }}</p>
            <p class="mt-1 text-xs text-gray-500">Used for account login and notifications</p>
          </div>

          <!-- Family Size -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Number of Children <span class="text-red-500">*</span>
            </label>
            <div class="grid grid-cols-6 gap-2">
              <button 
                v-for="size in familySizeOptions" 
                :key="size.value"
                type="button"
                @click="editForm.familySize = size.value"
                :class="[
                  'py-3 px-2 rounded-xl text-sm font-medium transition-all border-2',
                  editForm.familySize === size.value 
                    ? 'bg-primary-500 text-white border-primary-500 shadow-lg shadow-primary-200' 
                    : 'bg-gray-50 text-gray-700 border-gray-200 hover:bg-gray-100'
                ]"
              >
                {{ size.label }}
              </button>
            </div>
            <p v-if="fieldErrors.familySize" class="mt-1 text-xs text-red-500">{{ fieldErrors.familySize }}</p>
            <p class="mt-1 text-xs text-gray-500">
              Select the number of children in your family (maximum 5)
            </p>
          </div>

          <!-- Current Family Size Display -->
          <div class="bg-gray-50 rounded-xl p-4">
            <div class="flex items-center">
              <UsersIcon class="w-5 h-5 text-gray-400 mr-2" />
              <span class="text-sm font-medium text-gray-700">Current Selection:</span>
              <span class="ml-2 text-sm text-gray-900 font-semibold">
                {{ familySizeText }}
              </span>
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex gap-4 pt-4 border-t border-gray-100">
            <button 
              @click="goBack" 
              type="button" 
              class="flex-1 py-3 rounded-xl border border-gray-300 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-all"
              :disabled="isSubmitting"
            >
              Cancel
            </button>
            <button 
              type="submit" 
              :disabled="isSubmitting || !isFormValid" 
              class="flex-1 py-3 rounded-xl bg-primary-600 text-white text-sm font-medium hover:bg-primary-700 shadow-lg shadow-primary-200 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ isSubmitting ? 'Saving...' : 'Save Changes' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import { useAuthStore } from '@/stores/auth'
import { 
  ArrowLeftIcon,
  AlertCircleIcon,
  CheckCircleIcon,
  UsersIcon
} from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()

const isLoading = ref(true)
const isSubmitting = ref(false)
const errors = ref<string[]>([])
const fieldErrors = ref<Record<string, string>>({})
const successMessage = ref<string>('')

// Form data
const editForm = ref({
  fullName: '',
  email: '',
  familySize: 0
})

// Family size options (0-5 children)
const familySizeOptions = [
  { value: 0, label: '0' },
  { value: 1, label: '1' },
  { value: 2, label: '2' },
  { value: 3, label: '3' },
  { value: 4, label: '4' },
  { value: 5, label: '5' }
]

// Computed properties
const familySizeText = computed(() => {
  const size = editForm.value.familySize
  if (size === 0) return 'No children'
  if (size === 1) return '1 child'
  return `${size} children`
})

const isFormValid = computed(() => {
  return !!(
    editForm.value.fullName.trim() &&
    editForm.value.email.trim() &&
    editForm.value.familySize >= 0 &&
    editForm.value.familySize <= 5
  )
})

// Load current profile data
const loadProfile = async () => {
  try {
    isLoading.value = true
    
    // Use cached data if available
    if (authStore.user) {
      editForm.value = {
        fullName: authStore.user.fullName || '',
        email: authStore.user.email || '',
        familySize: authStore.user.familySize || 0
      }
    }
    
    // Fetch fresh data
    const userData = await authStore.fetchProfile()
    editForm.value = {
      fullName: userData.fullName || '',
      email: userData.email || '',
      familySize: userData.familySize || 0
    }
  } catch (error) {
    console.error('Failed to load profile:', error)
    errors.value = ['Failed to load profile information. Please try again.']
  } finally {
    isLoading.value = false
  }
}

// Form validation
const validateForm = (): boolean => {
  errors.value = []
  fieldErrors.value = {}

  // Family name validation
  if (!editForm.value.fullName.trim()) {
    fieldErrors.value.fullName = 'Family name is required'
    errors.value.push('Family name is required')
  } else if (editForm.value.fullName.trim().length < 2) {
    fieldErrors.value.fullName = 'Family name must be at least 2 characters'
    errors.value.push('Family name must be at least 2 characters')
  }

  // Email validation
  if (!editForm.value.email.trim()) {
    fieldErrors.value.email = 'Email address is required'
    errors.value.push('Email address is required')
  } else {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(editForm.value.email.trim())) {
      fieldErrors.value.email = 'Please enter a valid email address'
      errors.value.push('Please enter a valid email address')
    }
  }

  // Family size validation
  if (editForm.value.familySize < 0 || editForm.value.familySize > 5) {
    fieldErrors.value.familySize = 'Number of children must be between 0 and 5'
    errors.value.push('Number of children must be between 0 and 5')
  }

  return errors.value.length === 0
}

// Handle form submission
const handleSave = async () => {
  if (!validateForm()) {
    return
  }

  try {
    isSubmitting.value = true
    errors.value = []
    successMessage.value = ''
    
    const updateData = {
      fullName: editForm.value.fullName.trim(),
      email: editForm.value.email.trim(),
      familySize: editForm.value.familySize
    }

    const response = await authStore.updateProfile(updateData)
    
    successMessage.value = response.message || 'Profile updated successfully!'
    
    // Redirect back to profile after a short delay
    setTimeout(() => {
      router.push('/profile')
    }, 2000)
    
  } catch (error: any) {
    console.error('Failed to update profile:', error)
    
    // Handle validation errors from server
    if (error.response?.status === 400) {
      const errorData = error.response.data
      if (errorData.details) {
        // Handle field-specific errors
        errorData.details.forEach((detail: any) => {
          fieldErrors.value[detail.field] = detail.message
        })
        errors.value = errorData.details.map((detail: any) => detail.message)
      } else {
        errors.value = [errorData.error || 'Please check your input and try again']
      }
    } else if (error.response?.status === 409) {
      // Email already in use
      fieldErrors.value.email = 'This email address is already in use'
      errors.value = ['This email address is already in use by another account']
    } else {
      errors.value = ['Failed to update profile. Please try again.']
    }
  } finally {
    isSubmitting.value = false
  }
}

// Navigation
const goBack = () => {
  router.push('/profile')
}

// Clear field errors when user starts typing
watch(() => editForm.value.fullName, () => {
  if (fieldErrors.value.fullName) {
    delete fieldErrors.value.fullName
  }
  if (errors.value.length > 0) {
    errors.value = []
  }
})

watch(() => editForm.value.email, () => {
  if (fieldErrors.value.email) {
    delete fieldErrors.value.email
  }
  if (errors.value.length > 0) {
    errors.value = []
  }
})

watch(() => editForm.value.familySize, () => {
  if (fieldErrors.value.familySize) {
    delete fieldErrors.value.familySize
  }
  if (errors.value.length > 0) {
    errors.value = []
  }
})

// Initialize component
onMounted(() => {
  loadProfile()
})
</script>