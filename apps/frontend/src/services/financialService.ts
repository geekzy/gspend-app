import { financeApi } from './api/axios'
import { useLoadingStore } from '@/stores/loading'

export interface Category {
    id: string
    name: string
    type: 'income' | 'expense'
    icon: string
    color: string
    isSystem: boolean
    sortOrder: number
    usageCount?: number
}

export interface Income {
    id: string
    source: string
    amount: number
    frequency: string
    effectiveDate: string
}

export interface BudgetItem {
    id: string
    categoryId: string
    categoryName: string
    plannedAmount: number
    spentAmount: number
    notes: string
}

export interface Budget {
    id: string
    name: string
    periodType: string
    startDate: string
    endDate: string
    totalAmount: number
    items: BudgetItem[]
}

export interface Transaction {
    id: string
    categoryId: string
    categoryName: string
    budgetId?: string
    budgetName?: string
    type: 'income' | 'expense'
    amount: number
    description: string
    transactionDate: string
    paymentMethod: string
    notes: string
}

export interface TransactionFilter {
    startDate?: string
    endDate?: string
    categoryIds?: string[]
    type?: 'income' | 'expense' | 'all'
    page?: number
    perPage?: number
    sortBy?: string
    sortOrder?: 'asc' | 'desc'
}

export interface TransactionResponse {
    transactions: Transaction[]
    pagination: {
        page: number
        perPage: number
        total: number
        totalPages: number
    }
    filtersApplied: TransactionFilter
}

// Dashboard specific interfaces
export interface BudgetProgress {
    totalBudget: number
    totalSpent: number
    remainingBudget: number
    percentageUsed: number
}

export interface CategorySpending {
    categoryId: string
    categoryName: string
    amount: number
    percentage: number
}

export interface DashboardSummary {
    totalBalance: number
    monthlyIncome: number
    monthlyExpenses: number
    budgetProgress: BudgetProgress
    topCategories: CategorySpending[]
    recentTransactions: Transaction[]
}

// Report specific interfaces
export interface BudgetVsActualItem {
    categoryName: string
    budgeted: number
    actual: number
    variance: number
    percentageUsed: number
}

export interface BudgetVsActualReport {
    month: string
    categories: BudgetVsActualItem[]
    totalBudgeted: number
    totalSpent: number
    overallVariance: number
}

export interface SpendingByCategoryReport {
    startDate: string
    endDate: string
    categories: CategorySpending[]
    totalSpent: number
}

export interface MonthlySpending {
    month: string
    totalIncome: number
    totalExpenses: number
    netAmount: number
    topCategory: CategorySpending | null
}

export interface MonthlyTrendsReport {
    months: number
    monthlyData: MonthlySpending[]
    averageSpending: number
    trendDirection: 'increasing' | 'decreasing' | 'stable'
}

export const financialService = {
    // Dashboard
    async getDashboardSummary(): Promise<DashboardSummary> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.DASHBOARD_SUMMARY, async () => {
            const response = await financeApi.get('/dashboard/summary')
            // Backend returns wrapped response { success: true, data: ... }
            if (response.data && response.data.success && response.data.data) {
                return response.data.data
            }
            return response.data
        })
    },

    // Categories
    async getCategories(type?: string): Promise<Category[]> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.CATEGORIES_LIST, async () => {
            const params = type ? { type } : {}
            const response = await financeApi.get('/categories', { params })
            // Handle wrapped response
            if (response.data && response.data.success && response.data.data) {
                return response.data.data
            }
            return response.data
        })
    },

    async createCategory(category: Partial<Category>): Promise<Category> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.CATEGORY_CREATE, async () => {
            const response = await financeApi.post('/categories', category)
            return response.data
        })
    },

    async updateCategory(id: string, category: Partial<Category>): Promise<Category> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.CATEGORY_UPDATE, async () => {
            const response = await financeApi.put(`/categories/${id}`, category)
            return response.data
        })
    },

    async deleteCategory(id: string): Promise<void> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.CATEGORY_DELETE, async () => {
            await financeApi.delete(`/categories/${id}`)
        })
    },

    // Incomes
    async getIncomes(): Promise<Income[]> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.INCOME_LIST, async () => {
            const response = await financeApi.get('/incomes')
            return response.data
        })
    },

    async createIncome(income: Partial<Income>): Promise<Income> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.INCOME_CREATE, async () => {
            const response = await financeApi.post('/incomes', income)
            return response.data
        })
    },

    async updateIncome(id: string, income: Partial<Income>): Promise<Income> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.INCOME_UPDATE, async () => {
            const response = await financeApi.put(`/incomes/${id}`, income)
            return response.data
        })
    },

    async deleteIncome(id: string): Promise<void> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.INCOME_DELETE, async () => {
            await financeApi.delete(`/incomes/${id}`)
        })
    },

    // Budgets
    async getBudgets(): Promise<Budget[]> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.BUDGETS_LIST, async () => {
            const response = await financeApi.get('/budgets')
            if (response.data && response.data.success && response.data.data) {
                return response.data.data
            }
            return response.data
        })
    },

    async getActiveBudget(date?: string): Promise<Budget> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.BUDGET_ACTIVE, async () => {
            const params = date ? { date } : {}
            const response = await financeApi.get('/budgets/active', { params })
            if (response.data && response.data.success && response.data.data) {
                return response.data.data
            }
            return response.data
        })
    },

    async createBudget(budget: Partial<Budget>): Promise<Budget> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.BUDGET_CREATE, async () => {
            const response = await financeApi.post('/budgets', budget)
            return response.data
        })
    },

    async updateBudget(id: string, budget: Partial<Budget>): Promise<Budget> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.BUDGET_UPDATE, async () => {
            const response = await financeApi.put(`/budgets/${id}`, budget)
            return response.data
        })
    },

    async deleteBudget(id: string): Promise<void> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.BUDGET_DELETE, async () => {
            await financeApi.delete(`/budgets/${id}`)
        })
    },

    // Transactions
    async getTransactions(filter?: TransactionFilter): Promise<TransactionResponse> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.TRANSACTIONS_LIST, async () => {
            const response = await financeApi.get('/transactions', { params: filter })
            if (response.data && response.data.success && response.data.data) {
                return response.data.data
            }
            return response.data
        })
    },

    async createTransaction(transaction: Partial<Transaction>): Promise<Transaction> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.TRANSACTION_CREATE, async () => {
            const response = await financeApi.post('/transactions', transaction)
            return response.data
        })
    },

    async updateTransaction(id: string, transaction: Partial<Transaction>): Promise<Transaction> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.TRANSACTION_UPDATE, async () => {
            const response = await financeApi.put(`/transactions/${id}`, transaction)
            return response.data
        })
    },

    async deleteTransaction(id: string): Promise<void> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.TRANSACTION_DELETE, async () => {
            await financeApi.delete(`/transactions/${id}`)
        })
    },

    // Reports
    async getBudgetVsActualReport(month?: string): Promise<BudgetVsActualReport> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.REPORT_BUDGET_VS_ACTUAL, async () => {
            const params = month ? { month } : {}
            const response = await financeApi.get('/reports/budget-vs-actual', { params })
            if (response.data && response.data.success && response.data.data) {
                return response.data.data
            }
            return response.data
        })
    },

    async getSpendingByCategoryReport(startDate?: string, endDate?: string): Promise<SpendingByCategoryReport> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.REPORT_SPENDING_BY_CATEGORY, async () => {
            const params: any = {}
            if (startDate) params.startDate = startDate
            if (endDate) params.endDate = endDate
            const response = await financeApi.get('/reports/spending-by-category', { params })
            if (response.data && response.data.success && response.data.data) {
                return response.data.data
            }
            return response.data
        })
    },

    async getMonthlyTrendsReport(months?: number): Promise<MonthlyTrendsReport> {
        const loadingStore = useLoadingStore()
        return loadingStore.withLoading(loadingStore.LOADING_KEYS.REPORT_MONTHLY_TRENDS, async () => {
            const params = months ? { months } : {}
            const response = await financeApi.get('/reports/monthly-trends', { params })
            if (response.data && response.data.success && response.data.data) {
                return response.data.data
            }
            return response.data
        })
    }
}
