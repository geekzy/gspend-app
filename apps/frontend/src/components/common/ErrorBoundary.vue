<template>
  <div v-if="hasError" class="min-h-[400px] flex items-center justify-center p-8">
    <div class="max-w-md w-full bg-white rounded-2xl shadow-lg border border-gray-100 p-8 text-center">
      <div class="w-16 h-16 mx-auto mb-6 bg-red-100 rounded-full flex items-center justify-center">
        <AlertTriangleIcon class="w-8 h-8 text-red-600" />
      </div>
      
      <h3 class="text-lg font-bold text-gray-900 mb-2">
        {{ error?.title || 'Something went wrong' }}
      </h3>
      
      <p class="text-sm text-gray-600 mb-6">
        {{ error?.message || 'An unexpected error occurred. Please try refreshing the page.' }}
      </p>
      
      <div v-if="error?.suggestion" class="bg-blue-50 rounded-xl p-4 mb-6">
        <p class="text-sm text-blue-800">
          <strong>Suggestion:</strong> {{ error.suggestion }}
        </p>
      </div>
      
      <div class="flex flex-col sm:flex-row gap-3 justify-center">
        <button
          @click="retry"
          class="px-4 py-2 bg-blue-600 text-white rounded-xl hover:bg-blue-700 transition-colors font-medium"
        >
          Try Again
        </button>
        
        <button
          v-if="error?.action"
          @click="error.action.handler"
          class="px-4 py-2 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition-colors font-medium"
        >
          {{ error.action.label }}
        </button>
        
        <button
          @click="goHome"
          class="px-4 py-2 text-gray-600 hover:text-gray-800 transition-colors font-medium"
        >
          Go Home
        </button>
      </div>
      
      <details v-if="showDetails" class="mt-6 text-left">
        <summary class="cursor-pointer text-sm text-gray-500 hover:text-gray-700">
          Technical Details
        </summary>
        <pre class="mt-2 text-xs bg-gray-100 rounded p-3 overflow-auto text-gray-700">{{ errorDetails }}</pre>
      </details>
    </div>
  </div>
  
  <slot v-else />
</template>

<script setup lang="ts">
import { ref, onErrorCaptured, provide } from 'vue'
import { useRouter } from 'vue-router'
import { AlertTriangleIcon } from 'lucide-vue-next'
import type { UserFriendlyError } from '@/utils/errorHandling'
import { transformError } from '@/utils/errorHandling'

interface Props {
  showDetails?: boolean
  fallbackComponent?: any
}

const props = withDefaults(defineProps<Props>(), {
  showDetails: false
})

const emit = defineEmits<{
  error: [error: any]
  retry: []
}>()

const router = useRouter()
const hasError = ref(false)
const error = ref<UserFriendlyError | null>(null)
const errorDetails = ref<string>('')

// Capture errors from child components
onErrorCaptured((err: any, _instance: any, info: string) => {
  console.error('Error boundary caught error:', err, 'Info:', info)
  
  hasError.value = true
  error.value = transformError(err)
  errorDetails.value = `${err.message}\n\nStack: ${err.stack}\n\nComponent: ${info}`
  
  emit('error', err)
  
  // Prevent the error from propagating further
  return false
})

// Provide error handling to child components
provide('errorBoundary', {
  reportError: (err: any) => {
    hasError.value = true
    error.value = transformError(err)
    errorDetails.value = err.stack || err.message || 'Unknown error'
    emit('error', err)
  }
})

const retry = () => {
  hasError.value = false
  error.value = null
  errorDetails.value = ''
  emit('retry')
}

const goHome = () => {
  router.push('/')
}

// Reset error state when component is unmounted or route changes
const reset = () => {
  hasError.value = false
  error.value = null
  errorDetails.value = ''
}

defineExpose({
  reset,
  reportError: (err: any) => {
    hasError.value = true
    error.value = transformError(err)
    errorDetails.value = err.stack || err.message || 'Unknown error'
  }
})
</script>