export interface UserFriendlyError {
    title: string
    message: string
    suggestion?: string
    action?: {
        label: string
        handler: () => void
        style?: 'primary' | 'secondary'
    }
    type: 'error' | 'warning' | 'info'
}

export interface ApiError {
    response?: {
        status: number
        data?: {
            message?: string
            error?: string
            details?: Array<{
                field: string
                message: string
            }>
        }
    }
    message?: string
    code?: string
}

/**
 * Transform technical errors into user-friendly messages
 */
export function transformError(error: ApiError): UserFriendlyError {
    // Network errors
    if (!error.response) {
        const isOffline = typeof navigator !== 'undefined' && !navigator.onLine

        if (isOffline) {
            return {
                title: 'No Internet Connection',
                message: 'You seem to be offline. Please check your internet connection.',
                suggestion: 'We will try to reconnect automatically when you are back online.',
                action: {
                    label: 'Retry Now',
                    handler: () => window.location.reload()
                },
                type: 'warning'
            }
        }

        return {
            title: 'Server Unavailable',
            message: 'Unable to connect to the server. The service might be down.',
            suggestion: 'Please try again in a few moments.',
            action: {
                label: 'Retry',
                handler: () => window.location.reload()
            },
            type: 'error'
        }
    }

    const status = error.response.status
    const data = error.response.data

    // Handle different HTTP status codes
    switch (status) {
        case 400:
            return handle400Error(data)
        case 401:
            return handle401Error()
        case 403:
            return handle403Error()
        case 404:
            return handle404Error()
        case 409:
            return handle409Error(data)
        case 422:
            return handle422Error(data)
        case 429:
            return handle429Error()
        case 500:
        case 502:
        case 503:
        case 504:
            return handle5xxError(status)
        default:
            return handleUnknownError(error)
    }
}

function handle400Error(data: any): UserFriendlyError {
    if (data?.details && Array.isArray(data.details)) {
        const fieldErrors = data.details.map((d: any) => d.message).join(', ')
        return {
            title: 'Invalid Information',
            message: `Please check the following: ${fieldErrors}`,
            suggestion: 'Review the highlighted fields and correct any errors.',
            type: 'error'
        }
    }

    return {
        title: 'Invalid Request',
        message: data?.message || 'The information provided is not valid.',
        suggestion: 'Please check your input and try again.',
        type: 'error'
    }
}

function handle401Error(): UserFriendlyError {
    return {
        title: 'Session Expired',
        message: 'Your session has expired. Please log in again.',
        suggestion: 'You will be redirected to the login page.',
        action: {
            label: 'Login',
            handler: () => {
                localStorage.removeItem('auth_token')
                window.location.href = '/login'
            }
        },
        type: 'warning'
    }
}

function handle403Error(): UserFriendlyError {
    return {
        title: 'Access Denied',
        message: 'You don\'t have permission to perform this action.',
        suggestion: 'Contact your administrator if you believe this is an error.',
        type: 'error'
    }
}

function handle404Error(): UserFriendlyError {
    return {
        title: 'Not Found',
        message: 'The requested information could not be found.',
        suggestion: 'It may have been deleted or moved. Try refreshing the page.',
        action: {
            label: 'Refresh',
            handler: () => window.location.reload()
        },
        type: 'error'
    }
}

function handle409Error(data: any): UserFriendlyError {
    return {
        title: 'Conflict',
        message: data?.message || 'This action conflicts with existing data.',
        suggestion: 'Please check for duplicate entries or try a different approach.',
        type: 'error'
    }
}

function handle422Error(data: any): UserFriendlyError {
    return {
        title: 'Validation Error',
        message: data?.message || 'The provided data failed validation.',
        suggestion: 'Please review your input and ensure all required fields are filled correctly.',
        type: 'error'
    }
}

function handle429Error(): UserFriendlyError {
    return {
        title: 'Too Many Requests',
        message: 'You\'re making requests too quickly.',
        suggestion: 'Please wait a moment before trying again.',
        type: 'warning'
    }
}

function handle5xxError(status: number): UserFriendlyError {
    const messages = {
        500: 'The server encountered an unexpected error.',
        502: 'The server is temporarily unavailable.',
        503: 'The service is temporarily unavailable.',
        504: 'The server took too long to respond.'
    }

    return {
        title: 'Server Error',
        message: messages[status as keyof typeof messages] || 'A server error occurred.',
        suggestion: 'This is usually temporary. Please try again in a few minutes.',
        action: {
            label: 'Try Again',
            handler: () => window.location.reload()
        },
        type: 'error'
    }
}

function handleUnknownError(_error: ApiError): UserFriendlyError {
    return {
        title: 'Unexpected Error',
        message: 'Something unexpected happened.',
        suggestion: 'Please try again or contact support if the problem persists.',
        action: {
            label: 'Retry',
            handler: () => window.location.reload()
        },
        type: 'error'
    }
}

/**
 * Get user-friendly messages for specific operations
 */
export const operationMessages = {
    // Transaction operations
    transaction: {
        create: {
            success: 'Transaction added successfully',
            error: 'Failed to add transaction',
            suggestion: 'Please check that all required fields are filled and amounts are valid.'
        },
        update: {
            success: 'Transaction updated successfully',
            error: 'Failed to update transaction',
            suggestion: 'Please verify the transaction still exists and try again.'
        },
        delete: {
            success: 'Transaction deleted successfully',
            error: 'Failed to delete transaction',
            suggestion: 'The transaction may have already been deleted or is being used elsewhere.'
        },
        load: {
            error: 'Failed to load transactions',
            suggestion: 'Please check your internet connection and try refreshing the page.'
        }
    },

    // Budget operations
    budget: {
        create: {
            success: 'Budget created successfully',
            error: 'Failed to create budget',
            suggestion: 'Please ensure all budget categories have valid amounts and dates are correct.'
        },
        update: {
            success: 'Budget updated successfully',
            error: 'Failed to update budget',
            suggestion: 'Please check that budget amounts are positive and dates are valid.'
        },
        delete: {
            success: 'Budget deleted successfully',
            error: 'Failed to delete budget',
            suggestion: 'Active budgets with transactions may not be deletable.'
        },
        load: {
            error: 'Failed to load budget information',
            suggestion: 'Please refresh the page to try loading your budget again.'
        }
    },

    // Income operations
    income: {
        create: {
            success: 'Income source added successfully',
            error: 'Failed to add income source',
            suggestion: 'Please ensure the amount is positive and all required fields are filled.'
        },
        update: {
            success: 'Income source updated successfully',
            error: 'Failed to update income source',
            suggestion: 'Please verify the income source still exists and amounts are valid.'
        },
        delete: {
            success: 'Income source deleted successfully',
            error: 'Failed to delete income source',
            suggestion: 'Income sources linked to transactions may not be deletable.'
        },
        load: {
            error: 'Failed to load income sources',
            suggestion: 'Please refresh the page to try loading your income information again.'
        }
    },

    // Category operations
    category: {
        create: {
            success: 'Category created successfully',
            error: 'Failed to create category',
            suggestion: 'Please ensure the category name is unique and not empty.'
        },
        update: {
            success: 'Category updated successfully',
            error: 'Failed to update category',
            suggestion: 'System categories cannot be modified. Only custom categories can be edited.'
        },
        delete: {
            success: 'Category deleted successfully',
            error: 'Failed to delete category',
            suggestion: 'System categories and categories with transactions cannot be deleted.'
        },
        load: {
            error: 'Failed to load categories',
            suggestion: 'Please refresh the page to try loading categories again.'
        }
    },

    // Auth operations
    auth: {
        login: {
            success: 'Welcome back!',
            error: 'Login failed',
            suggestion: 'Please check your email and password and try again.'
        },
        register: {
            success: 'Account created successfully',
            error: 'Registration failed',
            suggestion: 'Please check that your email is unique and password meets requirements.'
        }
    },

    // Profile operations
    profile: {
        update: {
            success: 'Profile updated successfully',
            error: 'Failed to update profile',
            suggestion: 'Please ensure your email is valid and family size is between 0-5 children.'
        },
        passwordChange: {
            success: 'Password changed successfully',
            error: 'Failed to change password',
            suggestion: 'Please ensure your current password is correct and the new password meets requirements.'
        },
        load: {
            error: 'Failed to load profile information',
            suggestion: 'Please refresh the page or try logging in again.'
        }
    },

    // Report operations
    report: {
        load: {
            error: 'Failed to generate report',
            suggestion: 'Please try selecting a different date range or refresh the page.'
        }
    },

    // Dashboard operations
    dashboard: {
        load: {
            error: 'Failed to load dashboard data',
            suggestion: 'Please refresh the page to try loading your financial summary again.'
        }
    }
}

/**
 * Get contextual error message for specific operations
 */
export function getOperationError(
    operation: keyof typeof operationMessages,
    action: string,
    error?: ApiError
): UserFriendlyError {
    const opMessages = operationMessages[operation] as any
    const actionMessages = opMessages?.[action]

    if (!actionMessages) {
        return transformError(error || {})
    }

    // If we have a specific API error, transform it
    if (error) {
        const transformed = transformError(error)
        return {
            ...transformed,
            suggestion: actionMessages.suggestion || transformed.suggestion
        }
    }

    // Otherwise use the predefined message
    return {
        title: 'Operation Failed',
        message: actionMessages.error,
        suggestion: actionMessages.suggestion,
        type: 'error'
    }
}

/**
 * Validation error messages
 */
export const validationMessages = {
    required: (field: string) => `${field} is required`,
    email: 'Please enter a valid email address',
    minLength: (field: string, min: number) => `${field} must be at least ${min} characters`,
    maxLength: (field: string, max: number) => `${field} must be no more than ${max} characters`,
    positive: (field: string) => `${field} must be a positive number`,
    range: (field: string, min: number, max: number) => `${field} must be between ${min} and ${max}`,
    password: 'Password must contain at least 8 characters with uppercase, lowercase, and number',
    passwordMatch: 'Passwords do not match',
    futureDate: 'Date cannot be in the future',
    pastDate: 'Date cannot be in the past',
    invalidDate: 'Please enter a valid date'
}