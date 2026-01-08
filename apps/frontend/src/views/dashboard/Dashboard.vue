<template>
  <MainLayout>
    <div class="px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto w-full">
      <div v-if="isLoading" class="space-y-8">
        <!-- Loading skeleton for stats -->
        <SkeletonLoader type="stats" />
        
        <!-- Loading skeleton for content -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div class="lg:col-span-2">
            <SkeletonLoader type="table" :rows="5" />
          </div>
          <div>
            <SkeletonLoader type="chart" />
          </div>
        </div>
      </div>

      <template v-else>
        <div class="md:flex md:items-center md:justify-between mb-8">
          <div class="flex-1 min-w-0">
            <h2 class="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
              Welcome back, {{ userName }}! 👋
            </h2>
            <p class="mt-1 text-sm text-gray-500">
              Here's what's happening with your family budget this month.
            </p>
          </div>
          <div class="mt-4 flex md:mt-0 md:ml-4 gap-3">
            <router-link to="/reports" class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-xl shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-all">
              Download Report
            </router-link>
            <router-link to="/transactions" class="inline-flex items-center px-4 py-2 border border-transparent rounded-xl shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-all shadow-md active:scale-95">
              Add Transaction
            </router-link>
          </div>
        </div>

        <!-- Stats Grid -->
        <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3 mb-10">
          <!-- Balance Card -->
          <div class="bg-white overflow-hidden shadow-sm hover:shadow-md transition-shadow duration-300 rounded-2xl border border-gray-100 flex flex-col">
            <div class="p-5 flex-1">
              <div class="flex items-center">
                <div class="flex-shrink-0 bg-primary-100 rounded-xl p-3">
                  <WalletIcon class="h-6 w-6 text-primary-600" />
                </div>
                <div class="ml-5 w-0 flex-1">
                  <dl>
                    <dt class="text-sm font-medium text-gray-500 truncate">Total Balance</dt>
                    <dd class="flex items-baseline">
                      <div class="text-2xl font-bold text-gray-900 leading-tight">{{ formatCurrency(dashboardData?.totalBalance, { decimals: 2 }) }}</div>
                      <div v-if="hasBalanceChange" class="ml-2 flex items-baseline text-sm font-semibold" :class="balanceChangeClass">
                        {{ balanceChangeText }}
                      </div>
                    </dd>
                  </dl>
                </div>
              </div>
            </div>
            <div class="bg-gray-50 px-5 py-3 border-t border-gray-100">
              <div class="text-sm">
              <router-link to="/reports" class="font-medium text-primary-600 hover:text-primary-500">View detailed stats &rarr;</router-link>
              </div>
            </div>
          </div>

          <!-- Budget Card -->
          <div class="bg-white overflow-hidden shadow-sm hover:shadow-md transition-shadow duration-300 rounded-2xl border border-gray-100 flex flex-col">
            <div class="p-5 flex-1">
              <div class="flex items-center">
                <div class="flex-shrink-0 bg-blue-100 rounded-xl p-3">
                  <BarChart3Icon class="h-6 w-6 text-blue-600" />
                </div>
                <div class="ml-5 w-0 flex-1">
                  <dl>
                    <dt class="text-sm font-medium text-gray-500 truncate">Monthly Budget</dt>
                    <dd class="flex items-baseline justify-between mb-1">
                      <div class="text-2xl font-bold text-gray-900 leading-tight">{{ formatCurrency(dashboardData?.budgetProgress?.totalSpent, { decimals: 2 }) }}</div>
                      <div v-if="dashboardData?.budgetProgress" class="text-xs text-gray-400 font-medium whitespace-nowrap">Plan: {{ formatCurrency(dashboardData.budgetProgress.totalBudget) }}</div>
                    </dd>
                  </dl>
                  <!-- Budget Progress Bar -->
                  <ProgressBar
                    v-if="dashboardData?.budgetProgress"
                    :label="'Budget Progress'"
                    :spent="dashboardData.budgetProgress.totalSpent"
                    :total="dashboardData.budgetProgress.totalBudget"
                    :show-amounts="false"
                  />
                </div>
              </div>
            </div>
            <div class="bg-gray-50 px-5 py-3 border-t border-gray-100">
              <div class="text-sm text-gray-500">
                <span class="font-medium text-gray-900">{{ budgetPercent }}%</span> used of your monthly target
              </div>
            </div>
          </div>

          <!-- Spending Card -->
          <div class="bg-white overflow-hidden shadow-sm hover:shadow-md transition-shadow duration-300 rounded-2xl border border-gray-100 flex flex-col">
            <div class="p-5 flex-1">
              <div class="flex items-center">
                <div class="flex-shrink-0 bg-red-100 rounded-xl p-3">
                  <ArrowDownIcon class="h-6 w-6 text-red-600" />
                </div>
                <div class="ml-5 w-0 flex-1">
                  <dl>
                    <dt class="text-sm font-medium text-gray-500 truncate">Monthly Expenses</dt>
                    <dd class="flex items-baseline">
                      <div class="text-2xl font-bold text-red-600 leading-tight">{{ formatCurrency(dashboardData?.monthlyExpenses, { decimals: 2 }) }}</div>
                    </dd>
                  </dl>
                </div>
              </div>
            </div>
            <div class="bg-gray-50 px-5 py-3 border-t border-gray-100">
              <div class="text-sm">
                <span class="text-gray-500">Income: </span>
                <span class="text-green-600 font-medium">{{ formatCurrency(dashboardData?.monthlyIncome) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Content Sections -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <!-- Recent Transactions -->
          <div class="lg:col-span-2 bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden flex flex-col">
            <div class="p-6 border-b border-gray-100 flex justify-between items-center bg-white sticky top-0">
              <h3 class="text-lg font-bold text-gray-900 leading-none">Recent Transactions</h3>
              <router-link to="/transactions" class="text-sm font-semibold text-primary-600 hover:text-primary-500">View All</router-link>
            </div>
            <div class="overflow-y-auto min-h-[300px]">
              <div v-if="!dashboardData?.recentTransactions?.length" class="flex flex-col items-center justify-center p-12 text-gray-400">
                <WalletIcon class="w-12 h-12 mb-4 opacity-20" />
                <p>No recent transactions found.</p>
              </div>
              <table v-else class="min-w-full divide-y divide-gray-100">
                <thead class="bg-gray-50">
                  <tr>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Transaction</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Category</th>
                    <th scope="col" class="px-6 py-3 text-right text-xs font-bold text-gray-500 uppercase tracking-wider">Amount</th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-100">
                  <tr v-for="tx in dashboardData.recentTransactions" :key="tx.id" class="hover:bg-gray-50 transition-colors pointer-cursor">
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="flex items-center">
                        <div class="flex-shrink-0 h-10 w-10 flex items-center justify-center rounded-xl bg-gray-50 text-xl shadow-inner border border-gray-100">
                          {{ tx.type === 'income' ? '💰' : '🛒' }}
                        </div>
                        <div class="ml-4">
                          <div class="text-sm font-bold text-gray-900">{{ tx.description }}</div>
                          <div class="text-xs text-gray-400">{{ formatDate(tx.transactionDate) }}</div>
                        </div>
                      </div>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-500">
                      <span class="px-2.5 py-1 rounded-full text-xs font-bold bg-gray-100 text-gray-600">
                        {{ tx.categoryName }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-bold" :class="tx.type === 'expense' ? 'text-red-600' : 'text-green-600'">
                      {{ tx.type === 'expense' ? '-' : '' }}{{ formatCurrency(tx.amount) }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- Sidebar Widgets -->
          <div class="space-y-8">
            <!-- Categories Breakdown -->
            <div class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
              <h3 class="text-lg font-bold text-gray-900 mb-6">Top Categories</h3>
              <div v-if="!dashboardData?.topCategories?.length" class="text-center py-8 text-gray-400 text-sm">
                No spending data available.
              </div>
              <div v-else>
                <!-- Pie Chart -->
                <div class="mb-6">
                  <PieChart 
                    :data="chartData" 
                    :width="250" 
                    :height="250"
                    :show-legend="false"
                  />
                </div>
                <!-- Category List -->
                <div class="space-y-3">
                  <div v-for="(cat, index) in dashboardData.topCategories.slice(0, 4)" :key="cat.categoryId">
                    <div class="flex items-center justify-between mb-1.5 px-0.5">
                      <div class="flex items-center">
                        <div 
                          class="w-3 h-3 rounded-full mr-2 flex-shrink-0"
                          :style="{ backgroundColor: chartColors[index] }"
                        ></div>
                        <span class="text-sm font-bold text-gray-700">{{ cat.categoryName }}</span>
                      </div>
                      <span class="text-sm font-bold text-gray-900">{{ formatCurrency(cat.amount) }}</span>
                    </div>
                  </div>
                </div>
              </div>
              <router-link to="/categories" class="block w-full mt-6 py-2.5 px-4 rounded-xl border border-gray-200 text-sm font-bold text-gray-600 hover:bg-gray-50 hover:text-gray-900 transition-all text-center">
                Manage Categories
              </router-link>
            </div>
          </div>
        </div>
      </template>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import ProgressBar from '@/components/charts/ProgressBar.vue'
import PieChart from '@/components/charts/PieChart.vue'

import SkeletonLoader from '@/components/common/SkeletonLoader.vue'
import { financialService, DashboardSummary } from '@/services/financialService'
import { formatCurrency } from '@/utils/currency'
import { useAuthStore } from '@/stores/auth'
import { useLoadingStore } from '@/stores/loading'
import { useNotificationStore } from '@/stores/notification'
import { 
  WalletIcon, 
  BarChart3Icon, 
  ArrowDownIcon 
} from 'lucide-vue-next'

const authStore = useAuthStore()
const loadingStore = useLoadingStore()
const notificationStore = useNotificationStore()
const dashboardData = ref<DashboardSummary | null>(null)
const error = ref<string | null>(null)

const isLoading = computed(() => loadingStore.isLoading(loadingStore.LOADING_KEYS.DASHBOARD_SUMMARY))

const userName = computed(() => authStore.user?.fullName || 'Family')

const budgetPercent = computed(() => {
  if (!dashboardData.value?.budgetProgress) return 0
  return Math.round(dashboardData.value.budgetProgress.percentageUsed)
})

const hasBalanceChange = computed(() => {
  // Only show percentage if there's actual balance data and transactions
  return dashboardData.value?.totalBalance !== undefined && 
         dashboardData.value.totalBalance > 0 && 
         dashboardData.value.recentTransactions?.length > 0
})

const balanceChangeClass = computed(() => {
  // For now, assume positive change - in real app this would be calculated
  return 'text-green-600'
})

const balanceChangeText = computed(() => {
  // Placeholder - in real app this would be calculated from historical data
  return '+2.5%'
})

const chartColors = ['#3B82F6', '#EF4444', '#10B981', '#F59E0B', '#8B5CF6', '#EC4899']

const chartData = computed(() => {
  if (!dashboardData.value?.topCategories?.length) return []
  
  return dashboardData.value.topCategories.slice(0, 6).map(cat => ({
    label: cat.categoryName,
    value: cat.amount
  }))
})

onMounted(async () => {
  try {
    error.value = null
    dashboardData.value = await financialService.getDashboardSummary()
  } catch (err) {
    console.error('Failed to fetch dashboard data:', err)
    error.value = 'Failed to load dashboard data. Please try again.'
    notificationStore.error('Failed to load dashboard data. Please refresh the page.')
    
    // Fallback to empty data structure to prevent UI errors
    dashboardData.value = {
      totalBalance: 0,
      monthlyIncome: 0,
      monthlyExpenses: 0,
      budgetProgress: {
        totalBudget: 0,
        totalSpent: 0,
        remainingBudget: 0,
        percentageUsed: 0
      },
      topCategories: [],
      recentTransactions: []
    }
  }
})

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })
}
</script>
