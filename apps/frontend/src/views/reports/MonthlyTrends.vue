<template>
  <div class="p-6 max-w-6xl mx-auto">
    <!-- Breadcrumb Navigation -->
    <Breadcrumb :breadcrumbs="breadcrumbs" />
    
    <!-- Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900 mb-2">Monthly Spending Trends</h1>
      <p class="text-gray-600">Track your spending patterns over time to identify trends</p>
    </div>

    <!-- Month Range Filter -->
    <div class="mb-6 bg-white rounded-lg shadow p-4">
      <div class="flex flex-col sm:flex-row sm:items-center gap-4">
        <div class="flex items-center gap-2">
          <label for="months-select" class="text-sm font-medium text-gray-700 whitespace-nowrap">
            Show last:
          </label>
          <select
            id="months-select"
            v-model="selectedMonths"
            @change="loadReport"
            class="border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 min-w-0 flex-1 sm:flex-initial"
          >
            <option v-for="option in monthOptions" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
        </div>
        
        <button
          @click="loadReport"
          :disabled="loading"
          class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-sm w-full sm:w-auto"
        >
          <span v-if="loading">Loading...</span>
          <span v-else>Refresh</span>
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
      <div class="flex">
        <div class="text-red-400">
          <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
        </div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">Error loading report</h3>
          <p class="text-sm text-red-700 mt-1">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- Report Content -->
    <div v-else-if="report" class="space-y-6">
      <!-- Summary Cards -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div class="bg-white rounded-lg shadow p-6">
          <div class="text-sm font-medium text-gray-500">Average Monthly Spending</div>
          <div class="text-2xl font-bold text-gray-900">${{ report.averageSpending.toLocaleString() }}</div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="text-sm font-medium text-gray-500">Trend Direction</div>
          <div class="text-2xl font-bold" :class="trendColor">
            <span class="flex items-center">
              <component :is="trendIcon" class="h-6 w-6 mr-2" />
              {{ trendLabel }}
            </span>
          </div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="text-sm font-medium text-gray-500">Months Analyzed</div>
          <div class="text-2xl font-bold text-gray-900">{{ report.months }}</div>
        </div>
      </div>

      <!-- Line Charts -->
      <div class="grid grid-cols-1 xl:grid-cols-2 gap-6">
        <!-- Expenses Trend -->
        <div class="bg-white rounded-lg shadow p-4 sm:p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Monthly Expenses Trend</h3>
          <div v-if="chartData.length > 0" class="h-64 sm:h-80">
            <LineChart 
              :data="chartData" 
              :width="400" 
              :height="300"
              :color="chartColor"
              label="Monthly Expenses"
            />
          </div>
          <div v-else class="text-center py-8 text-gray-500">
            No expense data available for the selected period
          </div>
        </div>

        <!-- Income vs Expenses -->
        <div class="bg-white rounded-lg shadow p-4 sm:p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Net Amount Trend</h3>
          <div v-if="chartDataNet.length > 0" class="h-64 sm:h-80">
            <LineChart 
              :data="chartDataNet" 
              :width="400" 
              :height="300"
              :color="netTrendColor"
              label="Net Amount"
            />
          </div>
          <div v-else class="text-center py-8 text-gray-500">
            No net amount data available for the selected period
          </div>
        </div>
      </div>

      <!-- Monthly Breakdown Table -->
      <div class="bg-white rounded-lg shadow overflow-hidden">
        <div class="px-4 sm:px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-medium text-gray-900">Monthly Breakdown</h3>
        </div>
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Month
                </th>
                <th class="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Income
                </th>
                <th class="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Expenses
                </th>
                <th class="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Net Amount
                </th>
                <th class="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider hidden sm:table-cell">
                  Top Category
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="month in report.monthlyData" :key="month.month">
                <td class="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  {{ formatMonth(month.month) }}
                </td>
                <td class="px-4 sm:px-6 py-4 whitespace-nowrap text-sm text-green-600">
                  ${{ month.totalIncome.toLocaleString() }}
                </td>
                <td class="px-4 sm:px-6 py-4 whitespace-nowrap text-sm text-red-600">
                  ${{ month.totalExpenses.toLocaleString() }}
                </td>
                <td class="px-4 sm:px-6 py-4 whitespace-nowrap text-sm" :class="netAmountColor(month.netAmount)">
                  {{ month.netAmount >= 0 ? '+' : '' }}${{ month.netAmount.toLocaleString() }}
                </td>
                <td class="px-4 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-900 hidden sm:table-cell">
                  <div v-if="month.topCategory">
                    {{ month.topCategory.categoryName }}
                    <span class="text-gray-500">(${{ month.topCategory.amount.toLocaleString() }})</span>
                  </div>
                  <span v-else class="text-gray-400">-</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- No Data State -->
      <div v-if="report.monthlyData.length === 0" class="text-center py-12">
        <div class="text-gray-400 mb-4">
          <svg class="mx-auto h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-900 mb-2">No trend data found</h3>
        <p class="text-gray-600">Add some transactions to see your spending trends over time.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { financialService, type MonthlyTrendsReport } from '@/services/financialService'
import LineChart from '@/components/charts/LineChart.vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'

const report = ref<MonthlyTrendsReport | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)
const selectedMonths = ref(6)

const breadcrumbs = [
  { name: 'Reports', href: '/reports' },
  { name: 'Monthly Trends' }
]

const monthOptions = [
  { value: 3, label: '3 months' },
  { value: 6, label: '6 months' },
  { value: 9, label: '9 months' },
  { value: 12, label: '12 months' }
]

const chartData = computed(() => {
  if (!report.value || report.value.monthlyData.length === 0) return []
  
  return report.value.monthlyData
    .slice()
    .reverse() // Show oldest to newest for trend visualization
    .map(month => ({
      label: formatMonth(month.month),
      value: month.totalExpenses
    }))
})

const chartDataNet = computed(() => {
  if (!report.value || report.value.monthlyData.length === 0) return []
  
  return report.value.monthlyData
    .slice()
    .reverse() // Show oldest to newest for trend visualization
    .map(month => ({
      label: formatMonth(month.month),
      value: month.netAmount
    }))
})

const trendColor = computed(() => {
  if (!report.value) return 'text-gray-900'
  
  switch (report.value.trendDirection) {
    case 'increasing':
      return 'text-red-600'
    case 'decreasing':
      return 'text-green-600'
    case 'stable':
      return 'text-blue-600'
    default:
      return 'text-gray-900'
  }
})

const trendLabel = computed(() => {
  if (!report.value) return 'Unknown'
  
  switch (report.value.trendDirection) {
    case 'increasing':
      return 'Increasing'
    case 'decreasing':
      return 'Decreasing'
    case 'stable':
      return 'Stable'
    default:
      return 'Unknown'
  }
})

const trendIcon = computed(() => {
  if (!report.value) return 'div'
  
  switch (report.value.trendDirection) {
    case 'increasing':
      return 'svg' // Will be replaced with actual arrow up icon
    case 'decreasing':
      return 'svg' // Will be replaced with actual arrow down icon
    case 'stable':
      return 'svg' // Will be replaced with actual arrow right icon
    default:
      return 'div'
  }
})

const chartColor = computed(() => {
  if (!report.value) return '#3B82F6'
  
  switch (report.value.trendDirection) {
    case 'increasing':
      return '#EF4444' // red
    case 'decreasing':
      return '#10B981' // green
    case 'stable':
      return '#3B82F6' // blue
    default:
      return '#3B82F6'
  }
})

const netTrendColor = computed(() => {
  if (!report.value || report.value.monthlyData.length === 0) return '#3B82F6'
  
  // Calculate if net trend is positive or negative
  const recentMonths = report.value.monthlyData.slice(-3)
  const avgNet = recentMonths.reduce((sum, month) => sum + month.netAmount, 0) / recentMonths.length
  
  if (avgNet > 0) return '#10B981' // green for positive
  if (avgNet < 0) return '#EF4444' // red for negative
  return '#6B7280' // gray for neutral
})

const loadReport = async () => {
  loading.value = true
  error.value = null
  
  try {
    report.value = await financialService.getMonthlyTrendsReport(selectedMonths.value)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load report'
    console.error('Error loading monthly trends report:', err)
  } finally {
    loading.value = false
  }
}

const formatMonth = (monthStr: string) => {
  const date = new Date(monthStr + '-01') // Add day to make it a valid date
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long' })
}

const netAmountColor = (amount: number) => {
  if (amount > 0) return 'text-green-600 font-semibold'
  if (amount < 0) return 'text-red-600 font-semibold'
  return 'text-gray-900'
}

onMounted(() => {
  loadReport()
})
</script>