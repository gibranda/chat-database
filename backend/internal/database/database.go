package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// BuildPostgresConnString builds PostgreSQL connection string
func BuildPostgresConnString(host string, port int, dbname, user, password, sslmode string) string {
	if port == 0 {
		port = 5432
	}
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		host, port, dbname, user, password, sslmode)
}

// BuildMySQLConnString builds MySQL connection string
func BuildMySQLConnString(host string, port int, dbname, user, password string) string {
	if port == 0 {
		port = 3306
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		user, password, host, port, dbname)
}

type Database struct {
	db     *sql.DB
	dbType string
}

type TableInfo struct {
	Name        string   `json:"name"`
	Columns     []Column `json:"columns"`
	RowCount    int64    `json:"row_count"`
	Description string   `json:"description,omitempty"`
}

type Column struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Nullable   bool   `json:"nullable"`
	PrimaryKey bool   `json:"primary_key"`
	ForeignKey string `json:"foreign_key,omitempty"`
}

type QueryResult struct {
	Columns []string                 `json:"columns"`
	Rows    []map[string]interface{} `json:"rows"`
	Count   int                      `json:"count"`
}

type SchemaInfo struct {
	Tables        []TableInfo           `json:"tables"`
	Relationships []TableRelationship   `json:"relationships"`
	Summary       string                `json:"summary"`
}

type TableRelationship struct {
	FromTable  string `json:"from_table"`
	FromColumn string `json:"from_column"`
	ToTable    string `json:"to_table"`
	ToColumn   string `json:"to_column"`
}

func New(dbType, connectionString string) (*Database, error) {
	db, err := sql.Open(dbType, connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{
		db:     db,
		dbType: dbType,
	}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetTables() ([]string, error) {
	var query string
	switch d.dbType {
	case "postgres":
		query = `
			SELECT table_name 
			FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_type = 'BASE TABLE'
			ORDER BY table_name
		`
	case "mysql":
		query = `
			SELECT table_name 
			FROM information_schema.tables 
			WHERE table_schema = DATABASE()
			AND table_type = 'BASE TABLE'
			ORDER BY table_name
		`
	case "sqlite3":
		query = `
			SELECT name 
			FROM sqlite_master 
			WHERE type='table' 
			AND name NOT LIKE 'sqlite_%'
			ORDER BY name
		`
	default:
		return nil, fmt.Errorf("unsupported database type: %s", d.dbType)
	}

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

func (d *Database) GetTableInfo(tableName string) (*TableInfo, error) {
	columns, err := d.getColumns(tableName)
	if err != nil {
		return nil, err
	}

	rowCount, err := d.getRowCount(tableName)
	if err != nil {
		rowCount = 0 // Non-critical error
	}

	return &TableInfo{
		Name:     tableName,
		Columns:  columns,
		RowCount: rowCount,
	}, nil
}

func (d *Database) getColumns(tableName string) ([]Column, error) {
	var query string
	switch d.dbType {
	case "postgres":
		query = fmt.Sprintf(`
			SELECT 
				c.column_name,
				c.data_type,
				c.is_nullable,
				CASE WHEN pk.column_name IS NOT NULL THEN true ELSE false END as is_primary,
				COALESCE(fk.foreign_table || '.' || fk.foreign_column, '') as foreign_key
			FROM information_schema.columns c
			LEFT JOIN (
				SELECT ku.column_name
				FROM information_schema.table_constraints tc
				JOIN information_schema.key_column_usage ku
					ON tc.constraint_name = ku.constraint_name
				WHERE tc.table_name = '%s'
				AND tc.constraint_type = 'PRIMARY KEY'
			) pk ON c.column_name = pk.column_name
			LEFT JOIN (
				SELECT
					kcu.column_name,
					ccu.table_name AS foreign_table,
					ccu.column_name AS foreign_column
				FROM information_schema.table_constraints AS tc
				JOIN information_schema.key_column_usage AS kcu
					ON tc.constraint_name = kcu.constraint_name
				JOIN information_schema.constraint_column_usage AS ccu
					ON ccu.constraint_name = tc.constraint_name
				WHERE tc.table_name = '%s'
				AND tc.constraint_type = 'FOREIGN KEY'
			) fk ON c.column_name = fk.column_name
			WHERE c.table_name = '%s'
			ORDER BY c.ordinal_position
		`, tableName, tableName, tableName)
	case "mysql":
		query = fmt.Sprintf(`
			SELECT 
				c.column_name,
				c.data_type,
				c.is_nullable,
				CASE WHEN c.column_key = 'PRI' THEN true ELSE false END as is_primary,
				COALESCE(
					CONCAT(
						k.referenced_table_name, 
						'.', 
						k.referenced_column_name
					), 
					''
				) as foreign_key
			FROM information_schema.columns c
			LEFT JOIN information_schema.key_column_usage k
				ON c.table_name = k.table_name
				AND c.column_name = k.column_name
				AND k.referenced_table_name IS NOT NULL
			WHERE c.table_name = '%s'
			AND c.table_schema = DATABASE()
			ORDER BY c.ordinal_position
		`, tableName)
	case "sqlite3":
		query = fmt.Sprintf("PRAGMA table_info(%s)", tableName)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", d.dbType)
	}

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns: %w", err)
	}
	defer rows.Close()

	var columns []Column
	if d.dbType == "sqlite3" {
		for rows.Next() {
			var cid int
			var name, colType string
			var notNull, pk int
			var dfltValue interface{}
			
			if err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
				return nil, err
			}
			
			columns = append(columns, Column{
				Name:       name,
				Type:       colType,
				Nullable:   notNull == 0,
				PrimaryKey: pk == 1,
			})
		}
	} else {
		for rows.Next() {
			var col Column
			var nullable string
			
			if err := rows.Scan(&col.Name, &col.Type, &nullable, &col.PrimaryKey, &col.ForeignKey); err != nil {
				return nil, err
			}
			
			col.Nullable = nullable == "YES"
			columns = append(columns, col)
		}
	}

	return columns, nil
}

func (d *Database) getRowCount(tableName string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	var count int64
	err := d.db.QueryRow(query).Scan(&count)
	return count, err
}

func (d *Database) ExecuteQuery(query string, maxResults int) (*QueryResult, error) {
	// Clean and prepare query
	query = strings.TrimSpace(query)
	
	// Remove trailing semicolon if present
	query = strings.TrimRight(query, ";")
	
	// Add LIMIT if not present
	upperQuery := strings.ToUpper(query)
	if !strings.Contains(upperQuery, "LIMIT") && maxResults > 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, maxResults)
	}
	
	// Add semicolon back at the end
	query = query + ";"

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	return &QueryResult{
		Columns: columns,
		Rows:    results,
		Count:   len(results),
	}, nil
}

func (d *Database) GetFullSchema() (*SchemaInfo, error) {
	tables, err := d.GetTables()
	if err != nil {
		return nil, err
	}

	var tableInfos []TableInfo
	var relationships []TableRelationship

	for _, tableName := range tables {
		info, err := d.GetTableInfo(tableName)
		if err != nil {
			continue // Skip tables we can't read
		}
		tableInfos = append(tableInfos, *info)

		// Extract relationships
		for _, col := range info.Columns {
			if col.ForeignKey != "" {
				parts := strings.Split(col.ForeignKey, ".")
				if len(parts) == 2 {
					relationships = append(relationships, TableRelationship{
						FromTable:  tableName,
						FromColumn: col.Name,
						ToTable:    parts[0],
						ToColumn:   parts[1],
					})
				}
			}
		}
	}

	summary := d.generateSchemaSummary(tableInfos, relationships)

	return &SchemaInfo{
		Tables:        tableInfos,
		Relationships: relationships,
		Summary:       summary,
	}, nil
}

func (d *Database) generateSchemaSummary(tables []TableInfo, relationships []TableRelationship) string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("Database contains %d tables:\n\n", len(tables)))
	
	for _, table := range tables {
		sb.WriteString(fmt.Sprintf("Table: %s (%d rows)\n", table.Name, table.RowCount))
		sb.WriteString("Columns:\n")
		for _, col := range table.Columns {
			markers := []string{}
			if col.PrimaryKey {
				markers = append(markers, "PK")
			}
			if col.ForeignKey != "" {
				markers = append(markers, fmt.Sprintf("FK->%s", col.ForeignKey))
			}
			if !col.Nullable {
				markers = append(markers, "NOT NULL")
			}
			
			markerStr := ""
			if len(markers) > 0 {
				markerStr = fmt.Sprintf(" [%s]", strings.Join(markers, ", "))
			}
			
			sb.WriteString(fmt.Sprintf("  - %s: %s%s\n", col.Name, col.Type, markerStr))
		}
		sb.WriteString("\n")
	}
	
	if len(relationships) > 0 {
		sb.WriteString(fmt.Sprintf("Relationships (%d):\n", len(relationships)))
		for _, rel := range relationships {
			sb.WriteString(fmt.Sprintf("  - %s.%s -> %s.%s\n",
				rel.FromTable, rel.FromColumn, rel.ToTable, rel.ToColumn))
		}
	}
	
	return sb.String()
}

func (d *Database) ValidateSQL(query string) error {
	// Basic SQL injection prevention
	upperQuery := strings.ToUpper(strings.TrimSpace(query))
	
	// Check for dangerous patterns
	dangerous := []string{
		"DROP TABLE",
		"DROP DATABASE",
		"TRUNCATE",
		"ALTER TABLE",
		"CREATE TABLE",
		"CREATE DATABASE",
	}
	
	for _, pattern := range dangerous {
		if strings.Contains(upperQuery, pattern) {
			return fmt.Errorf("potentially dangerous SQL pattern detected: %s", pattern)
		}
	}
	
	return nil
}

func (d *Database) IsReadOnlyQuery(query string) bool {
	upperQuery := strings.ToUpper(strings.TrimSpace(query))
	return strings.HasPrefix(upperQuery, "SELECT") || 
	       strings.HasPrefix(upperQuery, "SHOW") ||
	       strings.HasPrefix(upperQuery, "DESCRIBE") ||
	       strings.HasPrefix(upperQuery, "EXPLAIN")
}
