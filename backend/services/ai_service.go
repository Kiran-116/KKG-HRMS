package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"hrms/config"
	"hrms/repositories"

	"github.com/google/uuid"
)

type AIService interface {
	ProcessHRQuery(ctx context.Context, userID uuid.UUID, query string) (string, error)
}

type aiService struct {
	userRepo       repositories.UserRepository
	leaveRepo      repositories.LeaveRepository
	attendanceRepo repositories.AttendanceRepository
	salaryRepo     repositories.SalaryRepository
}

func NewAIService(
	userRepo repositories.UserRepository,
	leaveRepo repositories.LeaveRepository,
	attendanceRepo repositories.AttendanceRepository,
	salaryRepo repositories.SalaryRepository,
) AIService {
	return &aiService{
		userRepo:       userRepo,
		leaveRepo:      leaveRepo,
		attendanceRepo: attendanceRepo,
		salaryRepo:     salaryRepo,
	}
}

func (s *aiService) ProcessHRQuery(ctx context.Context, userID uuid.UUID, query string) (string, error) {
	query = strings.ToLower(query)

	// Simple query parsing and response
	if strings.Contains(query, "leave") || strings.Contains(query, "leaves") {
		leaves, _ := s.leaveRepo.GetByUserID(ctx, userID, 10, 0)
		approvedCount := 0
		pendingCount := 0
		for _, leave := range leaves {
			if leave.Status == "approved" {
				approvedCount++
			} else if leave.Status == "pending" {
				pendingCount++
			}
		}
		return fmt.Sprintf("You have %d approved leaves and %d pending leave requests.", approvedCount, pendingCount), nil
	}

	if strings.Contains(query, "attendance") {
		attendances, _ := s.attendanceRepo.GetByUserID(ctx, userID, 10, 0)
		presentCount := 0
		for _, att := range attendances {
			if att.Status == "present" || att.Status == "late" {
				presentCount++
			}
		}
		return fmt.Sprintf("You have %d attendance records with %d present days.", len(attendances), presentCount), nil
	}

	if strings.Contains(query, "salary") {
		salaries, _ := s.salaryRepo.GetByUserID(ctx, userID, 1, 0)
		if len(salaries) > 0 {
			latest := salaries[0]
			return fmt.Sprintf("Your latest salary is $%.2f for %d/%d.", latest.NetSalary, latest.Month, latest.Year), nil
		}
		return "No salary records found.", nil
	}

	// Use OpenAI API for complex queries
	if config.AppConfig.OpenAI.APIKey != "" {
		return s.callOpenAI(query)
	}

	return "I can help you with leaves, attendance, and salary information. Please ask a specific question.", nil
}

func (s *aiService) callOpenAI(query string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	payload := map[string]interface{}{
		"model": config.AppConfig.OpenAI.Model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are an HR assistant. Answer HR-related questions concisely.",
			},
			{
				"role":    "user",
				"content": query,
			},
		},
		"max_tokens": 150,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.AppConfig.OpenAI.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("failed to call AI service")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return content, nil
				}
			}
		}
	}

	return "I couldn't process that query. Please try again.", nil
}
