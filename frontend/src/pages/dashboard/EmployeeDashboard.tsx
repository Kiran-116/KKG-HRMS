import React, { useEffect, useState } from 'react';
import { dashboardService, EmployeeDashboard } from '../../services/dashboardService';
import StatCard from '../../components/dashboard/StatCard';
import AttendancePieChart from '../../components/dashboard/AttendancePieChart';
import AttendanceChart from '../../components/dashboard/AttendanceChart';

const EmployeeDashboardPage: React.FC = () => {
  const [dashboard, setDashboard] = useState<EmployeeDashboard | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadDashboard();
  }, []);

  const loadDashboard = async () => {
    try {
      const data = await dashboardService.getEmployeeDashboard();
      setDashboard(data);
    } catch (error) {
      console.error('Failed to load dashboard:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="text-center py-12">Loading dashboard...</div>;
  }

  if (!dashboard) {
    return <div className="text-center py-12 text-red-600">Failed to load dashboard</div>;
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">My Dashboard</h1>
        <p className="text-gray-600 mt-2">Your personal HR information</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatCard
          title="Present Days (This Month)"
          value={dashboard.attendance_summary?.present_days ?? 0}
          icon="✅"
          color="green"
        />
        <StatCard
          title="Absent Days (This Month)"
          value={dashboard.attendance_summary?.absent_days ?? 0}
          icon="❌"
          color="red"
        />
        <StatCard
          title="Leave Balance"
          value={dashboard.leave_balance ?? 0}
          icon="📅"
          color="blue"
        />
        <StatCard
          title="Unread Notifications"
          value={dashboard.unread_notifications ?? 0}
          icon="🔔"
          color="yellow"
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Attendance Overview (This Month)</h2>
          {dashboard.attendance_summary ? (
            <AttendancePieChart 
              present={dashboard.attendance_summary.present_days}
              absent={dashboard.attendance_summary.absent_days}
            />
          ) : (
            <p className="text-sm text-gray-500">No attendance data available</p>
          )}
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Attendance Trend (Last 7 Days)</h2>
          <AttendanceChart data={dashboard.attendance_trend || []} />
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Salary Summary</h2>
          {dashboard.salary_summary ? (
            <div className="space-y-2">
              <p className="text-sm text-gray-600">
                {dashboard.salary_summary.month}/{dashboard.salary_summary.year}
              </p>
              <p className="text-2xl font-bold text-gray-900">
                ${dashboard.salary_summary.net_salary.toLocaleString()}
              </p>
              <div className="mt-4 space-y-1 text-sm">
                <p className="text-gray-600">
                  Base: ${dashboard.salary_summary.base_salary.toLocaleString()}
                </p>
                <p className="text-gray-600">
                  Bonus: ${dashboard.salary_summary.bonus.toLocaleString()}
                </p>
                <p className="text-gray-600">
                  Deductions: ${dashboard.salary_summary.deductions.toLocaleString()}
                </p>
              </div>
            </div>
          ) : (
            <p className="text-sm text-gray-500">No salary records found</p>
          )}
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Attendance Summary</h2>
          {dashboard.attendance_summary ? (
            <div className="space-y-2">
              <p className="text-sm text-gray-600">
                Month: {dashboard.attendance_summary.month}/{dashboard.attendance_summary.year}
              </p>
              <div className="mt-4 space-y-2">
                <div className="flex justify-between">
                  <span className="text-gray-600">Present:</span>
                  <span className="font-semibold">{dashboard.attendance_summary.present_days} days</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Absent:</span>
                  <span className="font-semibold">{dashboard.attendance_summary.absent_days} days</span>
                </div>
              </div>
            </div>
          ) : (
            <p className="text-sm text-gray-500">No attendance data available</p>
          )}
        </div>
      </div>
    </div>
  );
};

export default EmployeeDashboardPage;
