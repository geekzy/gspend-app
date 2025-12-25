import { financeApi } from './api/axios'

export interface Category {
    id: string
    name: string
    type: 'income' | 'expense'
    icon: string
    color: string
    isSystem: boolean
    sortOrder: number
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

export const financialService = {
    // Categories
    async getCategories(type?: string): Promise<Category[]> {
        const params = type ? { type } : {}
        const response = await financeApi.get('/categories', { params })
        return response.data
    },

    async createCategory(category: Partial<Category>): Promise<Category> {
        const response = await financeApi.post('/categories', category)
        return response.data
    },

    // Incomes
    async getIncomes(): Promise<Income[]> {
        const response = await financeApi.get('/incomes')
        return response.data
    },

    async createIncome(income: Partial<Income>): Promise<Income> {
        const response = await financeApi.post('/incomes', income)
        return response.data
    },

    // Budgets
    async getBudgets(): Promise<Budget[]> {
        const response = await financeApi.get('/budgets')
        return response.data
    },

    async getActiveBudget(date?: string): Promise<Budget> {
        const params = date ? { date } : {}
        const response = await financeApi.get('/budgets/active', { params })
        return response.data
    },

    async createBudget(budget: Partial<Budget>): Promise<Budget> {
        const response = await financeApi.post('/budgets', budget)
        return response.data
    },

    // Transactions
    async getTransactions(filter?: any): Promise<Transaction[]> {
        const response = await financeApi.get('/transactions', { params: filter })
        return response.data
    },

    async createTransaction(transaction: Partial<Transaction>): Promise<Transaction> {
        const response = await financeApi.post('/transactions', transaction)
        return response.data
    }
}
