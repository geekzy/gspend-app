<template>
  <div class="bg-white rounded-2xl border border-gray-200 p-4 space-y-4">
    <div class="flex items-center justify-between">
      <h3 class="text-sm font-bold text-gray-700">Categories</h3>
      <button 
        v-if="hasActiveFilter"
        @click="clearFilter"
        class="text-xs text-gray-500 hover:text-red-500 transition-colors"
      >
        Clear
      </button>
    </div>

    <!-- Transaction Type Filter -->
    <div class="space-y-2">
      <label class="block text-xs font-medium text-gray-600">Transaction Type</label>
      <div class="grid grid-cols-3 gap-2">
        <button
          @click="setTransactionType('all')"
          :class="[
            'px-3 py-2 text-xs font-medium rounded-xl transition-all',
            selectedType === 'all'
              ? 'bg-gray-100 text-gray-700 border border-gray-300'
              : 'bg-gray-50 text-gray-600 hover:bg-gray-100 border border-transparent'
          ]"
        >
          All
        </button>
        <button
          @click="setTransactionType('income')"
          :class="[
            'px-3 py-2 text-xs font-medium rounded-xl transition-all',
            selectedType === 'income'
              ? 'bg-green-100 text-green-700 border border-green-300'
              : 'bg-gray-50 text-gray-600 hover:bg-gray-100 border border-transparent'
          ]"
        >
          Income
        </button>
        <button
          @click="setTransactionType('expense')"
          :class="[
            'px-3 py-2 text-xs font-medium rounded-xl transition-all',
            selectedType === 'expense'
              ? 'bg-red-100 text-red-700 border border-red-300'
              : 'bg-gray-50 text-gray-600 hover:bg-gray-100 border border-transparent'
          ]"
        >
          Expense
        </button>
      </div>
    </div>

    <!-- Category Selection -->
    <div class="space-y-2">
      <label class="block text-xs font-medium text-gray-600">Select Categories</label>
      
      <!-- Search Categories -->
      <div class="relative">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search categories..."
          class="w-full px-3 py-2 text-sm border border-gray-200 rounded-xl focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
        />
      </div>

      <!-- Category List -->
      <div class="max-h-48 overflow-y-auto space-y-1">
        <div v-if="filteredCategories.length === 0" class="text-xs text-gray-500 text-center py-4">
          No categories found
        </div>
        
        <label
          v-for="category in filteredCategories"
          :key="category.id"
          class="flex items-center space-x-3 p-2 rounded-xl hover:bg-gray-50 cursor-pointer transition-colors"
        >
          <input
            type="checkbox"
            :value="category.id"
            v-model="selectedCategories"
            @change="onCategoryChange"
            class="w-4 h-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
          />
          
          <div class="flex items-center space-x-2 flex-1 min-w-0">
            <span 
              class="text-sm flex-shrink-0"
              :style="{ color: category.color }"
            >
              {{ category.icon }}
            </span>
            <span class="text-sm font-medium text-gray-700 truncate">
              {{ category.name }}
            </span>
            <span 
              :class="[
                'text-xs px-2 py-0.5 rounded-full flex-shrink-0',
                category.type === 'income' 
                  ? 'bg-green-100 text-green-700' 
                  : 'bg-red-100 text-red-700'
              ]"
            >
              {{ category.type }}
            </span>
          </div>
        </label>
      </div>
    </div>

    <!-- Selected Categories Summary -->
    <div v-if="selectedCategories.length > 0" class="pt-2 border-t border-gray-100">
      <div class="flex items-center justify-between text-xs mb-2">
        <span class="text-gray-500">Selected:</span>
        <span class="font-medium text-gray-700">{{ selectedCategories.length }} categories</span>
      </div>
      
      <div class="flex flex-wrap gap-1">
        <span
          v-for="categoryId in selectedCategories.slice(0, 3)"
          :key="categoryId"
          class="inline-flex items-center space-x-1 px-2 py-1 bg-primary-100 text-primary-700 text-xs rounded-lg"
        >
          <span>{{ getCategoryById(categoryId)?.icon }}</span>
          <span>{{ getCategoryById(categoryId)?.name }}</span>
        </span>
        
        <span
          v-if="selectedCategories.length > 3"
          class="inline-flex items-center px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded-lg"
        >
          +{{ selectedCategories.length - 3 }} more
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'

export interface Category {
  id: string
  name: string
  type: 'income' | 'expense'
  icon: string
  color: string
  isSystem: boolean
}

export interface CategoryFilter {
  type: 'all' | 'income' | 'expense'
  categoryIds: string[]
}

interface Props {
  categories: Category[]
  modelValue?: CategoryFilter
}

interface Emits {
  (e: 'update:modelValue', value: CategoryFilter): void
  (e: 'change', value: CategoryFilter): void
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: () => ({ type: 'all', categoryIds: [] })
})

const emit = defineEmits<Emits>()

// Local state
const selectedType = ref<'all' | 'income' | 'expense'>('all')
const selectedCategories = ref<string[]>([])
const searchQuery = ref('')

// Computed properties
const filteredCategories = computed(() => {
  let categories = props.categories

  // Filter by type
  if (selectedType.value !== 'all') {
    categories = categories.filter(cat => cat.type === selectedType.value)
  }

  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    categories = categories.filter(cat => 
      cat.name.toLowerCase().includes(query)
    )
  }

  // Sort by system categories first, then alphabetically
  return categories.sort((a, b) => {
    if (a.isSystem && !b.isSystem) return -1
    if (!a.isSystem && b.isSystem) return 1
    return a.name.localeCompare(b.name)
  })
})

const hasActiveFilter = computed(() => {
  return selectedType.value !== 'all' || selectedCategories.value.length > 0
})

// Helper functions
function getCategoryById(id: string): Category | undefined {
  return props.categories.find(cat => cat.id === id)
}

function setTransactionType(type: 'all' | 'income' | 'expense') {
  selectedType.value = type
  
  // Clear selected categories when changing type
  if (type !== 'all') {
    selectedCategories.value = selectedCategories.value.filter(id => {
      const category = getCategoryById(id)
      return category?.type === type
    })
  }
  
  emitChange()
}

function onCategoryChange() {
  emitChange()
}

function emitChange() {
  const filter: CategoryFilter = {
    type: selectedType.value,
    categoryIds: [...selectedCategories.value]
  }
  
  emit('update:modelValue', filter)
  emit('change', filter)
}

function clearFilter() {
  selectedType.value = 'all'
  selectedCategories.value = []
  searchQuery.value = ''
  emitChange()
}

// Watch for external changes
watch(() => props.modelValue, (newValue) => {
  if (newValue) {
    selectedType.value = newValue.type
    selectedCategories.value = [...newValue.categoryIds]
  }
}, { immediate: true })
</script>