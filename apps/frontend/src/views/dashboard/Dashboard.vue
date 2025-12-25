<template>
  <MainLayout>
    <div class="px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto w-full">
      <div class="md:flex md:items-center md:justify-between mb-8">
        <div class="flex-1 min-w-0">
          <h2 class="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
            Welcome back, Demo Family! 👋
          </h2>
          <p class="mt-1 text-sm text-gray-500">
            Here's what's happening with your family budget this month.
          </p>
        </div>
        <div class="mt-4 flex md:mt-0 md:ml-4 gap-3">
          <button type="button" class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-xl shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-all">
            Download Report
          </button>
          <button type="button" class="inline-flex items-center px-4 py-2 border border-transparent rounded-xl shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-all shadow-md active:scale-95">
            Add Transaction
          </button>
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
                    <div class="text-2xl font-bold text-gray-900 leading-tight">$12,450.00</div>
                    <div class="ml-2 flex items-baseline text-sm font-semibold text-green-600">
                      <span class="sr-only">Increased by</span>
                      +2.5%
                    </div>
                  </dd>
                </dl>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 px-5 py-3 border-t border-gray-100">
            <div class="text-sm">
              <a href="#" class="font-medium text-primary-600 hover:text-primary-500">View detailed stats &rarr;</a>
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
                    <div class="text-2xl font-bold text-gray-900 leading-tight">$4,500.00</div>
                    <div class="text-xs text-gray-400 font-medium whitespace-nowrap">Plan: $5,000</div>
                  </dd>
                </dl>
                <!-- Budget Bar -->
                <div class="w-full bg-gray-100 h-2 rounded-full overflow-hidden mt-3">
                  <div class="bg-blue-600 h-full rounded-full transition-all duration-1000" style="width: 78%"></div>
                </div>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 px-5 py-3 border-t border-gray-100">
            <div class="text-sm text-gray-500">
              <span class="font-medium text-gray-900">78%</span> used of your monthly target
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
                  <dt class="text-sm font-medium text-gray-500 truncate">Monthly Spending</dt>
                  <dd class="flex items-baseline">
                    <div class="text-2xl font-bold text-red-600 leading-tight">$3,240.50</div>
                  </dd>
                </dl>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 px-5 py-3 border-t border-gray-100">
            <div class="text-sm">
              <span class="text-red-500 font-medium">↑ 12%</span> more than last month
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
          <div class="overflow-y-auto">
            <table class="min-w-full divide-y divide-gray-100">
              <thead class="bg-gray-50">
                <tr>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Transaction</th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Category</th>
                  <th scope="col" class="px-6 py-3 text-right text-xs font-bold text-gray-500 uppercase tracking-wider">Amount</th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-100">
                <tr v-for="tx in mockTransactions" :key="tx.id" class="hover:bg-gray-50 transition-colors pointer-cursor">
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="flex items-center">
                      <div class="flex-shrink-0 h-10 w-10 flex items-center justify-center rounded-xl bg-gray-50 text-xl shadow-inner border border-gray-100">
                        {{ tx.icon }}
                      </div>
                      <div class="ml-4">
                        <div class="text-sm font-bold text-gray-900">{{ tx.description }}</div>
                        <div class="text-xs text-gray-400">{{ tx.date }}</div>
                      </div>
                    </div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-500">
                    <span class="px-2.5 py-1 rounded-full text-xs font-bold bg-gray-100 text-gray-600">
                      {{ tx.category }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-bold" :class="tx.amount \u003c 0 ? 'text-red-600' : 'text-green-600'">
                    {{ tx.amount \u003c 0 ? '' : '+' }}{{ tx.formattedAmount }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Sidebar Widgets -->
        <div class="space-y-8">
          <!-- Categories Breakdown (Mini) -->
          <div class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
            <h3 class="text-lg font-bold text-gray-900 mb-6">Top Categories</h3>
            <div class="space-y-4">
              <div v-for="cat in mockCategories" :key="cat.name">
                <div class="flex items-center justify-between mb-1.5 px-0.5">
                  <div class="flex items-center">
                    <span class="mr-2 text-base">{{ cat.icon }}</span>
                    <span class="text-sm font-bold text-gray-700">{{ cat.name }}</span>
                  </div>
                  <span class="text-sm font-bold text-gray-900">${{ cat.spent }}</span>
                </div>
                <div class="w-full bg-gray-100 h-1.5 rounded-full overflow-hidden shadow-inner">
                  <div class="h-full rounded-full" :style="{ width: cat.percent + '%', backgroundColor: cat.color }"></div>
                </div>
              </div>
            </div>
            <button class="w-full mt-6 py-2.5 px-4 rounded-xl border border-gray-200 text-sm font-bold text-gray-600 hover:bg-gray-50 hover:text-gray-900 transition-all">
              Manage Categories
            </button>
          </div>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import MainLayout from '@/layouts/MainLayout.vue'
import { 
  WalletIcon, 
  BarChart3Icon, 
  ArrowDownIcon 
} from 'lucide-vue-next'

const mockTransactions = [
  { id: 1, description: 'Supermarket shopping', category: 'Groceries', amount: -125.50, formattedAmount: '-$125.50', date: 'Today, 10:45 AM', icon: '🛒' },
  { id: 2, description: 'School Books', category: 'Education', amount: -68.00, formattedAmount: '-$68.00', date: 'Yesterday', icon: '📚' },
  { id: 3, description: 'Salary Deposit', category: 'Income', amount: 5000.00, formattedAmount: '$5,000.00', date: 'Oct 15, 2024', icon: '💰' },
  { id: 4, description: 'Electricity Bill', category: 'Utilities', amount: -120.00, formattedAmount: '-$120.00', date: 'Oct 14, 2024', icon: '⚡' },
  { id: 5, description: 'Family Dinner', category: 'Dining Out', amount: -85.20, formattedAmount: '-$85.20', date: 'Oct 12, 2024', icon: '🍕' },
]

const mockCategories = [
  { name: 'Groceries', spent: 650, percent: 80, color: '#10b981', icon: '🛒' },
  { name: 'Dining Out', spent: 320, percent: 45, color: '#f59e0b', icon: '🍽️' },
  { name: 'Transport', spent: 210, percent: 30, color: '#3b82f6', icon: '🚗' },
  { name: 'Utilities', spent: 450, percent: 60, color: '#6366f1', icon: '⚡' }
]
</script>
