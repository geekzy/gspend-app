import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService, type UpdateProfileRequest, type ChangePasswordRequest } from '@/services/authService'

export const useAuthStore = defineStore('auth', () => {
    const token = ref(localStorage.getItem('auth_token'))
    const user = ref(JSON.parse(localStorage.getItem('user_profile') || 'null'))

    const isAuthenticated = computed(() => !!token.value)

    function setSession(newToken: string, newUser: any) {
        token.value = newToken
        user.value = newUser
        localStorage.setItem('auth_token', newToken)
        localStorage.setItem('user_profile', JSON.stringify(newUser))
    }

    async function fetchProfile() {
        try {
            const userData = await authService.getProfile()
            user.value = userData
            localStorage.setItem('user_profile', JSON.stringify(userData))
            return userData
        } catch (error) {
            logout()
            throw error
        }
    }

    async function updateProfile(profileData: UpdateProfileRequest) {
        try {
            const response = await authService.updateProfile(profileData)
            user.value = response.user
            localStorage.setItem('user_profile', JSON.stringify(response.user))
            return response
        } catch (error) {
            throw error
        }
    }

    async function changePassword(passwordData: ChangePasswordRequest) {
        try {
            const response = await authService.changePassword(passwordData)
            return response
        } catch (error) {
            throw error
        }
    }

    function logout() {
        token.value = null
        user.value = null
        localStorage.removeItem('auth_token')
        localStorage.removeItem('user_profile')
    }

    return {
        token,
        user,
        isAuthenticated,
        setSession,
        fetchProfile,
        updateProfile,
        changePassword,
        logout
    }
})
