import { defineStore } from 'pinia'

export interface Message {
  id: string
  role: 'user' | 'assistant'
  content: string
  sql?: string
  results?: QueryResult
  reasoning?: ReasoningStep[]
  timestamp: Date
  error?: string
}

export interface QueryResult {
  columns: string[]
  rows: Record<string, any>[]
  count: number
}

export interface ReasoningStep {
  step: number
  action: string
  observation: string
  thought: string
}

export interface QueryHistoryItem {
  id: string
  question: string
  sql: string
  timestamp: Date
  success: boolean
}

export const useChatStore = defineStore('chat', {
  state: () => ({
    messages: [] as Message[],
    queryHistory: [] as QueryHistoryItem[],
    isLoading: false,
    currentQuery: ''
  }),

  actions: {
    addMessage(message: Omit<Message, 'id' | 'timestamp'>) {
      this.messages.push({
        ...message,
        id: crypto.randomUUID(),
        timestamp: new Date()
      })
    },

    addToHistory(item: Omit<QueryHistoryItem, 'id' | 'timestamp'>) {
      this.queryHistory.unshift({
        ...item,
        id: crypto.randomUUID(),
        timestamp: new Date()
      })
      
      // Keep only last 50 items
      if (this.queryHistory.length > 50) {
        this.queryHistory = this.queryHistory.slice(0, 50)
      }
    },

    clearMessages() {
      this.messages = []
    },

    clearHistory() {
      this.queryHistory = []
    },

    setLoading(loading: boolean) {
      this.isLoading = loading
    },

    setCurrentQuery(query: string) {
      this.currentQuery = query
    }
  },

  getters: {
    lastMessage: (state) => state.messages[state.messages.length - 1],
    messageCount: (state) => state.messages.length,
    recentHistory: (state) => state.queryHistory.slice(0, 10)
  }
})
