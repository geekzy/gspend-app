<template>
  <div class="p-6 max-w-6xl mx-auto">
    <!-- Breadcrumb Navigation -->
    <Breadcrumb :breadcrumbs="breadcrumbs" />
    
    <!-- Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900 mb-2">Budget vs Actual Report</h1>
      <p class="text-gray-600">Compare your planned budget with actual spending by category</p>
    </div>

    <!-- Month Selection Filter -->
    <div class="mb-6 bg-white rounded-lg shadow p-4">
      <div class="flex items-center gap-4">
        <label for="month-select" class="text-sm font-medium text-gray-700">
          Select Month:
        </label>
        <select
          id="month-select"
          v-model="selectedMonth"
          @change="loadReport"
          class="border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
        >
          <option v-for="month in availableMonths" :key="month.value" :value="month.value">
            {{ month.label }}
          </option>
        </select>
        <button
          @click="loadReport"
          :disabled="loading"
          class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-sm"
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
          <div class="text-sm font-medium text-gray-500">Total Budgeted</div>
          <div class="text-2xl font-bold text-gray-900">${{ report.totalBudgeted.toLocaleString() }}</div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="text-sm font-medium text-gray-500">Total Spent</div>
          <div class="text-2xl font-bold text-gray-900">${{ report.totalSpent.toLocaleString() }}</div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="text-sm font-medium text-gray-500">Overall Variance</div>
          <div class="text-2xl font-bold" :class="varianceColor(report.overallVariance)">
            {{ report.overallVariance >= 0 ? '+' : '' }}${{ report.overallVariance.toLocaleString() }}
          </div>
        </div>
      </div>

      <!-- Budget vs Actual Chart -->
      <div class="bg-white rounded-lg shadow p-4 sm:p-6">
        <h3 class="text-lg font-medium text-gray-900 mb-4">Budget vs Actual Comparison</h3>
        <div v-if="chartData.length > 0" class="h-64 sm:h-80 lg:h-96">
          <BarChart 
            :data="chartData" 
            :width="800" 
            :height="400"
            type="comparison"
            :colors="['#10B981', '#EF4444']"
            :comparison-labels="{ primary: 'Budgeted', secondary: 'Actual' }"
          />
        </div>
        <div v-else class="text-center py-8 text-gray-500">
          No budget data available for comparison
        </div>
      </div>

      <!-- Category Breakdown Table -->
      <div class="bg-white rounded-lg shadow overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-medium text-gray-900">Category Breakdown</h3>
        </div>
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Category
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Budgeted
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actual
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Variance
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  % Used
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Progress
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="category in report.categories" :key="category.categoryName">
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  {{ category.categoryName }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  ${{ category.budgeted.toLocaleString() }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  ${{ category.actual.toLocaleString() }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm" :class="varianceColor(category.variance)">
                  {{ category.variance >= 0 ? '+' : '' }}${{ category.variance.toLocaleString() }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm" :class="percentageColor(category.percentageUsed)">
                  {{ category.percentageUsed.toFixed(1) }}%
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="w-full bg-gray-200 rounded-full h-2">
                    <div
                      class="h-2 rounded-full transition-all duration-300"
                      :class="progressBarColor(category.percentageUsed)"
                      :style="{ width: Math.min(category.percentageUsed, 100) + '%' }"
                    ></div>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- No Data State -->
      <div v-if="report.categories.length === 0" class="text-center py-12">
        <div class="text-gray-400 mb-4">
          <svg class="mx-auto h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-900 mb-2">No budget data found</h3>
        <p class="text-gray-600">Create a budget for this month to see the comparison report.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { financialService, type BudgetVsActualReport } from '@/services/financialService'
import BarChart from '@/components/charts/BarChart.vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'

const report = ref<BudgetVsActualReport | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)
const selectedMonth = ref('')

const breadcrumbs = [
  { name: 'Reports', href: '/reports' },
  { name: 'Budget vs Actual' }
]

// Generate available months (current month and previous 11 months)
const availableMonths = computed(() => {
  const months = []
  const now = new Date()
  
  for (let i = 0; i < 12; i++) {
    const date = new Date(now.getFullYear(), now.getMonth() - i, 1)
    const value = date.toISOString().substring(0, 7) // YYYY-MM format
    const label = date.toLocaleDateString('en-US', { year: 'numeric', month: 'long' })
    months.push({ value, label })
  }
  
  return months
})

const chartData = computed(() => {
  if (!report.value || report.value.categories.length === 0) return []
  
  return report.value.categories.map(category => ({
    label: category.categoryName,
    budgeted: category.budgeted,
    actual: category.actual
  }))
})

const loadReport = async () => {
  if (!selectedMonth.value) return
  
  loading.value = true
  error.value = null
  
  try {
    report.value = await financialService.getBudgetVsActualReport(selectedMonth.value)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load report'
    console.error('Error loading budget vs actual report:', err)
  } finally {
    loading.value = false
  }
}

const varianceColor = (variance: number) => {
  if (variance > 0) return 'text-red-600' // Over budget
  if (variance < 0) return 'text-green-600' // Under budget
  return 'text-gray-900' // Exactly on budget
}

const percentageColor = (percentage: number) => {
  if (percentage > 100) return 'text-red-600 font-semibold'
  if (percentage > 80) return 'text-yellow-600 font-semibold'
  return 'text-gray-900'
}

const progressBarColor = (percentage: number) => {
  if (percentage > 100) return 'bg-red-500'
  if (percentage > 80) return 'bg-yellow-500'
  return 'bg-green-500'
}

onMounted(() => {
  // Set current month as default
  selectedMonth.value = availableMonths.value[0].value
  loadReport()
})
</script>