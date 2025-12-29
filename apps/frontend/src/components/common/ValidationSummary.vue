<template>
  <div v-if="hasErrors" :class="[
    'p-4 rounded-2xl border',
    variant === 'error' ? 'bg-red-50 border-red-200' : 'bg-orange-50 border-orange-200'
  ]">
    <div class="flex items-center mb-2">
      <AlertCircleIcon :class="[
        'w-5 h-5 mr-2',
        variant === 'error' ? 'text-red-500' : 'text-orange-500'
      ]" />
      <span :class="[
        'text-sm font-bold',
        variant === 'error' ? 'text-red-700' : 'text-orange-700'
      ]">
        {{ title }}
      </span>
    </div>
    
    <div v-if="showList && errorList.length > 0">
      <ul :class="[
        'text-sm space-y-1',
        variant === 'error' ? 'text-red-600' : 'text-orange-600'
      ]">
        <li v-for="error in errorList" :key="error" class="flex items-start">
          <span class="mr-2">•</span>
          <span>{{ error }}</span>
        </li>
      </ul>
    </div>
    
    <div v-else-if="!showList && message">
      <p :class="[
        'text-sm',
        variant === 'error' ? 'text-red-600' : 'text-orange-600'
      ]">
        {{ message }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { AlertCircleIcon } from 'lucide-vue-next'

interface Props {
  errors?: Record<string, string> | string[]
  title?: string
  message?: string
  variant?: 'error' | 'warning'
  showList?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  errors: () => ({}),
  title: 'Please fix the following errors:',
  message: '',
  variant: 'error',
  showList: true
})

// Convert errors to array format
const errorList = computed(() => {
  if (Array.isArray(props.errors)) {
    return props.errors.filter(Boolean)
  }
  
  if (typeof props.errors === 'object' && props.errors !== null) {
    return Object.values(props.errors).filter(Boolean)
  }
  
  return []
})

// Check if there are any errors to display
const hasErrors = computed(() => {
  return errorList.value.length > 0 || (props.message && props.message.trim().length > 0)
})
</script>