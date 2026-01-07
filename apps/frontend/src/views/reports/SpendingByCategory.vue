<template>
  <div class="p-6 max-w-6xl mx-auto">
    <!-- Breadcrumb Navigation -->
    <Breadcrumb :breadcrumbs="breadcrumbs" />
    
    <!-- Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900 mb-2">Spending by Category</h1>
      <p class="text-gray-600">See how your money is distributed across different spending categories</p>
    </div>

    <!-- Date Range Filter -->
    <div class="mb-6 bg-white rounded-lg shadow p-4">
      <div class="flex flex-col gap-4">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div class="flex items-center gap-2">
            <label for="start-date" class="text-sm font-medium text-gray-700 whitespace-nowrap">From:</label>
            <input
              id="start-date"
              type="date"
              v-model="startDate"
              class="border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 flex-1"
            />
          </div>
          <div class="flex items-center gap-2">
            <label for="end-date" class="text-sm font-medium text-gray-700 whitespace-nowrap">To:</label>
            <input
              id="end-date"
              type="date"
              v-model="endDate"
              class="border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 flex-1"
            />
          </div>
        </div>
        
        <!-- Quick Filter Buttons -->
        <div class="flex flex-wrap gap-2">
          <button
            @click="setQuickFilter('thisMonth')"
            class="px-3 py-2 text-sm bg-gray-100 hover:bg-gray-200 rounded-md transition-colors"
          >
            This Month
          </button>
          <button
            @click="setQuickFilter('lastMonth')"
            class="px-3 py-2 text-sm bg-gray-100 hover:bg-gray-200 rounded-md transition-colors"
          >
            Last Month
          </button>
          <button
            @click="setQuickFilter('last3Months')"
            class="px-3 py-2 text-sm bg-gray-100 hover:bg-gray-200 rounded-md transition-colors"
          >
            Last 3 Months
          </button>
          
          <button
            @click="loadReport"
            :disabled="loading"
            class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-sm ml-auto"
          >
            <span v-if="loading">Loading...</span>
            <span v-else>Apply Filter</span>
          </button>
        </div>
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
      <!-- Summary -->
      <div class="bg-white rounded-lg shadow p-6">
        <div class="text-center">
          <div class="text-sm font-medium text-gray-500 mb-1">Total Spending</div>
          <div class="text-3xl font-bold text-gray-900">{{ formatCurrency(report.totalSpent) }}</div>
          <div class="text-sm text-gray-600 mt-1">
            {{ formatDateRange(report.startDate, report.endDate) }}
          </div>
        </div>
      </div>

      <!-- Chart and Breakdown -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Pie Chart -->
        <div class="bg-white rounded-lg shadow p-4 sm:p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Spending Distribution</h3>
          <div v-if="chartData.length > 0" class="flex justify-center">
            <div class="w-full max-w-sm h-64 sm:h-80">
              <PieChart 
                :data="chartData" 
                :width="300" 
                :height="300"
                :show-legend="false"
              />
            </div>
          </div>
          <div v-else class="text-center py-8 text-gray-500">
            No spending data available for the selected period
          </div>
        </div>

        <!-- Category List -->
        <div class="bg-white rounded-lg shadow p-4 sm:p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Category Breakdown</h3>
          <div class="space-y-3 max-h-80 overflow-y-auto">
            <div 
              v-for="(category, index) in report.categories" 
              :key="category.categoryId"
              class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div class="flex items-center min-w-0 flex-1">
                <div 
                  class="w-4 h-4 rounded-full mr-3 flex-shrink-0"
                  :style="{ backgroundColor: chartColors[index % chartColors.length] }"
                ></div>
                <div class="min-w-0 flex-1">
                  <div class="font-medium text-gray-900 truncate">{{ category.categoryName }}</div>
                  <div class="text-sm text-gray-600">{{ category.percentage.toFixed(1) }}% of total</div>
                </div>
              </div>
              <div class="text-right ml-4 flex-shrink-0">
                <div class="font-semibold text-gray-900 text-sm sm:text-base">{{ formatCurrency(category.amount) }}</div>
              </div>
            </div>
          </div>
          
          <!-- No Data State -->
          <div v-if="report.categories.length === 0" class="text-center py-8">
            <div class="text-gray-400 mb-4">
              <svg class="mx-auto h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <h3 class="text-lg font-medium text-gray-900 mb-2">No spending data</h3>
            <p class="text-gray-600">No transactions found for the selected date range.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { financialService, type SpendingByCategoryReport } from '@/services/financialService'
import PieChart from '@/components/charts/PieChart.vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'
import { formatCurrency } from '@/utils/currency'

const report = ref<SpendingByCategoryReport | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)
const startDate = ref('')
const endDate = ref('')

const breadcrumbs = [
  { name: 'Reports', href: '/reports' },
  { name: 'Spending by Category' }
]

const chartColors = [
  '#3B82F6', // blue
  '#EF4444', // red
  '#10B981', // green
  '#F59E0B', // yellow
  '#8B5CF6', // purple
  '#EC4899', // pink
  '#06B6D4', // cyan
  '#84CC16', // lime
  '#F97316', // orange
  '#6366F1', // indigo
]

const chartData = computed(() => {
  if (!report.value || report.value.categories.length === 0) return []
  
  return report.value.categories.map(category => ({
    label: category.categoryName,
    value: category.amount
  }))
})

const setQuickFilter = (period: 'thisMonth' | 'lastMonth' | 'last3Months') => {
  const now = new Date()
  const currentYear = now.getFullYear()
  const currentMonth = now.getMonth()
  
  switch (period) {
    case 'thisMonth':
      startDate.value = new Date(currentYear, currentMonth, 1).toISOString().split('T')[0]
      endDate.value = new Date(currentYear, currentMonth + 1, 0).toISOString().split('T')[0]
      break
    case 'lastMonth':
      startDate.value = new Date(currentYear, currentMonth - 1, 1).toISOString().split('T')[0]
      endDate.value = new Date(currentYear, currentMonth, 0).toISOString().split('T')[0]
      break
    case 'last3Months':
      startDate.value = new Date(currentYear, currentMonth - 2, 1).toISOString().split('T')[0]
      endDate.value = new Date(currentYear, currentMonth + 1, 0).toISOString().split('T')[0]
      break
  }
  
  loadReport()
}

const loadReport = async () => {
  if (!startDate.value || !endDate.value) {
    error.value = 'Please select both start and end dates'
    return
  }
  
  loading.value = true
  error.value = null
  
  try {
    report.value = await financialService.getSpendingByCategoryReport(startDate.value, endDate.value)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load report'
    console.error('Error loading spending by category report:', err)
  } finally {
    loading.value = false
  }
}

const formatDateRange = (start: string, end: string) => {
  const startDate = new Date(start)
  const endDate = new Date(end)
  
  const options: Intl.DateTimeFormatOptions = { 
    year: 'numeric', 
    month: 'short', 
    day: 'numeric' 
  }
  
  return `${startDate.toLocaleDateString('en-US', options)} - ${endDate.toLocaleDateString('en-US', options)}`
}

onMounted(() => {
  // Set default to current month
  setQuickFilter('thisMonth')
})
</script>