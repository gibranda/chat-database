<template>
  <div class="card">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-bold text-gray-800">Query History</h3>
      <button 
        @click="handleClear"
        class="text-sm text-red-600 hover:text-red-700 font-medium"
        :disabled="chatStore.queryHistory.length === 0"
      >
        Clear
      </button>
    </div>

    <div v-if="chatStore.queryHistory.length === 0" class="text-center py-8">
      <svg class="w-12 h-12 text-gray-400 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <p class="text-sm text-gray-600">No queries yet</p>
    </div>

    <div v-else class="space-y-2 max-h-96 overflow-y-auto">
      <div 
        v-for="item in chatStore.recentHistory" 
        :key="item.id"
        class="border border-gray-200 rounded-lg p-3 hover:border-primary-300 hover:bg-gray-50 transition-all cursor-pointer"
        @click="$emit('select-query', item.question)"
      >
        <div class="flex items-start justify-between mb-2">
          <p class="text-sm text-gray-800 font-medium line-clamp-2">{{ item.question }}</p>
          <span 
            class="flex-shrink-0 ml-2 px-2 py-0.5 rounded text-xs font-medium"
            :class="item.success ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'"
          >
            {{ item.success ? '✓' : '✗' }}
          </span>
        </div>
        <p class="text-xs text-gray-500">{{ formatTime(item.timestamp) }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const chatStore = useChatStore()
const api = useApi()

const emit = defineEmits<{
  'select-query': [query: string]
}>()

const formatTime = (date: Date) => {
  return new Date(date).toLocaleString('id-ID', {
    day: '2-digit',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const handleClear = async () => {
  if (confirm('Are you sure you want to clear all query history?')) {
    chatStore.clearHistory()
    await api.clearHistory()
  }
}
</script>
