<template>
  <div class="w-full">
    <div class="flex items-center justify-between mb-2">
      <span class="text-sm font-medium text-gray-700">{{ label }}</span>
      <span class="text-sm font-bold text-gray-900">{{ percentage }}%</span>
    </div>
    <div class="w-full bg-gray-200 rounded-full h-3 overflow-hidden">
      <div 
        class="h-full rounded-full transition-all duration-1000 ease-out"
        :class="progressColorClass"
        :style="{ width: `${Math.min(percentage, 100)}%` }"
      ></div>
    </div>
    <div v-if="showAmounts" class="flex items-center justify-between mt-1 text-xs text-gray-500">
      <span>${{ spent.toLocaleString() }} spent</span>
      <span>${{ total.toLocaleString() }} budget</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  label: string
  spent: number
  total: number
  showAmounts?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showAmounts: true
})

const percentage = computed(() => {
  if (props.total === 0) return 0
  return Math.round((props.spent / props.total) * 100)
})

const progressColorClass = computed(() => {
  const pct = percentage.value
  if (pct >= 100) return 'bg-red-500'
  if (pct >= 80) return 'bg-yellow-500'
  return 'bg-blue-500'
})
</script>