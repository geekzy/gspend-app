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
            <span class="ml-2 text-sm font-normal text-gray-500">({{ expenseCategories.length }})</span>
          </h3>
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            <div 
              v-for="cat in expenseCategories" 
              :key="cat.id"
              class="bg-white p-4 rounded-2xl border border-gray-100 shadow-sm hover:shadow-md transition-all group"
            >
              <div class="flex items-start justify-between mb-3">
                <div 
                  class="w-10 h-10 rounded-xl flex items-center justify-center text-xl shadow-inner border border-gray-50"
                  :style="{ backgroundColor: cat.color + '20', color: cat.color }"
                >
                  {{ cat.icon || '📦' }}
                </div>
                <div class="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity">
                  <button 
                    @click="openEditModal(cat)"
                    class="p-1.5 text-gray-300 hover:text-blue-500 transition-colors"
                    title="Edit category"
                  >
                    <EditIcon class="w-4 h-4" />
                  </button>
                  <button 
                    v-if="!cat.isSystem"
                    @click="confirmDelete(cat)"
                    class="p-1.5 text-gray-300 hover:text-red-500 transition-colors"
                    title="Delete category"
                  >
                    <Trash2Icon class="w-4 h-4" />
                  </button>
                </div>
              </div>
              <div>
                <p class="text-sm font-bold text-gray-900 truncate">{{ cat.name }}</p>
                <div class="flex items-center justify-between mt-1">
                  <p class="text-xs text-gray-400 capitalize">
                    {{ cat.isSystem ? 'System' : 'Custom' }}
                  </p>
                  <p class="text-xs text-gray-500" v-if="cat.usageCount !== undefined">
                    {{ cat.usageCount }} uses
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Income Categories -->
        <div>
          <h3 class="text-lg font-bold text-gray-900 mb-6 flex items-center">
            <ArrowUpIcon class="w-5 h-5 mr-2 text-green-500" />
            Income Categories
            <span class="ml-2 text-sm font-normal text-gray-500">({{ incomeCategories.length }})</span>
          </h3>
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            <div 
              v-for="cat in incomeCategories" 
              :key="cat.id"
              class="bg-white p-4 rounded-2xl border border-gray-100 shadow-sm hover:shadow-md transition-all group"
            >
              <div class="flex items-start justify-between mb-3">
                <div 
                  class="w-10 h-10 rounded-xl flex items-center justify-center text-xl shadow-inner border border-gray-50"
                  :style="{ backgroundColor: cat.color + '20', color: cat.color }"
                >
                  {{ cat.icon || '💰' }}
                </div>
                <div class="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity">
                  <button 
                    @click="openEditModal(cat)"
                    class="p-1.5 text-gray-300 hover:text-blue-500 transition-colors"
                    title="Edit category"
                  >
                    <EditIcon class="w-4 h-4" />
                  </button>
                  <button 
                    v-if="!cat.isSystem"
                    @click="confirmDelete(cat)"
                    class="p-1.5 text-gray-300 hover:text-red-500 transition-colors"
                    title="Delete category"
                  >
                    <Trash2Icon class="w-4 h-4" />
                  </button>
                </div>
              </div>
              <div>
                <p class="text-sm font-bold text-gray-900 truncate">{{ cat.name }}</p>
                <div class="flex items-center justify-between mt-1">
                  <p class="text-xs text-gray-400 capitalize">
                    {{ cat.isSystem ? 'System' : 'Custom' }}
                  </p>
                  <p class="text-xs text-gray-500" v-if="cat.usageCount !== undefined">
                    {{ cat.usageCount }} uses
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Add/Edit Modal -->
      <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm">
        <div class="bg-white rounded-3xl w-full max-w-md p-8 shadow-2xl animate-in fade-in zoom-in duration-200">
          <h3 class="text-xl font-bold text-gray-900 mb-6">
            {{ isEditing ? 'Edit Category' : 'New Category' }}
          </h3>
          <CategoryForm
            :category="currentCategory"
            :is-editing="isEditing"
            @save="handleSaveCategory"
            @cancel="closeModal"
          />
        </div>
      </div>

      <!-- Delete Confirmation Modal -->
      <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm">
        <div class="bg-white rounded-3xl w-full max-w-md p-8 shadow-2xl animate-in fade-in zoom-in duration-200">
          <div class="text-center">
            <div class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-red-100 mb-4">
              <Trash2Icon class="h-6 w-6 text-red-600" />
            </div>
            <h3 class="text-lg font-bold text-gray-900 mb-2">Delete Category</h3>
            <p class="text-sm text-gray-500 mb-6">
              Are you sure you want to delete "{{ categoryToDelete?.name }}"? 
              <span v-if="categoryToDelete?.usageCount && categoryToDelete.usageCount > 0" class="text-red-600 font-medium">
                This category is used in {{ categoryToDelete.usageCount }} transactions.
              </span>
              This action cannot be undone.
            </p>
            <div class="flex gap-3">
              <button 
                @click="showDeleteModal = false"
                type="button" 
                class="flex-1 py-3 rounded-2xl border border-gray-200 text-sm font-bold text-gray-700 hover:bg-gray-50 transition-all"
              >
                Cancel
              </button>
              <button 
                @click="handleDeleteCategory"
                :disabled="isDeleting"
                type="button" 
                class="flex-1 py-3 rounded-2xl bg-red-600 text-white text-sm font-bold hover:bg-red-700 transition-all disabled:opacity-50"
              >
                {{ isDeleting ? 'Deleting...' : 'Delete' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import CategoryForm from '@/components/forms/CategoryForm.vue'
import { financialService, Category } from '@/services/financialService'
import { 
  ArrowDownIcon, 
  ArrowUpIcon, 
  Trash2Icon,
  EditIcon
} from 'lucide-vue-next'

const categories = ref<Category[]>([])
const isLoading = ref(true)
const isDeleting = ref(false)
const showModal = ref(false)
const showDeleteModal = ref(false)
const isEditing = ref(false)
const currentCategory = ref<Partial<Category>>({})
const categoryToDelete = ref<Category | null>(null)

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
  currentCategory.value = {
    name: '',
    type: 'expense',
    icon: '🛒',
    color: '#34d399'
  }
  isEditing.value = false
  showModal.value = true
}

const openEditModal = (category: Category) => {
  currentCategory.value = { ...category }
  isEditing.value = true
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  currentCategory.value = {}
  isEditing.value = false
}

const handleSaveCategory = async (categoryData: Partial<Category>) => {
  try {
    if (isEditing.value && currentCategory.value.id) {
      await financialService.updateCategory(currentCategory.value.id, categoryData)
    } else {
      await financialService.createCategory(categoryData)
    }
    closeModal()
    await fetchCategories()
  } catch (err) {
    console.error('Failed to save category:', err)
    throw err // Let CategoryForm handle the error display
  }
}

const confirmDelete = (category: Category) => {
  if (category.isSystem) {
    return // System categories cannot be deleted
  }
  categoryToDelete.value = category
  showDeleteModal.value = true
}

const handleDeleteCategory = async () => {
  if (!categoryToDelete.value) return
  
  try {
    isDeleting.value = true
    await financialService.deleteCategory(categoryToDelete.value.id)
    showDeleteModal.value = false
    categoryToDelete.value = null
    await fetchCategories()
  } catch (err) {
    console.error('Failed to delete category:', err)
  } finally {
    isDeleting.value = false
  }
}

onMounted(fetchCategories)
</script>
