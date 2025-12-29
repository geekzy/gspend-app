import { ref } from 'vue'
import { useNotificationStore } from '@/stores/notification'
import { transformError, getOperationError, operationMessages, type ApiError, type UserFriendlyError } from '@/utils/errorHandling'

export function useErrorHandling() {
    const notificationStore = useNotificationStore()
    const lastError = ref<UserFriendlyError | null>(null)

    /**
     * Handle any error with user-friendly messaging
     */
    const handleError = (error: ApiError, context?: string) => {
        console.error('Error occurred:', error, 'Context:', context)

        const friendlyError = transformError(error)
        lastError.value = friendlyError

        // Show notification with action if available
        if (friendlyError.action) {
            notificationStore.error(friendlyError.message, 8000, {
                title: friendlyError.title,
                actions: [friendlyError.action]
            })
        } else {
            notificationStore.error(friendlyError.message, 6000, {
                title: friendlyError.title
            })
        }

        return friendlyError
    }

    /**
     * Handle operation-specific errors
     */
    const handleOperationError = (
        operation: keyof typeof operationMessages,
        action: string,
        error: ApiError,
        context?: string
    ) => {
        console.error(`${operation}.${action} failed:`, error, 'Context:', context)

        const friendlyError = getOperationError(operation, action, error)
        lastError.value = friendlyError

        if (friendlyError.action) {
            notificationStore.error(friendlyError.message, 8000, {
                title: friendlyError.title,
                actions: [friendlyError.action]
            })
        } else {
            notificationStore.error(friendlyError.message, 6000, {
                title: friendlyError.title
            })
        }

        return friendlyError
    }

    /**
     * Handle network errors specifically
     */
    const handleNetworkError = (context?: string) => {
        console.error('Network error occurred. Context:', context)

        const friendlyError: UserFriendlyError = {
            title: 'Connection Problem',
            message: 'Unable to connect to the server. Please check your internet connection.',
            suggestion: 'Try refreshing the page or check your network connection.',
            action: {
                label: 'Retry',
                handler: () => window.location.reload()
            },
            type: 'error'
        }

        lastError.value = friendlyError
        notificationStore.networkError()

        return friendlyError
    }

    /**
     * Handle validation errors
     */
    const handleValidationError = (errors: Record<string, string>, context?: string) => {
        console.error('Validation errors:', errors, 'Context:', context)

        const errorCount = Object.keys(errors).length
        const message = errorCount === 1
            ? Object.values(errors)[0]
            : `Please fix ${errorCount} validation errors`

        const friendlyError: UserFriendlyError = {
            title: 'Validation Error',
            message,
            suggestion: 'Please review the highlighted fields and correct any errors.',
            type: 'error'
        }

        lastError.value = friendlyError
        notificationStore.validationError(message)

        return friendlyError
    }

    /**
     * Wrapper for async operations with automatic error handling
     */
    const withErrorHandling = async <T>(
        operation: () => Promise<T>,
        options: {
            context?: string
            operationType?: keyof typeof operationMessages
            operationAction?: string
            onError?: (error: UserFriendlyError) => void
            suppressNotification?: boolean
        } = {}
    ): Promise<T | null> => {
        try {
            return await operation()
        } catch (error) {
            let friendlyError: UserFriendlyError

            if (options.operationType && options.operationAction) {
                friendlyError = handleOperationError(
                    options.operationType,
                    options.operationAction,
                    error as ApiError,
                    options.context
                )
            } else {
                friendlyError = handleError(error as ApiError, options.context)
            }

            if (options.onError) {
                options.onError(friendlyError)
            }

            return null
        }
    }

    /**
     * Clear the last error
     */
    const clearError = () => {
        lastError.value = null
    }

    /**
     * Show success message for operations
     */
    const showOperationSuccess = (
        operation: keyof typeof operationMessages,
        action: string,
        customMessage?: string
    ) => {
        const opMessages = operationMessages[operation] as any
        const actionMessages = opMessages?.[action]

        const message = customMessage || actionMessages?.success || `${action} completed successfully`
        notificationStore.success(message)
    }

    /**
     * Retry mechanism for failed operations
     */
    const createRetryHandler = (operation: () => Promise<any>, maxRetries = 3) => {
        let retryCount = 0

        const retry = async (): Promise<any> => {
            try {
                return await operation()
            } catch (error) {
                retryCount++

                if (retryCount < maxRetries) {
                    console.log(`Retrying operation (attempt ${retryCount + 1}/${maxRetries})`)
                    // Exponential backoff
                    await new Promise(resolve => setTimeout(resolve, Math.pow(2, retryCount) * 1000))
                    return retry()
                } else {
                    throw error
                }
            }
        }

        return retry
    }

    return {
        lastError,
        handleError,
        handleOperationError,
        handleNetworkError,
        handleValidationError,
        withErrorHandling,
        clearError,
        showOperationSuccess,
        createRetryHandler
    }
}