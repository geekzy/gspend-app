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
