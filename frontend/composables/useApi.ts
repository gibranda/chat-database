export const useApi = () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase

  const checkHealth = async () => {
    try {
      const response = await $fetch(`${apiBase}/health`)
      return { success: true, data: response }
    } catch (error) {
      return { success: false, error }
    }
  }

  const sendQuery = async (question: string) => {
    try {
      const response = await $fetch(`${apiBase}/query`, {
        method: 'POST',
        body: { question }
      })
      return { success: true, data: response }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.message || 'Failed to send query' 
      }
    }
  }

  const getSchema = async () => {
    try {
      const response = await $fetch(`${apiBase}/schema`)
      return { success: true, data: response }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.message || 'Failed to fetch schema' 
      }
    }
  }

  const refreshSchema = async () => {
    try {
      const response = await $fetch(`${apiBase}/schema/refresh`, {
        method: 'POST'
      })
      return { success: true, data: response }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.message || 'Failed to refresh schema' 
      }
    }
  }

  const getTables = async () => {
    try {
      const response = await $fetch(`${apiBase}/tables`)
      return { success: true, data: response }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.message || 'Failed to fetch tables' 
      }
    }
  }

  const getTableInfo = async (tableName: string) => {
    try {
      const response = await $fetch(`${apiBase}/tables/${tableName}`)
      return { success: true, data: response }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.message || 'Failed to fetch table info' 
      }
    }
  }

  const clearHistory = async () => {
    try {
      const response = await $fetch(`${apiBase}/history/clear`, {
        method: 'POST'
      })
      return { success: true, data: response }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.message || 'Failed to clear history' 
      }
    }
  }

  const testConnection = async (credentials: any) => {
    try {
      const response = await $fetch(`${apiBase}/connection/test`, {
        method: 'POST',
        body: credentials
      })
      return { success: true, data: response }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.message || 'Failed to test connection' 
      }
    }
  }

  const connect = async (credentials: any) => {
    try {
      const response = await $fetch(`${apiBase}/connection/connect`, {
        method: 'POST',
        body: credentials
      })
      return { success: true, data: response }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.message || 'Failed to connect' 
      }
    }
  }

  return {
    checkHealth,
    sendQuery,
    getSchema,
    refreshSchema,
    getTables,
    getTableInfo,
    clearHistory,
    testConnection,
    connect
  }
}
