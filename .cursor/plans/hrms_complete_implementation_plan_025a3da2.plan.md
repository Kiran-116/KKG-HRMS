---
name: HRMS Complete Implementation Plan
overview: Build a production-grade HR Management System following clean architecture principles, implementing all 12 core modules from the PRD with comprehensive testing, Docker setup, and full frontend/backend integration.
todos:
  - id: phase1-foundation
    content: "Phase 1: Set up project foundation - Go module, Docker setup, database migrations, configuration management, middleware (CORS, logging, recovery)"
    status: completed
  - id: phase2-auth
    content: "Phase 2: Implement authentication & RBAC - User model, JWT auth, password hashing, auth middleware, login/register endpoints, frontend auth pages"
    status: completed
    dependencies:
      - phase1-foundation
  - id: phase3-employee
    content: "Phase 3: Employee Management - Employee CRUD operations, admin employee management, employee profile views, frontend employee pages"
    status: completed
    dependencies:
      - phase2-auth
  - id: phase4-attendance
    content: "Phase 4: Attendance Management - Check-in/check-out, attendance history, admin attendance view, attendance tracking logic"
    status: completed
    dependencies:
      - phase3-employee
  - id: phase5-leave
    content: "Phase 5: Leave Management - Leave application, approval/rejection workflow, leave history, leave balance calculation"
    status: completed
    dependencies:
      - phase3-employee
  - id: phase6-payroll
    content: "Phase 6: Payroll Management - Salary records, payslip generation (PDF), salary history, admin payroll management"
    status: completed
    dependencies:
      - phase3-employee
  - id: phase7-documents
    content: "Phase 7: Document Management - File upload/download, local storage implementation, cloud-ready storage interface, document listing"
    status: completed
    dependencies:
      - phase3-employee
  - id: phase8-notifications
    content: "Phase 8: Notification System - In-app notifications, SMTP email service, notification creation on events, notification center UI"
    status: completed
    dependencies:
      - phase5-leave
      - phase6-payroll
      - phase7-documents
  - id: phase9-dashboard
    content: "Phase 9: Dashboard - Admin dashboard with widgets (employees, attendance, leaves, payroll), employee dashboard with personal stats"
    status: completed
    dependencies:
      - phase4-attendance
      - phase5-leave
      - phase6-payroll
  - id: phase10-audit
    content: "Phase 10: Audit Logging - Audit log model, automatic logging middleware, audit log viewer for admins, event tracking"
    status: completed
    dependencies:
      - phase2-auth
  - id: phase11-ai
    content: "Phase 11: AI HR Assistant - OpenAI integration, natural language query parsing, AI chat interface, HR query handling"
    status: completed
    dependencies:
      - phase4-attendance
      - phase5-leave
      - phase6-payroll
  - id: phase12-observability
    content: "Phase 12: Observability - New Relic integration, HTTP request tracing, database query monitoring, error tracking"
    status: completed
    dependencies:
      - phase1-foundation
  - id: phase13-frontend
    content: "Phase 13: Frontend Infrastructure - React routing, layout components, sidebar navigation, header, TailwindCSS setup, responsive design"
    status: completed
    dependencies:
      - phase2-auth
  - id: phase14-testing
    content: "Phase 14: Testing - Unit tests for services/repositories, integration tests for API endpoints, frontend component tests, test coverage >80%"
    status: completed
    dependencies:
      - phase11-ai
  - id: phase15-docs
    content: "Phase 15: Documentation - API documentation, architecture docs, README with setup instructions, environment variables guide"
    status: completed
    dependencies:
      - phase14-testing
---

# HR Management System - Complete Implementation Plan

## Architecture Overview

The system follows **clean architecture** with clear separation of concerns:

```
Frontend (React + TypeScript + TailwindCSS)
    ↓ REST API
Backend (Golang + Gin)
    ├── Controllers (HTTP handlers)
    ├── Services (Business logic)
    ├── Repositories (Data access)
    ├── Models (Domain entities)
    ├── Middleware (Auth, logging, observability)
    └── Utils (Helpers, validators)
    ↓
PostgreSQL Database
    ↓
External Services (New Relic, OpenAI, SMTP)
```

### Key Architecture Decisions

1. **Clean Architecture Layers**:

   - **Controllers**: Handle HTTP requests/responses, input validation
   - **Services**: Business logic, orchestration, transaction management
   - **Repositories**: Database operations, abstraction layer
   - **Models**: Domain entities with validation

2. **Security**:

   - JWT tokens with refresh token support
   - bcrypt password hashing (cost factor 10)
   - Role-based middleware for route protection
   - Input validation on all endpoints
   - Rate limiting on sensitive endpoints

3. **Observability**:

   - New Relic instrumentation for all HTTP handlers
   - Structured logging with context
   - Database query monitoring
   - Error tracking

4. **Document Storage**:

   - Local filesystem for development
   - Interface-based design for cloud storage (S3-compatible)
   - File validation (type, size limits)

5. **Testing Strategy**:

   - Unit tests for services and repositories
   - Integration tests for API endpoints
   - Test database for integration tests

## Implementation Phases

### Phase 1: Project Foundation & Infrastructure

**Files to create:**

- `backend/go.mod` - Go module configuration
- `backend/main.go` - Application entry point
- `backend/config/config.go` - Configuration management
- `backend/database/database.go` - Database connection and migration setup
- `backend/middleware/` - CORS, logging, recovery middleware
- `docker-compose.yml` - PostgreSQL, backend, frontend services
- `backend/Dockerfile` - Backend container
- `frontend/package.json` - Frontend dependencies
- `frontend/Dockerfile` - Frontend container
- `.env.example` - Environment variables template

**Database Migrations:**

- `backend/migrations/001_create_users.sql`
- `backend/migrations/002_create_attendance.sql`
- `backend/migrations/003_create_leaves.sql`
- `backend/migrations/004_create_salaries.sql`
- `backend/migrations/005_create_documents.sql`
- `backend/migrations/006_create_notifications.sql`
- `backend/migrations/007_create_audit_logs.sql`

### Phase 2: Authentication & RBAC

**Backend Files:**

- `backend/models/user.go` - User model
- `backend/repositories/user_repository.go` - User data access
- `backend/services/auth_service.go` - Authentication logic
- `backend/controllers/auth_controller.go` - Auth endpoints
- `backend/middleware/auth_middleware.go` - JWT validation
- `backend/middleware/rbac_middleware.go` - Role-based access control
- `backend/utils/jwt.go` - JWT token generation/validation
- `backend/utils/password.go` - Password hashing utilities

**Frontend Files:**

- `frontend/src/services/api.ts` - API client setup
- `frontend/src/services/authService.ts` - Auth API calls
- `frontend/src/contexts/AuthContext.tsx` - Auth state management
- `frontend/src/pages/Login.tsx` - Login page
- `frontend/src/pages/Register.tsx` - Registration page
- `frontend/src/components/ProtectedRoute.tsx` - Route protection

**APIs:**

- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /api/auth/me` - Get current user

### Phase 3: Employee Management

**Backend Files:**

- `backend/models/employee.go` - Employee model (extends user)
- `backend/repositories/employee_repository.go` - Employee data access
- `backend/services/employee_service.go` - Employee business logic
- `backend/controllers/employee_controller.go` - Employee endpoints
- `backend/utils/validators.go` - Input validation helpers

**Frontend Files:**

- `frontend/src/services/employeeService.ts` - Employee API calls
- `frontend/src/pages/employees/EmployeeList.tsx` - Admin employee list
- `frontend/src/pages/employees/EmployeeForm.tsx` - Create/Edit employee
- `frontend/src/pages/employees/EmployeeProfile.tsx` - Employee profile view

**APIs:**

- `GET /api/employees` - List employees (admin)
- `POST /api/employees` - Create employee (admin)
- `PUT /api/employees/:id` - Update employee (admin)
- `GET /api/employees/:id` - Get employee details
- `GET /api/employees/me` - Get own profile

### Phase 4: Attendance Management

**Backend Files:**

- `backend/models/attendance.go` - Attendance model
- `backend/repositories/attendance_repository.go` - Attendance data access
- `backend/services/attendance_service.go` - Attendance business logic
- `backend/controllers/attendance_controller.go` - Attendance endpoints

**Frontend Files:**

- `frontend/src/services/attendanceService.ts` - Attendance API calls
- `frontend/src/pages/attendance/AttendanceCheckIn.tsx` - Check-in/out interface
- `frontend/src/pages/attendance/AttendanceHistory.tsx` - Attendance history
- `frontend/src/pages/attendance/AdminAttendanceView.tsx` - Admin attendance dashboard

**APIs:**

- `POST /api/attendance/checkin` - Employee check-in
- `POST /api/attendance/checkout` - Employee check-out
- `GET /api/attendance/me` - Own attendance history
- `GET /api/attendance` - All attendance (admin, with filters)

### Phase 5: Leave Management

**Backend Files:**

- `backend/models/leave.go` - Leave model
- `backend/repositories/leave_repository.go` - Leave data access
- `backend/services/leave_service.go` - Leave business logic
- `backend/controllers/leave_controller.go` - Leave endpoints

**Frontend Files:**

- `frontend/src/services/leaveService.ts` - Leave API calls
- `frontend/src/pages/leaves/ApplyLeave.tsx` - Leave application form
- `frontend/src/pages/leaves/LeaveHistory.tsx` - Leave history
- `frontend/src/pages/leaves/AdminLeaveApproval.tsx` - Admin leave approval

**APIs:**

- `POST /api/leaves/apply` - Apply for leave
- `GET /api/leaves/me` - Own leave history
- `GET /api/leaves` - All leaves (admin)
- `PUT /api/leaves/:id/approve` - Approve leave (admin)
- `PUT /api/leaves/:id/reject` - Reject leave (admin)

### Phase 6: Payroll Management

**Backend Files:**

- `backend/models/salary.go` - Salary model
- `backend/repositories/salary_repository.go` - Salary data access
- `backend/services/salary_service.go` - Salary business logic
- `backend/controllers/salary_controller.go` - Salary endpoints
- `backend/services/payslip_service.go` - Payslip generation (PDF)

**Frontend Files:**

- `frontend/src/services/salaryService.ts` - Salary API calls
- `frontend/src/pages/payroll/SalaryHistory.tsx` - Employee salary view
- `frontend/src/pages/payroll/AdminPayroll.tsx` - Admin payroll management
- `frontend/src/pages/payroll/PayslipView.tsx` - Payslip display/download

**APIs:**

- `POST /api/salary` - Create salary record (admin)
- `GET /api/salary/me` - Own salary history
- `GET /api/salary/:userId` - Employee salary (admin)
- `GET /api/salary/:id/payslip` - Download payslip PDF

### Phase 7: Document Management

**Backend Files:**

- `backend/models/document.go` - Document model
- `backend/repositories/document_repository.go` - Document data access
- `backend/services/document_service.go` - Document business logic
- `backend/services/storage_service.go` - File storage interface (local + cloud-ready)
- `backend/services/local_storage.go` - Local filesystem implementation
- `backend/controllers/document_controller.go` - Document endpoints

**Frontend Files:**

- `frontend/src/services/documentService.ts` - Document API calls
- `frontend/src/pages/documents/DocumentUpload.tsx` - Document upload
- `frontend/src/pages/documents/DocumentList.tsx` - Document listing
- `frontend/src/components/FileUpload.tsx` - Reusable file upload component

**APIs:**

- `POST /api/documents` - Upload document
- `GET /api/documents/me` - Own documents
- `GET /api/documents/:userId` - Employee documents (admin)
- `DELETE /api/documents/:id` - Delete document

### Phase 8: Notification System

**Backend Files:**

- `backend/models/notification.go` - Notification model
- `backend/repositories/notification_repository.go` - Notification data access
- `backend/services/notification_service.go` - Notification business logic
- `backend/services/email_service.go` - SMTP email service
- `backend/controllers/notification_controller.go` - Notification endpoints

**Frontend Files:**

- `frontend/src/services/notificationService.ts` - Notification API calls
- `frontend/src/components/NotificationBell.tsx` - Notification indicator
- `frontend/src/pages/notifications/NotificationList.tsx` - Notification center

**APIs:**

- `GET /api/notifications` - Get notifications
- `PUT /api/notifications/:id/read` - Mark as read
- `GET /api/notifications/unread-count` - Unread count

### Phase 9: Dashboard

**Backend Files:**

- `backend/services/dashboard_service.go` - Dashboard data aggregation
- `backend/controllers/dashboard_controller.go` - Dashboard endpoints

**Frontend Files:**

- `frontend/src/pages/dashboard/AdminDashboard.tsx` - Admin dashboard with widgets
- `frontend/src/pages/dashboard/EmployeeDashboard.tsx` - Employee dashboard
- `frontend/src/components/dashboard/StatCard.tsx` - Reusable stat widget
- `frontend/src/components/dashboard/ActivityFeed.tsx` - Recent activities widget

**APIs:**

- `GET /api/dashboard/admin` - Admin dashboard data
- `GET /api/dashboard/employee` - Employee dashboard data

### Phase 10: Audit Logging

**Backend Files:**

- `backend/models/audit_log.go` - Audit log model
- `backend/repositories/audit_repository.go` - Audit log data access
- `backend/services/audit_service.go` - Audit logging service
- `backend/middleware/audit_middleware.go` - Automatic audit logging
- `backend/controllers/audit_controller.go` - Audit log endpoints

**Frontend Files:**

- `frontend/src/pages/audit/AuditLogs.tsx` - Audit log viewer (admin)

**APIs:**

- `GET /api/audit-logs` - Get audit logs (admin, with filters)

### Phase 11: AI HR Assistant

**Backend Files:**

- `backend/services/ai_service.go` - OpenAI integration
- `backend/services/query_parser.go` - Natural language query parsing
- `backend/controllers/ai_controller.go` - AI assistant endpoint

**Frontend Files:**

- `frontend/src/services/aiService.ts` - AI API calls
- `frontend/src/pages/ai/HRAssistant.tsx` - AI chat interface
- `frontend/src/components/ai/ChatMessage.tsx` - Chat message component

**APIs:**

- `POST /api/ai/hr-assistant` - AI query endpoint

### Phase 12: Observability (New Relic)

**Backend Files:**

- `backend/middleware/newrelic_middleware.go` - New Relic instrumentation
- `backend/utils/observability.go` - Observability helpers
- Update `backend/main.go` - Initialize New Relic

**Configuration:**

- New Relic agent initialization
- HTTP transaction tracking
- Database query monitoring
- Error tracking

### Phase 13: Frontend Infrastructure & UI

**Frontend Files:**

- `frontend/src/App.tsx` - Main app component with routing
- `frontend/src/components/Layout.tsx` - Main layout with sidebar
- `frontend/src/components/Sidebar.tsx` - Navigation sidebar
- `frontend/src/components/Header.tsx` - Top header bar
- `frontend/src/utils/constants.ts` - App constants
- `frontend/src/utils/helpers.ts` - Utility functions
- `frontend/tailwind.config.js` - Tailwind configuration
- `frontend/src/index.css` - Global styles

### Phase 14: Testing

**Backend Tests:**

- `backend/services/*_test.go` - Service unit tests
- `backend/repositories/*_test.go` - Repository tests
- `backend/controllers/*_test.go` - API integration tests
- `backend/test/test_helper.go` - Test utilities

**Frontend Tests:**

- `frontend/src/components/*.test.tsx` - Component tests
- `frontend/src/services/*.test.ts` - Service tests

### Phase 15: API Documentation

**Files:**

- `docs/API.md` - Complete API documentation
- `docs/ARCHITECTURE.md` - Architecture documentation
- `README.md` - Project setup and run instructions

## Database Schema Details

All tables include `created_at` and `updated_at` timestamps. Foreign keys with proper constraints.

**Users Table:**

- Primary key: `id` (UUID)
- Unique: `email`
- Indexes: `email`, `role`

**Attendance Table:**

- Composite unique: `(user_id, date)`
- Indexes: `user_id`, `date`

**Leaves Table:**

- Indexes: `user_id`, `status`, `start_date`

**Salaries Table:**

- Composite unique: `(user_id, month, year)`
- Indexes: `user_id`, `month`, `year`

## Security Considerations

1. **JWT Configuration:**

   - Access token: 15 minutes expiry
   - Refresh token: 7 days expiry
   - Secure cookie storage for refresh tokens

2. **Password Policy:**

   - Minimum 8 characters
   - Bcrypt cost factor: 10

3. **Rate Limiting:**

   - Auth endpoints: 5 requests/minute
   - General API: 100 requests/minute

4. **Input Validation:**

   - All inputs validated using struct tags
   - Sanitization for user inputs

## Environment Variables

Required environment variables will be documented in `.env.example`:

- Database connection
- JWT secrets
- New Relic license key
- OpenAI API key
- SMTP configuration
- Server port and host

## Testing Strategy

1. **Unit Tests:** Services and repositories with mocked dependencies
2. **Integration Tests:** Full API endpoints with test database
3. **Test Coverage Target:** >80% for business logic
4. **Test Database:** Separate PostgreSQL instance for testing

## Implementation Order

Following the PRD's development plan, but with infrastructure setup first:

1. Foundation & Infrastructure
2. Authentication & RBAC
3. Employee Management
4. Attendance Management
5. Leave Management
6. Payroll Management
7. Document Management
8. Notification System
9. Dashboard
10. Audit Logging
11. AI HR Assistant
12. Observability
13. Frontend UI
14. Testing
15. Documentation