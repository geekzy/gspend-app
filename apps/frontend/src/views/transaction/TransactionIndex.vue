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

      <div v-if="isLoading" class="flex items-center justify-center min-h-[400px]">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>

      <div v-else class="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden">
        <!-- Table Header/Filters -->
        <div class="p-6 border-b border-gray-100 flex flex-col sm:flex-row justify-between items-center gap-4">
          <div class="relative w-full sm:w-64">
            <SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
            <input 
              v-model="search" 
              type="text" 
              placeholder="Search description..." 
              class="w-full pl-10 pr-4 py-2 bg-gray-50 border border-gray-200 rounded-xl text-sm outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>
          <div class="flex gap-2 w-full sm:w-auto">
            <select v-model="filterType" class="flex-1 sm:flex-none px-3 py-2 bg-gray-50 border border-gray-200 rounded-xl text-sm outline-none">
              <option value="all">All Types</option>
              <option value="expense">Expense</option>
              <option value="income">Income</option>
            </select>
          </div>
        </div>

        <!-- Transactions Table -->
        <div class="overflow-x-auto">
          <div v-if="filteredTransactions.length === 0" class="p-20 text-center text-gray-400">
            <HistoryIcon class="w-12 h-12 mx-auto mb-4 opacity-20" />
            <p>No transactions found matching your criteria.</p>
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
              <tr v-for="tx in filteredTransactions" :key="tx.id" class="hover:bg-gray-50 transition-colors">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {{ formatDate(tx.transactionDate) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="text-sm font-bold text-gray-900">{{ tx.description }}</div>
                  <div class="text-xs text-gray-400 font-medium">{{ tx.paymentMethod }}</div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-bold bg-gray-100 text-gray-600">
                    {{ tx.categoryName }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-bold">
                  <span :class="tx.type === 'expense' ? 'text-red-500' : 'text-green-600'">
                    {{ tx.type === 'expense' ? '-' : '+' }}${{ tx.amount.toLocaleString(undefined, { minimumFractionDigits: 2 }) }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-center">
                  <button class="p-1.5 text-gray-400 hover:text-red-500 transition-colors">
                    <Trash2Icon class="w-4 h-4" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Add Transaction Modal -->
      <div v-if="showAddModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm overflow-y-auto">
        <div class="bg-white rounded-3xl w-full max-w-lg p-8 shadow-2xl animate-in fade-in zoom-in duration-200">
          <div class="flex justify-between items-center mb-6">
            <h3 class="text-xl font-black text-gray-900">New Transaction</h3>
            <button @click="showAddModal = false" class="p-2 text-gray-400 hover:text-gray-600">
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
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import { financialService, Transaction, Category } from '@/services/financialService'
import { 
  HistoryIcon, 
  SearchIcon, 
  Trash2Icon, 
  XIcon 
} from 'lucide-vue-next'

const transactions = ref<Transaction[]>([])
const categories = ref<Category[]>([])
const isLoading = ref(true)
const isSubmitting = ref(false)
const showAddModal = ref(false)
const search = ref('')
const filterType = ref('all')

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

const filteredTransactions = computed(() => {
  return transactions.value.filter(tx => {
    const matchesSearch = tx.description.toLowerCase().includes(search.value.toLowerCase())
    const matchesType = filterType.value === 'all' || tx.type === filterType.value
    return matchesSearch && matchesType
  }).sort((a, b) => new Date(b.transactionDate).getTime() - new Date(a.transactionDate).getTime())
})

const fetchInitialData = async () => {
  try {
    isLoading.value = true
    const [txData, catData] = await Promise.all([
      financialService.getTransactions(),
      financialService.getCategories()
    ])
    transactions.value = txData
    categories.value = catData
  } catch (err) {
    console.error('Failed to fetch transaction data:', err)
  } finally {
    isLoading.value = false
  }
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
    await fetchInitialData()
  } catch (err) {
    console.error('Failed to create transaction:', err)
  } finally {
    isSubmitting.value = false
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

onMounted(fetchInitialData)
</script>
