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
              Change Password
            </h1>
            <p class="mt-1 text-sm text-gray-500">
              Update your account password for security.
            </p>
          </div>
        </div>
      </div>

      <!-- Error Display -->
      <div v-if="errors.length > 0" class="mb-6 p-4 bg-red-50 border border-red-200 rounded-xl">
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

      <!-- Password Change Form -->
      <div class="bg-white shadow-sm rounded-2xl border border-gray-100 overflow-hidden">
        <div class="px-6 py-5 border-b border-gray-100">
          <h3 class="text-lg font-semibold text-gray-900">Security Settings</h3>
          <p class="mt-1 text-sm text-gray-500">Change your account password</p>
        </div>
        
        <form @submit.prevent="handlePasswordChange" class="px-6 py-6 space-y-6">
          <!-- Current Password -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Current Password <span class="text-red-500">*</span>
            </label>
            <div class="relative">
              <input 
                v-model="passwordForm.currentPassword" 
                :type="showCurrentPassword ? 'text' : 'password'"
                :class="[
                  'w-full px-4 py-3 pr-12 border rounded-xl focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none transition-all',
                  fieldErrors.currentPassword ? 'border-red-300 bg-red-50' : 'border-gray-200'
                ]"
                placeholder="Enter your current password" 
                required 
              />
              <button
                type="button"
                @click="showCurrentPassword = !showCurrentPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
              >
                <EyeIcon v-if="!showCurrentPassword" class="w-5 h-5" />
                <EyeOffIcon v-else class="w-5 h-5" />
              </button>
            </div>
            <p v-if="fieldErrors.currentPassword" class="mt-1 text-xs text-red-500">{{ fieldErrors.currentPassword }}</p>
          </div>

          <!-- New Password -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              New Password <span class="text-red-500">*</span>
            </label>
            <div class="relative">
              <input 
                v-model="passwordForm.newPassword" 
                :type="showNewPassword ? 'text' : 'password'"
                :class="[
                  'w-full px-4 py-3 pr-12 border rounded-xl focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none transition-all',
                  fieldErrors.newPassword ? 'border-red-300 bg-red-50' : 'border-gray-200'
                ]"
                placeholder="Enter your new password" 
                required 
                @input="validatePasswordStrength"
              />
              <button
                type="button"
                @click="showNewPassword = !showNewPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
              >
                <EyeIcon v-if="!showNewPassword" class="w-5 h-5" />
                <EyeOffIcon v-else class="w-5 h-5" />
              </button>
            </div>
            <p v-if="fieldErrors.newPassword" class="mt-1 text-xs text-red-500">{{ fieldErrors.newPassword }}</p>
            
            <!-- Password Strength Indicator -->
            <div v-if="passwordForm.newPassword" class="mt-3">
              <div class="flex items-center justify-between mb-2">
                <span class="text-xs font-medium text-gray-700">Password Strength</span>
                <span :class="[
                  'text-xs font-medium',
                  passwordStrength.score >= 4 ? 'text-green-600' : 
                  passwordStrength.score >= 2 ? 'text-yellow-600' : 'text-red-600'
                ]">
                  {{ passwordStrength.label }}
                </span>
              </div>
              <div class="w-full bg-gray-200 rounded-full h-2">
                <div 
                  :class="[
                    'h-2 rounded-full transition-all duration-300',
                    passwordStrength.score >= 4 ? 'bg-green-500' : 
                    passwordStrength.score >= 2 ? 'bg-yellow-500' : 'bg-red-500'
                  ]"
                  :style="{ width: `${(passwordStrength.score / 4) * 100}%` }"
                ></div>
              </div>
              
              <!-- Password Requirements -->
              <div class="mt-3 space-y-1">
                <div class="flex items-center text-xs">
                  <CheckIcon v-if="passwordRequirements.length" class="w-3 h-3 text-green-500 mr-2" />
                  <XIcon v-else class="w-3 h-3 text-red-500 mr-2" />
                  <span :class="passwordRequirements.length ? 'text-green-600' : 'text-red-600'">
                    At least 8 characters
                  </span>
                </div>
                <div class="flex items-center text-xs">
                  <CheckIcon v-if="passwordRequirements.uppercase" class="w-3 h-3 text-green-500 mr-2" />
                  <XIcon v-else class="w-3 h-3 text-red-500 mr-2" />
                  <span :class="passwordRequirements.uppercase ? 'text-green-600' : 'text-red-600'">
                    One uppercase letter
                  </span>
                </div>
                <div class="flex items-center text-xs">
                  <CheckIcon v-if="passwordRequirements.lowercase" class="w-3 h-3 text-green-500 mr-2" />
                  <XIcon v-else class="w-3 h-3 text-red-500 mr-2" />
                  <span :class="passwordRequirements.lowercase ? 'text-green-600' : 'text-red-600'">
                    One lowercase letter
                  </span>
                </div>
                <div class="flex items-center text-xs">
                  <CheckIcon v-if="passwordRequirements.number" class="w-3 h-3 text-green-500 mr-2" />
                  <XIcon v-else class="w-3 h-3 text-red-500 mr-2" />
                  <span :class="passwordRequirements.number ? 'text-green-600' : 'text-red-600'">
                    One number
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- Confirm New Password -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Confirm New Password <span class="text-red-500">*</span>
            </label>
            <div class="relative">
              <input 
                v-model="passwordForm.confirmPassword" 
                :type="showConfirmPassword ? 'text' : 'password'"
                :class="[
                  'w-full px-4 py-3 pr-12 border rounded-xl focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none transition-all',
                  fieldErrors.confirmPassword ? 'border-red-300 bg-red-50' : 'border-gray-200'
                ]"
                placeholder="Confirm your new password" 
                required 
              />
              <button
                type="button"
                @click="showConfirmPassword = !showConfirmPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
              >
                <EyeIcon v-if="!showConfirmPassword" class="w-5 h-5" />
                <EyeOffIcon v-else class="w-5 h-5" />
              </button>
            </div>
            <p v-if="fieldErrors.confirmPassword" class="mt-1 text-xs text-red-500">{{ fieldErrors.confirmPassword }}</p>
            
            <!-- Password Match Indicator -->
            <div v-if="passwordForm.confirmPassword" class="mt-2 flex items-center text-xs">
              <CheckIcon v-if="passwordsMatch" class="w-3 h-3 text-green-500 mr-2" />
              <XIcon v-else class="w-3 h-3 text-red-500 mr-2" />
              <span :class="passwordsMatch ? 'text-green-600' : 'text-red-600'">
                {{ passwordsMatch ? 'Passwords match' : 'Passwords do not match' }}
              </span>
            </div>
          </div>

          <!-- Security Notice -->
          <div class="bg-blue-50 border border-blue-200 rounded-xl p-4">
            <div class="flex items-start">
              <ShieldCheckIcon class="w-5 h-5 text-blue-500 mr-3 mt-0.5 flex-shrink-0" />
              <div class="text-sm text-blue-700">
                <p class="font-medium mb-1">Security Tips:</p>
                <ul class="space-y-1 text-xs">
                  <li>• Use a unique password that you don't use elsewhere</li>
                  <li>• Consider using a password manager</li>
                  <li>• You'll be logged out and need to sign in again</li>
                </ul>
              </div>
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
              {{ isSubmitting ? 'Changing Password...' : 'Change Password' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import { useAuthStore } from '@/stores/auth'
import { 
  ArrowLeftIcon,
  AlertCircleIcon,
  CheckCircleIcon,
  EyeIcon,
  EyeOffIcon,
  CheckIcon,
  XIcon,
  ShieldCheckIcon
} from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()

const isSubmitting = ref(false)
const errors = ref<string[]>([])
const fieldErrors = ref<Record<string, string>>({})
const successMessage = ref<string>('')

// Password visibility toggles
const showCurrentPassword = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)

// Form data
const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// Password strength analysis
const passwordRequirements = computed(() => {
  const password = passwordForm.value.newPassword
  return {
    length: password.length >= 8,
    uppercase: /[A-Z]/.test(password),
    lowercase: /[a-z]/.test(password),
    number: /\d/.test(password)
  }
})

const passwordStrength = computed(() => {
  const reqs = passwordRequirements.value
  const score = Object.values(reqs).filter(Boolean).length
  
  const labels = ['Very Weak', 'Weak', 'Fair', 'Good', 'Strong']
  return {
    score,
    label: labels[score] || 'Very Weak'
  }
})

const passwordsMatch = computed(() => {
  return passwordForm.value.newPassword === passwordForm.value.confirmPassword
})

const isFormValid = computed(() => {
  return !!(
    passwordForm.value.currentPassword &&
    passwordForm.value.newPassword &&
    passwordForm.value.confirmPassword &&
    passwordRequirements.value.length &&
    passwordRequirements.value.uppercase &&
    passwordRequirements.value.lowercase &&
    passwordRequirements.value.number &&
    passwordsMatch.value
  )
})

// Validate password strength in real-time
const validatePasswordStrength = () => {
  // Clear previous new password errors when user types
  if (fieldErrors.value.newPassword) {
    delete fieldErrors.value.newPassword
  }
  if (errors.value.length > 0) {
    errors.value = []
  }
}

// Form validation
const validateForm = (): boolean => {
  errors.value = []
  fieldErrors.value = {}

  // Current password validation
  if (!passwordForm.value.currentPassword.trim()) {
    fieldErrors.value.currentPassword = 'Current password is required'
    errors.value.push('Current password is required')
  }

  // New password validation
  if (!passwordForm.value.newPassword) {
    fieldErrors.value.newPassword = 'New password is required'
    errors.value.push('New password is required')
  } else {
    const reqs = passwordRequirements.value
    if (!reqs.length) {
      fieldErrors.value.newPassword = 'Password must be at least 8 characters long'
      errors.value.push('Password must be at least 8 characters long')
    } else if (!reqs.uppercase) {
      fieldErrors.value.newPassword = 'Password must contain at least one uppercase letter'
      errors.value.push('Password must contain at least one uppercase letter')
    } else if (!reqs.lowercase) {
      fieldErrors.value.newPassword = 'Password must contain at least one lowercase letter'
      errors.value.push('Password must contain at least one lowercase letter')
    } else if (!reqs.number) {
      fieldErrors.value.newPassword = 'Password must contain at least one number'
      errors.value.push('Password must contain at least one number')
    }
  }

  // Confirm password validation
  if (!passwordForm.value.confirmPassword) {
    fieldErrors.value.confirmPassword = 'Please confirm your new password'
    errors.value.push('Please confirm your new password')
  } else if (!passwordsMatch.value) {
    fieldErrors.value.confirmPassword = 'Passwords do not match'
    errors.value.push('Passwords do not match')
  }

  // Check if new password is same as current
  if (passwordForm.value.newPassword && passwordForm.value.currentPassword && 
      passwordForm.value.newPassword === passwordForm.value.currentPassword) {
    fieldErrors.value.newPassword = 'New password must be different from current password'
    errors.value.push('New password must be different from current password')
  }

  return errors.value.length === 0
}

// Handle form submission
const handlePasswordChange = async () => {
  if (!validateForm()) {
    return
  }

  try {
    isSubmitting.value = true
    errors.value = []
    successMessage.value = ''
    
    const changeData = {
      currentPassword: passwordForm.value.currentPassword,
      newPassword: passwordForm.value.newPassword
    }

    const response = await authStore.changePassword(changeData)
    
    successMessage.value = response.message || 'Password changed successfully! You will be logged out in a few seconds.'
    
    // Clear the form
    passwordForm.value = {
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    }
    
    // Log out user and redirect to login after password change
    setTimeout(() => {
      authStore.logout()
      router.push('/login')
    }, 3000)
    
  } catch (error: any) {
    console.error('Failed to change password:', error)
    
    // Handle different error scenarios
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
    } else if (error.response?.status === 401) {
      // Current password is incorrect
      fieldErrors.value.currentPassword = 'Current password is incorrect'
      errors.value = ['Current password is incorrect']
    } else {
      errors.value = ['Failed to change password. Please try again.']
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
watch(() => passwordForm.value.currentPassword, () => {
  if (fieldErrors.value.currentPassword) {
    delete fieldErrors.value.currentPassword
  }
  if (errors.value.length > 0) {
    errors.value = []
  }
})

watch(() => passwordForm.value.confirmPassword, () => {
  if (fieldErrors.value.confirmPassword) {
    delete fieldErrors.value.confirmPassword
  }
  if (errors.value.length > 0) {
    errors.value = []
  }
})
</script>