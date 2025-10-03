<template>
  <div class="card">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-bold text-gray-800">Database Schema</h3>
      <button 
        @click="handleRefresh"
        :disabled="schemaStore.isLoading"
        class="text-sm text-primary-600 hover:text-primary-700 font-medium disabled:opacity-50"
      >
        <svg 
          class="w-5 h-5" 
          :class="{ 'animate-spin': schemaStore.isLoading }"
          fill="none" 
          stroke="currentColor" 
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
      </button>
    </div>

    <div v-if="schemaStore.isLoading" class="text-center py-8">
      <div class="inline-block w-8 h-8 border-4 border-primary-500 border-t-transparent rounded-full animate-spin"></div>
      <p class="text-sm text-gray-600 mt-2">Loading schema...</p>
    </div>

    <div v-else-if="schemaStore.schema" class="space-y-4">
      <!-- Stats -->
      <div class="grid grid-cols-2 gap-3">
        <div class="bg-gradient-primary text-white rounded-lg p-3">
          <p class="text-xs opacity-90">Tables</p>
          <p class="text-2xl font-bold">{{ schemaStore.tableCount }}</p>
        </div>
        <div class="bg-gradient-to-r from-purple-500 to-pink-500 text-white rounded-lg p-3">
          <p class="text-xs opacity-90">Relationships</p>
          <p class="text-2xl font-bold">{{ schemaStore.relationshipCount }}</p>
        </div>
      </div>

      <!-- Tables List -->
      <div>
        <h4 class="text-sm font-semibold text-gray-700 mb-2">Tables</h4>
        <div class="space-y-2 max-h-96 overflow-y-auto">
          <div 
            v-for="table in schemaStore.schema.tables" 
            :key="table.name"
            class="border border-gray-200 rounded-lg overflow-hidden hover:border-primary-300 transition-colors"
          >
            <button
              @click="toggleTable(table.name)"
              class="w-full px-3 py-2 bg-gray-50 hover:bg-gray-100 flex items-center justify-between transition-colors"
            >
              <div class="flex items-center space-x-2">
                <svg class="w-4 h-4 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M3 14h18m-9-4v8m-7 0h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
                <span class="text-sm font-medium text-gray-800">{{ table.name }}</span>
                <span class="text-xs text-gray-500">({{ table.row_count }} rows)</span>
              </div>
              <svg 
                class="w-4 h-4 text-gray-400 transition-transform"
                :class="{ 'rotate-90': expandedTables.has(table.name) }"
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>

            <div v-if="expandedTables.has(table.name)" class="px-3 py-2 bg-white">
              <div class="space-y-1">
                <div 
                  v-for="column in table.columns" 
                  :key="column.name"
                  class="flex items-center justify-between text-xs py-1"
                >
                  <div class="flex items-center space-x-2">
                    <span class="font-mono text-gray-700">{{ column.name }}</span>
                    <span class="text-gray-500">{{ column.type }}</span>
                  </div>
                  <div class="flex items-center space-x-1">
                    <span v-if="column.primary_key" class="px-1.5 py-0.5 bg-yellow-100 text-yellow-700 rounded text-xs font-medium">PK</span>
                    <span v-if="column.foreign_key" class="px-1.5 py-0.5 bg-blue-100 text-blue-700 rounded text-xs font-medium">FK</span>
                    <span v-if="!column.nullable" class="px-1.5 py-0.5 bg-red-100 text-red-700 rounded text-xs font-medium">NOT NULL</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="text-center py-8">
      <svg class="w-12 h-12 text-gray-400 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
      </svg>
      <p class="text-sm text-gray-600">No schema loaded</p>
    </div>
  </div>
</template>

<script setup lang="ts">
const schemaStore = useSchemaStore()
const api = useApi()
const expandedTables = ref(new Set<string>())

const toggleTable = (tableName: string) => {
  if (expandedTables.value.has(tableName)) {
    expandedTables.value.delete(tableName)
  } else {
    expandedTables.value.add(tableName)
  }
}

const handleRefresh = async () => {
  schemaStore.setLoading(true)
  const result = await api.refreshSchema()
  if (result.success) {
    const schemaResult = await api.getSchema()
    if (schemaResult.success && schemaResult.data) {
      schemaStore.setSchema(schemaResult.data.schema)
    }
  }
  schemaStore.setLoading(false)
}

// Load schema on mount
onMounted(async () => {
  if (!schemaStore.isSchemaLoaded) {
    schemaStore.setLoading(true)
    const result = await api.getSchema()
    if (result.success && result.data) {
      schemaStore.setSchema(result.data.schema)
    }
    schemaStore.setLoading(false)
  }
})
</script>
