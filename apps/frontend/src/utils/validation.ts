/**
 * Comprehensive validation utilities for gSpend application
 * Provides real-time validation feedback and custom validation messages
 */

export interface ValidationRule {
  validate: (value: any, formData?: any) => boolean
  message: string
}

export interface ValidationResult {
  isValid: boolean
  errors: Record<string, string>
}

export interface FieldValidationResult {
  isValid: boolean
  error?: string
}

/**
 * Common validation rules
 */
export const ValidationRules = {
  // Required field validation
  required: (message = 'This field is required'): ValidationRule => ({
    validate: (value: any) => {
      if (typeof value === 'string') {
        return value.trim().length > 0
      }
      return value !== null && value !== undefined && value !== ''
    },
    message
  }),

  // Positive number validation for monetary amounts
  positiveAmount: (message = 'Amount must be greater than 0'): ValidationRule => ({
    validate: (value: any) => {
      const num = Number(value)
      return !isNaN(num) && num > 0
    },
    message
  }),

  // Minimum length validation
  minLength: (min: number, message?: string): ValidationRule => ({
    validate: (value: any) => {
      const str = String(value || '').trim()
      return str.length >= min
    },
    message: message || `Must be at least ${min} characters`
  }),

  // Maximum length validation
  maxLength: (max: number, message?: string): ValidationRule => ({
    validate: (value: any) => {
      const str = String(value || '').trim()
      return str.length <= max
    },
    message: message || `Must be no more than ${max} characters`
  }),

  // Email validation
  email: (message = 'Please enter a valid email address'): ValidationRule => ({
    validate: (value: any) => {
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      return emailRegex.test(String(value || '').trim())
    },
    message
  }),

  // Password strength validation
  passwordStrength: (message = 'Password must be at least 8 characters with uppercase, lowercase, and number'): ValidationRule => ({
    validate: (value: any) => {
      const password = String(value || '')
      return (
        password.length >= 8 &&
        /[A-Z]/.test(password) &&
        /[a-z]/.test(password) &&
        /\d/.test(password)
      )
    },
    message
  }),

  // Date validation (not in future for transactions)
  notFutureDate: (message = 'Date cannot be in the future'): ValidationRule => ({
    validate: (value: any) => {
      if (!value) return true // Let required rule handle empty dates
      const selectedDate = new Date(value)
      const today = new Date()
      today.setHours(23, 59, 59, 999)
      return selectedDate <= today
    },
    message
  }),

  // Date range validation (end date after start date)
  dateRange: (startDateField: string, message = 'End date must be after start date'): ValidationRule => ({
    validate: (value: any, formData: any) => {
      if (!value || !formData?.[startDateField]) return true
      const startDate = new Date(formData[startDateField])
      const endDate = new Date(value)
      return endDate > startDate
    },
    message
  }),

  // Password confirmation validation
  passwordConfirm: (passwordField: string, message = 'Passwords do not match'): ValidationRule => ({
    validate: (value: any, formData: any) => {
      return value === formData?.[passwordField]
    },
    message
  }),

  // Family size validation (0-5 children)
  familySize: (message = 'Number of children must be between 0 and 5'): ValidationRule => ({
    validate: (value: any) => {
      const num = Number(value)
      return !isNaN(num) && num >= 0 && num <= 5 && Number.isInteger(num)
    },
    message
  }),

  // Category selection validation
  categoryRequired: (message = 'Please select a category'): ValidationRule => ({
    validate: (value: any) => {
      return value && value.trim().length > 0
    },
    message
  }),

  // Hex color validation
  hexColor: (message = 'Please select a valid color'): ValidationRule => ({
    validate: (value: any) => {
      return /^#[0-9A-F]{6}$/i.test(String(value || ''))
    },
    message
  }),

  // Frequency validation
  frequency: (validFrequencies: string[], message = 'Please select a valid frequency'): ValidationRule => ({
    validate: (value: any) => {
      return validFrequencies.includes(String(value || ''))
    },
    message
  }),

  // Numeric validation
  numeric: (message = 'Must be a valid number'): ValidationRule => ({
    validate: (value: any) => {
      return !isNaN(Number(value)) && isFinite(Number(value))
    },
    message
  }),

  // Custom validation rule
  custom: (validator: (value: any, formData?: any) => boolean, message: string): ValidationRule => ({
    validate: validator,
    message
  })
}

/**
 * Field validator class for real-time validation
 */
export class FieldValidator {
  private rules: ValidationRule[] = []

  constructor(rules: ValidationRule[] = []) {
    this.rules = rules
  }

  addRule(rule: ValidationRule): FieldValidator {
    this.rules.push(rule)
    return this
  }

  validate(value: any, formData?: any): FieldValidationResult {
    for (const rule of this.rules) {
      if (!rule.validate(value, formData)) {
        return {
          isValid: false,
          error: rule.message
        }
      }
    }
    return { isValid: true }
  }
}

/**
 * Form validator class for comprehensive form validation
 */
export class FormValidator {
  private fieldValidators: Record<string, FieldValidator> = {}

  addField(fieldName: string, validator: FieldValidator): FormValidator {
    this.fieldValidators[fieldName] = validator
    return this
  }

  validateField(fieldName: string, value: any, formData?: any): FieldValidationResult {
    const validator = this.fieldValidators[fieldName]
    if (!validator) {
      return { isValid: true }
    }
    return validator.validate(value, formData)
  }

  validateForm(formData: Record<string, any>): ValidationResult {
    const errors: Record<string, string> = {}
    let isValid = true

    for (const [fieldName, validator] of Object.entries(this.fieldValidators)) {
      const result = validator.validate(formData[fieldName], formData)
      if (!result.isValid) {
        errors[fieldName] = result.error!
        isValid = false
      }
    }

    return { isValid, errors }
  }
}

/**
 * Pre-configured validators for common form types
 */
export const FormValidators = {
  // Transaction form validator
  transaction: () => new FormValidator()
    .addField('description', new FieldValidator([
      ValidationRules.required('Description is required'),
      ValidationRules.minLength(2, 'Description must be at least 2 characters'),
      ValidationRules.maxLength(200, 'Description must be less than 200 characters')
    ]))
    .addField('amount', new FieldValidator([
      ValidationRules.required('Amount is required'),
      ValidationRules.positiveAmount('Amount must be greater than 0')
    ]))
    .addField('categoryId', new FieldValidator([
      ValidationRules.categoryRequired('Please select a category')
    ]))
    .addField('transactionDate', new FieldValidator([
      ValidationRules.required('Transaction date is required'),
      ValidationRules.notFutureDate('Transaction date cannot be in the future')
    ])),

  // Budget form validator
  budget: () => new FormValidator()
    .addField('name', new FieldValidator([
      ValidationRules.required('Budget name is required'),
      ValidationRules.minLength(2, 'Budget name must be at least 2 characters'),
      ValidationRules.maxLength(100, 'Budget name must be less than 100 characters')
    ]))
    .addField('startDate', new FieldValidator([
      ValidationRules.required('Start date is required')
    ]))
    .addField('endDate', new FieldValidator([
      ValidationRules.required('End date is required'),
      ValidationRules.dateRange('startDate', 'End date must be after start date')
    ])),

  // Income form validator
  income: () => new FormValidator()
    .addField('source', new FieldValidator([
      ValidationRules.required('Source name is required'),
      ValidationRules.minLength(2, 'Source name must be at least 2 characters'),
      ValidationRules.maxLength(100, 'Source name must be less than 100 characters')
    ]))
    .addField('amount', new FieldValidator([
      ValidationRules.required('Amount is required'),
      ValidationRules.positiveAmount('Amount must be greater than 0')
    ]))
    .addField('frequency', new FieldValidator([
      ValidationRules.frequency(['one-time', 'weekly', 'bi-weekly', 'monthly', 'yearly'], 'Please select a valid frequency')
    ]))
    .addField('effectiveDate', new FieldValidator([
      ValidationRules.required('Effective date is required')
    ])),

  // Category form validator
  category: () => new FormValidator()
    .addField('name', new FieldValidator([
      ValidationRules.required('Category name is required'),
      ValidationRules.minLength(2, 'Category name must be at least 2 characters'),
      ValidationRules.maxLength(50, 'Category name must be less than 50 characters')
    ]))
    .addField('type', new FieldValidator([
      ValidationRules.required('Category type is required')
    ]))
    .addField('icon', new FieldValidator([
      ValidationRules.required('Please select an icon')
    ]))
    .addField('color', new FieldValidator([
      ValidationRules.required('Please select a color'),
      ValidationRules.hexColor('Please select a valid color')
    ])),

  // Profile form validator
  profile: () => new FormValidator()
    .addField('familyName', new FieldValidator([
      ValidationRules.required('Family name is required'),
      ValidationRules.minLength(2, 'Family name must be at least 2 characters'),
      ValidationRules.maxLength(100, 'Family name must be less than 100 characters')
    ]))
    .addField('email', new FieldValidator([
      ValidationRules.required('Email is required'),
      ValidationRules.email('Please enter a valid email address')
    ]))
    .addField('numberOfChildren', new FieldValidator([
      ValidationRules.required('Number of children is required'),
      ValidationRules.familySize('Number of children must be between 0 and 5')
    ])),

  // Password change form validator
  passwordChange: () => new FormValidator()
    .addField('currentPassword', new FieldValidator([
      ValidationRules.required('Current password is required')
    ]))
    .addField('newPassword', new FieldValidator([
      ValidationRules.required('New password is required'),
      ValidationRules.passwordStrength('Password must be at least 8 characters with uppercase, lowercase, and number')
    ]))
    .addField('confirmPassword', new FieldValidator([
      ValidationRules.required('Please confirm your password'),
      ValidationRules.passwordConfirm('newPassword', 'Passwords do not match')
    ]))
}

/**
 * Utility functions for validation
 */
export const ValidationUtils = {
  // Debounce validation for real-time feedback
  debounceValidation: (fn: Function, delay = 300) => {
    let timeoutId: number
    return (...args: any[]) => {
      clearTimeout(timeoutId)
      timeoutId = setTimeout(() => fn.apply(null, args), delay)
    }
  },

  // Format validation errors for display
  formatErrors: (errors: Record<string, string>): string[] => {
    return Object.values(errors).filter(Boolean)
  },

  // Check if form has any errors
  hasErrors: (errors: Record<string, string>): boolean => {
    return Object.values(errors).some(error => error && error.length > 0)
  },

  // Clear specific field error
  clearFieldError: (errors: Record<string, string>, fieldName: string): Record<string, string> => {
    const newErrors = { ...errors }
    delete newErrors[fieldName]
    return newErrors
  },

  // Validate single field with debouncing
  validateFieldDebounced: (
    validator: FieldValidator,
    value: any,
    formData: any,
    onResult: (result: FieldValidationResult) => void
  ) => {
    const debouncedValidate = ValidationUtils.debounceValidation((val: any, data: any) => {
      const result = validator.validate(val, data)
      onResult(result)
    }, 300)

    debouncedValidate(value, formData)
  }
}