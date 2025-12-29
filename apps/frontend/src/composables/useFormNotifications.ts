/**
 * Composable for integrating form validation with toast notifications
 * Provides automatic success/error feedback for form operations
 */

import { useNotificationStore } from '@/stores/notification'
import { useErrorHandling } from './useErrorHandling'
import type { ValidationResult } from '@/utils/validation'

export interface FormNotificationOptions {
  showValidationErrors?: boolean
  showSuccessMessages?: boolean
  successMessage?: string
  errorMessage?: string
}

export function useFormNotifications(options: FormNotificationOptions = {}) {
  const {
    showValidationErrors = true,
    showSuccessMessages = true,
    successMessage = 'Changes saved successfully',
    errorMessage = 'Please fix the validation errors'
  } = options

  const notificationStore = useNotificationStore()
  const { handleOperationError, handleValidationError } = useErrorHandling()

  /**
   * Handle validation result and show appropriate notifications
   */
  const handleValidationResult = (result: ValidationResult): boolean => {
    if (!result.isValid && showValidationErrors) {
      handleValidationError(result.errors)
    }

    return result.isValid
  }

  /**
   * Show success notification for form operations
   */
  const showSuccess = (message?: string) => {
    if (showSuccessMessages) {
      notificationStore.success(message || successMessage)
    }
  }

  /**
   * Show error notification for form operations
   */
  const showError = (message?: string) => {
    notificationStore.error(message || errorMessage)
  }

  /**
   * Show validation-specific error
   */
  const showValidationError = (message?: string) => {
    if (showValidationErrors) {
      notificationStore.validationError(message)
    }
  }

  /**
   * Handle API errors and show appropriate notifications
   */
  const handleApiError = (error: any, defaultMessage = 'An error occurred. Please try again.') => {
    console.error('API Error:', error)

    if (error.response?.status === 400) {
      const errorData = error.response.data
      if (errorData.details && Array.isArray(errorData.details)) {
        // Multiple field errors
        const errors = errorData.details.reduce((acc: any, detail: any) => {
          acc[detail.field] = detail.message
          return acc
        }, {})
        handleValidationError(errors)
      } else if (errorData.message) {
        // Single error message
        showError(errorData.message)
      } else {
        showError(defaultMessage)
      }
    } else if (error.response?.status === 401) {
      showError('You are not authorized to perform this action')
    } else if (error.response?.status === 403) {
      showError('Access denied')
    } else if (error.response?.status === 404) {
      showError('The requested resource was not found')
    } else if (error.response?.status >= 500) {
      notificationStore.serverError()
    } else if (!error.response) {
      notificationStore.networkError()
    } else {
      showError(defaultMessage)
    }
  }

  /**
   * Handle CRUD operations with automatic notifications
   */
  const handleCrudOperation = async <T>(
    operation: 'create' | 'update' | 'delete' | 'save' | 'load',
    resource: string,
    operationFn: () => Promise<T>,
    options: {
      successMessage?: string
      errorMessage?: string
      onSuccess?: (result: T) => void
      onError?: (error: any) => void
    } = {}
  ): Promise<T | null> => {
    try {
      const result = await operationFn()

      if (options.onSuccess) {
        options.onSuccess(result)
      }

      if (options.successMessage) {
        notificationStore.success(options.successMessage)
      } else {
        notificationStore.crudSuccess(operation, resource)
      }

      return result
    } catch (error) {
      if (options.onError) {
        options.onError(error)
      }

      // Use the enhanced error handling
      const resourceMap: Record<string, keyof typeof import('@/utils/errorHandling').operationMessages> = {
        'transaction': 'transaction',
        'budget': 'budget',
        'income': 'income',
        'category': 'category',
        'profile': 'profile'
      }

      const operationType = resourceMap[resource.toLowerCase()]
      if (operationType) {
        handleOperationError(operationType, operation, error as any)
      } else {
        handleApiError(error)
      }

      return null
    }
  }

  /**
   * Wrapper for form submission with automatic error handling
   */
  const handleFormSubmission = async <T>(
    submitFn: () => Promise<T>,
    options: {
      successMessage?: string
      errorMessage?: string
      onSuccess?: (result: T) => void
      onError?: (error: any) => void
    } = {}
  ): Promise<T | null> => {
    try {
      const result = await submitFn()

      if (options.onSuccess) {
        options.onSuccess(result)
      }

      showSuccess(options.successMessage)
      return result
    } catch (error) {
      if (options.onError) {
        options.onError(error)
      }

      handleApiError(error, options.errorMessage)
      return null
    }
  }

  return {
    handleValidationResult,
    showSuccess,
    showError,
    showValidationError,
    handleApiError,
    handleFormSubmission,
    handleCrudOperation
  }
}