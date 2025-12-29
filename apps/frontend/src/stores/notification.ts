import { defineStore } from 'pinia'
import { ref } from 'vue'

export type NotificationType = 'success' | 'error' | 'info' | 'warning'

export interface Notification {
    id: string
    message: string
    type: NotificationType
    duration?: number
    title?: string
    actions?: NotificationAction[]
}

export interface NotificationAction {
    label: string
    handler: () => void
    style?: 'primary' | 'secondary'
}

export const useNotificationStore = defineStore('notification', () => {
    const notifications = ref<Notification[]>([])

    const addNotification = (message: string, type: NotificationType = 'info', duration = 5000, options?: { title?: string, actions?: NotificationAction[] }) => {
        const id = Math.random().toString(36).substring(2, 9)
        notifications.value.push({
            id,
            message,
            type,
            duration,
            title: options?.title,
            actions: options?.actions
        })

        if (duration > 0) {
            setTimeout(() => {
                removeNotification(id)
            }, duration)
        }
    }

    const removeNotification = (id: string) => {
        notifications.value = notifications.value.filter((n) => n.id !== id)
    }

    const success = (message: string, duration = 3000, options?: { title?: string, actions?: NotificationAction[] }) =>
        addNotification(message, 'success', duration, options)
    const error = (message: string, duration = 5000, options?: { title?: string, actions?: NotificationAction[] }) =>
        addNotification(message, 'error', duration, options)
    const info = (message: string, duration = 4000, options?: { title?: string, actions?: NotificationAction[] }) =>
        addNotification(message, 'info', duration, options)
    const warning = (message: string, duration = 4000, options?: { title?: string, actions?: NotificationAction[] }) =>
        addNotification(message, 'warning', duration, options)

    // Validation-specific notifications
    const validationError = (message: string = 'Please fix the validation errors') => {
        addNotification(message, 'error', 4000, { title: 'Validation Error' })
    }

    const validationSuccess = (message: string = 'All fields are valid') => {
        addNotification(message, 'success', 2000, { title: 'Success' })
    }

    // CRUD operation notifications
    const crudSuccess = (operation: string, resource: string) => {
        const messages = {
            create: `${resource} created successfully`,
            update: `${resource} updated successfully`,
            delete: `${resource} deleted successfully`,
            save: `${resource} saved successfully`
        }
        addNotification(messages[operation as keyof typeof messages] || `${operation} completed successfully`, 'success', 3000)
    }

    const crudError = (operation: string, resource: string, error?: string) => {
        const messages = {
            create: `Failed to create ${resource}`,
            update: `Failed to update ${resource}`,
            delete: `Failed to delete ${resource}`,
            save: `Failed to save ${resource}`,
            load: `Failed to load ${resource}`
        }
        const message = error || messages[operation as keyof typeof messages] || `${operation} failed`
        addNotification(message, 'error', 5000, { title: 'Operation Failed' })
    }

    // Network and loading notifications
    const networkError = (message: string = 'Network error. Please check your connection and try again.') => {
        addNotification(message, 'error', 6000, {
            title: 'Connection Error',
            actions: [
                {
                    label: 'Retry',
                    handler: () => window.location.reload(),
                    style: 'primary'
                }
            ]
        })
    }

    const serverError = (message: string = 'Server error. Please try again later.') => {
        addNotification(message, 'error', 6000, { title: 'Server Error' })
    }

    return {
        notifications,
        addNotification,
        removeNotification,
        success,
        error,
        info,
        warning,
        validationError,
        validationSuccess,
        crudSuccess,
        crudError,
        networkError,
        serverError,
    }
})
