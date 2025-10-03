<template>
  <div v-if="show" class="fixed inset-0 z-50 overflow-y-auto">
    <!-- Backdrop -->
    <div class="fixed inset-0 bg-black bg-opacity-50 transition-opacity" @click="$emit('close')"></div>
    
    <!-- Modal -->
    <div class="flex min-h-screen items-center justify-center p-4">
      <div class="relative bg-white rounded-2xl shadow-2xl max-w-2xl w-full p-8 fade-in">
        <!-- Header -->
        <div class="text-center mb-8">
          <div class="w-16 h-16 bg-gradient-primary rounded-2xl flex items-center justify-center mx-auto mb-4">
            <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" />
            </svg>
          </div>
          <h2 class="text-3xl font-bold gradient-text mb-2">Connect to Database</h2>
          <p class="text-gray-600">Choose your database type and enter credentials</p>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleConnect" class="space-y-6">
          <!-- Database Type -->
          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-2">Database Type</label>
            <div class="grid grid-cols-3 gap-3">
              <button
                v-for="type in dbTypes"
                :key="type.value"
                type="button"
                @click="form.type = type.value"
                :class="[
                  'p-4 border-2 rounded-lg transition-all',
                  form.type === type.value
                    ? 'border-primary-500 bg-primary-50'
                    : 'border-gray-200 hover:border-primary-300'
                ]"
              >
                <div class="text-2xl mb-2">{{ type.icon }}</div>
                <div class="text-sm font-semibold">{{ type.label }}</div>
              </button>
            </div>
          </div>

          <!-- PostgreSQL / MySQL Fields -->
          <div v-if="form.type !== 'sqlite3'" class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-semibold text-gray-700 mb-2">Host</label>
                <input
                  v-model="form.host"
                  type="text"
                  placeholder="localhost"
                  class="input-field"
                  required
                />
              </div>
              <div>
                <label class="block text-sm font-semibold text-gray-700 mb-2">Port</label>
                <input
                  v-model.number="form.port"
                  type="number"
                  :placeholder="form.type === 'postgres' ? '5432' : '3306'"
                  class="input-field"
                />
              </div>
            </div>

            <div>
              <label class="block text-sm font-semibold text-gray-700 mb-2">Database Name</label>
              <input
                v-model="form.database"
                type="text"
                placeholder="my_database"
                class="input-field"
                required
              />
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-semibold text-gray-700 mb-2">Username</label>
                <input
                  v-model="form.user"
                  type="text"
                  placeholder="postgres"
                  class="input-field"
                  required
                />
              </div>
              <div>
                <label class="block text-sm font-semibold text-gray-700 mb-2">Password</label>
                <input
                  v-model="form.password"
                  type="password"
                  placeholder="••••••••"
                  class="input-field"
                />
              </div>
            </div>

            <div v-if="form.type === 'postgres'">
              <label class="block text-sm font-semibold text-gray-700 mb-2">SSL Mode</label>
              <select v-model="form.sslmode" class="input-field">
                <option value="disable">Disable</option>
                <option value="require">Require</option>
                <option value="verify-ca">Verify CA</option>
                <option value="verify-full">Verify Full</option>
              </select>
            </div>
          </div>

          <!-- SQLite Fields -->
          <div v-else>
            <label class="block text-sm font-semibold text-gray-700 mb-2">Database File Path</label>
            <input
              v-model="form.path"
              type="text"
              placeholder="./database.db"
              class="input-field"
              required
            />
            <p class="text-xs text-gray-500 mt-1">Enter the path to your SQLite database file</p>
          </div>

          <!-- Error Message -->
          <div v-if="error" class="bg-red-50 border-l-4 border-red-500 p-4 rounded">
            <div class="flex items-start">
              <svg class="w-5 h-5 text-red-500 mt-0.5 mr-2 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <div>
                <p class="text-sm font-semibold text-red-800">Connection Failed</p>
                <p class="text-sm text-red-600 mt-1">{{ error }}</p>
              </div>
            </div>
          </div>

          <!-- Success Message -->
          <div v-if="testSuccess" class="bg-green-50 border-l-4 border-green-500 p-4 rounded">
            <div class="flex items-start">
              <svg class="w-5 h-5 text-green-500 mt-0.5 mr-2 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <div>
                <p class="text-sm font-semibold text-green-800">Connection Successful!</p>
                <p class="text-sm text-green-600 mt-1">Found {{ testSuccess.tables }} tables in the database</p>
              </div>
            </div>
          </div>

          <!-- Actions -->
          <div class="flex items-center space-x-3 pt-4">
            <button
              type="button"
              @click="handleTest"
              :disabled="loading"
              class="flex-1 px-6 py-3 border-2 border-primary-500 text-primary-600 font-semibold rounded-lg hover:bg-primary-50 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
            >
              <span v-if="!testing">Test Connection</span>
              <span v-else class="flex items-center justify-center">
                <svg class="animate-spin h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Testing...
              </span>
            </button>
            <button
              type="submit"
              :disabled="loading || !testSuccess"
              class="flex-1 px-6 py-3 bg-gradient-primary text-white font-semibold rounded-lg hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed transition-all"
            >
              <span v-if="!connecting">Connect & Start</span>
              <span v-else class="flex items-center justify-center">
                <svg class="animate-spin h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Connecting...
              </span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface ConnectionTestResult {
  success: boolean
  message?: string
  tables?: number
}

const props = defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  close: []
  connected: [data: any]
}>()

const api = useApi()

const dbTypes = [
  { value: 'postgres', label: 'PostgreSQL', icon: '🐘' },
  { value: 'mysql', label: 'MySQL', icon: '🐬' },
  { value: 'sqlite3', label: 'SQLite', icon: '📁' }
]

const form = reactive({
  type: 'postgres',
  host: 'localhost',
  port: 0,
  database: '',
  user: '',
  password: '',
  sslmode: 'disable',
  path: ''
})

const testing = ref(false)
const connecting = ref(false)
const error = ref('')
const testSuccess = ref<ConnectionTestResult | null>(null)

const loading = computed(() => testing.value || connecting.value)

const handleTest = async () => {
  testing.value = true
  error.value = ''
  testSuccess.value = null

  const result = await api.testConnection(form)
  testing.value = false

  if (result.success && result.data) {
    const data = result.data as ConnectionTestResult
    if (data.success) {
      testSuccess.value = data
    } else {
      error.value = data.message || 'Connection test failed'
    }
  } else {
    error.value = (result.error as string) || 'Failed to test connection'
  }
}

const handleConnect = async () => {
  if (!testSuccess.value) {
    error.value = 'Please test the connection first'
    return
  }

  connecting.value = true
  error.value = ''

  const result = await api.connect(form)
  connecting.value = false

  if (result.success && result.data) {
    emit('connected', result.data)
  } else {
    error.value = (result.error as string) || 'Failed to connect'
  }
}
</script>
