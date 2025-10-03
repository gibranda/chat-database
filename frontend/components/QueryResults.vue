<template>
  <div class="bg-white border border-gray-200 rounded-lg overflow-hidden">
    <!-- Header with Summary -->
    <div class="px-4 py-3 bg-gradient-to-r from-primary-50 to-purple-50 border-b border-gray-200">
      <div class="flex items-center justify-between mb-2">
        <div class="flex items-center space-x-2">
          <svg class="w-5 h-5 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
          <h4 class="text-sm font-bold text-gray-800">
            Query Results
          </h4>
        </div>
        <button 
          @click="downloadCSV" 
          class="text-xs text-primary-600 hover:text-primary-700 font-medium flex items-center space-x-1 px-2 py-1 hover:bg-white rounded transition-colors"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
          <span>Export CSV</span>
        </button>
      </div>
      
      <!-- Data Summary -->
      <div class="flex items-center space-x-4 text-xs">
        <div class="flex items-center space-x-1">
          <span class="font-semibold text-gray-700">Total Rows:</span>
          <span class="px-2 py-0.5 bg-white rounded-full font-bold text-primary-600">{{ results.count }}</span>
        </div>
        <div class="flex items-center space-x-1">
          <span class="font-semibold text-gray-700">Columns:</span>
          <span class="px-2 py-0.5 bg-white rounded-full font-bold text-primary-600">{{ results.columns.length }}</span>
        </div>
      </div>
    </div>

    <!-- Table -->
    <div class="overflow-x-auto max-h-96">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50 sticky top-0">
          <tr>
            <th 
              v-for="column in results.columns" 
              :key="column"
              class="px-4 py-3 text-left text-xs font-bold text-gray-700 uppercase tracking-wider border-b-2 border-primary-200"
            >
              {{ column }}
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr 
            v-for="(row, idx) in results.rows" 
            :key="idx" 
            class="hover:bg-primary-50 transition-colors"
            :class="{ 'bg-gray-50': idx % 2 === 1 }"
          >
            <td 
              v-for="column in results.columns" 
              :key="column"
              class="px-4 py-3 text-sm text-gray-700"
            >
              <span :class="getCellClass(row[column])">
                {{ formatValue(row[column]) }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Footer with additional info -->
    <div v-if="results.count > 10" class="px-4 py-2 bg-gray-50 border-t border-gray-200">
      <p class="text-xs text-gray-600 text-center">
        Showing {{ Math.min(results.count, results.rows.length) }} of {{ results.count }} rows
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { QueryResult } from '~/stores/chat'

const props = defineProps<{
  results: QueryResult
}>()

const formatValue = (value: any) => {
  if (value === null || value === undefined) {
    return 'NULL'
  }
  if (typeof value === 'boolean') {
    return value ? 'true' : 'false'
  }
  if (typeof value === 'number') {
    // Format numbers with thousand separators
    return new Intl.NumberFormat('id-ID').format(value)
  }
  if (typeof value === 'object') {
    return JSON.stringify(value)
  }
  return String(value)
}

const getCellClass = (value: any) => {
  if (value === null || value === undefined) {
    return 'text-gray-400 italic'
  }
  if (typeof value === 'number') {
    return 'font-mono font-semibold text-gray-900'
  }
  if (typeof value === 'boolean') {
    return value ? 'text-green-600 font-semibold' : 'text-red-600 font-semibold'
  }
  return ''
}

const downloadCSV = () => {
  const headers = props.results.columns.join(',')
  const rows = props.results.rows.map(row => 
    props.results.columns.map(col => {
      const value = formatValue(row[col])
      // Escape quotes and wrap in quotes if contains comma
      return value.includes(',') ? `"${value.replace(/"/g, '""')}"` : value
    }).join(',')
  )
  
  const csv = [headers, ...rows].join('\n')
  const blob = new Blob([csv], { type: 'text/csv' })
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `query_results_${new Date().getTime()}.csv`
  a.click()
  window.URL.revokeObjectURL(url)
}
</script>
