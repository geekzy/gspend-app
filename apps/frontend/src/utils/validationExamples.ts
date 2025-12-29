/**
 * Examples and tests for the validation system
 * This file demonstrates how to use the validation utilities
 */

import {
  ValidationRules,
  FieldValidator,
  FormValidators
} from './validation'

// Example: Testing positive amount validation
export function testAmountValidation() {
  const amountValidator = new FieldValidator([
    ValidationRules.required('Amount is required'),
    ValidationRules.positiveAmount('Amount must be greater than 0')
  ])

  // Test cases
  const testCases = [
    { value: '', expected: false, description: 'Empty string' },
    { value: 0, expected: false, description: 'Zero' },
    { value: -10, expected: false, description: 'Negative number' },
    { value: 10.50, expected: true, description: 'Positive decimal' },
    { value: '25.99', expected: true, description: 'String number' },
    { value: 'abc', expected: false, description: 'Non-numeric string' }
  ]

  console.log('Amount Validation Tests:')
  testCases.forEach(testCase => {
    const result = amountValidator.validate(testCase.value)
    const passed = result.isValid === testCase.expected
    console.log(`${passed ? '✅' : '❌'} ${testCase.description}: ${testCase.value} -> ${result.isValid ? 'Valid' : result.error}`)
  })
}

// Example: Testing required field validation
export function testRequiredFieldValidation() {
  const requiredValidator = new FieldValidator([
    ValidationRules.required('This field is required')
  ])

  const testCases = [
    { value: '', expected: false, description: 'Empty string' },
    { value: '   ', expected: false, description: 'Whitespace only' },
    { value: 'Hello', expected: true, description: 'Valid text' },
    { value: null, expected: false, description: 'Null value' },
    { value: undefined, expected: false, description: 'Undefined value' },
    { value: 0, expected: true, description: 'Zero (valid for numbers)' }
  ]

  console.log('\nRequired Field Validation Tests:')
  testCases.forEach(testCase => {
    const result = requiredValidator.validate(testCase.value)
    const passed = result.isValid === testCase.expected
    console.log(`${passed ? '✅' : '❌'} ${testCase.description}: "${testCase.value}" -> ${result.isValid ? 'Valid' : result.error}`)
  })
}

// Example: Testing cross-field validation (password confirmation)
export function testPasswordConfirmValidation() {
  const passwordConfirmValidator = new FieldValidator([
    ValidationRules.required('Please confirm your password'),
    ValidationRules.passwordConfirm('password', 'Passwords do not match')
  ])

  const testCases = [
    {
      confirmPassword: '',
      formData: { password: 'test123' },
      expected: false,
      description: 'Empty confirmation'
    },
    {
      confirmPassword: 'test123',
      formData: { password: 'test123' },
      expected: true,
      description: 'Matching passwords'
    },
    {
      confirmPassword: 'different',
      formData: { password: 'test123' },
      expected: false,
      description: 'Non-matching passwords'
    }
  ]

  console.log('\nPassword Confirmation Validation Tests:')
  testCases.forEach(testCase => {
    const result = passwordConfirmValidator.validate(testCase.confirmPassword, testCase.formData)
    const passed = result.isValid === testCase.expected
    console.log(`${passed ? '✅' : '❌'} ${testCase.description} -> ${result.isValid ? 'Valid' : result.error}`)
  })
}

// Example: Testing complete form validation
export function testTransactionFormValidation() {
  const transactionValidator = FormValidators.transaction()

  const testCases = [
    {
      formData: {
        description: 'Grocery shopping',
        amount: 45.99,
        categoryId: 'groceries',
        transactionDate: '2024-01-15'
      },
      expected: true,
      description: 'Valid transaction'
    },
    {
      formData: {
        description: '',
        amount: 45.99,
        categoryId: 'groceries',
        transactionDate: '2024-01-15'
      },
      expected: false,
      description: 'Missing description'
    },
    {
      formData: {
        description: 'Grocery shopping',
        amount: -10,
        categoryId: 'groceries',
        transactionDate: '2024-01-15'
      },
      expected: false,
      description: 'Negative amount'
    },
    {
      formData: {
        description: 'Grocery shopping',
        amount: 45.99,
        categoryId: '',
        transactionDate: '2024-01-15'
      },
      expected: false,
      description: 'Missing category'
    }
  ]

  console.log('\nTransaction Form Validation Tests:')
  testCases.forEach(testCase => {
    const result = transactionValidator.validateForm(testCase.formData)
    const passed = result.isValid === testCase.expected
    console.log(`${passed ? '✅' : '❌'} ${testCase.description} -> ${result.isValid ? 'Valid' : Object.keys(result.errors).length + ' errors'}`)
    if (!result.isValid) {
      Object.entries(result.errors).forEach(([field, error]) => {
        console.log(`    - ${field}: ${error}`)
      })
    }
  })
}

// Example: Testing password strength validation
export function testPasswordStrengthValidation() {
  const passwordValidator = new FieldValidator([
    ValidationRules.required('Password is required'),
    ValidationRules.passwordStrength()
  ])

  const testCases = [
    { value: '', expected: false, description: 'Empty password' },
    { value: 'weak', expected: false, description: 'Too short' },
    { value: 'password123', expected: false, description: 'No uppercase' },
    { value: 'PASSWORD123', expected: false, description: 'No lowercase' },
    { value: 'Password', expected: false, description: 'No number' },
    { value: 'Password123', expected: true, description: 'Strong password' },
    { value: 'MySecure123!', expected: true, description: 'Very strong password' }
  ]

  console.log('\nPassword Strength Validation Tests:')
  testCases.forEach(testCase => {
    const result = passwordValidator.validate(testCase.value)
    const passed = result.isValid === testCase.expected
    console.log(`${passed ? '✅' : '❌'} ${testCase.description}: "${testCase.value}" -> ${result.isValid ? 'Valid' : result.error}`)
  })
}

// Run all validation tests
export function runAllValidationTests() {
  console.log('🧪 Running Validation System Tests\n')

  testAmountValidation()
  testRequiredFieldValidation()
  testPasswordConfirmValidation()
  testTransactionFormValidation()
  testPasswordStrengthValidation()

  console.log('\n✨ Validation tests completed!')
}

// Example usage in development
if (import.meta.env.DEV) {
  // Uncomment to run tests in development
  // runAllValidationTests()
}