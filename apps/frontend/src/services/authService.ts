import { authApi } from './api/axios'
import { useLoadingStore } from '@/stores/loading'

export interface LoginResponse {
    accessToken: string
    refreshToken: string
    user: {
        id: string
        email: string
        fullName: string
        familySize: number
    }
}

export interface UpdateProfileRequest {
    fullName: string
    familySize: number
    email: string
}

export interface ChangePasswordRequest {
    currentPassword: string
    newPassword: string
}

export interface ChangePasswordResponse {
    success: boolean
    message: string
}

export interface ProfileResponse {
    success: boolean
    message: string
    user: {
        id: string
        email: string
        fullName: string
        familySize: number
    }
}

export const authService = {
    async login(credentials: any): Promise<LoginResponse> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.AUTH_LOGIN, async () => {
            const response = await authApi.post('/login', credentials)
            return response.data
        })
    },

    async register(userData: any): Promise<LoginResponse> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.AUTH_REGISTER, async () => {
            const response = await authApi.post('/register', userData)
            return response.data
        })
    },

    async getProfile() {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.PROFILE_LOAD, async () => {
            const response = await authApi.get('/me')
            return response.data
        })
    },

    async updateProfile(profileData: UpdateProfileRequest): Promise<ProfileResponse> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.PROFILE_UPDATE, async () => {
            const response = await authApi.put('/me', profileData)
            return response.data
        })
    },

    async changePassword(passwordData: ChangePasswordRequest): Promise<ChangePasswordResponse> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.PROFILE_PASSWORD_CHANGE, async () => {
            const response = await authApi.put('/change-password', passwordData)
            return response.data
        })
    }
}
