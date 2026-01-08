<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm overflow-y-auto">
    <div class="bg-white rounded-3xl w-full max-w-lg p-8 shadow-2xl animate-in fade-in zoom-in duration-200">
      <div class="flex justify-between items-center mb-6">
        <h3 class="text-xl font-black text-gray-900">Edit Transaction</h3>
        <button @click="$emit('close')" class="p-2 text-gray-400 hover:text-gray-600" aria-label="Close modal">
          <XIcon class="w-6 h-6" />
        </button>
      </div>

      <!-- Error Display -->
      <div v-if="hasErrors" class="mb-4 p-4 bg-red-50 border border-red-200 rounded-2xl">
        <div class="flex items-center mb-2">
          <AlertCircleIcon class="w-5 h-5 text-red-500 mr-2" />
          <span class="text-sm font-bold text-red-700">Please fix the following errors:</span>
        </div>
        <ul class="text-sm text-red-600 space-y-1">
          <li v-for="error in errorList" :key="error">• {{ error }}</li>
        </ul>
      </div>

      <form @submit.prevent="handleSave" class="space-y-4">
        <!-- Transaction Type -->
        <div class="grid grid-cols-2 gap-3">
          <button 
            type="button"
            @click="editForm.type = 'expense'"
            :class="[
              editForm.type === 'expense' ? 'bg-red-500 text-white shadow-lg shadow-red-200' : 'bg-gray-50 text-gray-500',
              'py-3 rounded-2xl text-sm font-bold transition-all'
            ]"
            role="radio"
            :aria-checked="editForm.type === 'expense'"
          >
            Expense
          </button>
          <button 
            type="button"
            @click="editForm.type = 'income'"
            :class="[
              editForm.type === 'income' ? 'bg-green-500 text-white shadow-lg shadow-green-200' : 'bg-gray-50 text-gray-500',
              'py-3 rounded-2xl text-sm font-bold transition-all'
            ]"
            role="radio"
            :aria-checked="editForm.type === 'income'"
          >
            Income
          </button>
        </div>

        <!-- Description -->
        <div>
          <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">
            Description <span class="text-red-500">*</span>
          </label>
          <input 
            v-model="editForm.description" 
            type="text" 
            :class="getFieldClass('description', 'w-full px-4 py-3 border rounded-2xl focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none', 'border-red-300 bg-red-50')"
            placeholder="What was this for?" 
            required 
            @blur="handleFieldBlur('description')"
            :aria-invalid="hasFieldError('description')"
            :aria-describedby="hasFieldError('description') ? 'error-description' : undefined"
          />
          <p v-if="hasFieldError('description')" id="error-description" class="mt-1 text-xs text-red-500">{{ getFieldError('description') }}</p>
        </div>

        <!-- Amount and Date -->
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">
              Amount <span class="text-red-500">*</span>
            </label>
            <div class="relative">
              <span class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400">{{ currencySymbol }}</span>
              <input 
                v-model.number="editForm.amount" 
                type="number" 
                step="0.01" 
                min="0.01"
                :class="getFieldClass('amount', 'w-full pl-12 pr-4 py-3 border rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none', 'border-red-300 bg-red-50')"
                placeholder="0.00" 
                required 
                @blur="handleFieldBlur('amount')"
                :aria-invalid="hasFieldError('amount')"
                :aria-describedby="hasFieldError('amount') ? 'error-amount' : undefined"
              />
            </div>
            <p v-if="hasFieldError('amount')" id="error-amount" class="mt-1 text-xs text-red-500">{{ getFieldError('amount') }}</p>
          </div>
          <div class="grid grid-cols-2 gap-2">
            <div>
              <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">
                Date <span class="text-red-500">*</span>
              </label>
              <input 
                v-model="editForm.transactionDate" 
                type="date" 
                :class="getFieldClass('transactionDate', 'w-full px-3 py-3 border rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none text-sm', 'border-red-300 bg-red-50')"
                required 
                @blur="handleFieldBlur('transactionDate')"
                :aria-invalid="hasFieldError('transactionDate')"
                :aria-describedby="hasFieldError('transactionDate') ? 'error-date' : undefined"
              />
              <p v-if="hasFieldError('transactionDate')" id="error-date" class="mt-1 text-xs text-red-500">{{ getFieldError('transactionDate') }}</p>
            </div>
            <div>
              <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">
                Time <span class="text-red-500">*</span>
              </label>
              <input 
                v-model="editForm.transactionTime" 
                type="time" 
                class="w-full px-3 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none text-sm"
                required 
              />
            </div>
          </div>
        </div>

        <!-- Category and Payment Method -->
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">
              Category <span class="text-red-500">*</span>
            </label>
            <select 
              v-model="editForm.categoryId" 
              :class="getFieldClass('categoryId', 'w-full px-4 py-3 border rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none bg-white', 'border-red-300 bg-red-50')"
              required
              @blur="handleFieldBlur('categoryId')"
              :aria-invalid="hasFieldError('categoryId')"
              :aria-describedby="hasFieldError('categoryId') ? 'error-category' : undefined"
            >
              <option value="">Select Category</option>
              <option v-for="cat in availableCategories" :key="cat.id" :value="cat.id">
                {{ cat.icon }} {{ cat.name }}
              </option>
            </select>
            <p v-if="hasFieldError('categoryId')" id="error-category" class="mt-1 text-xs text-red-500">{{ getFieldError('categoryId') }}</p>
          </div>
          <div>
            <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Payment Method</label>
            <select 
              v-model="editForm.paymentMethod" 
              class="w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none bg-white"
            >
              <option value="cash">Cash</option>
              <option value="bank_transfer">Bank Transfer</option>
              <option value="credit_card">Credit Card</option>
              <option value="debit_card">Debit Card</option>
            </select>
          </div>
        </div>

        <!-- Notes -->
        <div>
          <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Notes (Optional)</label>
          <textarea 
            v-model="editForm.notes" 
            rows="2" 
            class="w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none resize-none" 
            placeholder="Any extra details..."
          ></textarea>
        </div>

        <!-- Action Buttons -->
        <div class="flex gap-3 mt-8">
          <button 
            @click="$emit('close')" 
            type="button" 
            class="flex-1 py-3.5 rounded-2xl border border-gray-200 text-sm font-bold text-gray-700 hover:bg-gray-50 transition-all"
            :disabled="isSubmitting"
          >
            Cancel
          </button>
          <button 
            type="submit" 
            :disabled="isSubmitting || !isFormValid" 
            class="flex-1 py-3.5 rounded-2xl bg-primary-600 text-white text-sm font-black hover:bg-primary-700 shadow-lg shadow-primary-200 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ isSubmitting ? 'Saving...' : 'Save Changes' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { XIcon, AlertCircleIcon } from 'lucide-vue-next'
import { financialService, type Transaction, type Category } from '@/services/financialService'
import { useFormValidation } from '@/composables/useFormValidation'
import { FormValidators } from '@/utils/validation'
import { useNotificationStore } from '@/stores/notification'
import { currencySymbol } from '@/utils/currency'

interface Props {
  transaction: Transaction
  categories: Category[]
}

interface Emits {
  (e: 'close'): void
  (e: 'saved', transaction: Transaction): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const notificationStore = useNotificationStore()

const isSubmitting = ref(false)

// Form data
const editForm = ref({
  type: props.transaction.type,
  description: props.transaction.description,
  amount: props.transaction.amount,
  transactionDate: props.transaction.transactionDate.split('T')[0], // Convert to YYYY-MM-DD format
  transactionTime: props.transaction.transactionDate.includes('T') 
    ? props.transaction.transactionDate.split('T')[1]?.slice(0, 5) || '12:00'
    : '12:00',
  categoryId: props.transaction.categoryId,
  paymentMethod: props.transaction.paymentMethod,
  notes: props.transaction.notes || ''
})

// Set up form validation
const {
  errors,
  errorList,
  isValid,
  hasErrors,
  validateForm,
  validateField,
  clearFieldError,
  touchField,
  getFieldError,
  hasFieldError,
  getFieldClass
} = useFormValidation({
  validator: FormValidators.transaction(),
  formData: editForm,
  realTimeValidation: true
})

// Available categories based on transaction type
const availableCategories = computed(() => {
  return props.categories.filter(c => c.type === editForm.value.type)
})

// Form validation - using the new validation system
const isFormValid = computed(() => {
  return isValid.value && !!(
    editForm.value.description.trim() &&
    editForm.value.amount > 0 &&
    editForm.value.transactionDate &&
    editForm.value.categoryId
  )
})

// Watch for type changes to reset category if needed
watch(() => editForm.value.type, (newType) => {
  const currentCategory = props.categories.find(c => c.id === editForm.value.categoryId)
  if (currentCategory && currentCategory.type !== newType) {
    editForm.value.categoryId = ''
    clearFieldError('categoryId')
  }
})

const handleSave = async () => {
  const validationResult = validateForm()
  if (!validationResult.isValid) {
    notificationStore.error('Please fix the validation errors before saving')
    return
  }

  try {
    isSubmitting.value = true
    
    // Prepare the update data - combine date and time
    const combinedDateTime = `${editForm.value.transactionDate}T${editForm.value.transactionTime}:00`
    const updateData = {
      type: editForm.value.type,
      description: editForm.value.description.trim(),
      amount: editForm.value.amount,
      transactionDate: combinedDateTime,
      categoryId: editForm.value.categoryId,
      paymentMethod: editForm.value.paymentMethod,
      notes: editForm.value.notes.trim()
    }

    const updatedTransaction = await financialService.updateTransaction(props.transaction.id, updateData)
    notificationStore.success('Transaction updated successfully')
    emit('saved', updatedTransaction)
  } catch (error: any) {
    console.error('Failed to update transaction:', error)
    
    // Handle validation errors from server
    if (error.response?.status === 400) {
      const errorData = error.response.data
      if (errorData.details) {
        // Handle field-specific errors from server
        errorData.details.forEach((detail: any) => {
          // Server errors will override client validation
          errors.value[detail.field] = detail.message
        })
        notificationStore.error('Please check your input and try again')
      } else {
        notificationStore.error(errorData.message || 'Please check your input and try again')
      }
    } else {
      notificationStore.error('Failed to save transaction. Please try again.')
    }
  } finally {
    isSubmitting.value = false
  }
}

// Handle field blur events for validation
const handleFieldBlur = (fieldName: string) => {
  touchField(fieldName)
  validateField(fieldName)
}

// Clear field errors when user starts typing (handled by real-time validation)
watch(() => editForm.value.description, () => clearFieldError('description'))
watch(() => editForm.value.amount, () => clearFieldError('amount'))
watch(() => editForm.value.transactionDate, () => clearFieldError('transactionDate'))
watch(() => editForm.value.categoryId, () => clearFieldError('categoryId'))
</script>