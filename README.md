# HR Management System (HRMS)

A production-grade HR Management System built with React, TypeScript, Golang, and PostgreSQL.

## Features

- **Authentication & Authorization**: JWT-based authentication with role-based access control (Admin/Employee)
- **Employee Management**: CRUD operations for employee data
- **Attendance Management**: Check-in/check-out tracking
- **Leave Management**: Leave application and approval workflow
- **Payroll Management**: Salary records and payslip generation
- **Document Management**: File upload and storage
- **Notification System**: In-app and email notifications
- **Dashboard**: Admin and employee dashboards with analytics
- **Audit Logging**: Comprehensive audit trail
- **AI HR Assistant**: Natural language HR queries
- **Observability**: New Relic integration for monitoring

## Tech Stack

### Frontend
- React 18
- TypeScript
- TailwindCSS
- React Router

### Backend
- Golang 1.22+
- Gin Framework
- PostgreSQL
- JWT Authentication
- Clean Architecture

### Infrastructure
- Docker & Docker Compose
- PostgreSQL Database

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.22+ (for local development)
- Node.js 18+ (for local development)

### Environment Variables

Copy `.env.example` to `.env` and configure:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=hrms

# JWT
JWT_ACCESS_SECRET=your-secret-key
JWT_REFRESH_SECRET=your-refresh-secret

# OpenAI (optional)
OPENAI_API_KEY=your-api-key

# SMTP (optional)
SMTP_ENABLED=false
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email
SMTP_PASSWORD=your-password

# New Relic (optional)
# Enable to instrument incoming Gin requests.
NEW_RELIC_ENABLED=false
NEW_RELIC_APP_NAME=HRMS
# Set your New Relic license key via env var (recommended for security).
NEW_RELIC_LICENSE_KEY=YOUR_NEW_RELIC_LICENSE_KEY

# AI Monitoring (optional)
NEW_RELIC_AI_MONITORING_ENABLED=false
NEW_RELIC_CUSTOM_INSIGHTS_EVENTS_MAX_SAMPLES=10000
```

### Running with Docker

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Local Development

#### Backend

```bash
cd backend
go mod download
go run main.go
```

#### Frontend

```bash
cd frontend
npm install
npm start
```

## API Documentation

### Authentication

- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user
- `GET /api/auth/me` - Get current user

### Employees

- `GET /api/employees` - List employees (Admin)
- `POST /api/employees` - Create employee (Admin)
- `GET /api/employees/:id` - Get employee details
- `PUT /api/employees/:id` - Update employee (Admin)
- `GET /api/employees/me` - Get own profile

### Attendance

- `POST /api/attendance/checkin` - Check in
- `POST /api/attendance/checkout` - Check out
- `GET /api/attendance/me` - Get own attendance
- `GET /api/attendance` - Get all attendance (Admin)

### Leaves

- `POST /api/leaves/apply` - Apply for leave
- `GET /api/leaves/me` - Get own leaves
- `GET /api/leaves` - Get all leaves (Admin)
- `PUT /api/leaves/:id/approve` - Approve leave (Admin)
- `PUT /api/leaves/:id/reject` - Reject leave (Admin)

### Payroll

- `POST /api/salary` - Create salary record (Admin)
- `GET /api/salary/me` - Get own salary history
- `GET /api/salary/:userId` - Get employee salary (Admin)

### Documents

- `POST /api/documents` - Upload document
- `GET /api/documents/me` - Get own documents
- `GET /api/documents/:userId` - Get employee documents (Admin)
- `DELETE /api/documents/:id` - Delete document

### Notifications

- `GET /api/notifications` - Get notifications
- `GET /api/notifications/unread-count` - Get unread count
- `PUT /api/notifications/:id/read` - Mark as read

### Dashboard

- `GET /api/dashboard/admin` - Admin dashboard (Admin)
- `GET /api/dashboard/employee` - Employee dashboard

### Audit Logs

- `GET /api/audit-logs` - Get audit logs (Admin)

### AI Assistant

- `POST /api/ai/hr-assistant` - Process HR query

## Database Schema

The system uses PostgreSQL with the following main tables:

- `users` - User accounts
- `attendance` - Attendance records
- `leaves` - Leave requests
- `salaries` - Salary records
- `documents` - Document metadata
- `notifications` - User notifications
- `audit_logs` - Audit trail

## Architecture

The system follows clean architecture principles:

```
Frontend (React)
    ↓
REST API (Gin)
    ↓
Services (Business Logic)
    ↓
Repositories (Data Access)
    ↓
PostgreSQL Database
```

## Security

- JWT token-based authentication
- Password hashing with bcrypt
- Role-based access control
- Input validation
- Rate limiting on sensitive endpoints

## Testing

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test
```

## License

MIT
