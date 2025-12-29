<template>
  <div class="flex items-center justify-center" :class="containerClass">
    <div 
      class="animate-spin rounded-full border-b-2"
      :class="spinnerClass"
      :style="{ width: size, height: size }"
    ></div>
    <span v-if="message" class="ml-3 text-sm font-medium text-gray-600">
      {{ message }}
    </span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  size?: string
  color?: 'primary' | 'secondary' | 'white'
  message?: string
  fullHeight?: boolean
  inline?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: '32px',
  color: 'primary',
  fullHeight: false,
  inline: false
})

const containerClass = computed(() => {
  const classes = []
  
  if (props.fullHeight) {
    classes.push('min-h-[200px]')
  }
  
  if (props.inline) {
    classes.push('inline-flex')
  }
  
  return classes.join(' ')
})

const spinnerClass = computed(() => {
  const colorClasses = {
    primary: 'border-blue-600',
    secondary: 'border-gray-600',
    white: 'border-white'
  }
  
  return colorClasses[props.color]
})
</script>