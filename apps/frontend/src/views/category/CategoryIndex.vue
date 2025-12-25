<template>
  <MainLayout>
    <div class="px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto w-full">
      <div class="md:flex md:items-center md:justify-between mb-8">
        <div class="flex-1 min-w-0">
          <h2 class="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
            Categories
          </h2>
          <p class="mt-1 text-sm text-gray-500">
            Organize your transactions with custom categories.
          </p>
        </div>
        <div class="mt-4 flex md:mt-0 md:ml-4">
          <button 
            @click="openAddModal"
            type="button" 
            class="inline-flex items-center px-4 py-2.5 border border-transparent rounded-2xl shadow-lg text-sm font-bold text-white bg-primary-600 hover:bg-primary-700 transition-all active:scale-95"
          >
            Add Category
          </button>
        </div>
      </div>

      <div v-if="isLoading" class="flex items-center justify-center min-h-[400px]">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>

      <div v-else class="space-y-12">
        <!-- Expense Categories -->
        <div>
          <h3 class="text-lg font-bold text-gray-900 mb-6 flex items-center">
            <ArrowDownIcon class="w-5 h-5 mr-2 text-red-500" />
            Expense Categories
          </h3>
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            <div 
              v-for="cat in expenseCategories" 
              :key="cat.id"
              class="bg-white p-4 rounded-2xl border border-gray-100 shadow-sm hover:shadow-md transition-all flex items-center group"
            >
              <div 
                class="w-10 h-10 rounded-xl flex items-center justify-center text-xl shadow-inner border border-gray-50 mr-4"
                :style="{ backgroundColor: cat.color + '20', color: cat.color }"
              >
                {{ cat.icon || '📦' }}
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-bold text-gray-900 truncate">{{ cat.name }}</p>
                <p class="text-xs text-gray-400 capitalize">{{ cat.isSystem ? 'System' : 'Custom' }}</p>
              </div>
              <button 
                v-if="!cat.isSystem"
                class="opacity-0 group-hover:opacity-100 p-1.5 text-gray-300 hover:text-red-500 transition-all"
              >
                <Trash2Icon class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>

        <!-- Income Categories -->
        <div>
          <h3 class="text-lg font-bold text-gray-900 mb-6 flex items-center">
            <ArrowUpIcon class="w-5 h-5 mr-2 text-green-500" />
            Income Categories
          </h3>
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            <div 
              v-for="cat in incomeCategories" 
              :key="cat.id"
              class="bg-white p-4 rounded-2xl border border-gray-100 shadow-sm hover:shadow-md transition-all flex items-center group"
            >
              <div 
                class="w-10 h-10 rounded-xl flex items-center justify-center text-xl shadow-inner border border-gray-50 mr-4"
                :style="{ backgroundColor: cat.color + '20', color: cat.color }"
              >
                {{ cat.icon || '💰' }}
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-bold text-gray-900 truncate">{{ cat.name }}</p>
                <p class="text-xs text-gray-400 capitalize">{{ cat.isSystem ? 'System' : 'Custom' }}</p>
              </div>
              <button 
                v-if="!cat.isSystem"
                class="opacity-0 group-hover:opacity-100 p-1.5 text-gray-300 hover:text-red-500 transition-all"
              >
                <Trash2Icon class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Add Modal -->
      <div v-if="showAddModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm">
        <div class="bg-white rounded-3xl w-full max-w-md p-8 shadow-2xl animate-in fade-in zoom-in duration-200">
          <h3 class="text-xl font-bold text-gray-900 mb-6">New Category</h3>
          <form @submit.prevent="handleAddCategory" class="space-y-4">
            <div>
              <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Type</label>
              <div class="grid grid-cols-2 gap-3">
                <button 
                  type="button"
                  @click="newCategory.type = 'expense'"
                  :class="[
                    newCategory.type === 'expense' ? 'bg-primary-600 text-white shadow-lg shadow-primary-200' : 'bg-gray-50 text-gray-500',
                    'py-3 rounded-2xl text-sm font-bold transition-all'
                  ]"
                >
                  Expense
                </button>
                <button 
                  type="button"
                  @click="newCategory.type = 'income'"
                  :class="[
                    newCategory.type === 'income' ? 'bg-primary-600 text-white shadow-lg shadow-primary-200' : 'bg-gray-50 text-gray-500',
                    'py-3 rounded-2xl text-sm font-bold transition-all'
                  ]"
                >
                  Income
                </button>
              </div>
            </div>
            <div>
              <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Name</label>
              <input v-model="newCategory.name" type="text" class="w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none" placeholder="e.g. Subscriptions" required />
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Icon</label>
                <select v-model="newCategory.icon" class="w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500 outline-none bg-white">
                  <option v-for="icon in iconOptions" :key="icon" :value="icon">{{ icon }}</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">Color</label>
                <input v-model="newCategory.color" type="color" class="w-full h-12 p-1 border border-gray-200 rounded-2xl bg-white cursor-pointer" />
              </div>
            </div>
            <div class="flex gap-3 mt-8">
              <button @click="showAddModal = false" type="button" class="flex-1 py-3.5 rounded-2xl border border-gray-200 text-sm font-bold text-gray-700 hover:bg-gray-50 transition-all">Cancel</button>
              <button type="submit" :disabled="isSubmitting" class="flex-1 py-3.5 rounded-2xl bg-primary-600 text-white text-sm font-bold hover:bg-primary-700 shadow-lg shadow-primary-200 transition-all disabled:opacity-50">
                {{ isSubmitting ? 'Creating...' : 'Create' }}
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
import { financialService, Category } from '@/services/financialService'
import { 
  ArrowDownIcon, 
  ArrowUpIcon, 
  Trash2Icon 
} from 'lucide-vue-next'

const categories = ref<Category[]>([])
const isLoading = ref(true)
const isSubmitting = ref(false)
const showAddModal = ref(false)

const iconOptions = ['🛒', '🍽️', '🚗', '⚡', '🏠', '🎬', '🏥', '👔', '🏫', '✈️', '🎁', '⚽', '💰', '📈', '🏢', '🛠️']

const newCategory = ref({
  name: '',
  type: 'expense' as 'income' | 'expense',
  icon: '🛒',
  color: '#34d399'
})

const expenseCategories = computed(() => categories.value.filter(c => c.type === 'expense'))
const incomeCategories = computed(() => categories.value.filter(c => c.type === 'income'))

const fetchCategories = async () => {
  try {
    isLoading.value = true
    categories.value = await financialService.getCategories()
  } catch (err) {
    console.error('Failed to fetch categories:', err)
  } finally {
    isLoading.value = false
  }
}

const openAddModal = () => {
  newCategory.value = {
    name: '',
    type: 'expense',
    icon: '🛒',
    color: '#34d399'
  }
  showAddModal.value = true
}

const handleAddCategory = async () => {
  try {
    isSubmitting.value = true
    await financialService.createCategory(newCategory.value)
    showAddModal.value = false
    await fetchCategories()
  } catch (err) {
    console.error('Failed to create category:', err)
  } finally {
    isSubmitting.value = false
  }
}

onMounted(fetchCategories)
</script>
