<template>
  <div v-if="error && show" class="mt-1 flex items-start gap-1.5">
    <AlertCircleIcon v-if="icon" class="w-4 h-4 text-red-500 mt-0.5 flex-shrink-0" />
    <p :class="[
      'text-xs leading-tight',
      variant === 'inline' ? 'text-red-500' : 'text-red-600'
    ]">
      {{ error }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { AlertCircleIcon } from 'lucide-vue-next'

interface Props {
  error?: string
  show?: boolean
  variant?: 'inline' | 'block'
  icon?: boolean
}

withDefaults(defineProps<Props>(), {
  error: '',
  show: true,
  variant: 'inline',
  icon: true
})
</script>