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
    <div class="overflow-x-auto overflow-y-auto max-h-[60vh]">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50 sticky top-0">
          <tr>
            <th 
              v-for="column in results.columns" 
              :key="column"
              :aria-sort="sortKey === column ? (sortDir === 'asc' ? 'ascending' : 'descending') : 'none'"
              @click="toggleSort(column)"
              class="px-4 py-3 text-left text-xs font-bold text-gray-700 uppercase tracking-wider border-b-2 border-primary-200 cursor-pointer select-none"
            >
              <span class="inline-flex items-center gap-1">
                {{ column }}
                <svg v-if="sortKey === column" :class="['w-3.5 h-3.5', sortDir === 'desc' ? 'rotate-180' : '']" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M19 9l-7 7-7-7" />
                </svg>
              </span>
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr 
            v-for="(row, idx) in displayedRows" 
            :key="startIndex + idx" 
            class="hover:bg-primary-50 transition-colors"
            :class="{ 'bg-gray-50': (startIndex + idx) % 2 === 1 }"
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

    <!-- Footer: Pagination Controls -->
    <div class="px-4 py-2 bg-gray-50 border-t border-gray-200 flex items-center justify-between flex-wrap gap-2">
      <div class="text-xs text-gray-600">
        Showing {{ results.rows.length === 0 ? 0 : (startIndex + 1) }}–{{ endIndex }} of {{ results.rows.length }} rows
      </div>
      <div class="flex items-center gap-3">
        <div class="flex items-center gap-2 text-xs">
          <span class="text-gray-600">Rows per page</span>
          <select v-model.number="pageSize" class="border border-gray-300 rounded px-2 py-1 text-xs bg-white">
            <option v-for="s in pageSizes" :key="s" :value="s">{{ s }}</option>
          </select>
        </div>
        <div class="flex items-center gap-1">
          <button @click="firstPage" :disabled="page <= 1" class="px-2 py-1 text-xs border rounded disabled:opacity-50">«</button>
          <button @click="prevPage" :disabled="page <= 1" class="px-2 py-1 text-xs border rounded disabled:opacity-50">‹</button>
          <span class="text-xs text-gray-700 px-2">Page {{ page }} of {{ totalPages }}</span>
          <button @click="nextPage" :disabled="page >= totalPages" class="px-2 py-1 text-xs border rounded disabled:opacity-50">›</button>
          <button @click="lastPage" :disabled="page >= totalPages" class="px-2 py-1 text-xs border rounded disabled:opacity-50">»</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { QueryResult } from '~/stores/chat'

const props = defineProps<{
  results: QueryResult
}>()

// Sorting state
const sortKey = ref<string | null>(null)
const sortDir = ref<'asc' | 'desc' | null>(null)

const toggleSort = (column: string) => {
  if (sortKey.value !== column) {
    sortKey.value = column
    sortDir.value = 'asc'
  } else if (sortDir.value === 'asc') {
    sortDir.value = 'desc'
  } else if (sortDir.value === 'desc') {
    // cycle to none
    sortKey.value = null
    sortDir.value = null
  } else {
    sortDir.value = 'asc'
  }
}

const compareValues = (a: any, b: any) => {
  // Nulls last in asc
  const isNullA = a === null || a === undefined
  const isNullB = b === null || b === undefined
  if (isNullA && !isNullB) return 1
  if (!isNullA && isNullB) return -1
  if (isNullA && isNullB) return 0

  // Numbers (include numeric strings)
  const numA = typeof a === 'number' ? a : (typeof a === 'string' && a.trim() !== '' && !isNaN(Number(a)) ? Number(a) : null)
  const numB = typeof b === 'number' ? b : (typeof b === 'string' && b.trim() !== '' && !isNaN(Number(b)) ? Number(b) : null)
  if (numA !== null && numB !== null) {
    return numA === numB ? 0 : (numA < numB ? -1 : 1)
  }

  // Dates
  const dateA = typeof a === 'string' ? Date.parse(a) : NaN
  const dateB = typeof b === 'string' ? Date.parse(b) : NaN
  if (!isNaN(dateA) && !isNaN(dateB)) {
    return dateA === dateB ? 0 : (dateA < dateB ? -1 : 1)
  }

  // Booleans
  if (typeof a === 'boolean' && typeof b === 'boolean') {
    return a === b ? 0 : (a ? 1 : -1)
  }

  // Fallback string compare
  const sa = String(a).toLocaleLowerCase()
  const sb = String(b).toLocaleLowerCase()
  return sa === sb ? 0 : (sa < sb ? -1 : 1)
}

const sortedRows = computed(() => {
  const rows = props.results?.rows || []
  if (!sortKey.value || !sortDir.value) return rows
  const key = sortKey.value
  const dir = sortDir.value
  const arr = rows.slice()
  arr.sort((r1, r2) => {
    const cmp = compareValues((r1 as any)[key], (r2 as any)[key])
    return dir === 'asc' ? cmp : -cmp
  })
  return arr
})

// Pagination state
const page = ref(1)
const pageSizes = [10, 25, 50, 100]
const pageSize = ref(25)

const rowCount = computed(() => props.results?.rows?.length || 0)
const totalPages = computed(() => Math.max(1, Math.ceil(rowCount.value / pageSize.value)))
const startIndex = computed(() => Math.min((page.value - 1) * pageSize.value, Math.max(0, rowCount.value - 1)))
const endIndex = computed(() => Math.min(startIndex.value + pageSize.value, rowCount.value))
const displayedRows = computed(() => sortedRows.value.slice(startIndex.value, endIndex.value))

// Navigation handlers
const firstPage = () => (page.value = 1)
const prevPage = () => (page.value = Math.max(1, page.value - 1))
const nextPage = () => (page.value = Math.min(totalPages.value, page.value + 1))
const lastPage = () => (page.value = totalPages.value)

// Reset page when results or page size changes
watch(() => props.results, () => { page.value = 1 })
watch(pageSize, () => { page.value = 1 })
watch([sortKey, sortDir], () => { page.value = 1 })

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
