package services

import (
	"database/sql"
	"time"

	"hrms/database"
	"hrms/models"

	"github.com/google/uuid"
)

type DashboardService interface {
	GetAdminDashboard() (map[string]interface{}, error)
	GetEmployeeDashboard(userID uuid.UUID) (map[string]interface{}, error)
}

type dashboardService struct {
	db *sql.DB
}

func NewDashboardService() DashboardService {
	return &dashboardService{
		db: database.DB,
	}
}

func (s *dashboardService) GetAdminDashboard() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Total employees
	var totalEmployees int
	s.db.QueryRow("SELECT COUNT(*) FROM users WHERE is_active = true").Scan(&totalEmployees)
	result["total_employees"] = totalEmployees

	// Present employees today
	today := time.Now().Format("2006-01-02")
	var presentToday int
	s.db.QueryRow(`
		SELECT COUNT(DISTINCT user_id) 
		FROM attendance 
		WHERE date = $1 AND status IN ('present', 'late', 'half_day')
	`, today).Scan(&presentToday)
	result["present_today"] = presentToday

	// Absent employees today
	result["absent_today"] = totalEmployees - presentToday

	// Pending leave requests
	var pendingLeaves int
	s.db.QueryRow("SELECT COUNT(*) FROM leaves WHERE status = 'pending'").Scan(&pendingLeaves)
	result["pending_leaves"] = pendingLeaves

	// Recent activities (last 10 audit logs)
	rows, _ := s.db.Query(`
		SELECT action, entity_type, created_at 
		FROM audit_logs 
		ORDER BY created_at DESC 
		LIMIT 10
	`)
	defer rows.Close()
	
	var activities []map[string]interface{}
	for rows.Next() {
		var action, entityType string
		var createdAt time.Time
		rows.Scan(&action, &entityType, &createdAt)
		activities = append(activities, map[string]interface{}{
			"action":      action,
			"entity_type": entityType,
			"created_at":  createdAt,
		})
	}
	result["recent_activities"] = activities

	// Payroll summary (current month)
	currentMonth := int(time.Now().Month())
	currentYear := time.Now().Year()
	var payrollTotal float64
	s.db.QueryRow(`
		SELECT COALESCE(SUM(net_salary), 0) 
		FROM salaries 
		WHERE month = $1 AND year = $2
	`, currentMonth, currentYear).Scan(&payrollTotal)
	result["payroll_summary"] = map[string]interface{}{
		"month": currentMonth,
		"year":  currentYear,
		"total": payrollTotal,
	}

	return result, nil
}

func (s *dashboardService) GetEmployeeDashboard(userID uuid.UUID) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Attendance summary (current month)
	currentMonth := int(time.Now().Month())
	currentYear := time.Now().Year()
	var presentDays, absentDays int
	s.db.QueryRow(`
		SELECT 
			COUNT(*) FILTER (WHERE status IN ('present', 'late', 'half_day')) as present,
			COUNT(*) FILTER (WHERE status = 'absent') as absent
		FROM attendance
		WHERE user_id = $1 AND EXTRACT(MONTH FROM date) = $2 AND EXTRACT(YEAR FROM date) = $3
	`, userID, currentMonth, currentYear).Scan(&presentDays, &absentDays)
	result["attendance_summary"] = map[string]interface{}{
		"present_days": presentDays,
		"absent_days":  absentDays,
		"month":        currentMonth,
		"year":         currentYear,
	}

	// Leave balance (pending + approved leaves this year)
	var leaveBalance int
	s.db.QueryRow(`
		SELECT COUNT(*) 
		FROM leaves 
		WHERE user_id = $1 AND status IN ('pending', 'approved') AND EXTRACT(YEAR FROM start_date) = $2
	`, userID, currentYear).Scan(&leaveBalance)
	result["leave_balance"] = 20 - leaveBalance // Assuming 20 days annual leave

	// Salary summary (latest)
	var latestSalary models.Salary
	err := s.db.QueryRow(`
		SELECT id, user_id, base_salary, bonus, deductions, net_salary, month, year, created_at, updated_at
		FROM salaries
		WHERE user_id = $1
		ORDER BY year DESC, month DESC
		LIMIT 1
	`, userID).Scan(
		&latestSalary.ID,
		&latestSalary.UserID,
		&latestSalary.BaseSalary,
		&latestSalary.Bonus,
		&latestSalary.Deductions,
		&latestSalary.NetSalary,
		&latestSalary.Month,
		&latestSalary.Year,
		&latestSalary.CreatedAt,
		&latestSalary.UpdatedAt,
	)
	if err == nil {
		result["salary_summary"] = latestSalary
	}

	// Upcoming holidays (placeholder - would come from holidays table)
	result["upcoming_holidays"] = []interface{}{}

	// Unread notifications count
	var unreadCount int
	s.db.QueryRow(`
		SELECT COUNT(*) 
		FROM notifications 
		WHERE user_id = $1 AND is_read = false
	`, userID).Scan(&unreadCount)
	result["unread_notifications"] = unreadCount

	return result, nil
}
