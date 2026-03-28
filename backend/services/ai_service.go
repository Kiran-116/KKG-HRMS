package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"hrms/config"
	"hrms/database"
	"hrms/repositories"

	"github.com/google/uuid"
	"github.com/openai/openai-go"
)

type AIService interface {
	ProcessHRQuery(ctx context.Context, userID uuid.UUID, query string) (*AIResponse, error)
}

type aiService struct {
	db             *sql.DB
	userRepo       repositories.UserRepository
	leaveRepo      repositories.LeaveRepository
	attendanceRepo repositories.AttendanceRepository
	salaryRepo     repositories.SalaryRepository
}

type SQLQueryResponse struct {
	Query string `json:"query"`
}

type AIResponse struct {
	Type    string                   `json:"type"`
	Message string                   `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}

func NewAIService(
	userRepo repositories.UserRepository,
	leaveRepo repositories.LeaveRepository,
	attendanceRepo repositories.AttendanceRepository,
	salaryRepo repositories.SalaryRepository,
) AIService {
	return &aiService{
		db:             database.DB,
		userRepo:       userRepo,
		leaveRepo:      leaveRepo,
		attendanceRepo: attendanceRepo,
		salaryRepo:     salaryRepo,
	}
}

func (s *aiService) generateSQL(query string, role string, userID uuid.UUID) (string, error) {
	systemPrompt := fmt.Sprintf(`
You are an AI that converts user questions into SQL queries.

Rules:
- Only generate SELECT queries
- Do NOT modify data
- Use table names: users, leaves, attendance, salaries
- If role = EMPLOYEE → add WHERE user_id = '%s'
- If role = ADMIN → no restriction

Return JSON:
{
  "query": "SELECT ..."
}
`, userID.String())

	// Configure API key via environment for the SDK in case options are unavailable
	_ = os.Setenv("OPENAI_API_KEY", config.AppConfig.OpenAI.APIKey)
	client := openai.NewClient()

	var content string
	// Light retry on potential rate limiting
	for attempt := 0; attempt < 3; attempt++ {
		resp, err := client.Chat.Completions.New(
			context.Background(),
			openai.ChatCompletionNewParams{
				Model: openai.ChatModel(config.AppConfig.OpenAI.Model),
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage(systemPrompt),
					openai.UserMessage(query),
				},
			},
		)
		if err != nil {
			// backoff on rate limit or overload
			if strings.Contains(strings.ToLower(err.Error()), "rate") && attempt < 2 {
				time.Sleep(time.Duration(250*(1<<attempt)) * time.Millisecond)
				continue
			}
			return "", err
		}
		if len(resp.Choices) == 0 {
			return "", errors.New("no response from AI")
		}
		content = resp.Choices[0].Message.Content
		break
	}

	var sqlResp SQLQueryResponse
	if err := json.Unmarshal([]byte(content), &sqlResp); err != nil {
		return "", fmt.Errorf("failed to parse SQL: %w; content: %s", err, content)
	}
	return sqlResp.Query, nil
}

func validateSQL(query string) error {
	q := strings.TrimSpace(strings.ToLower(query))
	if strings.Contains(q, ";") {
		return errors.New("multiple statements are not allowed")
	}
	if !strings.HasPrefix(q, "select") {
		return errors.New("only SELECT queries allowed")
	}
	if strings.Contains(q, " drop ") ||
		strings.Contains(q, " delete ") ||
		strings.Contains(q, " update ") ||
		strings.Contains(q, " insert ") ||
		strings.Contains(q, " alter ") ||
		strings.Contains(q, " truncate ") {
		return errors.New("unsafe query detected")
	}
	allowed := map[string]bool{"users": true, "leaves": true, "attendance": true, "salaries": true}
	lq := " " + q + " "
	extractAndCheck := func(keyword string) error {
		idx := 0
		for {
			pos := strings.Index(lq[idx:], " "+keyword+" ")
			if pos == -1 {
				return nil
			}
			pos = pos + idx + len(" "+keyword+" ")
			rest := lq[pos:]
			end := len(rest)
			for i, ch := range rest {
				if ch == ' ' || ch == '\n' || ch == '\r' || ch == '\t' || ch == ',' {
					end = i
					break
				}
			}
			token := strings.Trim(rest[:end], " \n\r\t")
			token = strings.Trim(token, `"`)
			parts := strings.Split(token, ".")
			token = parts[len(parts)-1]
			token = strings.Trim(token, ",")
			if token != "" && !allowed[token] {
				return fmt.Errorf("table %s is not allowed", token)
			}
			idx = pos
		}
	}
	if err := extractAndCheck("from"); err != nil {
		return err
	}
	if err := extractAndCheck("join"); err != nil {
		return err
	}
	return nil
}

func containsUserPredicate(sqlQuery string, userID uuid.UUID) bool {
	q := strings.ToLower(sqlQuery)
	uid := strings.ToLower(userID.String())
	return strings.Contains(q, "user_id = '"+uid+"'") ||
		strings.Contains(q, "user_id='"+uid+"'") ||
		strings.Contains(q, "employee_id = '"+uid+"'") ||
		strings.Contains(q, "employee_id='"+uid+"'")
}

func (s *aiService) executeSQL(ctx context.Context, query string) ([]map[string]interface{}, error) {
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	results := []map[string]interface{}{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		ptrs := make([]interface{}, len(columns))
		for i := range values {
			ptrs[i] = &values[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		row := map[string]interface{}{}
		for i, col := range columns {
			row[col] = normalizeDBValue(values[i])
		}
		results = append(results, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// normalizeDBValue converts driver-native values into JSON-friendly types to avoid base64 []byte in output.
func normalizeDBValue(v interface{}) interface{} {
	switch x := v.(type) {
	case nil:
		return nil
	case []byte:
		// NUMERIC/DECIMAL often arrive as []byte containing ASCII; stringify
		return string(x)
	case time.Time:
		return x.UTC().Format(time.RFC3339)
	case uuid.UUID:
		return x.String()
	// Common numeric/int types already JSON-safe
	case int64, int32, int16, int8, int:
		return x
	case float64, float32:
		return x
	case bool, string:
		return x
	default:
		// Fallback to fmt for any other driver-specific types
		return fmt.Sprintf("%v", x)
	}
}

func (s *aiService) formatResponse(query string, data interface{}) (string, error) {
	// Deterministic salary formatter to avoid AI dependency and produce precise values
	lq := strings.ToLower(query)
	if strings.Contains(lq, "salary") {
		if rows, ok := data.([]map[string]interface{}); ok && len(rows) > 0 {
			row := rows[0]
			monthAny := row["month"]
			yearAny := row["year"]
			netAny := row["net_salary"]
			monthStr := fmt.Sprintf("%v", monthAny)
			yearStr := fmt.Sprintf("%v", yearAny)
			netStr := fmt.Sprintf("%v", netAny)
			// Try to format net salary as currency with two decimals if numeric-like
			if f, err := parseFloatLoose(netStr); err == nil {
				netStr = fmt.Sprintf("$%.2f", f)
			}
			return fmt.Sprintf("Your salary details: Net Salary %s for %s/%s.", netStr, monthStr, yearStr), nil
		}
	}

	if config.AppConfig.OpenAI.APIKey == "" {
		switch rows := data.(type) {
		case []map[string]interface{}:
			return fmt.Sprintf("Found %d result(s) for your query.", len(rows)), nil
		default:
			return "Processed your request.", nil
		}
	}
	systemPrompt := `
You are an HR assistant.

Convert database results into human-friendly response.
Be short and clear.
`
	_ = os.Setenv("OPENAI_API_KEY", config.AppConfig.OpenAI.APIKey)
	client := openai.NewClient()

	var content string
	for attempt := 0; attempt < 3; attempt++ {
		resp, err := client.Chat.Completions.New(
			context.Background(),
			openai.ChatCompletionNewParams{
				Model: openai.ChatModel(config.AppConfig.OpenAI.Model),
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage(systemPrompt),
					openai.UserMessage(fmt.Sprintf("Query: %s\nData: %+v", query, data)),
				},
			},
		)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "rate") && attempt < 2 {
				time.Sleep(time.Duration(250*(1<<attempt)) * time.Millisecond)
				continue
			}
			// Graceful local fallback summary
			if rows, ok := data.([]map[string]interface{}); ok {
				return fmt.Sprintf("I found %d result(s).", len(rows)), nil
			}
			return "I processed your request.", nil
		}
		if len(resp.Choices) == 0 {
			return "No response", nil
		}
		content = resp.Choices[0].Message.Content
		break
	}
	return content, nil
}

// parseFloatLoose tries to parse various numeric string formats to float64
func parseFloatLoose(s string) (float64, error) {
	// strip commas
	clean := strings.ReplaceAll(s, ",", "")
	return strconv.ParseFloat(clean, 64)
}

func (s *aiService) ProcessHRQuery(ctx context.Context, userID uuid.UUID, query string) (*AIResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	role := strings.ToLower(user.Role)

	sqlQuery := ""
	if config.AppConfig.OpenAI.APIKey != "" {
		sqlQuery, err = s.generateSQL(query, role, userID)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("AI not configured")
	}
	if err := validateSQL(sqlQuery); err != nil {
		return nil, err
	}
	if role == "employee" {
		if !containsUserPredicate(sqlQuery, userID) {
			return nil, errors.New("unauthorized query")
		}
	}
	data, err := s.executeSQL(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	message, err := s.formatResponse(query, data)
	if err != nil {
		return nil, err
	}
	return &AIResponse{
		Type:    "ai_sql_response",
		Message: message,
		Data:    data,
	}, nil
}
