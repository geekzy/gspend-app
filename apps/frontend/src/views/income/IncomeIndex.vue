<template>
  <MainLayout>
    <div class="px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto w-full">
      <div class="md:flex md:items-center md:justify-between mb-8">
        <div class="flex-1 min-w-0">
          <h2 class="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
            Income Management
          </h2>
          <p class="mt-1 text-sm text-gray-500">
            Track and manage your family income sources.
          </p>
        </div>
        <div class="mt-4 flex md:mt-0 md:ml-4">
          <button 
            @click="showAddModal = true"
            type="button" 
            class="inline-flex items-center px-4 py-2.5 border border-transparent rounded-2xl shadow-lg text-sm font-bold text-white bg-primary-600 hover:bg-primary-700 transition-all active:scale-95"
          >
            Add Income Source
          </button>
        </div>
      </div>

      <div v-if="isLoading" class="flex items-center justify-center min-h-[400px]">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>

      <div v-else-if="incomes.length === 0" class="bg-white rounded-3xl p-12 text-center border-2 border-dashed border-gray-200">
        <div class="bg-primary-50 w-16 h-16 rounded-2xl flex items-center justify-center mx-auto mb-4">
          <WalletIcon class="w-8 h-8 text-primary-600" />
        </div>
        <h3 class="text-lg font-bold text-gray-900 mb-2">No income sources yet</h3>
        <p class="text-gray-500 mb-6">Start by adding your salary or other steady income sources.</p>
        <button 
          @click="showAddModal = true"
          class="inline-flex items-center px-6 py-3 border border-transparent rounded-2xl text-sm font-bold text-primary-600 bg-primary-50 hover:bg-primary-100 transition-all"
        >
          Add First Source
        </button>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div 
          v-for="income in incomes" 
          :key="income.id"
          class="bg-white rounded-3xl p-6 shadow-sm border border-gray-100 hover:shadow-md transition-shadow"
        >
          <div class="flex justify-between items-start mb-4">
            <div class="bg-green-50 p-3 rounded-2xl text-green-600">
              <DollarSignIcon class="w-6 h-6" />
            </div>
            <div class="text-right">
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-bold bg-green-100 text-green-700 uppercase">
                {{ income.frequency }}
              </span>
            </div>
          </div>
          <h3 class="text-lg font-bold text-gray-900 mb-1">{{ income.source }}</h3>
          <p class="text-3xl font-black text-gray-900 mb-4">${{ income.amount.toLocaleString() }}</p>
          <div class="flex items-center text-sm text-gray-400 mb-6">
            <CalendarIcon class="w-4 h-4 mr-2" />
            Effective: {{ formatDate(income.effectiveDate) }}
          </div>
          <div class="flex gap-2">
            <button class="flex-1 py-2.5 rounded-xl border border-gray-200 text-sm font-bold text-gray-600 hover:bg-gray-50 transition-all">Edit</button>
            <button class="px-3 rounded-xl border border-gray-200 text-red-500 hover:bg-red-50 hover:border-red-100 transition-all">
              <Trash2Icon class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

      <!-- Add Modal (Simple placeholder for now) -->
      <div v-if="showAddModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm">
        <div class="bg-white rounded-3xl w-full max-w-md p-8 shadow-2xl animate-in fade-in zoom-in duration-200">
          <h3 class="text-xl font-bold text-gray-900 mb-6">Add Income Source</h3>
          <form @submit.prevent="handleAddIncome" class="space-y-4">
            <div>
              <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Source Name</label>
              <input v-model="newIncome.source" type="text" class="w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none" placeholder="e.g. Monthly Salary" required />
            </div>
            <div>
              <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Amount</label>
              <div class="relative">
                <span class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400">$</span>
                <input v-model.number="newIncome.amount" type="number" step="0.01" class="w-full pl-8 pr-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none" placeholder="0.00" required />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Frequency</label>
                <select v-model="newIncome.frequency" class="w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none appearance-none bg-white">
                  <option value="monthly">Monthly</option>
                  <option value="bi-weekly">Bi-weekly</option>
                  <option value="one-time">One-time</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Date</label>
                <input v-model="newIncome.effectiveDate" type="date" class="w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none" required />
              </div>
            </div>
            <div class="flex gap-3 mt-8">
              <button @click="showAddModal = false" type="button" class="flex-1 py-3.5 rounded-2xl border border-gray-200 text-sm font-bold text-gray-700 hover:bg-gray-50 transition-all">Cancel</button>
              <button type="submit" :disabled="isSubmitting" class="flex-1 py-3.5 rounded-2xl bg-primary-600 text-white text-sm font-bold hover:bg-primary-700 shadow-lg shadow-primary-200 transition-all disabled:opacity-50">
                {{ isSubmitting ? 'Saving...' : 'Save Source' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import { financialService, Income } from '@/services/financialService'
import { 
  WalletIcon, 
  DollarSignIcon, 
  CalendarIcon, 
  Trash2Icon 
} from 'lucide-vue-next'

const incomes = ref<Income[]>([])
const isLoading = ref(true)
const isSubmitting = ref(false)
const showAddModal = ref(false)

const newIncome = ref({
  source: '',
  amount: 0,
  frequency: 'monthly',
  effectiveDate: new Date().toISOString().split('T')[0]
})

const fetchIncomes = async () => {
  try {
    isLoading.value = true
    incomes.value = await financialService.getIncomes()
  } catch (err) {
    console.error('Failed to fetch incomes:', err)
  } finally {
    isLoading.value = false
  }
}

const handleAddIncome = async () => {
  try {
    isSubmitting.value = true
    await financialService.createIncome(newIncome.value)
    showAddModal.value = false
    // Reset form
    newIncome.value = {
      source: '',
      amount: 0,
      frequency: 'monthly',
      effectiveDate: new Date().toISOString().split('T')[0]
    }
    await fetchIncomes()
  } catch (err) {
    console.error('Failed to add income:', err)
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

onMounted(fetchIncomes)
</script>
