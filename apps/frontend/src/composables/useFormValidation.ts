/**
 * Vue composable for form validation with real-time feedback
 * Provides reactive validation state and methods for form components
 */

import { ref, reactive, computed, watch, type Ref } from 'vue'
import {
  FormValidator,
  FieldValidator,
  ValidationUtils,
  type ValidationResult,
  type FieldValidationResult
} from '@/utils/validation'

export interface UseFormValidationOptions {
  validator: FormValidator
  formData: Ref<Record<string, any>>
  realTimeValidation?: boolean
  debounceDelay?: number
}

export interface FormValidationState {
  errors: Record<string, string>
  isValidating: boolean
  isValid: boolean
  hasErrors: boolean
  touchedFields: Set<string>
}

export function useFormValidation(options: UseFormValidationOptions) {
  const { validator, formData, realTimeValidation = true, debounceDelay = 300 } = options

  // Reactive validation state
  const state = reactive<FormValidationState>({
    errors: {},
    isValidating: false,
    isValid: false,
    hasErrors: false,
    touchedFields: new Set()
  })

  // Track field validation timeouts for debouncing
  const validationTimeouts = new Map<string, number>()

  // Computed properties
  const errorList = computed(() => ValidationUtils.formatErrors(state.errors))
  const fieldErrors = computed(() => state.errors)

  /**
   * Validate a single field
   */
  const validateField = (fieldName: string, showError = true): FieldValidationResult => {
    const result = validator.validateField(fieldName, formData.value[fieldName], formData.value)

    if (showError) {
      if (result.isValid) {
        delete state.errors[fieldName]
      } else {
        state.errors[fieldName] = result.error!
      }
      updateValidationState()
    }

    return result
  }

  /**
   * Validate a field with debouncing for real-time validation
   */
  const validateFieldDebounced = (fieldName: string) => {
    // Clear existing timeout
    const existingTimeout = validationTimeouts.get(fieldName)
    if (existingTimeout) {
      clearTimeout(existingTimeout)
    }

    // Set new timeout
    const timeout = setTimeout(() => {
      validateField(fieldName, true)
      validationTimeouts.delete(fieldName)
    }, debounceDelay)

    validationTimeouts.set(fieldName, timeout)
  }

  /**
   * Validate the entire form
   */
  const validateForm = (): ValidationResult => {
    state.isValidating = true

    const result = validator.validateForm(formData.value)
    state.errors = result.errors

    // Mark all fields as touched
    Object.keys(formData.value).forEach(field => {
      state.touchedFields.add(field)
    })

    updateValidationState()
    state.isValidating = false

    return result
  }

  /**
   * Clear validation errors
   */
  const clearErrors = (fieldNames?: string[]) => {
    if (fieldNames) {
      fieldNames.forEach(field => {
        delete state.errors[field]
      })
    } else {
      state.errors = {}
    }
    updateValidationState()
  }

  /**
   * Clear error for a specific field
   */
  const clearFieldError = (fieldName: string) => {
    delete state.errors[fieldName]
    updateValidationState()
  }

  /**
   * Mark a field as touched (for showing validation on blur)
   */
  const touchField = (fieldName: string) => {
    state.touchedFields.add(fieldName)
  }

  /**
   * Check if a field has been touched
   */
  const isFieldTouched = (fieldName: string): boolean => {
    return state.touchedFields.has(fieldName)
  }

  /**
   * Get error for a specific field
   */
  const getFieldError = (fieldName: string): string | undefined => {
    return state.errors[fieldName]
  }

  /**
   * Check if a specific field has an error
   */
  const hasFieldError = (fieldName: string): boolean => {
    return !!state.errors[fieldName]
  }

  /**
   * Update overall validation state
   */
  const updateValidationState = () => {
    state.hasErrors = ValidationUtils.hasErrors(state.errors)
    state.isValid = !state.hasErrors
  }

  /**
   * Reset validation state
   */
  const reset = () => {
    state.errors = {}
    state.touchedFields.clear()
    updateValidationState()

    // Clear all pending timeouts
    validationTimeouts.forEach(timeout => clearTimeout(timeout))
    validationTimeouts.clear()
  }

  /**
   * Set up real-time validation watchers
   */
  if (realTimeValidation) {
    // Watch for changes in form data and validate touched fields
    watch(
      formData,
      (newData, oldData) => {
        if (!oldData) return

        Object.keys(newData).forEach(fieldName => {
          if (newData[fieldName] !== oldData[fieldName] && state.touchedFields.has(fieldName)) {
            validateFieldDebounced(fieldName)
          }
        })
      },
      { deep: true }
    )
  }

  // Initialize validation state
  updateValidationState()

  return {
    // State
    state: state,
    errors: fieldErrors,
    errorList,
    isValid: computed(() => state.isValid),
    hasErrors: computed(() => state.hasErrors),
    isValidating: computed(() => state.isValidating),

    // Methods
    validateField,
    validateForm,
    clearErrors,
    clearFieldError,
    touchField,
    isFieldTouched,
    getFieldError,
    hasFieldError,
    reset,

    // Utilities
    getFieldClass: (fieldName: string, baseClass = '', errorClass = 'border-red-300 bg-red-50') => {
      return hasFieldError(fieldName) ? `${baseClass} ${errorClass}` : baseClass
    }
  }
}

/**
 * Composable for simple field validation without full form context
 */
export function useFieldValidation(fieldValidator: FieldValidator) {
  const error = ref<string>('')
  const isValid = ref(true)
  const isValidating = ref(false)

  const validate = (value: any, formData?: any): FieldValidationResult => {
    isValidating.value = true
    const result = fieldValidator.validate(value, formData)

    error.value = result.error || ''
    isValid.value = result.isValid
    isValidating.value = false

    return result
  }

  const validateDebounced = ValidationUtils.debounceValidation(validate, 300)

  const clear = () => {
    error.value = ''
    isValid.value = true
  }

  return {
    error: error,
    isValid: isValid,
    isValidating: isValidating,
    validate,
    validateDebounced,
    clear
  }
}