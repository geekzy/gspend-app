<template>
  <div class="p-6 max-w-4xl mx-auto">
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900 mb-2">Financial Reports</h1>
      <p class="text-gray-600">Analyze your family's spending patterns and financial trends</p>
    </div>

    <!-- Report Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <!-- Budget vs Actual Report -->
      <div class="bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow duration-200">
        <div class="p-6">
          <div class="flex items-center mb-4">
            <div class="bg-blue-100 rounded-lg p-3 mr-4">
              <svg class="h-6 w-6 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-gray-900">Budget vs Actual</h3>
              <p class="text-sm text-gray-600">Compare planned vs actual spending</p>
            </div>
          </div>
          <p class="text-gray-700 mb-4">
            See how your actual spending compares to your budget by category. Track variance and percentage used for each budget item.
          </p>
          <router-link
            to="/reports/budget-vs-actual"
            class="inline-flex items-center px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 transition-colors"
          >
            View Report
            <svg class="ml-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </router-link>
        </div>
      </div>

      <!-- Spending by Category Report -->
      <div class="bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow duration-200">
        <div class="p-6">
          <div class="flex items-center mb-4">
            <div class="bg-green-100 rounded-lg p-3 mr-4">
              <svg class="h-6 w-6 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 3.055A9.001 9.001 0 1020.945 13H11V3.055z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.488 9H15V3.512A9.025 9.025 0 0120.488 9z" />
              </svg>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-gray-900">Spending by Category</h3>
              <p class="text-sm text-gray-600">Breakdown of spending distribution</p>
            </div>
          </div>
          <p class="text-gray-700 mb-4">
            Visualize how your money is distributed across different spending categories with interactive pie charts and detailed breakdowns.
          </p>
          <router-link
            to="/reports/spending-by-category"
            class="inline-flex items-center px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 transition-colors"
          >
            View Report
            <svg class="ml-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </router-link>
        </div>
      </div>

      <!-- Monthly Trends Report -->
      <div class="bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow duration-200">
        <div class="p-6">
          <div class="flex items-center mb-4">
            <div class="bg-purple-100 rounded-lg p-3 mr-4">
              <svg class="h-6 w-6 text-purple-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z" />
              </svg>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-gray-900">Monthly Trends</h3>
              <p class="text-sm text-gray-600">Track spending patterns over time</p>
            </div>
          </div>
          <p class="text-gray-700 mb-4">
            Analyze your spending trends over multiple months to identify patterns and make informed financial decisions.
          </p>
          <router-link
            to="/reports/monthly-trends"
            class="inline-flex items-center px-4 py-2 bg-purple-600 text-white text-sm font-medium rounded-md hover:bg-purple-700 transition-colors"
          >
            View Report
            <svg class="ml-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </router-link>
        </div>
      </div>
    </div>

    <!-- Quick Stats Section -->
    <div class="mt-12 bg-gray-50 rounded-lg p-6">
      <h2 class="text-xl font-semibold text-gray-900 mb-4">Quick Insights</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div class="text-center">
          <div class="text-2xl font-bold text-blue-600">3</div>
          <div class="text-sm text-gray-600">Available Reports</div>
        </div>
        <div class="text-center">
          <div class="text-2xl font-bold text-green-600">Real-time</div>
          <div class="text-sm text-gray-600">Data Updates</div>
        </div>
        <div class="text-center">
          <div class="text-2xl font-bold text-purple-600">Interactive</div>
          <div class="text-sm text-gray-600">Charts & Filters</div>
        </div>
      </div>
    </div>

    <!-- Help Section -->
    <div class="mt-8 bg-blue-50 border border-blue-200 rounded-lg p-6">
      <div class="flex">
        <div class="text-blue-400">
          <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
          </svg>
        </div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-blue-800">Getting Started with Reports</h3>
          <div class="text-sm text-blue-700 mt-1">
            <p>Reports are generated from your transaction and budget data. Make sure you have:</p>
            <ul class="list-disc list-inside mt-2 space-y-1">
              <li>Added some transactions to see spending patterns</li>
              <li>Created budgets to compare against actual spending</li>
              <li>Categorized your transactions for better insights</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// No additional logic needed for this index page
</script>