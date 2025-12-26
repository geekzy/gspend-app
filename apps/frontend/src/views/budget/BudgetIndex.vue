<template>
  <MainLayout>
    <div class="px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto w-full">
      <div class="md:flex md:items-center md:justify-between mb-8">
        <div class="flex-1 min-w-0">
          <h2 class="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
            Budget Planning
          </h2>
          <p class="mt-1 text-sm text-gray-500">
            Plan your monthly spending and track progress.
          </p>
        </div>
        <div class="mt-4 flex md:mt-0 md:ml-4">
          <button 
            @click="openAddModal"
            type="button" 
            class="inline-flex items-center px-4 py-2.5 border border-transparent rounded-2xl shadow-lg text-sm font-bold text-white bg-primary-600 hover:bg-primary-700 transition-all active:scale-95"
          >
            Create New Budget
          </button>
        </div>
      </div>

      <div v-if="isLoading" class="flex items-center justify-center min-h-[400px]">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>

      <div v-else-if="!activeBudget" class="bg-white rounded-3xl p-12 text-center border-2 border-dashed border-gray-200">
        <div class="bg-blue-50 w-16 h-16 rounded-2xl flex items-center justify-center mx-auto mb-4">
          <BarChart3Icon class="w-8 h-8 text-blue-600" />
        </div>
        <h3 class="text-lg font-bold text-gray-900 mb-2">No active budget</h3>
        <p class="text-gray-500 mb-6">Setup a budget for this month to start tracking your goals.</p>
        <button 
          @click="openAddModal"
          class="inline-flex items-center px-6 py-3 border border-transparent rounded-2xl text-sm font-bold text-blue-600 bg-blue-50 hover:bg-blue-100 transition-all"
        >
          Create Budget
        </button>
      </div>

      <template v-else>
        <!-- Active Budget Overview -->
        <div class="bg-white rounded-3xl p-8 shadow-sm border border-gray-100 mb-8 overflow-hidden relative">
          <div class="absolute top-0 right-0 w-64 h-64 bg-primary-50 rounded-full -mr-32 -mt-32 opacity-50"></div>
          <div class="relative z-10 lg:flex items-center justify-between">
            <div class="mb-6 lg:mb-0">
              <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-bold bg-primary-100 text-primary-700 uppercase mb-4">
                Active Budget
              </span>
              <h3 class="text-3xl font-black text-gray-900 mb-2">{{ activeBudget.name }}</h3>
              <p class="text-gray-500 flex items-center">
                <CalendarIcon class="w-4 h-4 mr-2" />
                {{ formatDate(activeBudget.startDate) }} - {{ formatDate(activeBudget.endDate) }}
              </p>
            </div>
            <div class="flex flex-col items-end">
              <div class="text-sm font-bold text-gray-500 mb-1">Total Allocated</div>
              <div class="text-4xl font-black text-gray-900">${{ activeBudget.totalAmount.toLocaleString() }}</div>
            </div>
          </div>
          
          <div class="mt-8 border-t border-gray-100 pt-8">
            <div class="flex items-center justify-between mb-2">
              <span class="text-sm font-bold text-gray-600">Overall Spending Progress</span>
              <span class="text-sm font-bold text-primary-600">{{ overallProgress }}%</span>
            </div>
            <div class="w-full bg-gray-100 h-3 rounded-full overflow-hidden shadow-inner">
              <div class="h-full bg-primary-600 rounded-full transition-all duration-1000" :style="{ width: overallProgress + '%' }"></div>
            </div>
          </div>
        </div>

        <!-- Budget Items Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div 
            v-for="item in activeBudget.items" 
            :key="item.id"
            class="bg-white rounded-3xl p-6 shadow-sm border border-gray-100 hover:shadow-md transition-shadow"
          >
            <div class="flex justify-between items-center mb-6">
              <div class="flex items-center">
                <div class="bg-blue-50 p-2 rounded-xl text-blue-600 mr-4">
                  <PieChartIcon class="w-5 h-5" />
                </div>
                <h4 class="font-bold text-gray-900">{{ item.categoryName }}</h4>
              </div>
              <button class="p-1.5 text-gray-400 hover:text-gray-600 transition-all">
                <MoreVerticalIcon class="w-5 h-5" />
              </button>
            </div>

            <div class="flex justify-between items-end mb-2">
              <div>
                <p class="text-xs font-bold text-gray-400 uppercase tracking-tight">Spent</p>
                <p class="text-xl font-black text-gray-900">${{ item.spentAmount.toLocaleString() }}</p>
              </div>
              <div class="text-right">
                <p class="text-xs font-bold text-gray-400 uppercase tracking-tight text-right">Limit</p>
                <p class="text-sm font-bold text-gray-500">${{ item.plannedAmount.toLocaleString() }}</p>
              </div>
            </div>

            <div class="w-full bg-gray-50 h-2 rounded-full overflow-hidden mb-4">
              <div 
                class="h-full rounded-full transition-all duration-1000" 
                :class="getProgressBarClass(item.spentAmount, item.plannedAmount)"
                :style="{ width: getPercent(item.spentAmount, item.plannedAmount) + '%' }"
              ></div>
            </div>

            <p v-if="item.notes" class="text-xs text-gray-400 italic mt-4 border-t border-gray-50 pt-2 line-clamp-1">{{ item.notes }}</p>
          </div>
        </div>
      </template>

      <!-- Create Budget Modal -->
      <div v-if="showAddModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm overflow-y-auto">
        <div class="bg-white rounded-3xl w-full max-w-2xl p-8 shadow-2xl animate-in fade-in zoom-in duration-200 my-8">
          <div class="flex justify-between items-center mb-8">
            <h3 class="text-2xl font-black text-gray-900">New Budget Plan</h3>
            <button @click="showAddModal = false" class="p-2 text-gray-400 hover:text-gray-600">
              <XIcon class="w-6 h-6" />
            </button>
          </div>

          <form @submit.prevent="handleCreateBudget" class="space-y-8">
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <div class="space-y-4">
                <div>
                  <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Budget Title</label>
                  <input v-model="newBudget.name" type="text" class="w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none" placeholder="e.g. October 2024 Household" required />
                </div>
                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Start Date</label>
                    <input v-model="newBudget.startDate" type="date" class="w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none" required />
                  </div>
                  <div>
                    <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">End Date</label>
                    <input v-model="newBudget.endDate" type="date" class="w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none" required />
                  </div>
                </div>
              </div>

              <div class="bg-gray-50 p-6 rounded-3xl">
                <div class="flex items-center text-primary-600 mb-4">
                  <InfoIcon class="w-5 h-5 mr-2" />
                  <span class="text-sm font-bold uppercase tracking-tight">Summary</span>
                </div>
                <div class="space-y-2">
                  <div class="flex justify-between text-sm">
                    <span class="text-gray-500">Categories</span>
                    <span class="font-bold font-mono text-gray-900">{{ newBudget.items?.length || 0 }}</span>
                  </div>
                  <div class="flex justify-between text-lg font-black border-t border-gray-200 pt-2 mt-2">
                    <span class="text-gray-900">Total Plan</span>
                    <span class="text-primary-600">${{ totalPlannedAmount.toLocaleString() }}</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Budget Items -->
            <div>
              <div class="flex justify-between items-center mb-4">
                <h4 class="text-lg font-bold text-gray-900">Category Allocations</h4>
                <button type="button" @click="addItem" class="text-primary-600 text-sm font-bold hover:underline">+ Add Category</button>
              </div>
              <div class="space-y-4 max-h-[300px] overflow-y-auto pr-2 custom-scrollbar">
                <div v-for="(item, index) in newBudget.items" :key="index" class="grid grid-cols-12 gap-3 items-start bg-gray-50 p-4 rounded-2xl relative">
                  <div class="col-span-6">
                    <select v-model="item.categoryId" class="w-full px-4 py-2 bg-white border border-gray-200 rounded-xl outline-none text-sm" @change="updateCategoryName(index)">
                      <option value="">Select Category</option>
                      <option v-for="cat in expenseCategories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                    </select>
                  </div>
                  <div class="col-span-5 relative">
                    <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 text-xs">$</span>
                    <input v-model.number="item.plannedAmount" type="number" class="w-full pl-6 pr-4 py-2 bg-white border border-gray-200 rounded-xl outline-none text-sm" placeholder="Amount" />
                  </div>
                  <button type="button" @click="removeItem(index)" class="col-span-1 p-2 text-gray-400 hover:text-red-500">
                    <XCircleIcon class="w-5 h-5" />
                  </button>
                </div>
              </div>
            </div>

            <div class="flex gap-4 pt-4">
              <button @click="showAddModal = false" type="button" class="flex-1 py-4 rounded-2xl border-2 border-gray-100 text-sm font-bold text-gray-500 hover:bg-gray-50 transition-all">Cancel</button>
              <button type="submit" :disabled="isSubmitting || !isValid" class="flex-1 py-4 rounded-2xl bg-primary-600 text-white text-sm font-black hover:bg-primary-700 shadow-xl shadow-primary-200 transition-all disabled:opacity-50">
                {{ isSubmitting ? 'Creating Budget...' : 'Finalize Budget' }}
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
import { financialService, Budget, Category } from '@/services/financialService'
import { 
  BarChart3Icon, 
  CalendarIcon, 
  PieChartIcon, 
  MoreVerticalIcon, 
  XIcon,
  InfoIcon,
  XCircleIcon
} from 'lucide-vue-next'

const activeBudget = ref<Budget | null>(null)
const categories = ref<Category[]>([])
const isLoading = ref(true)
const isSubmitting = ref(false)
const showAddModal = ref(false)

const newBudget = ref({
  name: '',
  periodType: 'monthly',
  startDate: new Date().toISOString().split('T')[0],
  endDate: new Date(new Date().getFullYear(), new Date().getMonth() + 1, 0).toISOString().split('T')[0],
  totalAmount: 0,
  items: [] as any[]
})

const expenseCategories = computed(() => categories.value.filter(c => c.type === 'expense'))

const totalPlannedAmount = computed(() => {
  return newBudget.value.items.reduce((sum, item) => sum + (Number(item.plannedAmount) || 0), 0)
})

const isValid = computed(() => {
  return newBudget.value.name && 
         newBudget.value.items.length > 0 && 
         newBudget.value.items.every(item => item.categoryId && item.plannedAmount > 0)
})

const overallProgress = computed(() => {
  if (!activeBudget.value || activeBudget.value.totalAmount === 0) return 0
  const spent = activeBudget.value.items.reduce((sum, item) => sum + item.spentAmount, 0)
  return Math.min(Math.round((spent / activeBudget.value.totalAmount) * 100), 100)
})

const fetchInitialData = async () => {
  try {
    isLoading.value = true
    const [budget, cats] = await Promise.all([
      financialService.getActiveBudget(),
      financialService.getCategories()
    ])
    activeBudget.value = budget
    categories.value = cats
  } catch (err) {
    console.error('Failed to fetch budget data:', err)
  } finally {
    isLoading.value = false
  }
}

const openAddModal = () => {
  newBudget.value = {
    name: '',
    periodType: 'monthly',
    startDate: new Date().toISOString().split('T')[0],
    endDate: new Date(new Date().getFullYear(), new Date().getMonth() + 1, 0).toISOString().split('T')[0],
    totalAmount: 0,
    items: [{ categoryId: '', plannedAmount: 0, notes: '' }]
  }
  showAddModal.value = true
}

const addItem = () => {
  newBudget.value.items.push({ categoryId: '', plannedAmount: 0, notes: '' })
}

const removeItem = (index: number) => {
  newBudget.value.items.splice(index, 1)
}

const updateCategoryName = (index: number) => {
  const item = newBudget.value.items[index]
  const cat = categories.value.find(c => c.id === item.categoryId)
  if (cat) {
    item.categoryName = cat.name
  }
}

const handleCreateBudget = async () => {
  try {
    isSubmitting.value = true
    newBudget.value.totalAmount = totalPlannedAmount.value
    await financialService.createBudget(newBudget.value as any)
    showAddModal.value = false
    await fetchInitialData()
  } catch (err) {
    console.error('Failed to create budget:', err)
  } finally {
    isSubmitting.value = false
  }
}

const getProgressBarClass = (spent: number, planned: number) => {
  const percent = (spent / planned) * 100
  if (percent >= 100) return 'bg-red-500'
  if (percent >= 85) return 'bg-orange-400'
  return 'bg-blue-600'
}

const getPercent = (spent: number, planned: number) => {
  if (planned === 0) return 0
  return Math.min(Math.round((spent / planned) * 100), 100)
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

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #e5e7eb;
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #d1d5db;
}
</style>
