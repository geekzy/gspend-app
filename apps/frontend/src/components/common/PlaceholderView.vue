<template>
  <MainLayout>
    <div class="px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto w-full">
      <header class="flex items-center gap-4 mb-8">
        <h1 class="text-3xl font-extrabold text-gray-900 tracking-tight">{{ title }}</h1>
      </header>
      
      <div class="bg-white p-12 sm:p-20 rounded-3xl shadow-sm border border-gray-100 flex flex-col items-center justify-center text-center">
        <div class="bg-primary-50 w-24 h-24 rounded-3xl flex items-center justify-center text-5xl mb-6 shadow-inner animate-pulse">
           🏗️
        </div>
        <h2 class="text-2xl font-bold text-gray-900 mb-3">Module Under Construction</h2>
        <p class="text-gray-500 max-w-md text-lg leading-relaxed">
          The <span class="font-bold text-primary-600">{{ title }}</span> feature is currently being developed as part of the Phase 0 prototype. 
        </p>
        <div class="mt-10 flex gap-4">
          <button @click="$router.push('/')" class="px-8 py-3 bg-primary-600 text-white font-bold rounded-2xl shadow-lg hover:bg-primary-700 hover:shadow-primary-100 transition-all active:scale-95">
            Return to Dashboard
          </button>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import MainLayout from '@/layouts/MainLayout.vue'
defineProps<{ title: string }>()
</script>
