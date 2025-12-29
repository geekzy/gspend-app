<template>
  <div class="fixed top-4 right-4 z-[100] flex flex-col gap-3 pointer-events-none w-full max-w-sm">
    <TransitionGroup 
      enter-active-class="transform ease-out duration-300 transition"
      enter-from-class="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
      enter-to-class="translate-y-0 opacity-100 sm:translate-x-0"
      leave-active-class="transition ease-in duration-100"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div 
        v-for="notification in notifications" 
        :key="notification.id"
        class="bg-white rounded-2xl shadow-2xl border border-gray-100 p-4 pointer-events-auto flex items-start gap-4 animate-in slide-in-from-right-4"
        :class="getTypeClass(notification.type)"
      >
        <div class="flex-shrink-0 mt-0.5">
          <CheckCircleIcon v-if="notification.type === 'success'" class="w-5 h-5 text-green-500" />
          <XCircleIcon v-else-if="notification.type === 'error'" class="w-5 h-5 text-red-500" />
          <AlertCircleIcon v-else-if="notification.type === 'warning'" class="w-5 h-5 text-orange-500" />
          <InfoIcon v-else class="w-5 h-5 text-blue-500" />
        </div>
        <div class="flex-1">
          <p v-if="notification.title" class="text-xs font-bold text-gray-600 uppercase tracking-wide mb-1">
            {{ notification.title }}
          </p>
          <p class="text-sm font-bold text-gray-900 leading-tight">
            {{ notification.message }}
          </p>
          <div v-if="notification.actions && notification.actions.length" class="mt-3 flex gap-2">
            <button
              v-for="action in notification.actions"
              :key="action.label"
              @click="action.handler"
              class="text-xs font-bold px-3 py-1.5 rounded-lg transition-colors"
              :class="action.style === 'primary' 
                ? 'bg-blue-600 text-white hover:bg-blue-700' 
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'"
            >
              {{ action.label }}
            </button>
          </div>
        </div>
        <button 
          @click="removeNotification(notification.id)"
          class="flex-shrink-0 text-gray-400 hover:text-gray-600 transition-colors"
        >
          <XIcon class="w-4 h-4" />
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useNotificationStore } from '@/stores/notification'
import { 
  CheckCircleIcon, 
  XCircleIcon, 
  AlertCircleIcon, 
  InfoIcon,
  XIcon
} from 'lucide-vue-next'

const store = useNotificationStore()
const { notifications } = storeToRefs(store)
const { removeNotification } = store

const getTypeClass = (type: string) => {
  switch (type) {
    case 'success': return 'border-l-4 border-l-green-500'
    case 'error': return 'border-l-4 border-l-red-500'
    case 'warning': return 'border-l-4 border-l-orange-500'
    default: return 'border-l-4 border-l-blue-500'
  }
}
</script>
