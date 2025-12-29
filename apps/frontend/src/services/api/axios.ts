import axios from 'axios'
import { useNotificationStore } from '@/stores/notification'
import { transformError } from '@/utils/errorHandling'

const API_AUTH_URL = import.meta.env.VITE_AUTH_SERVICE_URL || 'http://localhost/api/v1/auth'
const API_FINANCE_URL = import.meta.env.VITE_FINANCIAL_SERVICE_URL || 'http://localhost/api/v1'

export const authApi = axios.create({
    baseURL: API_AUTH_URL,
    headers: {
        'Content-Type': 'application/json'
    }
})

export const financeApi = axios.create({
    baseURL: API_FINANCE_URL,
    headers: {
        'Content-Type': 'application/json'
    }
})

// Add a request interceptor to add the auth token to requests
const addAuthToken = (config: any) => {
    const token = localStorage.getItem('auth_token')
    if (token) {
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
}

authApi.interceptors.request.use(addAuthToken)
financeApi.interceptors.request.use(addAuthToken)

// Add a response interceptor to handle errors with user-friendly messages
const handleError = (error: any) => {
    // Don't show notifications for errors that will be handled by specific components
    // This prevents double notifications
    if (error.config?.skipGlobalErrorHandler) {
        return Promise.reject(error)
    }

    const notificationStore = useNotificationStore()

    if (error.response) {
        const friendlyError = transformError(error)

        if (error.response.status === 401) {
            // Only redirect if not already on login/register pages
            if (!window.location.pathname.includes('/login') && !window.location.pathname.includes('/register')) {
                localStorage.removeItem('auth_token')
                notificationStore.error('Your session has expired. Please log in again.', 5000, {
                    title: 'Session Expired',
                    actions: [
                        {
                            label: 'Login',
                            handler: () => window.location.href = '/login',
                            style: 'primary'
                        }
                    ]
                })
            }
        } else {
            // Show user-friendly error message
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
    } else {
        // Network error
        notificationStore.networkError()
    }

    return Promise.reject(error)
}

authApi.interceptors.response.use(response => response, handleError)
financeApi.interceptors.response.use(response => response, handleError)
