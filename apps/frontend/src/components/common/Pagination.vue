<template>
  <div class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6">
    <!-- Mobile pagination -->
    <div class="flex-1 flex justify-between sm:hidden">
      <button
        @click="goToPrevious"
        :disabled="currentPage <= 1"
        class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Previous
      </button>
      <button
        @click="goToNext"
        :disabled="currentPage >= totalPages"
        class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Next
      </button>
    </div>

    <!-- Desktop pagination -->
    <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
      <div class="flex items-center space-x-4">
        <!-- Results info -->
        <p class="text-sm text-gray-700">
          Showing
          <span class="font-medium">{{ startItem }}</span>
          to
          <span class="font-medium">{{ endItem }}</span>
          of
          <span class="font-medium">{{ total }}</span>
          results
        </p>

        <!-- Items per page selector -->
        <div class="flex items-center space-x-2">
          <label class="text-sm text-gray-700">Show:</label>
          <select
            :value="perPage"
            @change="onPerPageChange"
            class="px-2 py-1 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
          >
            <option value="10">10</option>
            <option value="20">20</option>
            <option value="50">50</option>
          </select>
          <span class="text-sm text-gray-700">per page</span>
        </div>
      </div>

      <!-- Page navigation -->
      <div>
        <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
          <!-- Previous button -->
          <button
            @click="goToPrevious"
            :disabled="currentPage <= 1"
            class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <ChevronLeftIcon class="h-5 w-5" />
          </button>

          <!-- Page numbers -->
          <template v-for="page in visiblePages" :key="page">
            <button
              v-if="page !== '...'"
              @click="goToPage(page as number)"
              :class="[
                'relative inline-flex items-center px-4 py-2 border text-sm font-medium',
                page === currentPage
                  ? 'z-10 bg-primary-50 border-primary-500 text-primary-600'
                  : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'
              ]"
            >
              {{ page }}
            </button>
            <span
              v-else
              class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700"
            >
              ...
            </span>
          </template>

          <!-- Next button -->
          <button
            @click="goToNext"
            :disabled="currentPage >= totalPages"
            class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <ChevronRightIcon class="h-5 w-5" />
          </button>
        </nav>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ChevronLeftIcon, ChevronRightIcon } from 'lucide-vue-next'

interface Props {
  currentPage: number
  totalPages: number
  total: number
  perPage: number
}

interface Emits {
  (e: 'page-change', page: number): void
  (e: 'per-page-change', perPage: number): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Computed properties
const startItem = computed(() => {
  return (props.currentPage - 1) * props.perPage + 1
})

const endItem = computed(() => {
  return Math.min(props.currentPage * props.perPage, props.total)
})

const visiblePages = computed(() => {
  const pages: (number | string)[] = []
  const totalPages = props.totalPages
  const current = props.currentPage

  if (totalPages <= 7) {
    // Show all pages if 7 or fewer
    for (let i = 1; i <= totalPages; i++) {
      pages.push(i)
    }
  } else {
    // Always show first page
    pages.push(1)

    if (current <= 4) {
      // Show pages 2-5 and ellipsis
      for (let i = 2; i <= 5; i++) {
        pages.push(i)
      }
      pages.push('...')
    } else if (current >= totalPages - 3) {
      // Show ellipsis and last 4 pages
      pages.push('...')
      for (let i = totalPages - 4; i <= totalPages - 1; i++) {
        pages.push(i)
      }
    } else {
      // Show ellipsis, current page area, ellipsis
      pages.push('...')
      for (let i = current - 1; i <= current + 1; i++) {
        pages.push(i)
      }
      pages.push('...')
    }

    // Always show last page (if not already shown)
    if (totalPages > 1) {
      pages.push(totalPages)
    }
  }

  return pages
})

// Event handlers
const goToPage = (page: number) => {
  if (page >= 1 && page <= props.totalPages && page !== props.currentPage) {
    emit('page-change', page)
  }
}

const goToPrevious = () => {
  if (props.currentPage > 1) {
    emit('page-change', props.currentPage - 1)
  }
}

const goToNext = () => {
  if (props.currentPage < props.totalPages) {
    emit('page-change', props.currentPage + 1)
  }
}

const onPerPageChange = (event: Event) => {
  const target = event.target as HTMLSelectElement
  const newPerPage = parseInt(target.value)
  emit('per-page-change', newPerPage)
}
</script>