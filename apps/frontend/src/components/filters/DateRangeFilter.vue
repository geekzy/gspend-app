<template>
  <div class="bg-white rounded-2xl border border-gray-200 p-4 space-y-4">
    <div class="flex items-center justify-between">
      <h3 class="text-sm font-bold text-gray-700">Date Range</h3>
      <button 
        v-if="hasActiveFilter"
        @click="clearFilter"
        class="text-xs text-gray-500 hover:text-red-500 transition-colors"
      >
        Clear
      </button>
    </div>

    <!-- Quick Filter Buttons -->
    <div class="grid grid-cols-3 gap-2">
      <button
        v-for="quickFilter in quickFilters"
        :key="quickFilter.key"
        @click="applyQuickFilter(quickFilter)"
        :class="[
          'px-3 py-2 text-xs font-medium rounded-xl transition-all',
          activeQuickFilter === quickFilter.key
            ? 'bg-primary-100 text-primary-700 border border-primary-200'
            : 'bg-gray-50 text-gray-600 hover:bg-gray-100 border border-transparent'
        ]"
      >
        {{ quickFilter.label }}
      </button>
    </div>

    <!-- Custom Date Range -->
    <div class="space-y-3">
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label class="block text-xs font-medium text-gray-600 mb-1">Start Date</label>
          <input
            v-model="localStartDate"
            type="date"
            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-xl focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
            @change="onDateChange"
          />
        </div>
        <div>
          <label class="block text-xs font-medium text-gray-600 mb-1">End Date</label>
          <input
            v-model="localEndDate"
            type="date"
            class="w-full px-3 py-2 text-sm border border-gray-200 rounded-xl focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
            @change="onDateChange"
          />
        </div>
      </div>
      
      <!-- Apply Custom Range Button -->
      <button
        v-if="hasCustomDates && !isCustomRangeActive"
        @click="applyCustomRange"
        class="w-full px-3 py-2 text-sm font-medium text-primary-600 bg-primary-50 hover:bg-primary-100 rounded-xl transition-colors"
      >
        Apply Custom Range
      </button>
    </div>

    <!-- Active Filter Display -->
    <div v-if="hasActiveFilter" class="pt-2 border-t border-gray-100">
      <div class="flex items-center justify-between text-xs">
        <span class="text-gray-500">Active Filter:</span>
        <span class="font-medium text-gray-700">{{ activeFilterLabel }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'

export interface DateRange {
  startDate: string | null
  endDate: string | null
}

interface QuickFilter {
  key: string
  label: string
  startDate: string
  endDate: string
}

interface Props {
  modelValue?: DateRange
}

interface Emits {
  (e: 'update:modelValue', value: DateRange): void
  (e: 'change', value: DateRange): void
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: () => ({ startDate: null, endDate: null })
})

const emit = defineEmits<Emits>()

// Local state
const localStartDate = ref<string>('')
const localEndDate = ref<string>('')
const activeQuickFilter = ref<string | null>(null)

// Quick filter options
const quickFilters: QuickFilter[] = [
  {
    key: 'thisMonth',
    label: 'This Month',
    startDate: getMonthStart(0),
    endDate: getMonthEnd(0)
  },
  {
    key: 'lastMonth',
    label: 'Last Month',
    startDate: getMonthStart(-1),
    endDate: getMonthEnd(-1)
  },
  {
    key: 'last3Months',
    label: 'Last 3 Months',
    startDate: getMonthStart(-2),
    endDate: getMonthEnd(0)
  }
]

// Computed properties
const hasCustomDates = computed(() => {
  return localStartDate.value && localEndDate.value
})

const isCustomRangeActive = computed(() => {
  return hasCustomDates.value && 
         localStartDate.value === props.modelValue?.startDate &&
         localEndDate.value === props.modelValue?.endDate &&
         !activeQuickFilter.value
})

const hasActiveFilter = computed(() => {
  return props.modelValue?.startDate && props.modelValue?.endDate
})

const activeFilterLabel = computed(() => {
  if (activeQuickFilter.value) {
    const filter = quickFilters.find(f => f.key === activeQuickFilter.value)
    return filter?.label || ''
  }
  
  if (props.modelValue?.startDate && props.modelValue?.endDate) {
    const start = new Date(props.modelValue.startDate).toLocaleDateString('en-US', { 
      month: 'short', 
      day: 'numeric' 
    })
    const end = new Date(props.modelValue.endDate).toLocaleDateString('en-US', { 
      month: 'short', 
      day: 'numeric' 
    })
    return `${start} - ${end}`
  }
  
  return ''
})

// Helper functions
function getMonthStart(monthOffset: number): string {
  const date = new Date()
  date.setMonth(date.getMonth() + monthOffset)
  date.setDate(1)
  return date.toISOString().split('T')[0]
}

function getMonthEnd(monthOffset: number): string {
  const date = new Date()
  date.setMonth(date.getMonth() + monthOffset + 1)
  date.setDate(0)
  return date.toISOString().split('T')[0]
}

// Event handlers
function applyQuickFilter(filter: QuickFilter) {
  activeQuickFilter.value = filter.key
  localStartDate.value = filter.startDate
  localEndDate.value = filter.endDate
  
  const dateRange: DateRange = {
    startDate: filter.startDate,
    endDate: filter.endDate
  }
  
  emit('update:modelValue', dateRange)
  emit('change', dateRange)
}

function applyCustomRange() {
  if (!hasCustomDates.value) return
  
  activeQuickFilter.value = null
  
  const dateRange: DateRange = {
    startDate: localStartDate.value,
    endDate: localEndDate.value
  }
  
  emit('update:modelValue', dateRange)
  emit('change', dateRange)
}

function onDateChange() {
  // Reset quick filter when custom dates are changed
  activeQuickFilter.value = null
  
  // Auto-apply if both dates are set
  if (hasCustomDates.value) {
    applyCustomRange()
  }
}

function clearFilter() {
  activeQuickFilter.value = null
  localStartDate.value = ''
  localEndDate.value = ''
  
  const dateRange: DateRange = {
    startDate: null,
    endDate: null
  }
  
  emit('update:modelValue', dateRange)
  emit('change', dateRange)
}

// Watch for external changes
watch(() => props.modelValue, (newValue) => {
  if (newValue?.startDate) {
    localStartDate.value = newValue.startDate
  }
  if (newValue?.endDate) {
    localEndDate.value = newValue.endDate
  }
}, { immediate: true })
</script>