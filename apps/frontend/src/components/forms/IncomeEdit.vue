<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm overflow-y-auto">
    <div class="bg-white rounded-3xl w-full max-w-lg p-8 shadow-2xl animate-in fade-in zoom-in duration-200">
      <div class="flex justify-between items-center mb-6">
        <h3 class="text-xl font-black text-gray-900">Edit Income Source</h3>
        <button @click="$emit('close')" class="p-2 text-gray-400 hover:text-gray-600">
          <XIcon class="w-6 h-6" />
        </button>
      </div>

      <!-- Error Display -->
      <ValidationSummary 
        v-if="hasErrors"
        :errors="errors"
        class="mb-4"
      />

      <form @submit.prevent="handleSave" class="space-y-4">
        <!-- Source Name -->
        <div>
          <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">
            Source Name <span class="text-red-500">*</span>
          </label>
          <input 
            v-model="editForm.source" 
            type="text" 
            :class="getFieldClass('source', 'w-full px-4 py-3 border rounded-2xl focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none', 'border-red-300 bg-red-50')"
            placeholder="e.g. Monthly Salary, Freelance Work" 
            required 
            @blur="handleFieldBlur('source')"
          />
          <ErrorMessage :error="getFieldError('source')" />
        </div>

        <!-- Amount -->
        <div>
          <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">
            Amount <span class="text-red-500">*</span>
          </label>
          <div class="relative">
            <span class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400">$</span>
            <input 
              v-model.number="editForm.amount" 
              type="number" 
              step="0.01" 
              min="0.01"
              :class="getFieldClass('amount', 'w-full pl-8 pr-4 py-3 border rounded-2xl focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none', 'border-red-300 bg-red-50')"
              placeholder="0.00" 
              required 
              @blur="handleFieldBlur('amount')"
            />
          </div>
          <ErrorMessage :error="getFieldError('amount')" />
        </div>

        <!-- Frequency and Effective Date -->
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">
              Frequency <span class="text-red-500">*</span>
            </label>
            <select 
              v-model="editForm.frequency" 
              :class="getFieldClass('frequency', 'w-full px-4 py-3 border rounded-2xl focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none appearance-none bg-white', 'border-red-300 bg-red-50')"
              required
              @blur="handleFieldBlur('frequency')"
            >
              <option value="one-time">One-time</option>
              <option value="weekly">Weekly</option>
              <option value="bi-weekly">Bi-weekly</option>
              <option value="monthly">Monthly</option>
              <option value="yearly">Yearly</option>
            </select>
            <ErrorMessage :error="getFieldError('frequency')" />
          </div>
          <div>
            <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">
              Effective Date <span class="text-red-500">*</span>
            </label>
            <input 
              v-model="editForm.effectiveDate" 
              type="date" 
              :class="getFieldClass('effectiveDate', 'w-full px-4 py-3 border rounded-2xl focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none', 'border-red-300 bg-red-50')"
              required 
              @blur="handleFieldBlur('effectiveDate')"
            />
            <ErrorMessage :error="getFieldError('effectiveDate')" />
          </div>
        </div>

        <!-- Frequency Information -->
        <div class="bg-blue-50 p-4 rounded-2xl">
          <div class="flex items-start">
            <InfoIcon class="w-5 h-5 text-blue-500 mr-2 mt-0.5 flex-shrink-0" />
            <div class="text-sm text-blue-700">
              <p class="font-bold mb-1">Frequency Guide:</p>
              <ul class="space-y-1 text-xs">
                <li><strong>One-time:</strong> Single payment or bonus</li>
                <li><strong>Weekly:</strong> Every 7 days (52 times/year)</li>
                <li><strong>Bi-weekly:</strong> Every 2 weeks (26 times/year)</li>
                <li><strong>Monthly:</strong> Once per month (12 times/year)</li>
                <li><strong>Yearly:</strong> Annual salary or bonus</li>
              </ul>
            </div>
          </div>
        </div>

        <!-- Monthly Equivalent Display -->
        <div v-if="monthlyEquivalent > 0" class="bg-green-50 p-4 rounded-2xl">
          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <CalculatorIcon class="w-5 h-5 text-green-500 mr-2" />
              <span class="text-sm font-bold text-green-700">Monthly Equivalent:</span>
            </div>
            <span class="text-lg font-black text-green-700">
              {{ formatCurrency(monthlyEquivalent, { decimals: 2 }) }}
            </span>
          </div>
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
import { XIcon, InfoIcon, CalculatorIcon } from 'lucide-vue-next'
import { financialService, type Income } from '@/services/financialService'
import { useFormValidation } from '@/composables/useFormValidation'
import { useFormNotifications } from '@/composables/useFormNotifications'
import { FormValidators } from '@/utils/validation'
import ValidationSummary from '@/components/common/ValidationSummary.vue'
import ErrorMessage from '@/components/common/ErrorMessage.vue'
import { formatCurrency } from '@/utils/currency'

interface Props {
  income: Income
}

interface Emits {
  (e: 'close'): void
  (e: 'saved', income: Income): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const isSubmitting = ref(false)

// Form data
const editForm = ref({
  source: props.income.source,
  amount: props.income.amount,
  frequency: props.income.frequency,
  effectiveDate: props.income.effectiveDate.split('T')[0] // Convert to YYYY-MM-DD format
})

// Set up form validation
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
  validator: FormValidators.income(),
  formData: editForm,
  realTimeValidation: true
})

// Set up form notifications
const { handleFormSubmission, handleValidationResult } = useFormNotifications({
  successMessage: 'Income source updated successfully',
  errorMessage: 'Failed to update income source'
})

// Form validation - using the new validation system
const isFormValid = computed(() => {
  return isValid.value && !!(
    editForm.value.source.trim() &&
    editForm.value.amount > 0 &&
    editForm.value.frequency &&
    editForm.value.effectiveDate
  )
})

// Calculate monthly equivalent for display
const monthlyEquivalent = computed(() => {
  const amount = editForm.value.amount || 0
  
  switch (editForm.value.frequency) {
    case 'weekly':
      return (amount * 52) / 12 // 52 weeks per year / 12 months
    case 'bi-weekly':
      return (amount * 26) / 12 // 26 bi-weekly periods per year / 12 months
    case 'monthly':
      return amount
    case 'yearly':
      return amount / 12
    case 'one-time':
      return 0 // One-time payments don't have a monthly equivalent
    default:
      return 0
  }
})

const handleSave = async () => {
  const validationResult = validateForm()
  if (!handleValidationResult(validationResult)) {
    return
  }

  const updateData = {
    source: editForm.value.source.trim(),
    amount: editForm.value.amount,
    frequency: editForm.value.frequency,
    effectiveDate: editForm.value.effectiveDate
  }

  await handleFormSubmission(
    () => financialService.updateIncome(props.income.id, updateData),
    {
      onSuccess: (updatedIncome) => emit('saved', updatedIncome)
    }
  )
}

// Handle field blur events for validation
const handleFieldBlur = (fieldName: string) => {
  touchField(fieldName)
  validateField(fieldName)
}

// Clear field errors when user starts typing (handled by real-time validation)
watch(() => editForm.value.source, () => clearFieldError('source'))
watch(() => editForm.value.amount, () => clearFieldError('amount'))
watch(() => editForm.value.frequency, () => clearFieldError('frequency'))
watch(() => editForm.value.effectiveDate, () => clearFieldError('effectiveDate'))
</script>