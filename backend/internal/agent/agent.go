package agent

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gibranda/chat-with-database/internal/database"
	"github.com/gibranda/chat-with-database/internal/llm"
)

type Agent struct {
	llm                   *llm.OllamaClient
	db                    *database.Database
	maxIterations         int
	enableQueryValidation bool
	readonlyMode          bool
	maxResults            int
	conversationHistory   []llm.ChatMessage
	schemaCache           *database.SchemaInfo
}

type AgentResponse struct {
	Success      bool                   `json:"success"`
	Answer       string                 `json:"answer"`
	SQL          string                 `json:"sql,omitempty"`
	Results      *database.QueryResult  `json:"results,omitempty"`
	Reasoning    []ReasoningStep        `json:"reasoning"`
	Error        string                 `json:"error,omitempty"`
}

type ReasoningStep struct {
	Step        int    `json:"step"`
	Action      string `json:"action"`
	Observation string `json:"observation"`
	Thought     string `json:"thought"`
}

type ToolCall struct {
	Tool      string                 `json:"tool"`
	Arguments map[string]interface{} `json:"arguments"`
}

func NewAgent(
	llmClient *llm.OllamaClient,
	db *database.Database,
	maxIterations int,
	enableQueryValidation bool,
	readonlyMode bool,
	maxResults int,
) *Agent {
	return &Agent{
		llm:                   llmClient,
		db:                    db,
		maxIterations:         maxIterations,
		enableQueryValidation: enableQueryValidation,
		maxResults:            maxResults,
		conversationHistory:   make([]llm.ChatMessage, 0),
	}
}

func (a *Agent) ProcessQuery(question string) (*AgentResponse, error) {
	response := &AgentResponse{
		Success:   false,
		Reasoning: make([]ReasoningStep, 0),
	}

	// Step 1: Get schema if not cached
	if a.schemaCache == nil {
		schema, err := a.db.GetFullSchema()
		if err != nil {
			return nil, fmt.Errorf("failed to get schema: %w", err)
		}
		a.schemaCache = schema
	}

	response.Reasoning = append(response.Reasoning, ReasoningStep{
		Step:        1,
		Action:      "analyze_schema",
		Observation: fmt.Sprintf("Found %d tables in database", len(a.schemaCache.Tables)),
		Thought:     "Understanding database structure",
	})

	// Check if question is about listing tables
	questionLower := strings.ToLower(question)
	if strings.Contains(questionLower, "show") && (strings.Contains(questionLower, "table") || strings.Contains(questionLower, "tabel")) ||
		strings.Contains(questionLower, "list") && (strings.Contains(questionLower, "table") || strings.Contains(questionLower, "tabel")) ||
		strings.Contains(questionLower, "what table") || strings.Contains(questionLower, "apa saja tabel") {
		
		// Direct response for table listing
		var tableNames []string
		for _, table := range a.schemaCache.Tables {
			tableNames = append(tableNames, table.Name)
		}
		
		answer := fmt.Sprintf("Database ini memiliki %d tabel:\n\n", len(tableNames))
		for i, name := range tableNames {
			answer += fmt.Sprintf("%d. **%s**\n", i+1, name)
		}
		answer += "\nAnda bisa bertanya tentang data di tabel-tabel ini. Misalnya: 'tampilkan data dari tabel students' atau 'berapa jumlah data di tabel schools'."
		
		response.Success = true
		response.Answer = answer
		response.SQL = "-- Schema query (no SQL execution needed)"
		
		response.Reasoning = append(response.Reasoning, ReasoningStep{
			Step:        2,
			Action:      "list_tables",
			Observation: fmt.Sprintf("Listed %d tables from schema", len(tableNames)),
			Thought:     "Providing table list from cached schema",
		})
		
		return response, nil
	}

	// Step 2: Plan the query
	planPrompt := a.buildPlanningPrompt(question)
	plan, err := a.llm.Generate(planPrompt, a.getSystemPrompt())
	if err != nil {
		return nil, fmt.Errorf("failed to generate plan: %w", err)
	}

	response.Reasoning = append(response.Reasoning, ReasoningStep{
		Step:        2,
		Action:      "analyze_question",
		Observation: plan,
		Thought:     "Planning query approach",
	})

	// Step 3: Generate SQL
	sqlPrompt := a.buildSQLPrompt(question, plan)
	sqlResponse, err := a.llm.Generate(sqlPrompt, a.getSystemPrompt())
	if err != nil {
		return nil, fmt.Errorf("failed to generate SQL: %w", err)
	}

	// Extract SQL from response
	sql := a.extractSQL(sqlResponse)
	response.SQL = sql

	// Log generated SQL for debugging
	fmt.Printf("Generated SQL: %s\n", sql)

	response.Reasoning = append(response.Reasoning, ReasoningStep{
		Step:        3,
		Action:      "generate_sql",
		Observation: sql,
		Thought:     "Generated SQL query",
	})

	// Step 4: Validate SQL
	if a.enableQueryValidation {
		if err := a.db.ValidateSQL(sql); err != nil {
			response.Error = fmt.Sprintf("SQL validation failed: %v", err)
			return response, nil
		}

		if a.readonlyMode {
			isReadOnly := a.db.IsReadOnlyQuery(sql)
			fmt.Printf("Readonly mode check: SQL=%s, IsReadOnly=%v\n", sql[:min(50, len(sql))], isReadOnly)
			if !isReadOnly {
				response.Error = "Only read-only queries are allowed in readonly mode"
				return response, nil
			}
		}

		response.Reasoning = append(response.Reasoning, ReasoningStep{
			Step:        4,
			Action:      "validate_sql",
			Observation: "SQL validation passed",
			Thought:     "Query is safe to execute",
		})
	}

	// Step 5: Execute query
	results, err := a.db.ExecuteQuery(sql, a.maxResults)
	if err != nil {
		// Try to fix the query
		fixedSQL, fixErr := a.fixQuery(sql, err.Error(), question)
		if fixErr == nil {
			results, err = a.db.ExecuteQuery(fixedSQL, a.maxResults)
			if err == nil {
				response.SQL = fixedSQL
				response.Reasoning = append(response.Reasoning, ReasoningStep{
					Step:        5,
					Action:      "fix_and_retry",
					Observation: "Query fixed and executed successfully",
					Thought:     "Corrected SQL syntax error",
				})
			}
		}
		
		if err != nil {
			response.Error = fmt.Sprintf("Query execution failed: %v", err)
			return response, nil
		}
	}

	response.Results = results
	response.Reasoning = append(response.Reasoning, ReasoningStep{
		Step:        len(response.Reasoning) + 1,
		Action:      "execute_query",
		Observation: fmt.Sprintf("Retrieved %d rows", results.Count),
		Thought:     "Query executed successfully",
	})

	// Step 6: Generate natural language answer
	answerPrompt := a.buildAnswerPrompt(question, sql, results)
	answer, err := a.llm.Generate(answerPrompt, a.getSystemPrompt())
	if err != nil {
		answer = "Query executed successfully. See results below."
	}

	response.Answer = answer
	response.Success = true

	response.Reasoning = append(response.Reasoning, ReasoningStep{
		Step:        len(response.Reasoning) + 1,
		Action:      "generate_answer",
		Observation: answer,
		Thought:     "Formulated natural language response",
	})

	// Update conversation history
	a.conversationHistory = append(a.conversationHistory,
		llm.ChatMessage{Role: "user", Content: question},
		llm.ChatMessage{Role: "assistant", Content: answer},
	)

	return response, nil
}

func (a *Agent) getSystemPrompt() string {
	return `You are a friendly and helpful AI database assistant named "DB Assistant". Your role is to help users explore and understand their data through natural conversation.

Your personality:
- Friendly, approachable, and patient
- Explain technical concepts in simple terms
- Proactive in offering insights and suggestions
- Always respectful and professional

Key responsibilities:
1. Understand database schemas and relationships intuitively
2. Convert natural language questions to accurate SQL queries
3. Provide meaningful insights from data, not just raw results
4. Explain findings in clear, conversational language
5. Suggest related queries that might be helpful

Guidelines for SQL generation:
- Always use proper SQL syntax for the database type
- Consider table relationships and foreign keys carefully
- Use JOINs when querying related tables
- Apply appropriate filters, aggregations, and sorting
- Limit results to reasonable amounts (default 100 rows)
- Handle edge cases and NULL values gracefully
- Optimize queries for performance

When explaining results:
- Start with a brief, friendly summary
- Highlight interesting patterns or anomalies
- Provide context and meaning to numbers
- Suggest follow-up questions or analyses
- Use analogies when helpful
- Keep explanations concise but informative

Response style:
- Use conversational, natural language
- Avoid overly technical jargon unless necessary
- Be encouraging and positive
- Show enthusiasm about interesting findings
- Acknowledge limitations honestly`
}

func (a *Agent) buildPlanningPrompt(question string) string {
	return fmt.Sprintf(`Given this database schema:

%s

User question: "%s"

As a helpful database assistant, analyze this question and create a brief plan for answering it.

Consider:
1. Which ACTUAL tables from the schema contain the relevant data? (Use table names like: students, schools, teachers, etc.)
2. What relationships exist between these tables?
3. What aggregations, calculations, or transformations are needed?
4. What filters or conditions should be applied?
5. How should the results be sorted or limited?

IMPORTANT: 
- Focus on the actual data tables shown in the schema
- Do NOT plan to use information_schema or system tables unless specifically asked
- If asked about "tables" or "what data exists", refer to the table list in the schema

Provide a clear, concise plan (2-3 sentences) that shows you understand the user's intent and which actual tables to query.`, a.schemaCache.Summary, question)
}

func (a *Agent) buildSQLPrompt(question, plan string) string {
	return fmt.Sprintf(`Database schema:

%s

User question: "%s"

Plan: %s

Generate a SQL query to answer this question. 

IMPORTANT RULES:
- Return ONLY the SQL query, no explanations or markdown
- Do NOT wrap the SQL in backticks or quotes
- Do NOT use markdown code blocks
- Just return the raw SQL query
- End the query with a semicolon

QUERY GUIDELINES:
- Use the actual table names from the schema above
- Do NOT use information_schema or system tables unless specifically asked
- When asked about "tables" or "data", query the actual data tables (students, schools, etc.)
- Use appropriate JOINs when querying related tables
- Add LIMIT clause to prevent returning too many rows (default: 100)
- Use meaningful column aliases for better readability

Examples:
- For "show tables": List the table names you see in the schema
- For "show data": SELECT * FROM actual_table_name LIMIT 10;
- For "count records": SELECT COUNT(*) FROM actual_table_name;`, a.schemaCache.Summary, question, plan)
}

func (a *Agent) buildAnswerPrompt(question string, sql string, results *database.QueryResult) string {
	// Limit result preview to first 5 rows
	previewRows := results.Rows
	if len(previewRows) > 5 {
		previewRows = previewRows[:5]
	}

	resultsJSON, _ := json.MarshalIndent(previewRows, "", "  ")

	return fmt.Sprintf(`User asked: "%s"

SQL query executed:
%s

Results (%d rows total, showing first %d):
%s

As a friendly database assistant, provide a helpful and insightful response in Indonesian language.

Your response should include:
1. **Ringkasan Singkat**: Jelaskan apa yang ditunjukkan oleh data ini (2-3 kalimat)
2. **Insight Utama**: Temuan atau pola menarik yang terlihat dari data
3. **Angka Penting**: Statistik atau angka kunci yang perlu diperhatikan
4. **Observasi Tambahan**: Hal menarik lainnya atau saran analisis lanjutan

Gaya penulisan:
- Gunakan bahasa yang ramah dan mudah dipahami
- Fokus pada insight yang actionable
- Berikan konteks pada angka-angka
- Jika ada yang menarik atau tidak biasa, sebutkan
- Akhiri dengan saran pertanyaan lanjutan jika relevan

Contoh format:
"Berdasarkan data yang saya temukan, [ringkasan]. Yang menarik adalah [insight]. Secara angka, [statistik]. Anda mungkin juga ingin melihat [saran]."`,
		question, sql, results.Count, len(previewRows), string(resultsJSON))
}

func (a *Agent) extractSQL(response string) string {
	// Remove any backticks that might wrap the SQL
	response = strings.ReplaceAll(response, "`", "")
	
	// Try to extract SQL from code blocks
	if strings.Contains(response, "```sql") {
		start := strings.Index(response, "```sql") + 6
		end := strings.Index(response[start:], "```")
		if end > 0 {
			sql := strings.TrimSpace(response[start : start+end])
			return strings.ReplaceAll(sql, "`", "")
		}
	}

	if strings.Contains(response, "```") {
		start := strings.Index(response, "```") + 3
		end := strings.Index(response[start:], "```")
		if end > 0 {
			sql := strings.TrimSpace(response[start : start+end])
			return strings.ReplaceAll(sql, "`", "")
		}
	}

	// If no code blocks, try to find SELECT statement
	lines := strings.Split(response, "\n")
	var sqlLines []string
	inSQL := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		upper := strings.ToUpper(trimmed)

		if strings.HasPrefix(upper, "SELECT") ||
			strings.HasPrefix(upper, "WITH") {
			inSQL = true
		}

		if inSQL {
			sqlLines = append(sqlLines, line)
			if strings.HasSuffix(trimmed, ";") {
				break
			}
		}
	}

	if len(sqlLines) > 0 {
		sql := strings.TrimSpace(strings.Join(sqlLines, "\n"))
		return strings.ReplaceAll(sql, "`", "")
	}

	sql := strings.TrimSpace(response)
	sql = strings.ReplaceAll(sql, "`", "")
	
	// Remove any trailing semicolons and re-add one
	sql = strings.TrimRight(sql, ";")
	sql = strings.TrimSpace(sql)
	
	// Add single semicolon at the end if not empty
	if sql != "" && !strings.HasSuffix(sql, ";") {
		sql += ";"
	}
	
	return sql
}

func (a *Agent) fixQuery(sql, errorMsg, originalQuestion string) (string, error) {
	fixPrompt := fmt.Sprintf(`The following SQL query failed with an error:

SQL:
%s

Error:
%s

Original question: "%s"

Database schema:
%s

Fix the SQL query to resolve this error. Return ONLY the corrected SQL query in triple backticks.`,
		sql, errorMsg, originalQuestion, a.schemaCache.Summary)

	response, err := a.llm.Generate(fixPrompt, a.getSystemPrompt())
	if err != nil {
		return "", err
	}

	return a.extractSQL(response), nil
}

func (a *Agent) GetSchema() (*database.SchemaInfo, error) {
	if a.schemaCache == nil {
		schema, err := a.db.GetFullSchema()
		if err != nil {
			return nil, err
		}
		a.schemaCache = schema
	}
	return a.schemaCache, nil
}

func (a *Agent) RefreshSchema() error {
	schema, err := a.db.GetFullSchema()
	if err != nil {
		return err
	}
	a.schemaCache = schema
	return nil
}

func (a *Agent) ClearHistory() {
	a.conversationHistory = make([]llm.ChatMessage, 0)
}
