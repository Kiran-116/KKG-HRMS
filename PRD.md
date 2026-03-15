# HR Management System (HRMS)
## Product Requirements Document (PRD)

---

# 1. Overview

## Product Name
HR Management System (HRMS)

## Product Type
Internal company portal for managing employees, attendance, payroll, leave, and HR data.

## Goal
Build a modern HR platform that allows:

- Admins to manage employees, payroll, attendance, documents, and leave approvals.
- Employees to access their personal HR data (profile, attendance, salary, documents).
- System administrators to monitor system health using observability tools.
- Users to interact with HR data using an AI assistant.

---

# 2. Key Objectives

1. Implement secure **authentication and role-based access control**
2. Provide **admin dashboard for HR management**
3. Provide **employee self-service portal**
4. Implement **audit logging for traceability**
5. Integrate **New Relic for observability and monitoring**
6. Integrate **AI assistant for HR queries**
7. Design system using **clean architecture and scalable modules**

---

# 3. Target Users

## Admin
HR managers or administrators responsible for managing employee data.

Capabilities:
- Manage employees
- Approve leaves
- Monitor attendance
- Generate payslips
- Access analytics dashboard

## Employee
Regular users of the system.

Capabilities:
- View profile
- Check attendance
- Apply leave
- View salary and documents
- Use AI assistant for HR queries

---

# 4. Tech Stack

## Frontend
- React
- TypeScript
- TailwindCSS
- REST API integration

## Backend
- Golang
- Gin framework
- Clean architecture

## Database
- PostgreSQL

## Observability
- New Relic

## AI Integration
- OpenAI API

## Infrastructure
- Docker

---

# 5. System Architecture
React Frontend
|
REST API Layer
|
Golang Backend (Gin)
|
PostgreSQL Database
|
External Services
| |
New Relic OpenAI


---

# 6. Core Modules

1. Authentication
2. Role Based Access Control (RBAC)
3. Employee Management
4. Attendance Management
5. Leave Management
6. Payroll Management
7. Document Management
8. Notification System
9. Dashboard
10. Audit Logs
11. Observability
12. AI HR Assistant

---

# 7. Authentication & Authorization

## Features

- User registration
- User login
- Password hashing (bcrypt)
- JWT authentication
- Role-based access control

## Roles

### Admin
Full access to all HR features.

### Employee
Access only to personal data.

## APIs
POST /auth/register
POST /auth/login
GET /auth/me


---

# 8. Dashboard

## Admin Dashboard

Displays company overview.

Widgets:

- Total employees
- Present employees today
- Absent employees today
- Pending leave requests
- Recent activities
- Payroll summary

API:
GET /dashboard/admin


---

## Employee Dashboard

Displays personal HR information.

Widgets:

- Attendance summary
- Leave balance
- Salary summary
- Upcoming holidays
- Notifications

API:
GET /dashboard/employee


---

# 9. Employee Management

Admins can manage employees.

## Features

- Create employee
- Update employee
- View employee details
- Deactivate employee

## Employee Profile Fields

- Name
- Email
- Department
- Designation
- Joining date
- Salary

## APIs
GET /employees
POST /employees
PUT /employees/:id
GET /employees/:id
GET /employees/me


---

# 10. Attendance Management

Employees record daily attendance.

## Employee Features

- Check-in
- Check-out
- View attendance history

## Admin Features

- View all attendance
- Filter attendance by employee/date

## APIs
POST /attendance/checkin
POST /attendance/checkout
GET /attendance/me
GET /attendance


---

# 11. Leave Management

Employees can request leave.

## Employee Features

- Apply leave
- View leave history

## Admin Features

- Approve leave
- Reject leave
- View leave requests

## APIs
POST /leaves/apply
GET /leaves/me
GET /leaves
PUT /leaves/:id/approve
PUT /leaves/:id/reject


---

# 12. Payroll Management

Admins manage salary records.

## Features

Admin:

- Set salary
- Generate payslips
- View payroll history

Employee:

- View salary history
- Download payslips

## APIs
POST /salary
GET /salary/me
GET /salary/:userId


---

# 13. Document Management

Employees upload documents.

Examples:

- ID proof
- Offer letter
- Certificates
- Payslips

## Features

Employee:

- Upload document
- View documents

Admin:

- View employee documents

## APIs
POST /documents
GET /documents/me
GET /documents/:userId


---

# 14. Notification System

The system generates notifications for important events.

Examples:

- Leave approved/rejected
- Salary credited
- Document uploaded

## Channels

- In-app notifications
- Email notifications

## API
GET /notifications
PUT /notifications/:id/read


---

# 15. Audit Logging

All critical system actions must be logged.

Tracked events:

- User login
- Employee creation
- Salary update
- Leave approval
- Document upload

## Audit Log Fields
id
user_id
action
entity_type
entity_id
metadata
timestamp


## API
GET /audit-logs


---

# 16. Observability (New Relic)

System performance monitoring using New Relic.

Track:

- API latency
- Error rates
- Database query time
- System resource usage

Instrumentation must include:

- HTTP request tracing
- Error logging
- Performance monitoring

---

# 17. AI HR Assistant

An AI-powered assistant that answers HR-related questions.

Users can ask HR questions using natural language.

Examples:
How many leaves do I have left?
Show my attendance for last week
When was my last salary credited?


The AI converts natural language queries into system API calls.

## API
POST /ai/hr-assistant
Request:
{
"query": "How many leaves do I have left?"
}
Response:
{
"answer": "You have 8 leave days remaining."
}


---

# 18. Database Schema

## Users
id
name
email
password_hash
role
department
designation
created_at


---

## Attendance
id
user_id
date
check_in
check_out
status


---

## Leaves
id
user_id
start_date
end_date
reason
status
approved_by
created_at


---

## Salaries
id
user_id
base_salary
bonus
deductions
month
year
created_at

---

## Documents
id
user_id
file_url
document_type
uploaded_at


---

## Notifications
id
user_id
title
message
is_read
created_at

---

## Audit Logs
id
user_id
action
entity_type
entity_id
metadata
created_at


---

# 19. Security Requirements

- Password hashing using bcrypt
- JWT-based authentication
- Role-based authorization
- Secure document storage
- API validation
- Rate limiting for sensitive endpoints

---

# 20. Non Functional Requirements

## Performance
- API response < 300ms

## Scalability
- Modular architecture
- Service-based design

## Observability
- Monitoring via New Relic

## Maintainability
- Clean architecture
- Structured logging
- Modular services

---

# 21. Future Enhancements

- Real-time notifications using WebSockets
- AI-powered HR analytics
- Employee performance tracking
- Multi-organization support
- Mobile application

---

# 22. Suggested Project Folder Structure

Backend:
backend
│
├ controllers
├ services
├ repositories
├ middleware
├ models
├ routes
├ utils
└ main.go

Frontend:
frontend
│
├ components
├ pages
├ services
├ hooks
├ utils
└ App.tsx


---

# 23. Development Plan

## Phase 1
Authentication + RBAC

## Phase 2
Employee management

## Phase 3
Attendance system

## Phase 4
Leave management

## Phase 5
Payroll

## Phase 6
Dashboard

## Phase 7
Notifications + Audit logs

## Phase 8
AI HR assistant

## Phase 9
Observability integration

---

# 24. Success Criteria

The system is considered successful if:

- Admin can manage employees and HR operations
- Employees can access their HR data
- System tracks audit logs for critical actions
- New Relic provides performance monitoring
- AI assistant answers HR-related queries