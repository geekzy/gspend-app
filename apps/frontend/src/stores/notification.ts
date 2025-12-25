import { defineStore } from 'pinia'
import { ref } from 'vue'

export type NotificationType = 'success' | 'error' | 'info' | 'warning'

export interface Notification {
    id: string
    message: string
    type: NotificationType
    duration?: number
}

export const useNotificationStore = defineStore('notification', () => {
    const notifications = ref<Notification[]>([])

    const addNotification = (message: string, type: NotificationType = 'info', duration = 5000) => {
        const id = Math.random().toString(36).substring(2, 9)
        notifications.value.push({ id, message, type, duration })

        if (duration > 0) {
            setTimeout(() => {
                removeNotification(id)
            }, duration)
        }
    }

    const removeNotification = (id: string) => {
        notifications.value = notifications.value.filter((n) => n.id !== id)
    }

    const success = (message: string) => addNotification(message, 'success')
    const error = (message: string) => addNotification(message, 'error')
    const info = (message: string) => addNotification(message, 'info')
    const warning = (message: string) => addNotification(message, 'warning')

    return {
        notifications,
        addNotification,
        removeNotification,
        success,
        error,
        info,
        warning,
    }
})
