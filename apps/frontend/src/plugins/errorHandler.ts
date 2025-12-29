import type { App } from 'vue'
import { useNotificationStore } from '@/stores/notification'
import { transformError } from '@/utils/errorHandling'

export function setupErrorHandler(app: App) {
    // Global error handler for uncaught errors
    app.config.errorHandler = (error: any, _instance: any, info: string) => {
        console.error('Global error handler:', error, 'Info:', info)

        // Don't show notifications for errors that are already handled
        if (error.handled) {
            return
        }

        const notificationStore = useNotificationStore()
        const friendlyError = transformError(error)

        // Show user-friendly error notification
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
    }

    // Handle unhandled promise rejections
    window.addEventListener('unhandledrejection', (event) => {
        console.error('Unhandled promise rejection:', event.reason)

        // Don't show notifications for errors that are already handled
        if (event.reason?.handled) {
            return
        }

        const notificationStore = useNotificationStore()
        const friendlyError = transformError(event.reason)

        // Show user-friendly error notification
        notificationStore.error(friendlyError.message, 6000, {
            title: friendlyError.title
        })

        // Prevent the default browser error handling
        event.preventDefault()
    })

    // Handle general JavaScript errors
    window.addEventListener('error', (event) => {
        console.error('Global JavaScript error:', event.error)

        // Don't show notifications for errors that are already handled
        if (event.error?.handled) {
            return
        }

        const notificationStore = useNotificationStore()
        const friendlyError = transformError(event.error)

        // Show user-friendly error notification
        notificationStore.error(friendlyError.message, 6000, {
            title: friendlyError.title
        })
    })
}

// Utility to mark errors as handled to prevent duplicate notifications
export function markErrorAsHandled(error: any) {
    if (error && typeof error === 'object') {
        error.handled = true
    }
    return error
}