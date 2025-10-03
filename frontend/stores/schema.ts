import { defineStore } from 'pinia'

export interface Column {
  name: string
  type: string
  nullable: boolean
  primary_key: boolean
  foreign_key?: string
}

export interface TableInfo {
  name: string
  columns: Column[]
  row_count: number
  description?: string
}

export interface TableRelationship {
  from_table: string
  from_column: string
  to_table: string
  to_column: string
}

export interface SchemaInfo {
  tables: TableInfo[]
  relationships: TableRelationship[]
  summary: string
}

export const useSchemaStore = defineStore('schema', {
  state: () => ({
    schema: null as SchemaInfo | null,
    isLoading: false,
    lastUpdated: null as Date | null,
    selectedTable: null as string | null
  }),

  actions: {
    setSchema(schema: SchemaInfo) {
      this.schema = schema
      this.lastUpdated = new Date()
    },

    setLoading(loading: boolean) {
      this.isLoading = loading
    },

    selectTable(tableName: string | null) {
      this.selectedTable = tableName
    },

    clearSchema() {
      this.schema = null
      this.lastUpdated = null
      this.selectedTable = null
    }
  },

  getters: {
    tableCount: (state) => state.schema?.tables.length || 0,
    relationshipCount: (state) => state.schema?.relationships.length || 0,
    
    getTable: (state) => (tableName: string) => {
      return state.schema?.tables.find(t => t.name === tableName)
    },

    tableNames: (state) => state.schema?.tables.map(t => t.name) || [],

    isSchemaLoaded: (state) => state.schema !== null
  }
})
