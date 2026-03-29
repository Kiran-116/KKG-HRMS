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
	// Generate last 6 months and join with salary data
	payrollTrendRows, _ := s.db.QueryContext(ctx, `
		WITH months AS (
			SELECT 
				EXTRACT(MONTH FROM generate_series(
					CURRENT_DATE - INTERVAL '5 months',
					CURRENT_DATE,
					'1 month'
				))::int AS month,
				EXTRACT(YEAR FROM generate_series(
					CURRENT_DATE - INTERVAL '5 months',
					CURRENT_DATE,
					'1 month'
				))::int AS year
		)
		SELECT 
			m.month,
			m.year,
			COALESCE(SUM(s.net_salary), 0) as total
		FROM months m
		LEFT JOIN salaries s ON s.month = m.month AND s.year = m.year
		GROUP BY m.month, m.year
		ORDER BY m.year ASC, m.month ASC
	`)
	defer payrollTrendRows.Close()

	var payrollTrend []map[string]interface{}
	for payrollTrendRows.Next() {
		var month, year int
		var total float64
		if err := payrollTrendRows.Scan(&month, &year, &total); err == nil {
			payrollTrend = append(payrollTrend, map[string]interface{}{
				"month": month,
				"year":  year,
				"total": total,
			})
		}
	}
	result["payroll_trend"] = payrollTrend

	return result, nil
}

func (s *dashboardService) GetEmployeeDashboard(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Attendance summary (current month)
	currentMonth := int(time.Now().Month())
	currentYear := time.Now().Year()
	var presentDays, leaveDays, businessDays int
	// Business days only, count up to TODAY; exclude approved leaves from absences
	s.db.QueryRowContext(ctx, `
		WITH month_days AS (
			SELECT generate_series(
				date_trunc('month', CURRENT_DATE)::date,
				CURRENT_DATE,
				interval '1 day'
			)::date AS d
		),
		workdays AS (
			SELECT d FROM month_days WHERE EXTRACT(ISODOW FROM d) < 6
		)
		SELECT
			COUNT(*) AS business_days,
			COALESCE(SUM(CASE WHEN a.status IN ('present','late','half_day') THEN 1 ELSE 0 END), 0) AS present_days,
			COALESCE(SUM(CASE WHEN l.user_id IS NOT NULL THEN 1 ELSE 0 END), 0) AS leave_days
		FROM workdays w
		LEFT JOIN attendance a
		  ON a.user_id = $1 AND a.date = w.d
		LEFT JOIN leaves l
		  ON l.user_id = $1 AND l.status = 'approved' AND w.d BETWEEN l.start_date AND l.end_date
	`, userID).Scan(&businessDays, &presentDays, &leaveDays)
	absentDays := businessDays - presentDays - leaveDays
	if absentDays < 0 {
		absentDays = 0
	}
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
		WITH dates AS (
			SELECT generate_series(CURRENT_DATE - interval '6 days', CURRENT_DATE, interval '1 day')::date AS d
		),
		flags AS (
			SELECT 
				d.d,
				(EXTRACT(ISODOW FROM d.d) < 6)::int AS is_workday
			FROM dates d
		)
		SELECT 
			f.d AS date,
			COALESCE(MAX(CASE WHEN a.status IN ('present','late','half_day') THEN 1 ELSE 0 END), 0) AS present,
			CASE 
				WHEN f.is_workday = 0 THEN 0
				WHEN EXISTS (
					SELECT 1 FROM leaves l 
					WHERE l.user_id = $1 AND l.status = 'approved' AND f.d BETWEEN l.start_date AND l.end_date
				) THEN 0
				ELSE (1 - COALESCE(MAX(CASE WHEN a.status IN ('present','late','half_day') THEN 1 ELSE 0 END), 0))
			END AS absent
		FROM flags f
		LEFT JOIN attendance a
		  ON a.user_id = $1 AND a.date = f.d
		GROUP BY f.d, f.is_workday
		ORDER BY f.d ASC
	`, userID)
	defer attendanceTrendRows.Close()

	var attendanceTrend []map[string]interface{}
	for attendanceTrendRows.Next() {
		var date time.Time
		var present, absent int
		attendanceTrendRows.Scan(&date, &present, &absent)
		attendanceTrend = append(attendanceTrend, map[string]interface{}{
			"date":    date.Format("2006-01-02"),
			"present": present,
			"absent":  absent,
		})
	}
	result["attendance_trend"] = attendanceTrend

	return result, nil
}
