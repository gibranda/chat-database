<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 via-white to-primary-50">
    <!-- Connection Modal -->
    <ConnectionModal 
      :show="!isConnected" 
      @close="() => {}" 
      @connected="handleConnected"
    />
    
    <!-- Header -->
    <header class="bg-white border-b border-gray-200 sticky top-0 z-10 shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-3">
            <div class="w-10 h-10 bg-gradient-primary rounded-lg flex items-center justify-center">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
              </svg>
            </div>
            <div>
              <h1 class="text-xl font-bold gradient-text">AI Database Agent</h1>
              <p class="text-xs text-gray-500">Natural Language to SQL</p>
            </div>
          </div>

          <div class="flex items-center space-x-4">
            <!-- Database Connection Info -->
            <div v-if="isConnected" class="flex items-center space-x-2 px-3 py-1.5 bg-green-50 rounded-lg border border-green-200">
              <div class="w-2 h-2 rounded-full bg-green-500"></div>
              <span class="text-sm text-green-700 font-medium">
                {{ connectionInfo?.database || 'Connected' }}
              </span>
            </div>

            <!-- Backend Status -->
            <div v-else class="flex items-center space-x-2">
              <div 
                class="w-2 h-2 rounded-full"
                :class="backendConnected ? 'bg-yellow-500' : 'bg-red-500 pulse-slow'"
              ></div>
              <span class="text-sm text-gray-600">
                {{ backendConnected ? 'No Database' : 'Backend Offline' }}
              </span>
            </div>

            <!-- Change Connection Button -->
            <button
              v-if="isConnected"
              @click="handleChangeConnection"
              class="text-sm text-primary-600 hover:text-primary-800 font-medium flex items-center space-x-1"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
              </svg>
              <span>Change</span>
            </button>

            <!-- Disconnect Button -->
            <button
              v-if="isConnected"
              @click="handleDisconnect"
              class="text-sm text-red-600 hover:text-red-800 font-medium flex items-center space-x-1"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
              <span>Disconnect</span>
            </button>

            <!-- Clear Chat Button -->
            <button
              @click="handleClearChat"
              class="text-sm text-gray-600 hover:text-gray-800 font-medium"
              :disabled="chatStore.messageCount === 0"
            >
              Clear Chat
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
        <!-- Sidebar -->
        <aside class="lg:col-span-1 space-y-4">
          <SchemaExplorer />
          <QueryHistory @select-query="handleSelectQuery" />
        </aside>

        <!-- Chat Area -->
        <main class="lg:col-span-3">
          <div class="card flex flex-col">
            <!-- Welcome Message -->
            <div v-if="chatStore.messageCount === 0" class="flex-1 flex items-center justify-center">
              <div class="text-center max-w-2xl px-4">
                <div class="w-20 h-20 bg-gradient-primary rounded-2xl flex items-center justify-center mx-auto mb-6">
                  <svg class="w-12 h-12 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
                  </svg>
                </div>
                <h2 class="text-3xl font-bold gradient-text mb-4">
                  Welcome to AI Database Agent
                </h2>
                <p class="text-gray-600 mb-8">
                  Ask questions about your database in natural language. I'll convert them to SQL and execute them for you.
                </p>

                <!-- Example Questions -->
                <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                  <button
                    v-for="example in exampleQuestions"
                    :key="example"
                    @click="handleSelectQuery(example)"
                    class="text-left p-4 border-2 border-gray-200 rounded-lg hover:border-primary-300 hover:bg-primary-50 transition-all group"
                  >
                    <div class="flex items-start space-x-3">
                      <svg class="w-5 h-5 text-primary-500 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                      </svg>
                      <p class="text-sm text-gray-700 group-hover:text-gray-900">{{ example }}</p>
                    </div>
                  </button>
                </div>
              </div>
            </div>

            <!-- Messages -->
            <div v-else class="flex-1 mb-4 space-y-4 px-2" ref="messagesContainer">
              <ChatMessage 
                v-for="message in chatStore.messages" 
                :key="message.id"
                :message="message"
              />

              <!-- Loading Indicator -->
              <div v-if="chatStore.isLoading" class="flex items-start space-x-3">
                <div class="flex-shrink-0">
                  <div class="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center">
                    <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                    </svg>
                  </div>
                </div>
                <div class="bg-white border border-gray-200 rounded-2xl rounded-tl-none px-4 py-3">
                  <div class="flex items-center space-x-2">
                    <div class="w-2 h-2 bg-primary-500 rounded-full animate-bounce"></div>
                    <div class="w-2 h-2 bg-primary-500 rounded-full animate-bounce" style="animation-delay: 0.1s"></div>
                    <div class="w-2 h-2 bg-primary-500 rounded-full animate-bounce" style="animation-delay: 0.2s"></div>
                    <span class="text-sm text-gray-600 ml-2">Thinking...</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Input Area -->
            <div class="border-t border-gray-200 pt-4 flex-shrink-0">
              <form @submit.prevent="handleSubmit" class="flex items-start space-x-2">
                <div class="flex-1">
                  <textarea
                    v-model="inputMessage"
                    @keydown.enter.exact.prevent="handleSubmit"
                    placeholder="Ask a question about your database..."
                    rows="2"
                    class="input-field resize-none text-sm"
                    :disabled="chatStore.isLoading || !backendConnected"
                  ></textarea>
                  <p class="text-xs text-gray-500 mt-1">Press Enter to send, Shift+Enter for new line</p>
                </div>
                <button
                  type="submit"
                  :disabled="!inputMessage.trim() || chatStore.isLoading || !backendConnected"
                  class="px-4 py-2 bg-gradient-primary text-white font-semibold rounded-lg hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 flex items-center space-x-1 h-[42px]"
                >
                  <span class="text-sm">Send</span>
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
                  </svg>
                </button>
              </form>
            </div>
          </div>
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const chatStore = useChatStore()
const schemaStore = useSchemaStore()
const api = useApi()
const inputMessage = ref('')
const backendConnected = ref(false)
const isConnected = ref(false)
const connectionInfo = ref<any>(null)
const messagesContainer = ref<HTMLElement | null>(null)

const exampleQuestions = [
  'Show me all tables in the database',
  'What are the top 10 customers by total orders?',
  'How many products do we have in each category?',
  'Show me sales trends for the last 6 months'
]

// Check backend health
const checkBackend = async () => {
  const result = await api.checkHealth()
  backendConnected.value = result.success
}

// Handle successful connection
const handleConnected = async (data: any) => {
  connectionInfo.value = data
  isConnected.value = true
  
  // Check backend health after connection
  await checkBackend()
  
  // Load schema
  schemaStore.setLoading(true)
  const schemaResult = await api.getSchema()
  if (schemaResult.success && schemaResult.data) {
    schemaStore.setSchema((schemaResult.data as any).schema)
  }
  schemaStore.setLoading(false)
}

// Handle query submission
const handleSubmit = async () => {
  if (!inputMessage.value.trim() || chatStore.isLoading) return

  const question = inputMessage.value.trim()
  inputMessage.value = ''

  // Add user message
  chatStore.addMessage({
    role: 'user',
    content: question
  })

  // Scroll to bottom
  nextTick(() => {
    scrollToBottom()
  })

  // Send query
  chatStore.setLoading(true)
  const result = await api.sendQuery(question)
  chatStore.setLoading(false)

  if (result.success && result.data) {
    const data = result.data as any
    
    // Check if query was successful or has error
    if (data.success === false || data.error) {
      // Query failed on backend
      chatStore.addMessage({
        role: 'assistant',
        content: '',
        error: data.error || 'Query execution failed'
      })
      
      chatStore.addToHistory({
        question,
        sql: data.sql || '',
        success: false
      })
    } else {
      // Query successful
      chatStore.addMessage({
        role: 'assistant',
        content: data.answer || 'Query executed successfully.',
        sql: data.sql,
        results: data.results,
        reasoning: data.reasoning
      })

      // Add to history
      chatStore.addToHistory({
        question,
        sql: data.sql || '',
        success: true
      })
    }
  } else {
    // API call failed
    chatStore.addMessage({
      role: 'assistant',
      content: '',
      error: result.error as string || 'Failed to connect to backend'
    })

    chatStore.addToHistory({
      question,
      sql: '',
      success: false
    })
  }

  // Scroll to bottom
  nextTick(() => {
    scrollToBottom()
  })
}

const handleSelectQuery = (query: string) => {
  inputMessage.value = query
  nextTick(() => {
    handleSubmit()
  })
}

const handleClearChat = () => {
  if (confirm('Are you sure you want to clear the chat?')) {
    chatStore.clearMessages()
  }
}

const handleDisconnect = async () => {
  if (!confirm('Are you sure you want to disconnect from the database?')) {
    return
  }

  const result = await api.disconnect()
  if (result.success) {
    // Reset state
    isConnected.value = false
    connectionInfo.value = null
    
    // Clear stores
    chatStore.clearMessages()
    schemaStore.clearSchema()
    
    // Show success message
    console.log('Disconnected successfully')
  } else {
    alert('Failed to disconnect: ' + result.error)
  }
}

const handleChangeConnection = () => {
  if (confirm('Disconnect from current database and connect to a new one?')) {
    // Reset connection state to show modal
    isConnected.value = false
    connectionInfo.value = null
    
    // Clear stores
    chatStore.clearMessages()
    schemaStore.clearSchema()
  }
}

const scrollToBottom = () => {
  const el = messagesContainer.value
  if (el) {
    // If the container has its own scrollbar, use it
    if (el.scrollHeight > el.clientHeight) {
      el.scrollTop = el.scrollHeight
      return
    }
  }
  // Otherwise, scroll the window (page-level scroll)
  if (typeof window !== 'undefined') {
    window.scrollTo({ top: document.documentElement.scrollHeight, behavior: 'smooth' })
  }
}

// Check backend on mount
onMounted(() => {
  checkBackend()
  // Don't start interval until connected
  watch(isConnected, (connected) => {
    if (connected) {
      setInterval(checkBackend, 10000) // Check every 10 seconds
    }
  })
})

useHead({
  title: 'Chat'
})
</script>
