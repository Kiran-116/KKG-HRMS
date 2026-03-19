package services

import (
	"context"
	"database/sql"
	"time"

	"hrms/database"
	"hrms/models"

	"github.com/google/uuid"
)

type DashboardService interface {
	GetAdminDashboard(ctx context.Context, rangeParam string) (map[string]interface{}, error)
	GetEmployeeDashboard(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
}

type dashboardService struct {
	db *sql.DB
}

func NewDashboardService() DashboardService {
	return &dashboardService{
		db: database.DB,
	}
}

func (s *dashboardService) GetAdminDashboard(ctx context.Context, rangeParam string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Total employees
	var totalEmployees int
	s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE is_active = true").Scan(&totalEmployees)
	result["total_employees"] = totalEmployees

	// Present employees today
	today := time.Now().Format("2006-01-02")
	var presentToday int
	s.db.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT user_id) 
		FROM attendance 
		WHERE date = $1 AND status IN ('present', 'late', 'half_day')
	`, today).Scan(&presentToday)
	result["present_today"] = presentToday

	// Absent employees today
	result["absent_today"] = totalEmployees - presentToday

	// Pending leave requests
	var pendingLeaves int
	s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM leaves WHERE status = 'pending'").Scan(&pendingLeaves)
	result["pending_leaves"] = pendingLeaves

	// Recent activities (last 10 audit logs)
	rows, _ := s.db.QueryContext(ctx, `
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
	s.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(net_salary), 0) 
		FROM salaries 
		WHERE month = $1 AND year = $2
	`, currentMonth, currentYear).Scan(&payrollTotal)
	result["payroll_summary"] = map[string]interface{}{
		"month": currentMonth,
		"year":  currentYear,
		"total": payrollTotal,
	}

	// Determine range for attendance trend
	days := 30
	switch rangeParam {
	case "day":
		days = 7
	case "year":
		days = 365
	}

	// Attendance trend over selected range
	// Use a date series to include days with no attendance rows.
	attendanceTrendRows, _ := s.db.QueryContext(ctx, `
		WITH dates AS (
			SELECT generate_series(CURRENT_DATE - ($1::int - 1) * INTERVAL '1 day', CURRENT_DATE, INTERVAL '1 day')::date AS date
		)
		SELECT 
			d.date,
			COALESCE(COUNT(DISTINCT a.user_id) FILTER (WHERE a.status IN ('present', 'late', 'half_day')), 0) AS present
		FROM dates d
		LEFT JOIN attendance a ON a.date = d.date
		GROUP BY d.date
		ORDER BY d.date ASC
	`, days)
	defer attendanceTrendRows.Close()

	var attendanceTrend []map[string]interface{}
	for attendanceTrendRows.Next() {
		var date time.Time
		var present int
		attendanceTrendRows.Scan(&date, &present)
		absent := 0
		if totalEmployees > present {
			absent = totalEmployees - present
		}
		attendanceTrend = append(attendanceTrend, map[string]interface{}{
			"date":    date.Format("2006-01-02"),
			"present": present,
			"absent":  absent,
		})
	}
	result["attendance_trend"] = attendanceTrend

	// Monthly payroll (last 6 months)
	startMonth := currentMonth - 5
	startYear := currentYear
	if startMonth <= 0 {
		startMonth += 12
		startYear--
	}
	payrollTrendRows, _ := s.db.QueryContext(ctx, `
		SELECT month, year, COALESCE(SUM(net_salary), 0) as total
		FROM salaries
		WHERE (year = $1 AND month >= $2) OR (year = $3 AND month <= $4)
		GROUP BY month, year
		ORDER BY year ASC, month ASC
		LIMIT 6
	`, currentYear, startMonth, startYear, currentMonth)
	defer payrollTrendRows.Close()

	var payrollTrend []map[string]interface{}
	for payrollTrendRows.Next() {
		var month, year int
		var total float64
		payrollTrendRows.Scan(&month, &year, &total)
		payrollTrend = append(payrollTrend, map[string]interface{}{
			"month": month,
			"year":  year,
			"total": total,
		})
	}
	result["payroll_trend"] = payrollTrend

	return result, nil
}

func (s *dashboardService) GetEmployeeDashboard(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Attendance summary (current month)
	currentMonth := int(time.Now().Month())
	currentYear := time.Now().Year()
	var presentDays, absentDays int
	s.db.QueryRowContext(ctx, `
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
	s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM leaves 
		WHERE user_id = $1 AND status IN ('pending', 'approved') AND EXTRACT(YEAR FROM start_date) = $2
	`, userID, currentYear).Scan(&leaveBalance)
	result["leave_balance"] = 20 - leaveBalance // Assuming 20 days annual leave

	// Salary summary (latest)
	var latestSalary models.Salary
	err := s.db.QueryRowContext(ctx, `
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
	s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM notifications 
		WHERE user_id = $1 AND is_read = false
	`, userID).Scan(&unreadCount)
	result["unread_notifications"] = unreadCount

	// Attendance trend (last 7 days)
	attendanceTrendRows, _ := s.db.QueryContext(ctx, `
		SELECT date, 
		       CASE WHEN status IN ('present', 'late', 'half_day') THEN 1 ELSE 0 END as present
		FROM attendance
		WHERE user_id = $1 AND date >= CURRENT_DATE - INTERVAL '7 days'
		ORDER BY date ASC
	`, userID)
	defer attendanceTrendRows.Close()

	var attendanceTrend []map[string]interface{}
	for attendanceTrendRows.Next() {
		var date time.Time
		var present int
		attendanceTrendRows.Scan(&date, &present)
		attendanceTrend = append(attendanceTrend, map[string]interface{}{
			"date":    date.Format("2006-01-02"),
			"present": present,
		})
	}
	result["attendance_trend"] = attendanceTrend

	return result, nil
}
