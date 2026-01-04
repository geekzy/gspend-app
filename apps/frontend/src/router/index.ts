import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/Dashboard.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/auth/Login.vue')
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('@/views/auth/Register.vue')
    },
    {
        path: '/forgot-password',
        name: 'ForgotPassword',
        component: () => import('@/views/auth/ForgotPassword.vue')
    },
    {
        path: '/reset-password',
        name: 'ResetPassword',
        component: () => import('@/views/auth/ResetPassword.vue')
    },
    {
        path: '/verify-email',
        name: 'VerifyEmail',
        component: () => import('@/views/auth/VerifyEmail.vue')
    },
    {
        path: '/income',
        name: 'Income',
        component: () => import('@/views/income/IncomeIndex.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/budget',
        name: 'Budget',
        component: () => import('@/views/budget/BudgetIndex.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/transactions',
        name: 'Transactions',
        component: () => import('@/views/transaction/TransactionIndex.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/profile',
        name: 'Profile',
        component: () => import('@/views/profile/Profile.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/profile/edit',
        name: 'ProfileEdit',
        component: () => import('@/views/profile/ProfileEdit.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/profile/password',
        name: 'PasswordChange',
        component: () => import('@/views/profile/PasswordChange.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/reports',
        name: 'Reports',
        component: () => import('@/views/reports/ReportsIndex.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/reports/budget-vs-actual',
        name: 'BudgetVsActual',
        component: () => import('@/views/reports/BudgetVsActual.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/reports/spending-by-category',
        name: 'SpendingByCategory',
        component: () => import('@/views/reports/SpendingByCategory.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/reports/monthly-trends',
        name: 'MonthlyTrends',
        component: () => import('@/views/reports/MonthlyTrends.vue'),
        meta: { requiresAuth: true }
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// Simple mock auth guard
router.beforeEach((to, _from, next) => {
    const isAuthenticated = localStorage.getItem('auth_token')
    if (to.meta.requiresAuth && !isAuthenticated) {
        next('/login')
    } else {
        next()
    }
})

export default router
