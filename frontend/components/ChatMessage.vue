<template>
  <div class="fade-in mb-6">
    <div :class="messageClass">
      <!-- User Message -->
      <div v-if="message.role === 'user'" class="flex items-start space-x-3">
        <div class="flex-shrink-0">
          <div class="w-8 h-8 rounded-full bg-gradient-primary flex items-center justify-center">
            <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
          </div>
        </div>
        <div class="flex-1">
          <div class="bg-gradient-primary text-white rounded-2xl rounded-tl-none px-4 py-3">
            <p class="text-sm">{{ message.content }}</p>
          </div>
          <p class="text-xs text-gray-400 mt-1 ml-2">{{ formatTime(message.timestamp) }}</p>
        </div>
      </div>

      <!-- Assistant Message -->
      <div v-else class="flex items-start space-x-3">
        <div class="flex-shrink-0">
          <div class="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center">
            <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
          </div>
        </div>
        <div class="flex-1 max-w-4xl">
          <!-- Error Message -->
          <div v-if="message.error" class="bg-red-50 border-l-4 border-red-500 rounded-lg px-4 py-3 mb-3">
            <div class="flex items-start">
              <svg class="w-5 h-5 text-red-500 mt-0.5 mr-2 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <div class="flex-1">
                <p class="text-sm font-semibold text-red-800">Error</p>
                <p class="text-sm text-red-600 mt-1">{{ message.error }}</p>
              </div>
            </div>
          </div>

          <!-- Success Response -->
          <div v-else>
            <!-- Insight/Summary -->
            <div v-if="message.content" class="bg-white border border-gray-200 rounded-2xl rounded-tl-none px-4 py-3 mb-3">
              <div class="flex items-start space-x-2">
                <svg class="w-5 h-5 text-primary-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <div class="flex-1">
                  <h4 class="text-sm font-semibold text-gray-800 mb-1">Insight</h4>
                  <p class="text-sm text-gray-700 whitespace-pre-wrap">{{ message.content }}</p>
                </div>
              </div>
            </div>

            <!-- Results Data -->
            <div v-if="message.results && message.results.count > 0" class="mb-3">
              <QueryResults :results="message.results" />
            </div>

            <!-- Collapsible SQL Query -->
            <div v-if="message.sql" class="mb-3">
              <button 
                @click="showSQL = !showSQL"
                class="w-full flex items-center justify-between px-4 py-2 bg-gray-50 hover:bg-gray-100 border border-gray-200 rounded-lg transition-colors"
              >
                <div class="flex items-center space-x-2">
                  <svg 
                    class="w-4 h-4 text-gray-600 transition-transform" 
                    :class="{ 'rotate-90': showSQL }"
                    fill="none" 
                    stroke="currentColor" 
                    viewBox="0 0 24 24"
                  >
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                  </svg>
                  <span class="text-sm font-semibold text-gray-700">Generated SQL</span>
                </div>
                <button 
                  v-if="showSQL"
                  @click.stop="copySQL" 
                  class="text-xs text-primary-600 hover:text-primary-700 font-medium px-2 py-1 hover:bg-primary-50 rounded"
                >
                  {{ copied ? '✓ Copied' : 'Copy' }}
                </button>
              </button>
              
              <div v-if="showSQL" class="mt-2 code-block">
                <pre class="text-xs">{{ message.sql }}</pre>
              </div>
            </div>

            <!-- Collapsible Reasoning Steps -->
            <div v-if="message.reasoning && message.reasoning.length > 0">
              <ReasoningSteps :steps="message.reasoning" />
            </div>
          </div>

          <p class="text-xs text-gray-400 mt-2">{{ formatTime(message.timestamp) }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Message } from '~/stores/chat'

const props = defineProps<{
  message: Message
}>()

const copied = ref(false)
const showSQL = ref(false)

const messageClass = computed(() => {
  return props.message.role === 'user' ? 'flex justify-end' : 'flex justify-start'
})

const formatTime = (date: Date) => {
  return new Date(date).toLocaleTimeString('id-ID', { 
    hour: '2-digit', 
    minute: '2-digit' 
  })
}

const copySQL = async () => {
  if (props.message.sql) {
    await navigator.clipboard.writeText(props.message.sql)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  }
}
</script>
