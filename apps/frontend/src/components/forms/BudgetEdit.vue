<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm overflow-y-auto">
    <div class="bg-white rounded-3xl w-full max-w-4xl p-8 shadow-2xl animate-in fade-in zoom-in duration-200 my-8">
      <div class="flex justify-between items-center mb-8">
        <h3 class="text-2xl font-black text-gray-900">Edit Budget</h3>
        <button @click="$emit('close')" class="p-2 text-gray-400 hover:text-gray-600">
          <XIcon class="w-6 h-6" />
        </button>
      </div>

      <!-- Error Display -->
      <ValidationSummary 
        v-if="hasErrors || Object.keys(itemFieldErrors).length > 0"
        :errors="{ ...errors, ...Object.values(itemFieldErrors).reduce((acc, errs) => ({ ...acc, ...errs }), {}) }"
        class="mb-6"
      />

      <form @submit.prevent="handleSave" class="space-y-8">
        <!-- Budget Basic Info -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">
                Budget Title <span class="text-red-500">*</span>
              </label>
              <input 
                v-model="editForm.name" 
                type="text" 
                :class="getFieldClass('name', 'w-full px-4 py-3 border rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none', 'border-red-300 bg-red-50')"
                placeholder="e.g. October 2024 Household" 
                required 
                @blur="handleFieldBlur('name')"
              />
              <ErrorMessage :error="getFieldError('name')" />
            </div>
            
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">
                  Start Date <span class="text-red-500">*</span>
                </label>
                <input 
                  v-model="editForm.startDate" 
                  type="date" 
                  :class="getFieldClass('startDate', 'w-full px-4 py-3 border rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none', 'border-red-300 bg-red-50')"
                  required 
                  @blur="handleFieldBlur('startDate')"
                />
                <ErrorMessage :error="getFieldError('startDate')" />
              </div>
              <div>
                <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">
                  End Date <span class="text-red-500">*</span>
                </label>
                <input 
                  v-model="editForm.endDate" 
                  type="date" 
                  :class="getFieldClass('endDate', 'w-full px-4 py-3 border rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none', 'border-red-300 bg-red-50')"
                  required 
                  @blur="handleFieldBlur('endDate')"
                />
                <ErrorMessage :error="getFieldError('endDate')" />
              </div>
            </div>
          </div>

          <!-- Budget Summary -->
          <div class="bg-gray-50 p-6 rounded-3xl">
            <div class="flex items-center text-primary-600 mb-4">
              <InfoIcon class="w-5 h-5 mr-2" />
              <span class="text-sm font-bold uppercase tracking-tight">Budget Summary</span>
            </div>
            <div class="space-y-2">
              <div class="flex justify-between text-sm">
                <span class="text-gray-500">Categories</span>
                <span class="font-bold font-mono text-gray-900">{{ editForm.items.length }}</span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-gray-500">Original Total</span>
                <span class="font-bold font-mono text-gray-600">${{ originalTotalAmount.toLocaleString() }}</span>
              </div>
              <div class="flex justify-between text-lg font-black border-t border-gray-200 pt-2 mt-2">
                <span class="text-gray-900">New Total</span>
                <span :class="totalPlannedAmount !== originalTotalAmount ? 'text-orange-600' : 'text-primary-600'">
                  ${{ totalPlannedAmount.toLocaleString() }}
                </span>
              </div>
              <div v-if="totalPlannedAmount !== originalTotalAmount" class="text-xs text-orange-600 font-medium">
                {{ totalPlannedAmount > originalTotalAmount ? '+' : '' }}${{ (totalPlannedAmount - originalTotalAmount).toLocaleString() }} change
              </div>
            </div>
          </div>
        </div>

        <!-- Budget Items -->
        <div>
          <div class="flex justify-between items-center mb-4">
            <h4 class="text-lg font-bold text-gray-900">Category Allocations</h4>
            <button 
              type="button" 
              @click="addItem" 
              class="inline-flex items-center px-3 py-2 text-primary-600 text-sm font-bold hover:bg-primary-50 rounded-xl transition-colors"
            >
              <PlusIcon class="w-4 h-4 mr-1" />
              Add Category
            </button>
          </div>
          
          <div class="space-y-4 max-h-[400px] overflow-y-auto pr-2 custom-scrollbar">
            <div 
              v-for="(item, index) in editForm.items" 
              :key="item.id || index" 
              class="grid grid-cols-12 gap-3 items-start bg-gray-50 p-4 rounded-2xl relative"
            >
              <!-- Category Selection -->
              <div class="col-span-5">
                <select 
                  v-model="item.categoryId" 
                  :class="[
                    'w-full px-4 py-2 bg-white border rounded-xl outline-none text-sm',
                    getItemFieldError(index, 'categoryId') ? 'border-red-300 bg-red-50' : 'border-gray-200'
                  ]"
                  @change="updateCategoryName(index)"
                  @blur="handleItemFieldBlur(index, 'categoryId')"
                >
                  <option value="">Select Category</option>
                  <option v-for="cat in expenseCategories" :key="cat.id" :value="cat.id">
                    {{ cat.icon }} {{ cat.name }}
                  </option>
                </select>
                <ErrorMessage :error="getItemFieldError(index, 'categoryId')" />
              </div>
              
              <!-- Planned Amount -->
              <div class="col-span-3 relative">
                <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 text-xs">$</span>
                <input 
                  v-model.number="item.plannedAmount" 
                  type="number" 
                  step="0.01"
                  min="0.01"
                  :class="[
                    'w-full pl-6 pr-4 py-2 bg-white border rounded-xl outline-none text-sm',
                    getItemFieldError(index, 'plannedAmount') ? 'border-red-300 bg-red-50' : 'border-gray-200'
                  ]"
                  placeholder="Amount" 
                  @blur="handleItemFieldBlur(index, 'plannedAmount')"
                />
                <ErrorMessage :error="getItemFieldError(index, 'plannedAmount')" />
              </div>
              
              <!-- Spent Amount (Read-only) -->
              <div class="col-span-3">
                <div class="px-4 py-2 bg-gray-100 border border-gray-200 rounded-xl text-sm text-gray-600">
                  ${{ (item.spentAmount || 0).toLocaleString() }} spent
                </div>
              </div>
              
              <!-- Remove Button -->
              <button 
                type="button" 
                @click="removeItem(index)" 
                class="col-span-1 p-2 text-gray-400 hover:text-red-500 transition-colors"
                title="Remove category"
              >
                <XCircleIcon class="w-5 h-5" />
              </button>
              
              <!-- Notes (Full Width) -->
              <div class="col-span-12 mt-2">
                <textarea 
                  v-model="item.notes" 
                  rows="1" 
                  class="w-full px-4 py-2 bg-white border border-gray-200 rounded-xl outline-none text-sm resize-none" 
                  placeholder="Notes for this category (optional)..."
                ></textarea>
              </div>
            </div>
          </div>
          
          <div v-if="editForm.items.length === 0" class="text-center py-8 text-gray-400">
            <PlusCircleIcon class="w-12 h-12 mx-auto mb-2 opacity-20" />
            <p>No categories added yet. Click "Add Category" to get started.</p>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex gap-4 pt-4 border-t border-gray-100">
          <button 
            @click="$emit('close')" 
            type="button" 
            class="flex-1 py-4 rounded-2xl border-2 border-gray-100 text-sm font-bold text-gray-500 hover:bg-gray-50 transition-all"
            :disabled="isSubmitting"
          >
            Cancel
          </button>
          <button 
            type="submit" 
            :disabled="isSubmitting || !isFormValid" 
            class="flex-1 py-4 rounded-2xl bg-primary-600 text-white text-sm font-black hover:bg-primary-700 shadow-xl shadow-primary-200 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ isSubmitting ? 'Saving Budget...' : 'Save Changes' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { 
  XIcon, 
  InfoIcon, 
  PlusIcon, 
  XCircleIcon, 
  PlusCircleIcon 
} from 'lucide-vue-next'
import { financialService, type Budget, type Category, type BudgetItem } from '@/services/financialService'
import { useFormValidation } from '@/composables/useFormValidation'
import { useFormNotifications } from '@/composables/useFormNotifications'
import { FormValidators, ValidationRules, FieldValidator } from '@/utils/validation'
import ValidationSummary from '@/components/common/ValidationSummary.vue'
import ErrorMessage from '@/components/common/ErrorMessage.vue'

interface Props {
  budget: Budget
  categories: Category[]
}

interface Emits {
  (e: 'close'): void
  (e: 'saved', budget: Budget): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const isSubmitting = ref(false)

// Form data
const editForm = ref({
  name: props.budget.name,
  periodType: props.budget.periodType,
  startDate: props.budget.startDate.split('T')[0], // Convert to YYYY-MM-DD format
  endDate: props.budget.endDate.split('T')[0],
  totalAmount: props.budget.totalAmount,
  items: [...props.budget.items] // Create a copy of items
})

const originalTotalAmount = props.budget.totalAmount

// Set up form validation with custom budget validator
const budgetValidator = FormValidators.budget()
  .addField('totalAmount', new FieldValidator([
    ValidationRules.required('Total amount is required'),
    ValidationRules.positiveAmount('Total amount must be greater than 0')
  ]))

const {
  errors,
  isValid,
  hasErrors,
  validateForm,
  validateField,
  clearFieldError,
  touchField,
  getFieldError,
  getFieldClass
} = useFormValidation({
  validator: budgetValidator,
  formData: editForm,
  realTimeValidation: true
})

// Set up form notifications
const { handleFormSubmission, handleValidationResult } = useFormNotifications({
  successMessage: 'Budget updated successfully',
  errorMessage: 'Failed to update budget'
})

// Available expense categories
const expenseCategories = computed(() => {
  return props.categories.filter(c => c.type === 'expense')
})

// Calculate total planned amount
const totalPlannedAmount = computed(() => {
  return editForm.value.items.reduce((sum, item) => sum + (Number(item.plannedAmount) || 0), 0)
})

// Form validation - enhanced with item validation
const isFormValid = computed(() => {
  return isValid.value && !!(
    editForm.value.name.trim() &&
    editForm.value.startDate &&
    editForm.value.endDate &&
    editForm.value.items.length > 0 &&
    editForm.value.items.every(item => item.categoryId && item.plannedAmount > 0)
  )
})

// Item validation state
const itemFieldErrors = ref<Record<number, Record<string, string>>>({})

// Helper function to get item field errors
const getItemFieldError = (index: number, field: string): string | undefined => {
  return itemFieldErrors.value[index]?.[field]
}

// Validate individual budget items
const validateBudgetItem = (index: number, item: BudgetItem): boolean => {
  const itemErrors: Record<string, string> = {}
  
  if (!item.categoryId) {
    itemErrors.categoryId = 'Category is required'
  }
  
  if (!item.plannedAmount || item.plannedAmount <= 0) {
    itemErrors.plannedAmount = 'Amount must be greater than 0'
  }
  
  if (Object.keys(itemErrors).length > 0) {
    itemFieldErrors.value[index] = itemErrors
    return false
  } else {
    delete itemFieldErrors.value[index]
    return true
  }
}

// Validate all budget items
const validateAllItems = (): boolean => {
  let allValid = true
  
  editForm.value.items.forEach((item, index) => {
    if (!validateBudgetItem(index, item)) {
      allValid = false
    }
  })
  
  // Check for duplicate categories
  const categoryIds = editForm.value.items.map(item => item.categoryId).filter(Boolean)
  const duplicates = categoryIds.filter((id, index) => categoryIds.indexOf(id) !== index)
  
  if (duplicates.length > 0) {
    editForm.value.items.forEach((item, index) => {
      if (duplicates.includes(item.categoryId)) {
        if (!itemFieldErrors.value[index]) itemFieldErrors.value[index] = {}
        itemFieldErrors.value[index].categoryId = 'Category already used'
      }
    })
    allValid = false
  }
  
  return allValid
}

// Budget item management
const addItem = () => {
  editForm.value.items.push({
    id: '', // Will be generated by backend
    categoryId: '',
    categoryName: '',
    plannedAmount: 0,
    spentAmount: 0,
    notes: ''
  } as BudgetItem)
}

const removeItem = (index: number) => {
  editForm.value.items.splice(index, 1)
  // Clear any errors for this item
  if (itemFieldErrors.value[index]) {
    delete itemFieldErrors.value[index]
  }
}

const updateCategoryName = (index: number) => {
  const item = editForm.value.items[index]
  const cat = props.categories.find(c => c.id === item.categoryId)
  if (cat) {
    item.categoryName = cat.name
  }
  
  // Clear category error when valid category is selected
  if (itemFieldErrors.value[index]?.categoryId) {
    delete itemFieldErrors.value[index].categoryId
  }
  
  // Validate the item
  validateBudgetItem(index, item)
}

const handleSave = async () => {
  // Validate form and items
  const formValidation = validateForm()
  const itemsValid = validateAllItems()
  
  if (!handleValidationResult(formValidation) || !itemsValid) {
    return
  }

  const updateData = {
    name: editForm.value.name.trim(),
    periodType: editForm.value.periodType,
    startDate: editForm.value.startDate,
    endDate: editForm.value.endDate,
    totalAmount: totalPlannedAmount.value,
    items: editForm.value.items.map(item => ({
      id: item.id || '',
      categoryId: item.categoryId,
      categoryName: item.categoryName,
      plannedAmount: item.plannedAmount,
      spentAmount: item.spentAmount || 0,
      notes: item.notes?.trim() || ''
    }))
  }

  await handleFormSubmission(
    () => financialService.updateBudget(props.budget.id, updateData),
    {
      onSuccess: (updatedBudget) => emit('saved', updatedBudget)
    }
  )
}

// Handle field blur events
const handleFieldBlur = (fieldName: string) => {
  touchField(fieldName)
  validateField(fieldName)
}

const handleItemFieldBlur = (index: number, _fieldName: string) => {
  const item = editForm.value.items[index]
  validateBudgetItem(index, item)
}

// Clear field errors when user starts typing (handled by real-time validation)
watch(() => editForm.value.name, () => clearFieldError('name'))
watch(() => editForm.value.startDate, () => clearFieldError('startDate'))
watch(() => editForm.value.endDate, () => clearFieldError('endDate'))

// Clear item field errors when user makes changes
editForm.value.items.forEach((_, index) => {
  watch(() => editForm.value.items[index]?.plannedAmount, () => {
    if (itemFieldErrors.value[index]?.plannedAmount) {
      delete itemFieldErrors.value[index].plannedAmount
      validateBudgetItem(index, editForm.value.items[index])
    }
  })
  
  watch(() => editForm.value.items[index]?.categoryId, () => {
    if (itemFieldErrors.value[index]?.categoryId) {
      delete itemFieldErrors.value[index].categoryId
      validateBudgetItem(index, editForm.value.items[index])
    }
  })
})
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #e5e7eb;
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #d1d5db;
}
</style>