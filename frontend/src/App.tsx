import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ToastContainer } from 'react-toastify';
import { AuthProvider } from './contexts/AuthContext';
import { WebSocketProvider } from './contexts/WebSocketContext';
import ProtectedRoute from './components/ProtectedRoute';
import Layout from './components/Layout';
import Login from './pages/Login';
import Register from './pages/Register';

// Dashboard
import Dashboard from './pages/Dashboard';

// Employees
import EmployeeList from './pages/employees/EmployeeList';
import EmployeeForm from './pages/employees/EmployeeForm';

// Attendance
import AttendanceCheckIn from './pages/attendance/AttendanceCheckIn';
import AttendanceHistory from './pages/attendance/AttendanceHistory';
import AdminAttendanceView from './pages/attendance/AdminAttendanceView';

// Leaves
import ApplyLeave from './pages/leaves/ApplyLeave';
import LeaveHistory from './pages/leaves/LeaveHistory';
import AdminLeaveApproval from './pages/leaves/AdminLeaveApproval';

// Payroll
import SalaryHistory from './pages/payroll/SalaryHistory';
import AdminPayroll from './pages/payroll/AdminPayroll';

// Documents
import DocumentList from './pages/documents/DocumentList';
import DocumentUpload from './pages/documents/DocumentUpload';

// Notifications
import NotificationList from './pages/notifications/NotificationList';

// AI Assistant
import HRAssistant from './pages/ai/HRAssistant';

// Audit
import AuditLogs from './pages/audit/AuditLogs';

// 404
import NotFound from './pages/NotFound';

function App() {
  return (
    <AuthProvider>
      <WebSocketProvider>
        <Router>
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            
            <Route
              path="/"
              element={
                <ProtectedRoute>
                  <Layout />
                </ProtectedRoute>
              }
            >
            {/* Dashboard */}
            <Route
              path="dashboard"
              element={
                <ProtectedRoute>
                  <Dashboard />
                </ProtectedRoute>
              }
            />

            {/* Employees - Admin Only */}
            <Route
              path="employees"
              element={
                <ProtectedRoute requireAdmin>
                  <EmployeeList />
                </ProtectedRoute>
              }
            />
            <Route
              path="employees/new"
              element={
                <ProtectedRoute requireAdmin>
                  <EmployeeForm />
                </ProtectedRoute>
              }
            />
            <Route
              path="employees/:id"
              element={
                <ProtectedRoute requireAdmin>
                  <EmployeeForm />
                </ProtectedRoute>
              }
            />

            {/* Attendance */}
            <Route
              path="attendance"
              element={
                <ProtectedRoute requireAdmin>
                  <AdminAttendanceView />
                </ProtectedRoute>
              }
            />
            <Route
              path="attendance/checkin"
              element={
                <ProtectedRoute>
                  <AttendanceCheckIn />
                </ProtectedRoute>
              }
            />
            <Route
              path="attendance/history"
              element={
                <ProtectedRoute>
                  <AttendanceHistory />
                </ProtectedRoute>
              }
            />

            {/* Leaves */}
            <Route
              path="leaves"
              element={
                <ProtectedRoute requireAdmin>
                  <AdminLeaveApproval />
                </ProtectedRoute>
              }
            />
            <Route
              path="leaves/me"
              element={
                <ProtectedRoute>
                  <LeaveHistory />
                </ProtectedRoute>
              }
            />
            <Route
              path="leaves/apply"
              element={
                <ProtectedRoute>
                  <ApplyLeave />
                </ProtectedRoute>
              }
            />

            {/* Payroll */}
            <Route
              path="payroll"
              element={
                <ProtectedRoute requireAdmin>
                  <AdminPayroll />
                </ProtectedRoute>
              }
            />
            <Route
              path="salary/me"
              element={
                <ProtectedRoute>
                  <SalaryHistory />
                </ProtectedRoute>
              }
            />

            {/* Documents */}
            <Route
              path="documents"
              element={
                <ProtectedRoute>
                  <DocumentList />
                </ProtectedRoute>
              }
            />
            <Route
              path="documents/upload"
              element={
                <ProtectedRoute>
                  <DocumentUpload />
                </ProtectedRoute>
              }
            />

            {/* Notifications */}
            <Route
              path="notifications"
              element={
                <ProtectedRoute>
                  <NotificationList />
                </ProtectedRoute>
              }
            />

            {/* AI Assistant */}
            <Route
              path="ai-assistant"
              element={
                <ProtectedRoute>
                  <HRAssistant />
                </ProtectedRoute>
              }
            />

            {/* Audit Logs - Admin Only */}
            <Route
              path="audit-logs"
              element={
                <ProtectedRoute requireAdmin>
                  <AuditLogs />
                </ProtectedRoute>
              }
            />
          </Route>

          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          
          {/* 404 - Catch all unmatched routes */}
          <Route path="*" element={<NotFound />} />
          </Routes>
        </Router>
        <ToastContainer
          position="top-right"
          autoClose={3000}
          hideProgressBar={false}
          newestOnTop={false}
          closeOnClick
          rtl={false}
          pauseOnFocusLoss
          draggable
          pauseOnHover
          theme="light"
        />
      </WebSocketProvider>
    </AuthProvider>
  );
}

export default App;
