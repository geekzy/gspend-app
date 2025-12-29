import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface LoadingState {
    [key: string]: boolean
}

export const useLoadingStore = defineStore('loading', () => {
    const loadingStates = ref<LoadingState>({})

    // Set loading state for a specific operation
    const setLoading = (key: string, isLoading: boolean) => {
        if (isLoading) {
            loadingStates.value[key] = true
        } else {
            delete loadingStates.value[key]
        }
    }

    // Check if a specific operation is loading
    const isLoading = (key: string): boolean => {
        return loadingStates.value[key] || false
    }

    // Check if any operation is loading
    const isAnyLoading = computed(() => {
        return Object.keys(loadingStates.value).length > 0
    })

    // Get all currently loading operations
    const getLoadingOperations = computed(() => {
        return Object.keys(loadingStates.value)
    })

    // Wrapper function to automatically handle loading states
    const withLoading = async <T>(
        key: string,
        operation: () => Promise<T>
    ): Promise<T> => {
        try {
            setLoading(key, true)
            const result = await operation()
            return result
        } finally {
            setLoading(key, false)
        }
    }

    // Predefined loading keys for common operations
    const LOADING_KEYS = {
        // Dashboard
        DASHBOARD_SUMMARY: 'dashboard.summary',

        // Transactions
        TRANSACTIONS_LIST: 'transactions.list',
        TRANSACTION_CREATE: 'transaction.create',
        TRANSACTION_UPDATE: 'transaction.update',
        TRANSACTION_DELETE: 'transaction.delete',

        // Budgets
        BUDGETS_LIST: 'budgets.list',
        BUDGET_CREATE: 'budget.create',
        BUDGET_UPDATE: 'budget.update',
        BUDGET_DELETE: 'budget.delete',
        BUDGET_ACTIVE: 'budget.active',

        // Income
        INCOME_LIST: 'income.list',
        INCOME_CREATE: 'income.create',
        INCOME_UPDATE: 'income.update',
        INCOME_DELETE: 'income.delete',

        // Categories
        CATEGORIES_LIST: 'categories.list',
        CATEGORY_CREATE: 'category.create',
        CATEGORY_UPDATE: 'category.update',
        CATEGORY_DELETE: 'category.delete',

        // Reports
        REPORT_BUDGET_VS_ACTUAL: 'report.budgetVsActual',
        REPORT_SPENDING_BY_CATEGORY: 'report.spendingByCategory',
        REPORT_MONTHLY_TRENDS: 'report.monthlyTrends',

        // Profile
        PROFILE_LOAD: 'profile.load',
        PROFILE_UPDATE: 'profile.update',
        PROFILE_PASSWORD_CHANGE: 'profile.passwordChange',

        // Auth
        AUTH_LOGIN: 'auth.login',
        AUTH_REGISTER: 'auth.register',
        AUTH_LOGOUT: 'auth.logout'
    }

    return {
        loadingStates,
        setLoading,
        isLoading,
        isAnyLoading,
        getLoadingOperations,
        withLoading,
        LOADING_KEYS
    }
})