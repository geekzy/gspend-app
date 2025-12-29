<template>
  <form @submit.prevent="handleSubmit" class="space-y-4">
    <!-- Type Selection -->
    <div>
      <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Type</label>
      <div class="grid grid-cols-2 gap-3">
        <button 
          type="button"
          @click="formData.type = 'expense'"
          :class="[
            formData.type === 'expense' ? 'bg-primary-600 text-white shadow-lg shadow-primary-200' : 'bg-gray-50 text-gray-500',
            'py-3 rounded-2xl text-sm font-bold transition-all'
          ]"
        >
          Expense
        </button>
        <button 
          type="button"
          @click="formData.type = 'income'"
          :class="[
            formData.type === 'income' ? 'bg-primary-600 text-white shadow-lg shadow-primary-200' : 'bg-gray-50 text-gray-500',
            'py-3 rounded-2xl text-sm font-bold transition-all'
          ]"
        >
          Income
        </button>
      </div>
      <p v-if="errors.type" class="mt-1 text-xs text-red-500 ml-1">{{ errors.type }}</p>
    </div>

    <!-- Name Input -->
    <div>
      <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Name</label>
      <input 
        v-model="formData.name" 
        type="text" 
        :class="[
          'w-full px-4 py-3 border rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none transition-colors',
          errors.name ? 'border-red-300 bg-red-50' : 'border-gray-200'
        ]"
        placeholder="e.g. Groceries, Salary, Entertainment" 
        required 
      />
      <p v-if="errors.name" class="mt-1 text-xs text-red-500 ml-1">{{ errors.name }}</p>
    </div>

    <!-- Icon and Color Selection -->
    <div class="grid grid-cols-2 gap-4">
      <!-- Icon Picker -->
      <div>
        <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Icon</label>
        <div class="relative">
          <button
            type="button"
            @click="showIconPicker = !showIconPicker"
            class="w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none bg-white flex items-center justify-between"
          >
            <span class="flex items-center">
              <span class="text-xl mr-2">{{ formData.icon }}</span>
              <span class="text-sm text-gray-600">Select icon</span>
            </span>
            <ChevronDownIcon class="w-4 h-4 text-gray-400" />
          </button>
          
          <!-- Icon Picker Dropdown -->
          <div 
            v-if="showIconPicker" 
            class="absolute z-10 mt-1 w-full bg-white border border-gray-200 rounded-2xl shadow-lg p-3 max-h-48 overflow-y-auto"
          >
            <div class="grid grid-cols-6 gap-2">
              <button
                v-for="icon in familyFriendlyIcons"
                :key="icon"
                type="button"
                @click="selectIcon(icon)"
                :class="[
                  'p-2 rounded-xl text-xl hover:bg-gray-100 transition-colors',
                  formData.icon === icon ? 'bg-primary-100 ring-2 ring-primary-500' : ''
                ]"
              >
                {{ icon }}
              </button>
            </div>
          </div>
        </div>
        <p v-if="errors.icon" class="mt-1 text-xs text-red-500 ml-1">{{ errors.icon }}</p>
      </div>

      <!-- Color Picker -->
      <div>
        <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Color</label>
        <div class="space-y-2">
          <!-- Color Preview -->
          <div 
            class="w-full h-12 rounded-2xl border border-gray-200 flex items-center justify-center text-white font-bold text-sm"
            :style="{ backgroundColor: formData.color }"
          >
            Preview
          </div>
          
          <!-- Predefined Color Palette -->
          <div class="grid grid-cols-4 gap-2">
            <button
              v-for="color in colorPalette"
              :key="color"
              type="button"
              @click="formData.color = color"
              :class="[
                'w-8 h-8 rounded-xl border-2 transition-all',
                formData.color === color ? 'border-gray-400 scale-110' : 'border-gray-200'
              ]"
              :style="{ backgroundColor: color }"
              :title="color"
            />
          </div>
          
          <!-- Custom Color Input -->
          <input 
            v-model="formData.color" 
            type="color" 
            class="w-full h-8 border border-gray-200 rounded-xl bg-white cursor-pointer" 
          />
        </div>
        <p v-if="errors.color" class="mt-1 text-xs text-red-500 ml-1">{{ errors.color }}</p>
      </div>
    </div>

    <!-- Form Actions -->
    <div class="flex gap-3 mt-8">
      <button 
        @click="$emit('cancel')" 
        type="button" 
        class="flex-1 py-3.5 rounded-2xl border border-gray-200 text-sm font-bold text-gray-700 hover:bg-gray-50 transition-all"
      >
        Cancel
      </button>
      <button 
        type="submit" 
        :disabled="isSubmitting" 
        class="flex-1 py-3.5 rounded-2xl bg-primary-600 text-white text-sm font-bold hover:bg-primary-700 shadow-lg shadow-primary-200 transition-all disabled:opacity-50"
      >
        {{ isSubmitting ? (isEditing ? 'Updating...' : 'Creating...') : (isEditing ? 'Update' : 'Create') }}
      </button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { ref, reactive, watch, onMounted, onUnmounted } from 'vue'
import { ChevronDownIcon } from 'lucide-vue-next'
import type { Category } from '@/services/financialService'

interface Props {
  category?: Partial<Category>
  isEditing?: boolean
}

interface Emits {
  (e: 'save', category: Partial<Category>): void
  (e: 'cancel'): void
}

const props = withDefaults(defineProps<Props>(), {
  category: () => ({}),
  isEditing: false
})

const emit = defineEmits<Emits>()

const isSubmitting = ref(false)
const showIconPicker = ref(false)
const errors = reactive<Record<string, string>>({})

// Family-friendly icon options organized by category
const familyFriendlyIcons = [
  // Housing & Utilities
  '🏠', '💡', '🔥', '💧', '📱', '📺', '🌐',
  
  // Food & Groceries
  '🛒', '🍽️', '🥗', '🍕', '☕', '🧊', '🍎',
  
  // Transportation
  '🚗', '⛽', '🚌', '🚇', '✈️', '🚲', '🛴',
  
  // Children & Family
  '👶', '🎒', '📚', '🎨', '⚽', '🎮', '🧸',
  '👕', '👟', '🎂', '🎪', '🏊', '🎵', '📖',
  
  // Healthcare & Personal
  '🏥', '💊', '🦷', '👓', '💄', '🧴', '✂️',
  
  // Work & Income
  '💼', '💰', '📈', '🏢', '🛠️', '💻', '📊',
  
  // Entertainment & Hobbies
  '🎬', '🎭', '🎪', '🎨', '📷', '🎸', '🏃',
  
  // Shopping & Miscellaneous
  '🛍️', '🎁', '💳', '🏪', '📦', '🔧', '🧹',
  
  // Savings & Investment
  '🏦', '💎', '📊', '💹', '🔒', '🎯', '📈'
]

// Predefined color palette optimized for categories
const colorPalette = [
  '#EF4444', '#F97316', '#F59E0B', '#EAB308', // Reds & Oranges
  '#84CC16', '#22C55E', '#10B981', '#14B8A6', // Greens
  '#06B6D4', '#0EA5E9', '#3B82F6', '#6366F1', // Blues
  '#8B5CF6', '#A855F7', '#D946EF', '#EC4899', // Purples & Pinks
  '#6B7280', '#374151', '#1F2937', '#111827'  // Grays
]

const formData = reactive<Partial<Category>>({
  name: '',
  type: 'expense',
  icon: '🛒',
  color: '#34d399'
})

// Initialize form data from props
const initializeForm = () => {
  Object.assign(formData, {
    name: props.category?.name || '',
    type: props.category?.type || 'expense',
    icon: props.category?.icon || '🛒',
    color: props.category?.color || '#34d399'
  })
}

// Validate form data
const validateForm = (): boolean => {
  Object.keys(errors).forEach(key => delete errors[key])
  
  if (!formData.name?.trim()) {
    errors.name = 'Category name is required'
  } else if (formData.name.trim().length < 2) {
    errors.name = 'Category name must be at least 2 characters'
  } else if (formData.name.trim().length > 50) {
    errors.name = 'Category name must be less than 50 characters'
  }
  
  if (!formData.type) {
    errors.type = 'Category type is required'
  }
  
  if (!formData.icon) {
    errors.icon = 'Please select an icon'
  }
  
  if (!formData.color) {
    errors.color = 'Please select a color'
  } else if (!/^#[0-9A-F]{6}$/i.test(formData.color)) {
    errors.color = 'Please select a valid color'
  }
  
  return Object.keys(errors).length === 0
}

const selectIcon = (icon: string) => {
  formData.icon = icon
  showIconPicker.value = false
}

const handleSubmit = async () => {
  if (!validateForm()) {
    return
  }
  
  try {
    isSubmitting.value = true
    await emit('save', { ...formData })
  } catch (err) {
    console.error('Failed to save category:', err)
    // Handle API errors here if needed
  } finally {
    isSubmitting.value = false
  }
}

// Close icon picker when clicking outside
const handleClickOutside = (event: Event) => {
  const target = event.target as HTMLElement
  if (!target.closest('.relative')) {
    showIconPicker.value = false
  }
}

// Watch for prop changes
watch(() => props.category, initializeForm, { immediate: true, deep: true })

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>