<template>
  <div class="mt-3">
    <button 
      @click="expanded = !expanded"
      class="w-full flex items-center justify-between px-4 py-2 bg-gray-50 hover:bg-gray-100 border border-gray-200 rounded-lg transition-colors"
    >
      <div class="flex items-center space-x-2">
        <svg 
          class="w-4 h-4 text-gray-600 transition-transform" 
          :class="{ 'rotate-90': expanded }"
          fill="none" 
          stroke="currentColor" 
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
        <svg class="w-4 h-4 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
        </svg>
        <span class="text-sm font-semibold text-gray-700">Agent Reasoning Process</span>
      </div>
      <span class="text-xs text-gray-500 px-2 py-0.5 bg-white rounded-full">{{ steps.length }} steps</span>
    </button>

    <div v-if="expanded" class="mt-2 space-y-2">
      <div 
        v-for="step in steps" 
        :key="step.step"
        class="bg-gradient-to-r from-primary-50 to-purple-50 border-l-4 border-primary-500 rounded-r-lg p-3"
      >
        <div class="flex items-start space-x-3">
          <div class="flex-shrink-0">
            <div class="w-6 h-6 rounded-full bg-primary-500 text-white flex items-center justify-center text-xs font-bold">
              {{ step.step }}
            </div>
          </div>
          <div class="flex-1 min-w-0">
            <h5 class="text-sm font-semibold text-gray-800 capitalize">{{ step.action.replace(/_/g, ' ') }}</h5>
            <p class="text-xs text-gray-600 mt-1">
              <strong class="text-gray-700">Thought:</strong> {{ step.thought }}
            </p>
            <p class="text-xs text-gray-500 mt-1">
              <strong class="text-gray-600">Observation:</strong> 
              <span class="line-clamp-2">{{ step.observation }}</span>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ReasoningStep } from '~/stores/chat'

defineProps<{
  steps: ReasoningStep[]
}>()

const expanded = ref(false)
</script>
