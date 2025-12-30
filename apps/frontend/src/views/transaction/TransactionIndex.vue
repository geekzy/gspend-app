<template>
  <MainLayout>
    <div class="px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto w-full">
      <div class="md:flex md:items-center md:justify-between mb-8">
        <div class="flex-1 min-w-0">
          <h2 class="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
            Transactions
          </h2>
          <p class="mt-1 text-sm text-gray-500">
            Keep track of every dollar in and out.
          </p>
        </div>
        <div class="mt-4 flex md:mt-0 md:ml-4">
          <button 
            @click="openAddModal"
            type="button" 
            class="inline-flex items-center px-4 py-2.5 border border-transparent rounded-2xl shadow-lg text-sm font-bold text-white bg-primary-600 hover:bg-primary-700 transition-all active:scale-95"
          >
            Add Transaction
          </button>
        </div>
      </div>

      <!-- Filters Section -->
      <div class="grid grid-cols-1 lg:grid-cols-4 gap-6 mb-6">
        <!-- Date Range Filter -->
        <div class="lg:col-span-1">
          <DateRangeFilter 
            v-model="dateFilter"
            @change="onDateFilterChange"
          />
        </div>
        
        <!-- Category Filter -->
        <div class="lg:col-span-1">
          <CategoryFilter 
            :categories="categories"
            v-model="categoryFilter"
            @change="onCategoryFilterChange"
          />
        </div>
        
        <!-- Search and Quick Filters -->
        <div class="lg:col-span-2 space-y-4">
          <!-- Search -->
          <div class="bg-white rounded-2xl border border-gray-200 p-4">
            <label class="block text-sm font-bold text-gray-700 mb-2">Search</label>
            <div class="relative">
              <SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <input 
                v-model="search" 
                type="text" 
                placeholder="Search description..." 
                class="w-full pl-10 pr-4 py-2 bg-gray-50 border border-gray-200 rounded-xl text-sm outline-none focus:ring-2 focus:ring-primary-500"
                @input="onSearchChange"
              />
            </div>
          </div>
          
          <!-- Applied Filters -->
          <div v-if="hasActiveFilters" class="bg-white rounded-2xl border border-gray-200 p-4">
            <div class="flex items-center justify-between mb-2">
              <span class="text-sm font-bold text-gray-700">Active Filters</span>
              <button 
                @click="clearAllFilters"
                class="text-xs text-red-500 hover:text-red-600 transition-colors"
              >
                Clear All
              </button>
            </div>
            <div class="flex flex-wrap gap-2">
              <span
                v-if="dateFilter.startDate && dateFilter.endDate"
                class="inline-flex items-center space-x-1 px-2 py-1 bg-blue-100 text-blue-700 text-xs rounded-lg"
              >
                <span>📅</span>
                <span>{{ formatDateRange(dateFilter.startDate, dateFilter.endDate) }}</span>
              </span>
              
              <span
                v-if="categoryFilter.type !== 'all'"
                class="inline-flex items-center space-x-1 px-2 py-1 bg-purple-100 text-purple-700 text-xs rounded-lg"
              >
                <span>🏷️</span>
                <span>{{ categoryFilter.type === 'income' ? 'Income' : 'Expense' }}</span>
              </span>
              
              <span
                v-for="categoryId in categoryFilter.categoryIds.slice(0, 2)"
                :key="categoryId"
                class="inline-flex items-center space-x-1 px-2 py-1 bg-green-100 text-green-700 text-xs rounded-lg"
              >
                <span>{{ getCategoryById(categoryId)?.icon }}</span>
                <span>{{ getCategoryById(categoryId)?.name }}</span>
              </span>
              
              <span
                v-if="categoryFilter.categoryIds.length > 2"
                class="inline-flex items-center px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded-lg"
              >
                +{{ categoryFilter.categoryIds.length - 2 }} more categories
              </span>
              
              <span
                v-if="search"
                class="inline-flex items-center space-x-1 px-2 py-1 bg-yellow-100 text-yellow-700 text-xs rounded-lg"
              >
                <span>🔍</span>
                <span>"{{ search }}"</span>
              </span>
            </div>
          </div>
        </div>
      </div>

      <div v-if="isLoading" class="flex items-center justify-center min-h-[400px]">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>

      <div v-else class="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden">
        <!-- Results Summary -->
        <div class="p-6 border-b border-gray-100 flex flex-col sm:flex-row justify-between items-center gap-4">
          <div class="text-sm text-gray-600">
            Showing {{ (currentPage - 1) * perPage + 1 }} to {{ Math.min(currentPage * perPage, totalTransactions) }} of {{ totalTransactions }} transactions
            <span v-if="hasActiveFilters" class="text-primary-600 font-medium">(filtered)</span>
          </div>
          <div class="flex items-center gap-2">
            <label class="text-sm text-gray-600">Sort by:</label>
            <select 
              v-model="sortBy" 
              @change="onSortChange"
              class="px-3 py-1 bg-gray-50 border border-gray-200 rounded-lg text-sm outline-none"
            >
              <option value="transactionDate">Date</option>
              <option value="amount">Amount</option>
              <option value="description">Description</option>
            </select>
            <button
              @click="toggleSortOrder"
              class="p-1 text-gray-400 hover:text-gray-600 transition-colors"
              aria-label="Toggle sort order"
            >
              <ArrowUpDownIcon class="w-4 h-4" />
            </button>
          </div>
        </div>

        <!-- Transactions Table -->
        <div class="overflow-x-auto">
          <div v-if="transactions.length === 0" class="p-20 text-center text-gray-400">
            <HistoryIcon class="w-12 h-12 mx-auto mb-4 opacity-20" />
            <p v-if="hasActiveFilters">No transactions found matching your criteria.</p>
            <p v-else>No transactions found. Add your first transaction to get started!</p>
          </div>
          <table v-else class="min-w-full divide-y divide-gray-100">
            <thead class="bg-gray-50">
              <tr>
                <th scope="col" class="px-6 py-4 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Date</th>
                <th scope="col" class="px-6 py-4 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Description</th>
                <th scope="col" class="px-6 py-4 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Category</th>
                <th scope="col" class="px-6 py-4 text-right text-xs font-bold text-gray-500 uppercase tracking-wider">Amount</th>
                <th scope="col" class="px-6 py-4 text-center text-xs font-bold text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-100">
              <tr v-for="tx in transactions" :key="tx.id" class="hover:bg-gray-50 transition-colors">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {{ formatDate(tx.transactionDate) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="text-sm font-bold text-gray-900">{{ tx.description }}</div>
                  <div class="text-xs text-gray-400 font-medium">{{ tx.paymentMethod }}</div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="inline-flex items-center space-x-1 px-2.5 py-0.5 rounded-full text-xs font-bold bg-gray-100 text-gray-600">
                    <span>{{ getCategoryById(tx.categoryId)?.icon }}</span>
                    <span>{{ tx.categoryName }}</span>
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-bold">
                  <span :class="tx.type === 'expense' ? 'text-red-500' : 'text-green-600'">
                    {{ tx.type === 'expense' ? '-' : '+' }}${{ tx.amount.toLocaleString(undefined, { minimumFractionDigits: 2 }) }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-center">
                  <div class="flex items-center justify-center space-x-2">
                    <button 
                      @click="openEditModal(tx)"
                      class="p-1.5 text-gray-400 hover:text-blue-500 transition-colors"
                      title="Edit transaction"
                      aria-label="Edit transaction"
                    >
                      <EditIcon class="w-4 h-4" />
                    </button>
                    <button 
                      @click="deleteTransaction(tx.id)"
                      class="p-1.5 text-gray-400 hover:text-red-500 transition-colors"
                      title="Delete transaction"
                      aria-label="Delete transaction"
                    >
                      <Trash2Icon class="w-4 h-4" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <Pagination
          v-if="totalPages > 1"
          :current-page="currentPage"
          :total-pages="totalPages"
          :total="totalTransactions"
          :per-page="perPage"
          @page-change="onPageChange"
          @per-page-change="onPerPageChange"
        />
      </div>

      <!-- Add Transaction Modal -->
      <div v-if="showAddModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm overflow-y-auto">
        <div class="bg-white rounded-3xl w-full max-w-lg p-8 shadow-2xl animate-in fade-in zoom-in duration-200">
          <div class="flex justify-between items-center mb-6">
            <h3 class="text-xl font-black text-gray-900">New Transaction</h3>
            <button @click="showAddModal = false" class="p-2 text-gray-400 hover:text-gray-600" aria-label="Close modal">
              <XIcon class="w-6 h-6" />
            </button>
          </div>

          <form @submit.prevent="handleAddTransaction" class="space-y-4">
            <div class="grid grid-cols-2 gap-3">
              <button 
                type="button"
                @click="newTransaction.type = 'expense'"
                :class="[
                  newTransaction.type === 'expense' ? 'bg-red-500 text-white shadow-lg shadow-red-200' : 'bg-gray-50 text-gray-500',
                  'py-3 rounded-2xl text-sm font-bold transition-all'
                ]"
                role="radio"
                :aria-checked="newTransaction.type === 'expense'"
              >
                Expense
              </button>
              <button 
                type="button"
                @click="newTransaction.type = 'income'"
                :class="[
                  newTransaction.type === 'income' ? 'bg-green-500 text-white shadow-lg shadow-green-200' : 'bg-gray-50 text-gray-500',
                  'py-3 rounded-2xl text-sm font-bold transition-all'
                ]"
                role="radio"
                :aria-checked="newTransaction.type === 'income'"
              >
                Income
              </button>
            </div>

            <div>
              <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Description</label>
              <input v-model="newTransaction.description" type="text" class="w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none" placeholder="What was this for?" required />
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Amount</label>
                <div class="relative">
                  <span class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400">$</span>
                  <input v-model.number="newTransaction.amount" type="number" step="0.01" class="w-full pl-8 pr-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none" placeholder="0.00" required />
                </div>
              </div>
              <div>
                <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Date</label>
                <input v-model="newTransaction.transactionDate" type="date" class="w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none" required />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Category</label>
                <select v-model="newTransaction.categoryId" class="w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none bg-white" required>
                  <option value="">Select Category</option>
                  <option v-for="cat in availableCategories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Payment Method</label>
                <select v-model="newTransaction.paymentMethod" class="w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none bg-white">
                  <option value="cash">Cash</option>
                  <option value="bank_transfer">Bank Transfer</option>
                  <option value="credit_card">Credit Card</option>
                  <option value="debit_card">Debit Card</option>
                </select>
              </div>
            </div>

            <div>
              <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Notes (Optional)</label>
              <textarea v-model="newTransaction.notes" rows="2" class="w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none resize-none" placeholder="Any extra details..."></textarea>
            </div>

            <div class="flex gap-3 mt-8">
              <button @click="showAddModal = false" type="button" class="flex-1 py-3.5 rounded-2xl border border-gray-200 text-sm font-bold text-gray-700 hover:bg-gray-50 transition-all">Cancel</button>
              <button type="submit" :disabled="isSubmitting" class="flex-1 py-3.5 rounded-2xl bg-primary-600 text-white text-sm font-black hover:bg-primary-700 shadow-lg shadow-primary-200 transition-all disabled:opacity-50">
                {{ isSubmitting ? 'Recording...' : 'Record Transaction' }}
              </button>
            </div>
          </form>
        </div>
      </div>
      <!-- Edit Transaction Modal -->
      <TransactionEdit
        v-if="showEditModal && editingTransaction"
        :transaction="editingTransaction"
        :categories="categories"
        @close="closeEditModal"
        @saved="onTransactionSaved"
      />
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import DateRangeFilter, { type DateRange } from '@/components/filters/DateRangeFilter.vue'
import CategoryFilter, { type CategoryFilter as CategoryFilterType } from '@/components/filters/CategoryFilter.vue'
import Pagination from '@/components/common/Pagination.vue'
import TransactionEdit from '@/components/forms/TransactionEdit.vue'
import { 
  financialService, 
  Transaction, 
  Category, 
  TransactionResponse 
} from '@/services/financialService'
import { 
  HistoryIcon, 
  SearchIcon, 
  Trash2Icon, 
  XIcon,
  ArrowUpDownIcon,
  EditIcon
} from 'lucide-vue-next'

const transactions = ref<Transaction[]>([])
const categories = ref<Category[]>([])
const totalTransactions = ref(0)
const isLoading = ref(true)
const isSubmitting = ref(false)
const showAddModal = ref(false)
const showEditModal = ref(false)
const editingTransaction = ref<Transaction | null>(null)

// Filter state
const search = ref('')
const dateFilter = ref<DateRange>({ startDate: null, endDate: null })
const categoryFilter = ref<CategoryFilterType>({ type: 'all', categoryIds: [] })
const sortBy = ref('transactionDate')
const sortOrder = ref<'asc' | 'desc'>('desc')

// Pagination state
const currentPage = ref(1)
const perPage = ref(20)
const totalPages = computed(() => Math.ceil(totalTransactions.value / perPage.value))

const newTransaction = ref({
  type: 'expense' as 'income' | 'expense',
  description: '',
  amount: 0,
  transactionDate: new Date().toISOString().split('T')[0],
  categoryId: '',
  paymentMethod: 'bank_transfer',
  notes: ''
})

const availableCategories = computed(() => {
  return categories.value.filter(c => c.type === newTransaction.value.type)
})

const hasActiveFilters = computed(() => {
  return !!(
    search.value ||
    (dateFilter.value.startDate && dateFilter.value.endDate) ||
    categoryFilter.value.type !== 'all' ||
    categoryFilter.value.categoryIds.length > 0
  )
})

// Helper functions
const getCategoryById = (id: string): Category | undefined => {
  return categories.value.find(cat => cat.id === id)
}

const formatDateRange = (startDate: string, endDate: string): string => {
  const start = new Date(startDate).toLocaleDateString('en-US', { 
    month: 'short', 
    day: 'numeric' 
  })
  const end = new Date(endDate).toLocaleDateString('en-US', { 
    month: 'short', 
    day: 'numeric' 
  })
  return `${start} - ${end}`
}

const fetchTransactions = async () => {
  try {
    isLoading.value = true
    
    // For now, we'll fetch all and filter client-side since backend may not support all filters yet
    const response = await financialService.getTransactions()
    
    // Handle both old format (array) and new format (object with pagination)
    let allTransactions: Transaction[]
    if (Array.isArray(response)) {
      allTransactions = response
    } else {
      allTransactions = (response as TransactionResponse).transactions
    }
    
    // Apply client-side filtering as fallback
    let filtered = allTransactions

    // Search filter
    if (search.value) {
      const searchLower = search.value.toLowerCase()
      filtered = filtered.filter(tx => 
        tx.description.toLowerCase().includes(searchLower) ||
        tx.categoryName.toLowerCase().includes(searchLower)
      )
    }

    // Date filter
    if (dateFilter.value.startDate && dateFilter.value.endDate) {
      const startDate = new Date(dateFilter.value.startDate)
      const endDate = new Date(dateFilter.value.endDate)
      filtered = filtered.filter(tx => {
        const txDate = new Date(tx.transactionDate)
        return txDate >= startDate && txDate <= endDate
      })
    }

    // Type filter
    if (categoryFilter.value.type !== 'all') {
      filtered = filtered.filter(tx => tx.type === categoryFilter.value.type)
    }

    // Category filter
    if (categoryFilter.value.categoryIds.length > 0) {
      filtered = filtered.filter(tx => 
        categoryFilter.value.categoryIds.includes(tx.categoryId)
      )
    }

    // Update total count after filtering
    totalTransactions.value = filtered.length

    // Sort
    filtered.sort((a, b) => {
      let aVal: any, bVal: any
      
      switch (sortBy.value) {
        case 'amount':
          aVal = a.amount
          bVal = b.amount
          break
        case 'description':
          aVal = a.description.toLowerCase()
          bVal = b.description.toLowerCase()
          break
        case 'transactionDate':
        default:
          aVal = new Date(a.transactionDate)
          bVal = new Date(b.transactionDate)
          break
      }

      if (sortOrder.value === 'asc') {
        return aVal > bVal ? 1 : -1
      } else {
        return aVal < bVal ? 1 : -1
      }
    })

    // Apply pagination
    const startIndex = (currentPage.value - 1) * perPage.value
    const endIndex = startIndex + perPage.value
    transactions.value = filtered.slice(startIndex, endIndex)
  } catch (err) {
    console.error('Failed to fetch transaction data:', err)
  } finally {
    isLoading.value = false
  }
}

const fetchCategories = async () => {
  try {
    const catData = await financialService.getCategories()
    categories.value = catData
  } catch (err) {
    console.error('Failed to fetch categories:', err)
  }
}

const fetchInitialData = async () => {
  await Promise.all([
    fetchTransactions(),
    fetchCategories()
  ])
}

// Event handlers
const onDateFilterChange = (newDateFilter: DateRange) => {
  dateFilter.value = newDateFilter
  currentPage.value = 1 // Reset to first page
  fetchTransactions()
}

const onCategoryFilterChange = (newCategoryFilter: CategoryFilterType) => {
  categoryFilter.value = newCategoryFilter
  currentPage.value = 1 // Reset to first page
  fetchTransactions()
}

const onSearchChange = () => {
  currentPage.value = 1 // Reset to first page
  fetchTransactions()
}

const onSortChange = () => {
  fetchTransactions()
}

const toggleSortOrder = () => {
  sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  fetchTransactions()
}

const clearAllFilters = () => {
  search.value = ''
  dateFilter.value = { startDate: null, endDate: null }
  categoryFilter.value = { type: 'all', categoryIds: [] }
  currentPage.value = 1
  fetchTransactions()
}

const onPageChange = (page: number) => {
  currentPage.value = page
  fetchTransactions()
}

const onPerPageChange = (newPerPage: number) => {
  perPage.value = newPerPage
  currentPage.value = 1 // Reset to first page when changing page size
  fetchTransactions()
}

const openAddModal = () => {
  newTransaction.value = {
    type: 'expense',
    description: '',
    amount: 0,
    transactionDate: new Date().toISOString().split('T')[0],
    categoryId: '',
    paymentMethod: 'bank_transfer',
    notes: ''
  }
  showAddModal.value = true
}

const handleAddTransaction = async () => {
  try {
    isSubmitting.value = true
    await financialService.createTransaction(newTransaction.value as any)
    showAddModal.value = false
    await fetchTransactions()
  } catch (err) {
    console.error('Failed to create transaction:', err)
  } finally {
    isSubmitting.value = false
  }
}

const openEditModal = (transaction: Transaction) => {
  editingTransaction.value = transaction
  showEditModal.value = true
}

const closeEditModal = () => {
  showEditModal.value = false
  editingTransaction.value = null
}

const onTransactionSaved = async () => {
  showEditModal.value = false
  editingTransaction.value = null
  await fetchTransactions()
}

const deleteTransaction = async (id: string) => {
  if (!confirm('Are you sure you want to delete this transaction?')) {
    return
  }
  
  try {
    await financialService.deleteTransaction(id)
    await fetchTransactions()
  } catch (err) {
    console.error('Failed to delete transaction:', err)
  }
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString(undefined, { 
    month: 'short', 
    day: 'numeric', 
    year: 'numeric' 
  })
}

// Debounce search to avoid too many API calls
let searchTimeout: number
watch(search, () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    onSearchChange()
  }, 300) as unknown as number
})

onMounted(fetchInitialData)
</script>
