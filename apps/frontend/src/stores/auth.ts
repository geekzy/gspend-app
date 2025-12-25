import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
    const token = ref(localStorage.getItem('auth_token'))
    const user = ref(JSON.parse(localStorage.getItem('user_profile') || 'null'))

    const isAuthenticated = computed(() => !!token.value)

    function login(mockToken: string, mockUser: any) {
        token.value = mockToken
        user.value = mockUser
        localStorage.setItem('auth_token', mockToken)
        localStorage.setItem('user_profile', JSON.stringify(mockUser))
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
        login,
        logout
    }
})
