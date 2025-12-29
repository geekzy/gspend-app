<template>
  <button
    :disabled="isLoading || disabled"
    :class="buttonClass"
    @click="handleClick"
  >
    <LoadingSpinner
      v-if="isLoading"
      size="16px"
      :color="spinnerColor"
      inline
    />
    <slot v-else name="icon" />
    
    <span :class="{ 'ml-2': isLoading || $slots.icon }">
      <slot>{{ isLoading ? loadingText : text }}</slot>
    </span>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import LoadingSpinner from './LoadingSpinner.vue'

interface Props {
  isLoading?: boolean
  disabled?: boolean
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
  size?: 'sm' | 'md' | 'lg'
  text?: string
  loadingText?: string
  fullWidth?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isLoading: false,
  disabled: false,
  variant: 'primary',
  size: 'md',
  loadingText: 'Loading...'
})

const emit = defineEmits<{
  click: [event: MouseEvent]
}>()

const buttonClass = computed(() => {
  const baseClasses = [
    'inline-flex items-center justify-center font-medium rounded-xl transition-all duration-200',
    'focus:outline-none focus:ring-2 focus:ring-offset-2',
    'disabled:opacity-50 disabled:cursor-not-allowed'
  ]

  // Size classes
  const sizeClasses = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2 text-sm',
    lg: 'px-6 py-3 text-base'
  }

  // Variant classes
  const variantClasses = {
    primary: 'bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500 shadow-sm',
    secondary: 'bg-gray-100 text-gray-900 hover:bg-gray-200 focus:ring-gray-500 border border-gray-300',
    danger: 'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500 shadow-sm',
    ghost: 'text-gray-700 hover:bg-gray-100 focus:ring-gray-500'
  }

  const classes = [
    ...baseClasses,
    sizeClasses[props.size],
    variantClasses[props.variant]
  ]

  if (props.fullWidth) {
    classes.push('w-full')
  }

  return classes.join(' ')
})

const spinnerColor = computed(() => {
  return props.variant === 'secondary' || props.variant === 'ghost' ? 'secondary' : 'white'
})

const handleClick = (event: MouseEvent) => {
  if (!props.isLoading && !props.disabled) {
    emit('click', event)
  }
}
</script>