import { authApi } from './api/axios'

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

export const authService = {
    async login(credentials: any): Promise<LoginResponse> {
        const response = await authApi.post('/login', credentials)
        return response.data
    },

    async register(userData: any): Promise<LoginResponse> {
        const response = await authApi.post('/register', userData)
        return response.data
    },

    async getProfile() {
        const response = await authApi.get('/me')
        return response.data
    }
}
