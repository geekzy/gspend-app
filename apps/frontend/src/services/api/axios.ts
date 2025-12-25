import axios from 'axios'
import { useNotificationStore } from '@/stores/notification'

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

// Add a response interceptor to handle errors
const handleError = (error: any) => {
    const notificationStore = useNotificationStore()

    if (error.response) {
        if (error.response.status === 401) {
            // Only redirect if not already on login/register pages
            if (!window.location.pathname.includes('/login') && !window.location.pathname.includes('/register')) {
                localStorage.removeItem('auth_token')
                window.location.href = '/login'
            }
        } else {
            // Show toast for other errors
            const message = error.response.data?.message || 'An unexpected error occurred'
            notificationStore.error(message)
        }
    } else {
        notificationStore.error('Network error or server unreachable')
    }

    return Promise.reject(error)
}

authApi.interceptors.response.use(response => response, handleError)
financeApi.interceptors.response.use(response => response, handleError)
